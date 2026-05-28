package middleware

import (
	"net/http"
	"testing"
	"time"

	"github.com/educabot/team-ai-toolkit/web"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// newFakeClock returns a clock function whose time can be advanced by mutating
// the pointed-at value. Safe to use across a single test (not concurrent).
func newFakeClock(t *testing.T) *time.Time {
	t.Helper()
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	return &now
}

func clockFn(ptr *time.Time) func() time.Time {
	return func() time.Time { return *ptr }
}

// TestRateLimiter_AllowsUpToCapacityThenBlocks verifies that a bucket with
// capacity N allows exactly N requests in the same window and blocks the (N+1)th.
func TestRateLimiter_AllowsUpToCapacityThenBlocks(t *testing.T) {
	const capacity = 3
	ts := newFakeClock(t)
	rl := NewRateLimiter(capacity, clockFn(ts))
	key := "org-a"

	for i := range capacity {
		assert.True(t, rl.Allow(key), "request %d/%d: expected allowed, got blocked", i+1, capacity)
	}

	assert.False(t, rl.Allow(key), "expected request to be blocked after capacity exhausted")
}

// TestRateLimiter_RefillsAfterElapsedTime verifies that advancing the fake clock
// by a full hour causes the bucket to refill to capacity.
func TestRateLimiter_RefillsAfterElapsedTime(t *testing.T) {
	const capacity = 5
	ts := newFakeClock(t)
	rl := NewRateLimiter(capacity, clockFn(ts))
	key := "org-b"

	for range capacity {
		rl.Allow(key)
	}
	assert.False(t, rl.Allow(key), "expected blocked after draining")

	*ts = ts.Add(time.Hour)

	for i := range capacity {
		assert.True(t, rl.Allow(key), "after refill: request %d/%d expected allowed", i+1, capacity)
	}

	assert.False(t, rl.Allow(key), "expected blocked after re-draining refilled bucket")
}

// TestRateLimiter_PartialRefill verifies that advancing the clock by half an hour
// refills exactly half the bucket tokens.
func TestRateLimiter_PartialRefill(t *testing.T) {
	const capacity = 10
	ts := newFakeClock(t)
	rl := NewRateLimiter(capacity, clockFn(ts))
	key := "org-partial"

	for range capacity {
		rl.Allow(key)
	}

	*ts = ts.Add(30 * time.Minute)

	const expected = 5
	for i := range expected {
		assert.True(t, rl.Allow(key), "partial refill: request %d/%d expected allowed", i+1, expected)
	}
	assert.False(t, rl.Allow(key), "expected blocked after consuming partially-refilled tokens")
}

// TestRateLimiter_DifferentKeysHaveIndependentBuckets verifies that exhausting
// the bucket for one key does not affect a different key's bucket.
func TestRateLimiter_DifferentKeysHaveIndependentBuckets(t *testing.T) {
	const capacity = 2
	ts := newFakeClock(t)
	rl := NewRateLimiter(capacity, clockFn(ts))
	keyA := "org-x"
	keyB := "org-y"

	for range capacity {
		rl.Allow(keyA)
	}

	assert.False(t, rl.Allow(keyA), "keyA: expected blocked after exhaustion")

	for i := range capacity {
		assert.True(t, rl.Allow(keyB), "keyB: request %d/%d expected allowed (independent bucket)", i+1, capacity)
	}
	assert.False(t, rl.Allow(keyB), "keyB: expected blocked after its own exhaustion")
}

// TestRateLimiter_UnlimitedWhenMaxPerHourIsZero verifies that maxPerHour <= 0
// puts the limiter in unlimited mode where Allow always returns true.
func TestRateLimiter_UnlimitedWhenMaxPerHourIsZero(t *testing.T) {
	ts := newFakeClock(t)
	rl := NewRateLimiter(0, clockFn(ts))
	key := "org-unlimited"

	const iterations = 1000
	for i := range iterations {
		assert.True(t, rl.Allow(key), "unlimited mode: request %d expected allowed", i+1)
	}
}

// TestRateLimiter_UnlimitedWhenMaxPerHourIsNegative verifies the same unlimited
// behaviour for negative maxPerHour values.
func TestRateLimiter_UnlimitedWhenMaxPerHourIsNegative(t *testing.T) {
	ts := newFakeClock(t)
	rl := NewRateLimiter(-10, clockFn(ts))
	key := "org-neg"

	for i := range 100 {
		assert.True(t, rl.Allow(key), "unlimited mode (negative): request %d expected allowed", i+1)
	}
}

// TestRateLimitMiddleware_AllowsRequestWithNoOrgID verifies that a request with
// no org context is allowed through (deferred to auth/tenant middleware).
func TestRateLimitMiddleware_AllowsRequestWithNoOrgID(t *testing.T) {
	interceptor := RateLimitMiddleware(10)
	req := web.NewMockRequest()

	resp := interceptor(req)

	assert.Equal(t, 0, resp.Status)
}

// TestRateLimitMiddleware_AllowsWithinLimit verifies that requests within the
// rate limit are passed through.
func TestRateLimitMiddleware_AllowsWithinLimit(t *testing.T) {
	interceptor := RateLimitMiddleware(100)
	orgID := uuid.New()
	req := web.NewMockRequest()
	req.Values[OrgIDKey] = orgID

	resp := interceptor(req)

	assert.Equal(t, 0, resp.Status)
}

// TestRateLimitMiddleware_BlocksWhenLimitExceeded verifies that once the org
// bucket is exhausted, subsequent requests receive HTTP 429.
func TestRateLimitMiddleware_BlocksWhenLimitExceeded(t *testing.T) {
	ts := newFakeClock(t)
	limiter := NewRateLimiter(1, clockFn(ts))

	orgID := uuid.New()

	interceptor := func(req web.Request) web.Response {
		oid := OrgID(req)
		if oid == uuid.Nil {
			return web.Response{}
		}
		if !limiter.Allow(oid.String()) {
			return web.Err(http.StatusTooManyRequests, "rate_limited", "rate limit exceeded")
		}
		return web.Response{}
	}

	req := web.NewMockRequest()
	req.Values[OrgIDKey] = orgID

	first := interceptor(req)
	second := interceptor(req)

	assert.Equal(t, 0, first.Status)
	assert.Equal(t, http.StatusTooManyRequests, second.Status)
}
