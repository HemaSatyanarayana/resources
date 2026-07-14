# 01 — Goroutines & WaitGroup

The two most fundamental tools: starting concurrent work, and waiting for it.

## Concepts

### `go f()` — launch a goroutine

```go
go doWork()          // starts doWork() concurrently, returns immediately
go func() { ... }()  // same, with an anonymous function
```

A goroutine is a function running independently, scheduled by the Go runtime
onto OS threads. It's cheap — a few KB of stack — so thousands are fine. But:

- You get **no return value**. The only way out is a shared variable or a channel.
- There is **no automatic join**. `main` returning kills every goroutine
  instantly, finished or not. You must wait explicitly.

### `sync.WaitGroup` — wait for a set of goroutines

A `WaitGroup` is a concurrency-safe counter:

```go
var wg sync.WaitGroup
wg.Add(1)              // +1 to the counter, BEFORE launching
go func() {
    defer wg.Done()    // -1 when this goroutine finishes
    // ...work...
}()
wg.Wait()              // blocks until the counter hits 0
```

Rules that will bite you if broken:
- Call `wg.Add` **before** `go`, in the launching goroutine — never inside the
  goroutine (it may not have run yet when `Wait` is reached).
- `defer wg.Done()` is the safe idiom — it runs even if the goroutine panics.
- Its zero value is ready to use. Don't copy a `WaitGroup` after use (pass a
  `*sync.WaitGroup` if you must pass it).

### Race-free parallel writes

You do **not** always need a lock. If each goroutine writes to a *distinct*
element of a pre-sized slice, there's no race — no two goroutines touch the same
memory:

```go
out := make([]int, len(nums))   // pre-size it
for i, n := range nums {
    wg.Add(1)
    go func() {
        defer wg.Done()
        out[i] = n * n           // i and n are per-iteration (Go 1.22+), safe
    }()
}
```

> **Loop variables (Go 1.22+):** each iteration gets a fresh `i` and `n`, so
> capturing them in the closure is safe. On Go ≤1.21 all goroutines shared one
> variable and you had to pass them as arguments. This module is on 1.25, so
> capture freely — but know the history, it's a classic interview gotcha.

## Your task

Create `goroutines.go` (package `goroutines`) and implement, **from scratch**:

```go
// ParallelMap applies f to every element of nums, running each call in its own
// goroutine, and returns the results in the SAME ORDER as nums.
func ParallelMap(nums []int, f func(int) int) []int

// WaitAll runs every function in fns in its own goroutine and returns only
// after all of them have finished.
func WaitAll(fns []func())
```

The test file `goroutines_test.go` is the spec — match those names exactly.

## Run

```bash
go test -race -v ./01-goroutines/
```

## Hints

- `ParallelMap`: pre-size `out := make([]int, len(nums))`, one goroutine per
  index writing `out[i] = f(n)`, a `WaitGroup` to join, then `return out`.
- `WaitAll`: `wg.Add(len(fns))`, loop launching each with `defer wg.Done()`,
  then `wg.Wait()`.
- If `-race` complains, you're sharing a variable you shouldn't — check that
  each goroutine writes its own slot.

<details>
<summary>Reference solution</summary>

```go
package goroutines

import "sync"

func ParallelMap(nums []int, f func(int) int) []int {
	out := make([]int, len(nums))
	var wg sync.WaitGroup
	wg.Add(len(nums))
	for i, n := range nums {
		go func() {
			defer wg.Done()
			out[i] = f(n)
		}()
	}
	wg.Wait()
	return out
}

func WaitAll(fns []func()) {
	var wg sync.WaitGroup
	wg.Add(len(fns))
	for _, fn := range fns {
		go func() {
			defer wg.Done()
			fn()
		}()
	}
	wg.Wait()
}
```

</details>
