package crawler

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeFetcher is an in-memory site graph. It never touches the network, records
// how many times each URL was fetched, and tracks the peak number of concurrent
// Fetch calls so tests can assert the crawler bounds its parallelism.
type fakeFetcher struct {
	graph map[string][]string
	delay time.Duration

	mu    sync.Mutex
	calls map[string]int

	inFlight int32
	peak     int32
}

func newFetcher(graph map[string][]string, delay time.Duration) *fakeFetcher {
	return &fakeFetcher{graph: graph, delay: delay, calls: map[string]int{}}
}

func (f *fakeFetcher) Fetch(ctx context.Context, url string) ([]string, error) {
	cur := atomic.AddInt32(&f.inFlight, 1)
	for {
		old := atomic.LoadInt32(&f.peak)
		if cur <= old || atomic.CompareAndSwapInt32(&f.peak, old, cur) {
			break
		}
	}
	defer atomic.AddInt32(&f.inFlight, -1)

	f.mu.Lock()
	f.calls[url]++
	f.mu.Unlock()

	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	links, ok := f.graph[url]
	if !ok {
		return nil, fmt.Errorf("404: %s", url)
	}
	return links, nil
}

func (f *fakeFetcher) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, c := range f.calls {
		n += c
	}
	return n
}

func TestCrawlVisitsAllReachable(t *testing.T) {
	f := newFetcher(map[string][]string{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {"d"},
		"d": {},
	}, 0)

	got := Crawl(context.Background(), "a", 10, 4, f)
	want := []string{"a", "b", "c", "d"}
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Crawl = %v, want %v", got, want)
	}
	if n := f.callCount(); n != 4 {
		t.Errorf("fetched %d times, want 4 (each URL exactly once)", n)
	}
}

func TestCrawlRespectsDepth(t *testing.T) {
	f := newFetcher(map[string][]string{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {"d"},
		"d": {},
	}, 0)

	// depth 0 = seed only; links found at depth d are crawled at depth d+1,
	// crawled only while d+1 <= maxDepth.
	got := Crawl(context.Background(), "a", 1, 4, f)
	want := []string{"a", "b", "c"} // d is at depth 2, beyond maxDepth 1
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Crawl(maxDepth=1) = %v, want %v", got, want)
	}
}

func TestCrawlSeedOnly(t *testing.T) {
	f := newFetcher(map[string][]string{"a": {"b", "c"}, "b": {}, "c": {}}, 0)
	got := Crawl(context.Background(), "a", 0, 4, f)
	if want := []string{"a"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Crawl(maxDepth=0) = %v, want %v", got, want)
	}
}

func TestCrawlDedupesCycles(t *testing.T) {
	f := newFetcher(map[string][]string{
		"a": {"b", "a"}, // self-link
		"b": {"a"},      // back-link (cycle)
	}, 0)

	got := Crawl(context.Background(), "a", 100, 4, f)
	want := []string{"a", "b"}
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Crawl = %v, want %v", got, want)
	}
	if n := f.callCount(); n != 2 {
		t.Errorf("fetched %d times, want 2 — cycle not de-duplicated (or infinite loop)", n)
	}
}

func TestCrawlSkipsErrors(t *testing.T) {
	// "b" is not in the graph, so fetching it returns an error; it must not
	// appear in the result, but siblings must still be crawled.
	f := newFetcher(map[string][]string{
		"a": {"b", "c"},
		"c": {},
	}, 0)
	got := Crawl(context.Background(), "a", 10, 4, f)
	want := []string{"a", "c"}
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Crawl = %v, want %v", got, want)
	}
}

func TestCrawlBoundsConcurrency(t *testing.T) {
	links := make([]string, 20)
	graph := map[string][]string{}
	for i := range links {
		u := fmt.Sprintf("leaf-%02d", i)
		links[i] = u
		graph[u] = []string{}
	}
	graph["root"] = links

	f := newFetcher(graph, 15*time.Millisecond)
	const workers = 4
	Crawl(context.Background(), "root", 5, workers, f)

	if peak := atomic.LoadInt32(&f.peak); peak > workers {
		t.Errorf("peak concurrency %d, want <= %d", peak, workers)
	}
	if peak := atomic.LoadInt32(&f.peak); peak < 2 {
		t.Errorf("peak concurrency %d — crawler is not parallel", peak)
	}
}

func TestCrawlIsConcurrent(t *testing.T) {
	links := make([]string, 20)
	graph := map[string][]string{}
	for i := range links {
		u := fmt.Sprintf("leaf-%02d", i)
		links[i] = u
		graph[u] = []string{}
	}
	graph["root"] = links

	f := newFetcher(graph, 10*time.Millisecond)
	start := time.Now()
	Crawl(context.Background(), "root", 5, 10, f)
	// 21 fetches * 10ms sequential = 210ms. With 10 workers it should be far less.
	if elapsed := time.Since(start); elapsed > 120*time.Millisecond {
		t.Errorf("took %v — crawler does not look concurrent", elapsed)
	}
}

func TestCrawlContextCancel(t *testing.T) {
	f := newFetcher(map[string][]string{"a": {"b", "c"}, "b": {}, "c": {}}, 0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	got := Crawl(ctx, "a", 10, 4, f)
	if len(got) != 0 {
		t.Errorf("Crawl with cancelled ctx = %v, want empty", got)
	}
	if n := f.callCount(); n != 0 {
		t.Errorf("fetched %d times with cancelled ctx, want 0", n)
	}
}

func TestCrawlWorkersFloor(t *testing.T) {
	f := newFetcher(map[string][]string{"a": {"b"}, "b": {}}, 0)
	// workers <= 0 must be treated as 1, not deadlock or panic.
	got := Crawl(context.Background(), "a", 5, 0, f)
	want := []string{"a", "b"}
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Crawl(workers=0) = %v, want %v", got, want)
	}
}
