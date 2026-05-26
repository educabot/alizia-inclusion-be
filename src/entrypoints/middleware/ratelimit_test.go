package middleware

import (
	"net/http"
	"testing"
	"time"

	"github.com/educabot/team-ai-toolkit/web"
	"github.com/google/uuid"
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
	// Arrange
	const capacity = 3
	ts := newFakeClock(t)
	rl := NewRateLimiter(capacity, clockFn(ts))
	key := "org-a"

	// Act + Assert — first N requests must be allowed
	for i := range capacity {
		if !rl.Allow(key) {
			t.Fatalf("request %d/%d: expected allowed, got blocked", i+1, capacity)
		}
	}

	// (N+1)th request in the same window must be blocked
	if rl.Allow(key) {
		t.Fatal("expected request to be blocked after capacity exhausted, but it was allowed")
	}
}

// TestRateLimiter_RefillsAfterElapsedTime verifies that advancing the fake clock
// by a full hour causes the bucket to refill to capacity.
func TestRateLimiter_RefillsAfterElapsedTime(t *testing.T) {
	// Arrange
	const capacity = 5
	ts := newFakeClock(t)
	rl := NewRateLimiter(capacity, clockFn(ts))
	key := "org-b"

	// Drain the bucket completely.
	for range capacity {
		rl.Allow(key)
	}
	if rl.Allow(key) {
		t.Fatal("expected blocked after draining, got allowed")
	}

	// Act — advance clock by exactly one hour (full refill).
	*ts = ts.Add(time.Hour)

	// Assert — bucket should be full again; all capacity requests must be allowed.
	for i := range capacity {
		if !rl.Allow(key) {
			t.Fatalf("after refill: request %d/%d expected allowed, got blocked", i+1, capacity)
		}
	}

	// One beyond capacity must block again.
	if rl.Allow(key) {
		t.Fatal("expected blocked after re-draining refilled bucket, got allowed")
	}
}

// TestRateLimiter_PartialRefill verifies that advancing the clock by half an hour
// refills exactly half the bucket tokens.
func TestRateLimiter_PartialRefill(t *testing.T) {
	// Arrange: capacity 10, drain all, advance by 30 min → expect 5 tokens back.
	const capacity = 10
	ts := newFakeClock(t)
	rl := NewRateLimiter(capacity, clockFn(ts))
	key := "org-partial"

	for range capacity {
		rl.Allow(key)
	}

	// Act — half an hour passes.
	*ts = ts.Add(30 * time.Minute)

	// Assert — exactly 5 more requests should be allowed.
	const expected = 5
	for i := range expected {
		if !rl.Allow(key) {
			t.Fatalf("partial refill: request %d/%d expected allowed, got blocked", i+1, expected)
		}
	}
	if rl.Allow(key) {
		t.Fatal("expected blocked after consuming partially-refilled tokens, got allowed")
	}
}

// TestRateLimiter_DifferentKeysHaveIndependentBuckets verifies that exhausting
// the bucket for one key does not affect a different key's bucket.
func TestRateLimiter_DifferentKeysHaveIndependentBuckets(t *testing.T) {
	// Arrange
	const capacity = 2
	ts := newFakeClock(t)
	rl := NewRateLimiter(capacity, clockFn(ts))
	keyA := "org-x"
	keyB := "org-y"

	// Act — drain keyA completely.
	for range capacity {
		rl.Allow(keyA)
	}

	// Assert — keyA is blocked but keyB is still open.
	if rl.Allow(keyA) {
		t.Fatal("keyA: expected blocked after exhaustion, got allowed")
	}

	for i := range capacity {
		if !rl.Allow(keyB) {
			t.Fatalf("keyB: request %d/%d expected allowed (independent bucket), got blocked", i+1, capacity)
		}
	}
	if rl.Allow(keyB) {
		t.Fatal("keyB: expected blocked after its own exhaustion, got allowed")
	}
}

// TestRateLimiter_UnlimitedWhenMaxPerHourIsZero verifies that maxPerHour <= 0
// puts the limiter in unlimited mode where Allow always returns true.
func TestRateLimiter_UnlimitedWhenMaxPerHourIsZero(t *testing.T) {
	// Arrange
	ts := newFakeClock(t)
	rl := NewRateLimiter(0, clockFn(ts))
	key := "org-unlimited"

	// Act + Assert — a large number of requests must all be allowed.
	const iterations = 1000
	for i := range iterations {
		if !rl.Allow(key) {
			t.Fatalf("unlimited mode: request %d expected allowed, got blocked", i+1)
		}
	}
}

// TestRateLimiter_UnlimitedWhenMaxPerHourIsNegative verifies the same unlimited
// behaviour for negative maxPerHour values.
func TestRateLimiter_UnlimitedWhenMaxPerHourIsNegative(t *testing.T) {
	// Arrange
	ts := newFakeClock(t)
	rl := NewRateLimiter(-10, clockFn(ts))
	key := "org-neg"

	// Act + Assert
	for i := range 100 {
		if !rl.Allow(key) {
			t.Fatalf("unlimited mode (negative): request %d expected allowed, got blocked", i+1)
		}
	}
}

// TestRateLimitMiddleware_AllowsRequestWithNoOrgID verifies that a request with
// no org context is allowed through (deferred to auth/tenant middleware).
func TestRateLimitMiddleware_AllowsRequestWithNoOrgID(t *testing.T) {
	// Arrange
	interceptor := RateLimitMiddleware(10)
	req := web.NewMockRequest()
	// No OrgID set — OrgID(req) returns uuid.Nil.

	// Act
	resp := interceptor(req)

	// Assert
	if resp.Status != 0 {
		t.Errorf("expected pass-through (status 0) for missing org, got %d", resp.Status)
	}
}

// TestRateLimitMiddleware_AllowsWithinLimit verifies that requests within the
// rate limit are passed through.
func TestRateLimitMiddleware_AllowsWithinLimit(t *testing.T) {
	// Arrange
	interceptor := RateLimitMiddleware(100)
	orgID := uuid.New()
	req := web.NewMockRequest()
	req.Values[OrgIDKey] = orgID

	// Act
	resp := interceptor(req)

	// Assert
	if resp.Status != 0 {
		t.Errorf("expected pass-through, got status %d", resp.Status)
	}
}

// TestRateLimitMiddleware_BlocksWhenLimitExceeded verifies that once the org
// bucket is exhausted, subsequent requests receive HTTP 429.
func TestRateLimitMiddleware_BlocksWhenLimitExceeded(t *testing.T) {
	// Arrange: limit of 1 per hour so the second call must fail.
	// We build a limiter directly and embed it via a closure to keep time injectable.
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

	// Act — first request consumes the only token.
	first := interceptor(req)

	// Act — second request in the same window should be blocked.
	second := interceptor(req)

	// Assert
	if first.Status != 0 {
		t.Errorf("first request: expected pass-through (0), got %d", first.Status)
	}
	if second.Status != http.StatusTooManyRequests {
		t.Errorf("second request: expected %d, got %d", http.StatusTooManyRequests, second.Status)
	}
}
