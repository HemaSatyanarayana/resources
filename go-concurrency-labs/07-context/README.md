# 07 — Context

`context.Context` is the standard way to carry **cancellation**, **deadlines**,
and request-scoped values across API boundaries and goroutines. It's the
`done`-channel pattern from Lab 06, standardized and everywhere in Go's stdlib
(`net/http`, `database/sql`, gRPC…).

## Concepts

### What a Context gives you

```go
ctx.Done()  // <-chan struct{} — closed when the context is canceled
ctx.Err()   // nil while active; context.Canceled or context.DeadlineExceeded after
```

Watch `ctx.Done()` in a `select`, exactly like a `done` channel.

### Creating contexts

```go
ctx := context.Background()                             // root, never canceled

ctx, cancel := context.WithCancel(parent)               // cancel() to stop
defer cancel()

ctx, cancel := context.WithTimeout(parent, 2*time.Second) // auto-cancels after 2s
defer cancel()

ctx, cancel := context.WithDeadline(parent, someTime)
defer cancel()
```

**Always `defer cancel()`** — even for `WithTimeout`. `cancel` releases the
context's resources; skipping it leaks a timer/goroutine until the deadline.

### The cancellation idiom

```go
func work(ctx context.Context) error {
    select {
    case <-time.After(5 * time.Second):
        return nil                 // finished the work
    case <-ctx.Done():
        return ctx.Err()           // canceled or timed out — report why
    }
}
```

- Cancellation **propagates down**: canceling a parent cancels all children.
- A function that gets a `ctx` should return promptly when it's done, returning
  `ctx.Err()`.
- Convention: `ctx` is the **first** parameter, named `ctx`. Never store it in a
  struct; pass it explicitly.

### First-error fan-out (mini errgroup)

A common pattern: run N tasks concurrently, and if **any** fails, cancel the rest
and return that error. You derive a cancelable context, hand it to every task,
and the first failure calls `cancel()` so siblings watching `ctx.Done()` stop
early. (The real `golang.org/x/sync/errgroup` packages this up; here you build
it by hand.)

## Your task

Create `ctxwork.go` (package `ctxlab`) and implement, **from scratch**:

```go
// SleepOrCancel sleeps for d, but returns early if ctx is canceled first.
// Returns nil if it slept the full duration, or ctx.Err() if canceled.
func SleepOrCancel(ctx context.Context, d time.Duration) error

// FetchAll calls fetch for every url concurrently, passing ctx to each. It
// returns the results in URL ORDER on success. If any fetch returns an error,
// FetchAll cancels the remaining fetches and returns the first error observed
// (results is nil in that case).
func FetchAll(ctx context.Context, urls []string,
    fetch func(ctx context.Context, url string) (string, error)) ([]string, error)
```

## Run

```bash
go test -race -v ./07-context/
```

## Hints

- `SleepOrCancel`: `select { case <-time.After(d): return nil; case <-ctx.Done(): return ctx.Err() }`.
- `FetchAll`:
  - Derive `ctx, cancel := context.WithCancel(ctx)` and `defer cancel()`.
  - `results := make([]string, len(urls))`; each goroutine writes its own index.
  - Collect errors safely — a shared `error` needs a `sync.Mutex`, or use a
    buffered error channel and read the first. On first error, call `cancel()`.
  - `WaitGroup` to join; if any error, return `nil, firstErr`, else `results, nil`.
- Pass the **derived** `ctx` into each `fetch` so cancellation reaches them.

<details>
<summary>Reference solution</summary>

```go
package ctxlab

import (
	"context"
	"sync"
	"time"
)

func SleepOrCancel(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func FetchAll(ctx context.Context, urls []string,
	fetch func(ctx context.Context, url string) (string, error)) ([]string, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]string, len(urls))

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
	)

	wg.Add(len(urls))
	for i, url := range urls {
		go func() {
			defer wg.Done()
			v, err := fetch(ctx, url)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
					cancel() // stop the siblings
				}
				mu.Unlock()
				return
			}
			results[i] = v
		}()
	}
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	return results, nil
}
```

</details>
