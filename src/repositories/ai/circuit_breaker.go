package ai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type circuitState int

const (
	stateClosed circuitState = iota
	stateOpen
	stateHalfOpen
)

// CircuitBreaker wraps a providers.AIClient and applies the circuit breaker
// pattern to protect against cascading failures. It is a drop-in decorator
// that itself implements providers.AIClient.
type CircuitBreaker struct {
	client           providers.AIClient
	failureThreshold int
	cooldown         time.Duration
	now              func() time.Time

	mu       sync.Mutex
	state    circuitState
	failures int
	openedAt time.Time
}

// NewCircuitBreaker returns a CircuitBreaker wrapping client.
// failureThreshold consecutive failures open the circuit.
// cooldown is the minimum duration the circuit stays open before a trial call
// is allowed (half-open).
// now is the clock function used to determine elapsed time; pass nil to use
// time.Now.
func NewCircuitBreaker(
	client providers.AIClient,
	failureThreshold int,
	cooldown time.Duration,
	now func() time.Time,
) *CircuitBreaker {
	if now == nil {
		now = time.Now
	}
	return &CircuitBreaker{
		client:           client,
		failureThreshold: failureThreshold,
		cooldown:         cooldown,
		now:              now,
		state:            stateClosed,
	}
}

// gate checks whether the call should be allowed and transitions state as
// needed. It returns an error if the circuit is open and the cooldown has not
// yet elapsed. The caller must hold no lock when calling gate; gate acquires
// and releases mu internally.
func (cb *CircuitBreaker) gate() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case stateClosed:
		return nil
	case stateHalfOpen:
		// A trial call is already in flight (handled by the caller).
		return nil
	case stateOpen:
		if cb.now().Sub(cb.openedAt) >= cb.cooldown {
			cb.state = stateHalfOpen
			return nil
		}
		return fmt.Errorf("circuit open: %w", providers.ErrServiceUnavailable)
	}
	return nil
}

// record updates the circuit state based on whether the call succeeded or
// failed. It must be called after every pass-through call.
func (cb *CircuitBreaker) record(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err == nil {
		// Success: close the circuit and reset the failure counter.
		cb.state = stateClosed
		cb.failures = 0
		return
	}

	// Failure path.
	if cb.state == stateHalfOpen {
		// Trial call failed: re-open and restart the cooldown.
		cb.state = stateOpen
		cb.openedAt = cb.now()
		return
	}

	cb.failures++
	if cb.failures >= cb.failureThreshold {
		cb.state = stateOpen
		cb.openedAt = cb.now()
	}
}

// Generate passes through to the wrapped client unless the circuit is open.
func (cb *CircuitBreaker) Generate(ctx context.Context, params providers.GenerateParams) (string, error) {
	if err := cb.gate(); err != nil {
		return "", err
	}
	result, err := cb.client.Generate(ctx, params)
	cb.record(err)
	return result, err
}

// Chat passes through to the wrapped client unless the circuit is open.
func (cb *CircuitBreaker) Chat(ctx context.Context, messages []providers.ChatMessage) (*providers.ChatResponse, error) {
	if err := cb.gate(); err != nil {
		return nil, err
	}
	result, err := cb.client.Chat(ctx, messages)
	cb.record(err)
	return result, err
}

// ChatWithTools passes through to the wrapped client unless the circuit is open.
func (cb *CircuitBreaker) ChatWithTools(ctx context.Context, messages []providers.ChatMessage, tools []providers.ToolDefinition) (*providers.ChatResponse, error) {
	if err := cb.gate(); err != nil {
		return nil, err
	}
	result, err := cb.client.ChatWithTools(ctx, messages, tools)
	cb.record(err)
	return result, err
}
