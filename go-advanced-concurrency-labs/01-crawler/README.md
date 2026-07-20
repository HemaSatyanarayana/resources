# 01 — Concurrent Web Crawler

Your first end-to-end project. A crawler looks trivial ("fetch a page, follow
its links, repeat") but it hides the hardest problem in concurrent programming:
**knowing when you're finished**. The set of pages to visit isn't fixed — it
grows as you discover links — so you can't just `wg.Add(n); wg.Wait()` over a
known count.

## The system

```
seed ─▶ fetch ─▶ links ─┬─▶ fetch ─▶ links ─┬─▶ ...
                        ├─▶ fetch ─▶ links ─┘
                        └─▶ fetch ─▶ (already visited → skip)
        ▲                                   │
        └──────── bounded by `workers` ─────┘
```

Three constraints turn this into a real exercise:

1. **Bounded parallelism.** At most `workers` fetches run at once — you don't
   spawn 10,000 goroutines hammering a site. A **semaphore channel** of capacity
   `workers` enforces this.
2. **Dedup under a race.** Two goroutines may discover the same link at the same
   instant. The *check* ("have we seen this?") and the *mark* ("now we have")
   must happen under **one** lock hold, or you'll fetch it twice.
3. **Termination.** The frontier grows. Track *outstanding* work with a
   `WaitGroup`: `Add(1)` the moment you schedule a URL, `Done()` when its crawl
   returns. `wg.Wait()` unblocks exactly when the last discovered URL finishes.

### Why a semaphore, not a worker pool?

A classic worker pool has N goroutines `range`-ing a jobs channel. That works,
but closing the jobs channel at the right time is fiddly when the producer is
*also* every worker. The semaphore pattern — "spawn a goroutine per URL, but
each must grab one of `workers` tokens before fetching" — sidesteps channel-close
choreography entirely. Goroutines are cheap; blocked-on-a-semaphore goroutines
are practically free.

### Depth semantics (read carefully — a test pins this)

- The seed is depth `0`.
- Links found on a depth-`d` page are crawled at depth `d+1`, **only if**
  `d < maxDepth`.
- So `maxDepth == 0` crawls the seed alone; `maxDepth == 1` crawls the seed and
  its direct links; etc.

## Your task

Implement, from scratch, in `crawler.go` (package `crawler`):

```go
type Fetcher interface {
	Fetch(ctx context.Context, url string) (urls []string, err error)
}

func Crawl(ctx context.Context, seed string, maxDepth, workers int, f Fetcher) []string
```

Return the URLs fetched **without error** (the seed included, if it fetched).
Order is irrelevant — tests sort. Treat `workers <= 0` as `1`. Stop early and
return whatever you have if `ctx` is cancelled.

## Run

```bash
go test -race -v ./01-crawler/
```

## Hints

- Two sets behind one `sync.Mutex`: `visited` (scheduled) and `fetched` (succeeded).
- Semaphore: `sem := make(chan struct{}, workers)`. Acquire with
  `select { case <-ctx.Done(): return; case sem <- struct{}{}: }`, release with
  `defer func(){ <-sem }()`.
- Mark the seed visited **before** launching its crawl, so a page linking back to
  it doesn't re-schedule it.
- `wg.Add(1)` immediately before every `go crawl(...)`, `defer wg.Done()` first
  thing inside `crawl`.
- The context-cancel test uses an already-cancelled context and expects **zero**
  fetches — so check `ctx.Done()` *before* fetching, not after.

<details>
<summary>Reference solution</summary>

```go
package crawler

import (
	"context"
	"sort"
	"sync"
)

func Crawl(ctx context.Context, seed string, maxDepth, workers int, f Fetcher) []string {
	if workers <= 0 {
		workers = 1
	}

	var mu sync.Mutex
	visited := map[string]bool{} // scheduled (never schedule twice)
	fetched := map[string]bool{} // succeeded

	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	var crawl func(url string, depth int)
	crawl = func(url string, depth int) {
		defer wg.Done()

		// Acquire a worker slot (or bail if cancelled).
		select {
		case <-ctx.Done():
			return
		case sem <- struct{}{}:
		}
		defer func() { <-sem }()

		links, err := f.Fetch(ctx, url)
		if err != nil {
			return
		}
		mu.Lock()
		fetched[url] = true
		mu.Unlock()

		if depth >= maxDepth {
			return
		}
		for _, link := range links {
			mu.Lock()
			if !visited[link] {
				visited[link] = true
				wg.Add(1)
				go crawl(link, depth+1)
			}
			mu.Unlock()
		}
	}

	mu.Lock()
	visited[seed] = true
	mu.Unlock()
	wg.Add(1)
	go crawl(seed, 0)
	wg.Wait()

	out := make([]string, 0, len(fetched))
	for u := range fetched {
		out = append(out, u)
	}
	sort.Strings(out)
	return out
}
```

Why it's race-free: every read/write of `visited` and `fetched` is under `mu`;
the check-and-mark of `visited` is a single critical section, so a URL is
scheduled exactly once. `wg` counts scheduled-but-unfinished crawls, so `Wait`
returns precisely when the frontier is exhausted.

</details>
