# 08 — Concurrency

This is Go's signature feature. Goroutines are cheap (a few KB of stack), and the runtime multiplexes thousands of them onto a handful of OS threads. This chapter covers the two coordination primitives from `sync`; the next covers channels.

## Concepts

- **`go f()`** launches a goroutine. It runs concurrently; you get no return value and no automatic "join".
- **`sync.WaitGroup`** waits for a set of goroutines: `wg.Add(n)`, `defer wg.Done()` inside each, `wg.Wait()` in the caller.
- **`sync.Mutex`** protects shared mutable state: `mu.Lock()` / `defer mu.Unlock()`. Exactly one goroutine holds it at a time.
- **Data race** = two goroutines touching the same memory concurrently, at least one writing. Undefined behavior. **Always run `go test -race`.**
- **Loop-variable capture:** since Go 1.22 each iteration has a fresh loop variable, so `go func(){ use(i) }()` is safe. On older Go you had to pass `i` as an argument.
- **Worker pool:** a fixed number of goroutines pulling jobs off a channel — bounds concurrency and reuses goroutines.

## Your task

Implement everything in [`concurrency.go`](concurrency.go):

| Function | Primitive |
|----------|-----------|
| `ParallelSquare` | `WaitGroup` + independent index writes (no mutex needed) |
| `SafeCounter` | `Mutex` |
| `ConcurrentSum` | worker pool over channels |

## Run

Always with the race detector:

```bash
go test -race -v ./exercises/08-concurrency/
```

## Hints

- In `ParallelSquare`, capture the index: `for i, n := range nums { wg.Add(1); go func(){ defer wg.Done(); out[i] = n*n }() }`. Writing distinct indices of a pre-sized slice is race-free.
- `SafeCounter.Value` must also take the lock — reading an `int` concurrently with a write is still a race.
- For `ConcurrentSum`: make a `jobs := make(chan int)` and `results := make(chan int)`. Start `workers` goroutines. Feed jobs in a goroutine and `close(jobs)` when done so workers' `range jobs` loops end. Sum `workers` values off `results`.

<details>
<summary>Reference solution</summary>

```go
package concurrency

import "sync"

func ParallelSquare(nums []int) []int {
	out := make([]int, len(nums))
	var wg sync.WaitGroup
	for i, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out[i] = n * n
		}()
	}
	wg.Wait()
	return out
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

func ConcurrentSum(nums []int, workers int) int {
	if workers <= 0 {
		workers = 1
	}
	jobs := make(chan int)
	results := make(chan int)

	for w := 0; w < workers; w++ {
		go func() {
			partial := 0
			for n := range jobs {
				partial += n
			}
			results <- partial
		}()
	}

	go func() {
		for _, n := range nums {
			jobs <- n
		}
		close(jobs)
	}()

	total := 0
	for w := 0; w < workers; w++ {
		total += <-results
	}
	return total
}
```

</details>
