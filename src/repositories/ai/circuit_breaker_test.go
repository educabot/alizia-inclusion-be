package ai_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	ai "github.com/educabot/alizia-inclusion-be/src/repositories/ai"
)

// fakeAIClient is a controllable test double for providers.AIClient.
// Set generateErr / chatErr to make calls fail; leave nil for success.
// calls tracks how many times any method was invoked.
type fakeAIClient struct {
	generateErr error
	chatErr     error
	calls       int
}

func (f *fakeAIClient) Generate(_ context.Context, _ providers.GenerateParams) (string, error) {
	f.calls++
	if f.generateErr != nil {
		return "", f.generateErr
	}
	return "ok", nil
}

func (f *fakeAIClient) Chat(_ context.Context, _ []providers.ChatMessage) (*providers.ChatResponse, error) {
	f.calls++
	if f.chatErr != nil {
		return nil, f.chatErr
	}
	return &providers.ChatResponse{Content: "ok"}, nil
}

func (f *fakeAIClient) ChatWithTools(_ context.Context, _ []providers.ChatMessage, _ []providers.ToolDefinition) (*providers.ChatResponse, error) {
	f.calls++
	if f.chatErr != nil {
		return nil, f.chatErr
	}
	return &providers.ChatResponse{Content: "ok"}, nil
}

// fakeClock is a controllable clock that can be advanced manually.
type fakeClock struct {
	t time.Time
}

func (fc *fakeClock) Now() time.Time { return fc.t }

func (fc *fakeClock) Advance(d time.Duration) { fc.t = fc.t.Add(d) }

const (
	threshold = 3
	cooldown  = 10 * time.Second
)

func newBreaker(client providers.AIClient, clock *fakeClock) *ai.CircuitBreaker {
	return ai.NewCircuitBreaker(client, threshold, cooldown, clock.Now)
}

// TestCircuitBreaker_PassThrough_WhenClosed verifies that calls reach the
// wrapped client and the response is returned unchanged while the circuit is
// closed.
func TestCircuitBreaker_PassThrough_WhenClosed(t *testing.T) {
	fake := &fakeAIClient{}
	clock := &fakeClock{t: time.Now()}
	cb := newBreaker(fake, clock)

	result, err := cb.Generate(context.Background(), providers.GenerateParams{UserPrompt: "hello"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "ok" {
		t.Errorf("expected result %q, got %q", "ok", result)
	}
	if fake.calls != 1 {
		t.Errorf("expected 1 call, got %d", fake.calls)
	}
}

// TestCircuitBreaker_OpensAfterThresholdFailures verifies that after
// failureThreshold consecutive failures the circuit opens and subsequent calls
// are short-circuited (the wrapped client is NOT called again).
func TestCircuitBreaker_OpensAfterThresholdFailures(t *testing.T) {
	upstreamErr := fmt.Errorf("upstream down")
	fake := &fakeAIClient{generateErr: upstreamErr}
	clock := &fakeClock{t: time.Now()}
	cb := newBreaker(fake, clock)

	for i := 0; i < threshold; i++ {
		if _, err := cb.Generate(context.Background(), providers.GenerateParams{}); err == nil {
			t.Fatalf("call %d: expected error, got nil", i+1)
		}
	}
	callsAfterOpen := fake.calls

	_, err := cb.Generate(context.Background(), providers.GenerateParams{})

	if err == nil {
		t.Fatal("expected error when circuit open")
	}
	if !errors.Is(err, providers.ErrServiceUnavailable) {
		t.Errorf("expected ErrServiceUnavailable, got: %v", err)
	}
	if fake.calls != callsAfterOpen {
		t.Errorf("wrapped client should not be called when circuit is open: got %d, want %d", fake.calls, callsAfterOpen)
	}
}

// TestCircuitBreaker_SuccessResetsFailureCounter verifies that a success resets
// the consecutive-failure counter so that the circuit does not open prematurely
// when failures stay below the threshold between successes.
func TestCircuitBreaker_SuccessResetsFailureCounter(t *testing.T) {
	fake := &fakeAIClient{generateErr: fmt.Errorf("transient")}
	clock := &fakeClock{t: time.Now()}
	cb := newBreaker(fake, clock)

	for i := 0; i < threshold-1; i++ {
		_, _ = cb.Generate(context.Background(), providers.GenerateParams{})
	}

	fake.generateErr = nil
	if _, err := cb.Generate(context.Background(), providers.GenerateParams{}); err != nil {
		t.Fatalf("call after partial failures should succeed: %v", err)
	}

	fake.generateErr = fmt.Errorf("transient again")
	for i := 0; i < threshold-1; i++ {
		_, _ = cb.Generate(context.Background(), providers.GenerateParams{})
	}

	fake.generateErr = nil
	_, finalErr := cb.Generate(context.Background(), providers.GenerateParams{})

	if finalErr != nil {
		t.Fatalf("circuit must remain closed when failures never reach threshold consecutively: %v", finalErr)
	}
}

// TestCircuitBreaker_TrialCallClosesCircuit verifies that after the cooldown
// elapses, one trial call is allowed; if it succeeds the circuit closes and
// normal traffic resumes.
func TestCircuitBreaker_TrialCallClosesCircuit(t *testing.T) {
	fake := &fakeAIClient{chatErr: fmt.Errorf("down")}
	clock := &fakeClock{t: time.Now()}
	cb := newBreaker(fake, clock)

	for i := 0; i < threshold; i++ {
		_, _ = cb.Chat(context.Background(), nil)
	}
	callsWhenOpen := fake.calls

	if _, err := cb.Chat(context.Background(), nil); !errors.Is(err, providers.ErrServiceUnavailable) {
		t.Fatalf("expected ErrServiceUnavailable while open, got: %v", err)
	}
	if fake.calls != callsWhenOpen {
		t.Errorf("no calls expected while open: got %d, want %d", fake.calls, callsWhenOpen)
	}

	clock.Advance(cooldown)
	fake.chatErr = nil
	if _, err := cb.Chat(context.Background(), nil); err != nil {
		t.Fatalf("trial call should succeed: %v", err)
	}

	callsAfterClose := fake.calls
	if _, err := cb.Chat(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error after close: %v", err)
	}
	if fake.calls != callsAfterClose+1 {
		t.Errorf("circuit must be closed after successful trial: got %d, want %d", fake.calls, callsAfterClose+1)
	}
}

// TestCircuitBreaker_FailingTrialReopensCircuit verifies that a trial call that
// fails re-opens the circuit and restarts the cooldown.
func TestCircuitBreaker_FailingTrialReopensCircuit(t *testing.T) {
	fake := &fakeAIClient{chatErr: fmt.Errorf("down")}
	clock := &fakeClock{t: time.Now()}
	cb := newBreaker(fake, clock)

	for i := 0; i < threshold; i++ {
		_, _ = cb.Chat(context.Background(), nil)
	}

	clock.Advance(cooldown)
	if _, err := cb.Chat(context.Background(), nil); err == nil {
		t.Fatal("trial call must propagate the upstream error")
	}

	callsAfterTrial := fake.calls
	if _, err := cb.Chat(context.Background(), nil); !errors.Is(err, providers.ErrServiceUnavailable) {
		t.Errorf("circuit must be open after failing trial, got: %v", err)
	}
	if fake.calls != callsAfterTrial {
		t.Errorf("wrapped client must not be called when re-opened: got %d, want %d", fake.calls, callsAfterTrial)
	}

	clock.Advance(cooldown)
	fake.chatErr = nil
	if _, err := cb.Chat(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error after second trial: %v", err)
	}
}
