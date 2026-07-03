# 09 — Channels & Pipelines

Channels are how goroutines communicate. Go's motto: **"Don't communicate by sharing memory; share memory by communicating."** This chapter builds composable pipeline stages and fans multiple channels into one.

## Concepts

- **`ch := make(chan int)`** is unbuffered — a send blocks until a receive is ready (a rendezvous). `make(chan int, n)` buffers `n` values.
- **Directional types** document intent: `<-chan int` is receive-only, `chan<- int` is send-only. Pipeline stages return `<-chan int`.
- **Close to signal "no more values":** `close(ch)`. **Only the sender closes**, never the receiver. `for v := range ch` loops until the channel is closed.
- **Receiving from a closed channel** yields the zero value immediately with `ok == false`: `v, ok := <-ch`.
- **`select`** waits on multiple channel operations; whichever is ready first wins (random among ties).
- **Fan-out / fan-in:** spread work across goroutines, then merge their outputs back into one channel.

## Your task

Implement the stages in [`pipeline.go`](pipeline.go). Delete the `var _ = sync.WaitGroup{}` line once you actually use `sync.WaitGroup` in `Merge`.

| Function | Skill |
|----------|-------|
| `Gen` | Producer goroutine + `close` |
| `Square` | A stage: read → transform → emit → close |
| `Merge` | Fan-in with `WaitGroup` |
| `Collect` | `range` a channel to completion |
| `FirstOf` | `select` |

## Run

```bash
go test -race -v ./exercises/09-channels-pipeline/
```

## Hints

- Every stage follows the shape: `out := make(chan int); go func(){ defer close(out); for v := range in { out <- transform(v) } }(); return out`.
- In `Merge`, launch one goroutine per input that copies it into `out`; `wg.Add(1)` before each, `wg.Done()` when its input drains. A **separate** goroutine does `wg.Wait(); close(out)` — you can't close `out` until all forwarders finish.
- Forgetting to `close` leaves a `range` blocked forever (a deadlock or leaked goroutine).
- `FirstOf` is a two-case `select` with `case v := <-a:` / `case v := <-b:`.

<details>
<summary>Reference solution</summary>

```go
package pipeline

import "sync"

func Gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * v
		}
	}()
	return out
}

func Merge(ins ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(ins))
	for _, in := range ins {
		go func() {
			defer wg.Done()
			for v := range in {
				out <- v
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func Collect(in <-chan int) []int {
	var out []int
	for v := range in {
		out = append(out, v)
	}
	return out
}

func FirstOf(a, b <-chan int) int {
	select {
	case v := <-a:
		return v
	case v := <-b:
		return v
	}
}
```

</details>
