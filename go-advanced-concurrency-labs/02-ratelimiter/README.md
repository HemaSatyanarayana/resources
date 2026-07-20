# 02 — Token-Bucket Rate Limiter

Every client of a downstream API needs one of these. The **token bucket** is the
canonical algorithm: a bucket holds up to `burst` tokens; each request spends
one; tokens drip back in at a fixed rate. It permits short bursts (up to the
bucket size) while capping the long-run average.

The elegant Go insight: **a buffered channel *is* a token bucket.** Its capacity
is the burst; the number of elements in it is the tokens currently available;
send/receive are already goroutine-safe. You write almost no locking.

## The system

```
                  refill goroutine (ticker)
                        │  every `interval`: add up to `refill`
                        ▼            (non-blocking — drop if full)
 Allow()  ── receive ──▶ [ ●●●○○ ]  capacity = burst
 Wait()   ── receive ──▶  bucket
```

- **`Allow()`** — a *non-blocking* receive. Got a token? Return `true`. Empty?
  Return `false` instantly. This is the "try, and give up if rate-limited" path.
- **`Wait(ctx)`** — a *blocking* receive, raced against `ctx.Done()`. Either you
  get a token (return `nil`) or the context expires first (return `ctx.Err()`).
  This is the "back-pressure: slow down and wait your turn" path.
- **Refill goroutine** — a `time.Ticker`; on each tick it tops up the bucket with
  a **non-blocking** send so it never blocks when the bucket is already full.
- **`Stop()`** — closes a `done` channel to end the goroutine; wrapped in a
  `sync.Once` so calling it twice doesn't panic on a double `close`.

### The two subtleties that bite people

1. **Refill must not block.** If you do a plain `bucket <- struct{}{}` and the
   bucket is full, the ticker goroutine blocks forever and stops ticking. Always
   `select { case bucket <- struct{}{}: default: }`.
2. **No `sync.Mutex` needed for the count.** Resist adding a counter + lock — the
   channel already gives you an atomic, blocking-capable, concurrency-safe count.
   The only lock-like thing you need is the `sync.Once` in `Stop`.

## Your task

Implement, from scratch, in `ratelimiter.go` (package `ratelimiter`):

```go
type Limiter struct { /* your fields */ }

func NewTokenBucket(burst, refill int, interval time.Duration) *Limiter
func (l *Limiter) Allow() bool
func (l *Limiter) Wait(ctx context.Context) error
func (l *Limiter) Stop()
```

Bucket starts **full**. `burst < 1` or `refill < 1` → treat as `1`. `Stop` must be
idempotent and must actually end the refill goroutine (no leaks).

## Run

```bash
go test -race -v ./02-ratelimiter/
```

## Hints

- Pre-fill: `for i := 0; i < burst; i++ { bucket <- struct{}{} }`.
- Refill tick: loop `refill` times doing a non-blocking send; `break` out of the
  loop the first time the default case fires (bucket is full).
- `Wait`: `select { case <-l.bucket: return nil; case <-ctx.Done(): return ctx.Err() }`.
- `Stop`: `l.once.Do(func() { close(l.done) })`.
- The refill goroutine's loop: `select { case <-l.done: return; case <-ticker.C: /* top up */ }`.

<details>
<summary>Reference solution</summary>

```go
package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type Limiter struct {
	bucket chan struct{}
	done   chan struct{}
	once   sync.Once
}

func NewTokenBucket(burst, refill int, interval time.Duration) *Limiter {
	if burst < 1 {
		burst = 1
	}
	if refill < 1 {
		refill = 1
	}
	l := &Limiter{
		bucket: make(chan struct{}, burst),
		done:   make(chan struct{}),
	}
	for i := 0; i < burst; i++ { // start full
		l.bucket <- struct{}{}
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-l.done:
				return
			case <-ticker.C:
				for i := 0; i < refill; i++ {
					select {
					case l.bucket <- struct{}{}:
					default: // bucket full — stop topping up
						i = refill
					}
				}
			}
		}
	}()
	return l
}

func (l *Limiter) Allow() bool {
	select {
	case <-l.bucket:
		return true
	default:
		return false
	}
}

func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-l.bucket:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *Limiter) Stop() {
	l.once.Do(func() { close(l.done) })
}
```

</details>
