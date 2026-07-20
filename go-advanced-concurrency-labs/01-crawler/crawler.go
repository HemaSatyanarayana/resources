// Package crawler — Lab 01: a bounded, concurrent web crawler.
//
// Read README.md first; crawler_test.go is the spec. Fill in every TODO.
// Run: go test -race -v ./01-crawler/
package crawler

import "context"

// Fetcher retrieves the links found on a page. Tests supply an in-memory fake,
// so you never touch the network — you only orchestrate concurrency.
type Fetcher interface {
	// Fetch returns the URLs linked from url, or an error (e.g. a 404).
	Fetch(ctx context.Context, url string) (urls []string, err error)
}

// Crawl starts at seed and concurrently crawls every reachable URL, using at
// most `workers` simultaneous fetches. Each URL is fetched at most once. Links
// discovered at depth d are crawled at depth d+1, only while d+1 <= maxDepth
// (so maxDepth 0 fetches the seed alone). URLs whose Fetch errors are skipped.
// The crawl stops early if ctx is cancelled.
//
// It returns the URLs that were fetched successfully (order does not matter;
// tests sort before comparing). Treat workers <= 0 as 1.
func Crawl(ctx context.Context, seed string, maxDepth, workers int, f Fetcher) []string {
	// TODO: implement.
	//
	// The hard part is not fan-out — it's knowing when you're DONE. The frontier
	// grows as you discover links, so a plain WaitGroup over a fixed set won't
	// work. Sketch:
	//
	//   1. Guard workers <= 0 -> 1.
	//   2. Protect a `visited` set (URLs already scheduled) with a sync.Mutex,
	//      and a `fetched` set (URLs that returned no error).
	//   3. Bound concurrency with a semaphore channel: sem := make(chan struct{}, workers).
	//   4. A recursive closure crawl(url, depth) that:
	//        - returns immediately if ctx is done (select on ctx.Done() vs. sem<-),
	//        - acquires the semaphore, fetches, releases it,
	//        - records fetched[url] on success,
	//        - for each link, if depth < maxDepth and not visited: mark visited,
	//          wg.Add(1), and `go crawl(link, depth+1)`.
	//      Check-and-mark `visited` under the SAME lock hold, or you'll fetch
	//      duplicates in a race.
	//   5. Track outstanding work with a sync.WaitGroup: Add(1) before each
	//      `go crawl`, Done() (defer) inside it. wg.Wait() then collect fetched.
	panic("TODO: implement Crawl")
}
