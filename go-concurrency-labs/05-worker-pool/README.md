# 05 — Worker Pools

Launching a goroutine per item is fine for a handful of tasks. For thousands —
or when each task hits a rate-limited API or uses lots of memory — you want a
**fixed** number of workers pulling from a shared queue. That's a worker pool:
it **bounds concurrency** and reuses goroutines.

## Concepts

### The shape

```
              ┌────────┐
 jobs ───────▶│ worker │──┐
        │     └────────┘  │
        │     ┌────────┐  ├──▶ results
        ├────▶│ worker │──┤
        │     └────────┘  │
        │     ┌────────┐  │
        └────▶│ worker │──┘
              └────────┘
```

- One `jobs` channel. N worker goroutines each `range` over it.
- Each worker sends to a `results` channel.
- A feeder goroutine pushes all jobs then **closes** `jobs` — that's what ends
  every worker's `range` loop.
- A `WaitGroup` tracks the workers; when all are done, **close** `results` so the
  collector's `range` ends.

### Why close in that order

- Closing `jobs` tells workers "no more work" → their `for j := range jobs`
  loops exit → each worker returns → `WaitGroup` drains → close `results`.
- If you close `results` too early, a worker still sending panics. So:
  `wg.Wait()` **then** `close(results)`, and do that wait in its **own**
  goroutine so the collector can start receiving immediately (otherwise: deadlock
  on an unbuffered `results`).

### Preserving order

Results arrive in whatever order workers finish — not input order. To return
results aligned to inputs, tag each job with its index and write the result to
`out[index]`:

```go
type job struct{ idx, val int }
type res struct{ idx, out int }
```

Then `out[r.idx] = r.out`. Distinct indices ⇒ no lock needed on `out`.

## Your task

Create `pool.go` (package `workerpool`) and implement, **from scratch**:

```go
// Map applies f to every input using exactly `workers` worker goroutines and
// returns the results in INPUT ORDER. If workers <= 0, treat it as 1.
func Map(inputs []int, workers int, f func(int) int) []int
```

## Run

```bash
go test -race -v ./05-worker-pool/
```

## Hints

- Define small structs to carry the index alongside the value in and out.
- Feeder goroutine: send every `{idx,val}` then `close(jobs)`.
- Start `workers` goroutines (guard `workers <= 0 → 1`), each: `wg.Add(1)` before
  launch, `for j := range jobs { results <- res{j.idx, f(j.val)} }`, `wg.Done()`.
- Closer goroutine: `wg.Wait(); close(results)`.
- Collector (main goroutine): `out := make([]int, len(inputs))`, then
  `for r := range results { out[r.idx] = r.out }`, then `return out`.
- Handle empty input: return an empty (non-nil) slice without hanging.

<details>
<summary>Reference solution</summary>

```go
package workerpool

import "sync"

func Map(inputs []int, workers int, f func(int) int) []int {
	if workers <= 0 {
		workers = 1
	}
	out := make([]int, len(inputs))
	if len(inputs) == 0 {
		return out
	}

	type job struct{ idx, val int }
	type res struct{ idx, out int }

	jobs := make(chan job)
	results := make(chan res)

	// Feeder.
	go func() {
		defer close(jobs)
		for i, v := range inputs {
			jobs <- job{i, v}
		}
	}()

	// Workers.
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				results <- res{j.idx, f(j.val)}
			}
		}()
	}

	// Closer.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collector.
	for r := range results {
		out[r.idx] = r.out
	}
	return out
}
```

</details>
