package ai_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
	ai "github.com/educabot/alizia-inclusion-be/src/repositories/ai"
)

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
	client := &mockproviders.MockAIClient{}
	client.On("Generate", mock.Anything, mock.Anything).Return("ok", nil)
	clock := &fakeClock{t: time.Now()}
	cb := newBreaker(client, clock)

	result, err := cb.Generate(context.Background(), providers.GenerateParams{UserPrompt: "hello"})

	require.NoError(t, err)
	assert.Equal(t, "ok", result)
	client.AssertNumberOfCalls(t, "Generate", 1)
}

// TestCircuitBreaker_OpensAfterThresholdFailures verifies that after
// failureThreshold consecutive failures the circuit opens and subsequent calls
// are short-circuited (the wrapped client is NOT called again).
func TestCircuitBreaker_OpensAfterThresholdFailures(t *testing.T) {
	upstreamErr := fmt.Errorf("upstream down")
	client := &mockproviders.MockAIClient{}
	client.On("Generate", mock.Anything, mock.Anything).Return("", upstreamErr)
	clock := &fakeClock{t: time.Now()}
	cb := newBreaker(client, clock)

	for i := 0; i < threshold; i++ {
		_, err := cb.Generate(context.Background(), providers.GenerateParams{})
		require.Error(t, err, "call %d: expected error, got nil", i+1)
	}
	callsAfterOpen := len(client.Calls)

	_, err := cb.Generate(context.Background(), providers.GenerateParams{})

	require.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrServiceUnavailable)
	assert.Len(t, client.Calls, callsAfterOpen, "wrapped client must not be called once open")
}

// TestCircuitBreaker_SuccessResetsFailureCounter verifies that a success resets
// the consecutive-failure counter so that the circuit does not open prematurely
// when failures stay below the threshold between successes.
func TestCircuitBreaker_SuccessResetsFailureCounter(t *testing.T) {
	client := &mockproviders.MockAIClient{}
	// Sequence: 2 failures, 1 success (resets), 2 failures, 1 success.
	// Failures never reach the threshold consecutively, so the circuit stays closed.
	client.On("Generate", mock.Anything, mock.Anything).Return("", fmt.Errorf("transient")).Times(threshold - 1)
	client.On("Generate", mock.Anything, mock.Anything).Return("ok", nil).Once()
	client.On("Generate", mock.Anything, mock.Anything).Return("", fmt.Errorf("transient again")).Times(threshold - 1)
	client.On("Generate", mock.Anything, mock.Anything).Return("ok", nil).Once()
	clock := &fakeClock{t: time.Now()}
	cb := newBreaker(client, clock)

	for i := 0; i < threshold-1; i++ {
		_, _ = cb.Generate(context.Background(), providers.GenerateParams{})
	}
	_, err := cb.Generate(context.Background(), providers.GenerateParams{})
	require.NoError(t, err, "call after partial failures should succeed")

	for i := 0; i < threshold-1; i++ {
		_, _ = cb.Generate(context.Background(), providers.GenerateParams{})
	}
	_, finalErr := cb.Generate(context.Background(), providers.GenerateParams{})

	assert.NoError(t, finalErr, "circuit must remain closed when failures never reach threshold consecutively")
	client.AssertExpectations(t)
}

// TestCircuitBreaker_TrialCallClosesCircuit verifies that after the cooldown
// elapses, one trial call is allowed; if it succeeds the circuit closes and
// normal traffic resumes.
func TestCircuitBreaker_TrialCallClosesCircuit(t *testing.T) {
	client := &mockproviders.MockAIClient{}
	// The first `threshold` calls fail and open the circuit; every later call that
	// reaches the client succeeds (the trial call and subsequent traffic).
	client.On("Chat", mock.Anything, mock.Anything).Return((*providers.ChatResponse)(nil), fmt.Errorf("down")).Times(threshold)
	client.On("Chat", mock.Anything, mock.Anything).Return(&providers.ChatResponse{Content: "ok"}, nil)
	clock := &fakeClock{t: time.Now()}
	cb := newBreaker(client, clock)

	for i := 0; i < threshold; i++ {
		_, _ = cb.Chat(context.Background(), nil)
	}
	callsWhenOpen := len(client.Calls)

	_, err := cb.Chat(context.Background(), nil)
	assert.ErrorIs(t, err, providers.ErrServiceUnavailable)
	assert.Len(t, client.Calls, callsWhenOpen, "wrapped client must not be called while open")

	clock.Advance(cooldown)
	_, err = cb.Chat(context.Background(), nil)
	require.NoError(t, err, "trial call should succeed")

	callsAfterClose := len(client.Calls)
	_, err = cb.Chat(context.Background(), nil)
	require.NoError(t, err, "unexpected error after close")
	assert.Len(t, client.Calls, callsAfterClose+1, "circuit must be closed after successful trial")
}

// TestCircuitBreaker_FailingTrialReopensCircuit verifies that a trial call that
// fails re-opens the circuit and restarts the cooldown.
func TestCircuitBreaker_FailingTrialReopensCircuit(t *testing.T) {
	client := &mockproviders.MockAIClient{}
	// First `threshold` calls fail (open the circuit) plus the failing trial call
	// after the cooldown; the final trial call succeeds.
	client.On("Chat", mock.Anything, mock.Anything).Return((*providers.ChatResponse)(nil), fmt.Errorf("down")).Times(threshold + 1)
	client.On("Chat", mock.Anything, mock.Anything).Return(&providers.ChatResponse{Content: "ok"}, nil)
	clock := &fakeClock{t: time.Now()}
	cb := newBreaker(client, clock)

	for i := 0; i < threshold; i++ {
		_, _ = cb.Chat(context.Background(), nil)
	}

	clock.Advance(cooldown)
	_, err := cb.Chat(context.Background(), nil)
	require.Error(t, err, "trial call must propagate the upstream error")

	callsAfterTrial := len(client.Calls)
	_, err = cb.Chat(context.Background(), nil)
	assert.ErrorIs(t, err, providers.ErrServiceUnavailable, "circuit must be open after failing trial")
	assert.Len(t, client.Calls, callsAfterTrial, "wrapped client must not be called when re-opened")

	clock.Advance(cooldown)
	_, err = cb.Chat(context.Background(), nil)
	assert.NoError(t, err, "unexpected error after second trial")
}
