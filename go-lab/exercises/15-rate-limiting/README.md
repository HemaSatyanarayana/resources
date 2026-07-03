# 15 — Rate Limiting

Rate limiting protects a backend from abuse and overload. The **token bucket** is the workhorse algorithm: a bucket holds up to N tokens, refills at a steady rate, and each request spends one token — spare capacity ("burst") is allowed, sustained rate is capped. You'll build a concurrency-safe limiter and wire it into HTTP middleware that returns **429 Too Many Requests**.

## Concepts

- **Token bucket:** `capacity` = max burst; `ratePerSec` = sustained rate. Refill lazily — compute how many tokens to add from the elapsed time on each call, rather than running a background ticker.
- **Lazy refill formula:** `tokens = min(capacity, tokens + elapsed.Seconds() * ratePerSec)`, then update `last`.
- **Thread safety:** an HTTP server calls `Allow()` from many goroutines at once. A `sync.Mutex` around read-modify-write of `tokens` prevents over-issuing. (The concurrency test proves exactly 100 of 500 racing requests get through.)
- **Injectable clock:** production uses `time.Now`; tests swap in a fake `now func() time.Time` so behavior is deterministic — no `time.Sleep` in tests.
- **HTTP 429** is the standard "slow down" status; often paired with a `Retry-After` header (not required here).

## Your task

Implement everything in [`ratelimit.go`](ratelimit.go):

| Function | Skill |
|----------|-------|
| `NewTokenBucket` | Start full, wire the real clock |
| `refill` | Lazy time-based refill, capped |
| `Allow` | Locked read-modify-write |
| `Middleware` | 429 vs pass-through |

## Run

```bash
go test -race -v ./exercises/15-rate-limiting/
```

## Hints

- In `NewTokenBucket`, set `tokens` and `capacity` to `float64(capacity)`, `last` to `time.Now()`, and `now` to `time.Now`.
- `refill` (caller holds the lock): `elapsed := b.now().Sub(b.last); b.last = b.now(); b.tokens = math.Min(b.capacity, b.tokens + elapsed.Seconds()*b.ratePerSec)`.
- `Allow`: `b.mu.Lock(); defer b.mu.Unlock(); b.refill(); if b.tokens >= 1 { b.tokens--; return true }; return false`.
- `Middleware`: if `!b.Allow()`, `http.Error(w, "rate limited", http.StatusTooManyRequests)` and `return`; else call `next.ServeHTTP`.

<details>
<summary>Reference solution</summary>

```go
package ratelimit

import (
	"math"
	"net/http"
	"time"
)

func NewTokenBucket(capacity int, ratePerSec float64) *TokenBucket {
	return &TokenBucket{
		tokens:     float64(capacity),
		capacity:   float64(capacity),
		ratePerSec: ratePerSec,
		last:       time.Now(),
		now:        time.Now,
	}
}

func (b *TokenBucket) refill() {
	now := b.now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now
	b.tokens = math.Min(b.capacity, b.tokens+elapsed*b.ratePerSec)
}

func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

func Middleware(b *TokenBucket, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !b.Allow() {
			http.Error(w, "rate limited", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

</details>
