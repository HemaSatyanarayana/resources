// Package ratelimit drills a classic backend building block: a concurrency-safe
// token-bucket rate limiter, plus the HTTP middleware that returns 429 when a
// caller exceeds their rate.
package ratelimit

import (
	"net/http"
	"sync"
	"time"
)

// TokenBucket is a thread-safe token-bucket limiter. It holds up to `capacity`
// tokens and refills at `ratePerSec` tokens per second. Each allowed request
// consumes one token.
//
// The `now` field lets tests inject a fake clock; production code uses time.Now.
type TokenBucket struct {
	mu         sync.Mutex
	tokens     float64
	capacity   float64
	ratePerSec float64
	last       time.Time
	now        func() time.Time
}

// NewTokenBucket returns a bucket that starts FULL (capacity tokens), refilling
// at ratePerSec tokens per second.
func NewTokenBucket(capacity int, ratePerSec float64) *TokenBucket {
	panic("TODO: implement NewTokenBucket")
}

// refill adds tokens for the time elapsed since `last`, capped at capacity.
// Callers must hold b.mu. Implement this helper first — Allow uses it.
func (b *TokenBucket) refill() {
	panic("TODO: implement refill")
}

// Allow reports whether one request may proceed right now. If a token is
// available it consumes one and returns true; otherwise it returns false.
// It must be safe to call from many goroutines at once.
func (b *TokenBucket) Allow() bool {
	panic("TODO: implement Allow")
}

// Middleware wraps next so that requests exceeding the limiter get status 429
// with body "rate limited"; allowed requests pass through.
func Middleware(b *TokenBucket, next http.Handler) http.Handler {
	panic("TODO: implement Middleware")
}
