# 13 — Context

`context.Context` is how Go carries **deadlines, cancellation signals, and request-scoped values** across API boundaries and goroutines. Every serious backend threads a `ctx` through its handlers, DB calls, and outbound requests. Master it and you can cancel work cleanly instead of leaking goroutines.

## Concepts

- **The rule:** `ctx` is the **first parameter** of a function, named `ctx`. Never store it in a struct; pass it explicitly.
- **`ctx.Done()`** returns a channel that's closed when the context is cancelled or its deadline passes. **`ctx.Err()`** then reports why: `context.Canceled` or `context.DeadlineExceeded`.
- **Deriving contexts:** `context.WithCancel(parent)`, `context.WithTimeout(parent, d)`, `context.WithDeadline`. Each returns a `cancel` func — **always `defer cancel()`** to release resources.
- **Values:** `context.WithValue(ctx, key, val)` for request-scoped data (request IDs, auth subject). Use an **unexported key type** to avoid collisions — never a bare string.
- **Cancellation is cooperative:** your code must actually *check* `ctx.Done()` (usually via `select`) for cancellation to take effect.

## Your task

Implement everything in [`contextx.go`](contextx.go):

| Function | Skill |
|----------|-------|
| `SumWithContext` | `select` over `ctx.Done()` vs work |
| `WithRequestID` / `RequestID` | Typed context values |
| `Race` | Result channel vs `ctx.Done()` |

## Run

```bash
go test -v ./exercises/13-context/
```

## Hints

- The core loop:
  ```go
  for _, n := range nums {
      select {
      case <-ctx.Done():
          return sum, ctx.Err()
      case <-time.After(perItem):
          sum += n
      }
  }
  ```
- `WithValue`/`Value` use the **same** unexported `requestIDKey`. The stored value comes back as `any`, so type-assert: `id, ok := ctx.Value(requestIDKey).(string)`.
- For `Race`, use a **buffered** channel (`make(chan int, 1)`) so the worker goroutine can send even after the caller has already returned on `ctx.Done()` — otherwise it blocks forever (a goroutine leak).

<details>
<summary>Reference solution</summary>

```go
package contextx

import (
	"context"
	"time"
)

func SumWithContext(ctx context.Context, nums []int, perItem time.Duration) (int, error) {
	sum := 0
	for _, n := range nums {
		select {
		case <-ctx.Done():
			return sum, ctx.Err()
		case <-time.After(perItem):
			sum += n
		}
	}
	return sum, nil
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

func Race(ctx context.Context, work func() int) (int, error) {
	done := make(chan int, 1)
	go func() { done <- work() }()
	select {
	case v := <-done:
		return v, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
```

</details>
