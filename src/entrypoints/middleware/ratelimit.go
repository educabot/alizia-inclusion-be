package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/educabot/team-ai-toolkit/web"
	"github.com/google/uuid"
)

// RateLimiter is an in-memory token-bucket rate limiter keyed by arbitrary strings.
// It is safe for concurrent use.
//
// Note: X-RateLimit-* headers are intentionally not emitted. The web.Response
// abstraction only exposes Status and Body — there is no mechanism to set custom
// HTTP response headers through this layer.
type RateLimiter struct {
	mu         sync.Mutex
	maxTokens  float64
	refillRate float64 // tokens per nanosecond
	now        func() time.Time
	buckets    map[string]*tokenBucket
}

type tokenBucket struct {
	tokens   float64
	lastSeen time.Time
}

// NewRateLimiter creates a RateLimiter allowing up to maxPerHour requests per key per hour.
// The now function is the time source — inject time.Now in production, a fake clock in tests.
//
// If maxPerHour <= 0, the limiter is in unlimited mode: Allow always returns true.
// Document this contract so callers can use 0 / negative values to disable rate limiting.
func NewRateLimiter(maxPerHour int, now func() time.Time) *RateLimiter {
	var refill float64
	var capacity float64
	if maxPerHour > 0 {
		capacity = float64(maxPerHour)
		refill = capacity / float64(time.Hour) // tokens per nanosecond
	}
	return &RateLimiter{
		maxTokens:  capacity,
		refillRate: refill,
		now:        now,
		buckets:    make(map[string]*tokenBucket),
	}
}

// Allow reports whether the request identified by key should be allowed through.
// It consumes one token from the bucket. If the bucket is empty, it returns false.
// When maxPerHour <= 0 (unlimited mode), Allow always returns true.
func (rl *RateLimiter) Allow(key string) bool {
	// Unlimited mode: maxPerHour was <= 0 at construction time.
	if rl.maxTokens <= 0 {
		return true
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := rl.now()

	b, exists := rl.buckets[key]
	if !exists {
		// First request for this key: start with a full bucket minus the token we're consuming.
		rl.buckets[key] = &tokenBucket{
			tokens:   rl.maxTokens - 1,
			lastSeen: now,
		}
		return true
	}

	// Refill based on elapsed time since last request.
	elapsed := now.Sub(b.lastSeen)
	if elapsed > 0 {
		refilled := float64(elapsed) * rl.refillRate
		b.tokens += refilled
		if b.tokens > rl.maxTokens {
			b.tokens = rl.maxTokens
		}
		b.lastSeen = now
	}

	if b.tokens < 1 {
		return false
	}

	b.tokens--
	return true
}

// RateLimitMiddleware returns a web.Interceptor that enforces per-org token-bucket
// rate limiting. The limiter is shared across all requests handled by the returned
// interceptor.
//
// If maxPerHour <= 0, the middleware is a no-op (unlimited).
// If the request has no org ID (uuid.Nil), the request is allowed through —
// auth/tenant middleware upstream is responsible for rejecting unauthenticated calls.
func RateLimitMiddleware(maxPerHour int) web.Interceptor {
	limiter := NewRateLimiter(maxPerHour, time.Now)

	return func(req web.Request) web.Response {
		orgID := OrgID(req)
		if orgID == uuid.Nil {
			// No org context yet — let auth/tenant middleware handle it.
			return web.Response{}
		}

		if !limiter.Allow(orgID.String()) {
			return web.Err(http.StatusTooManyRequests, "rate_limited", "rate limit exceeded")
		}

		return web.Response{}
	}
}
