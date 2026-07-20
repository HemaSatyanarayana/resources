package cache

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetCachesValue(t *testing.T) {
	var calls int32
	c := New(time.Minute, func(key string) (int, error) {
		return int(atomic.AddInt32(&calls, 1)), nil
	})
	defer c.Close()

	v1, err := c.Get("k")
	if err != nil || v1 != 1 {
		t.Fatalf("first Get = (%d, %v), want (1, nil)", v1, err)
	}
	v2, err := c.Get("k") // must be served from cache, not reloaded
	if err != nil || v2 != 1 {
		t.Fatalf("second Get = (%d, %v), want (1, nil)", v2, err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("loader called %d times, want 1 (second Get should hit the cache)", got)
	}
}

func TestSingleflightDedupsConcurrentMisses(t *testing.T) {
	var calls int32
	release := make(chan struct{})
	c := New(time.Minute, func(key string) (int, error) {
		atomic.AddInt32(&calls, 1)
		<-release // hold the load open so every caller piles onto this one call
		return 42, nil
	})
	defer c.Close()

	const n = 30
	results := make([]int, n)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			v, err := c.Get("k")
			if err != nil {
				t.Errorf("Get = err %v", err)
			}
			results[i] = v
		}(i)
	}
	close(start)   // fire all callers
	close(release) // let the single in-flight load complete
	wg.Wait()

	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("loader called %d times, want exactly 1 — concurrent misses were not de-duplicated", got)
	}
	for i, v := range results {
		if v != 42 {
			t.Errorf("caller %d got %d, want 42 (all should share the one load)", i, v)
		}
	}
}

func TestTTLExpiryReloads(t *testing.T) {
	var calls int32
	c := New(30*time.Millisecond, func(key string) (int, error) {
		return int(atomic.AddInt32(&calls, 1)), nil
	})
	defer c.Close()

	if v, _ := c.Get("k"); v != 1 {
		t.Fatalf("Get #1 = %d, want 1", v)
	}
	if v, _ := c.Get("k"); v != 1 {
		t.Fatalf("Get #2 = %d, want 1 (still fresh)", v)
	}

	time.Sleep(50 * time.Millisecond) // let the entry expire

	if v, _ := c.Get("k"); v != 2 {
		t.Errorf("Get after TTL = %d, want 2 (entry should have been reloaded)", v)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("loader called %d times, want 2", got)
	}
}

func TestDifferentKeysLoadIndependently(t *testing.T) {
	var calls int32
	c := New(time.Minute, func(key string) (string, error) {
		atomic.AddInt32(&calls, 1)
		return "val:" + key, nil
	})
	defer c.Close()

	if a, _ := c.Get("a"); a != "val:a" {
		t.Errorf(`Get("a") = %q, want "val:a"`, a)
	}
	if b, _ := c.Get("b"); b != "val:b" {
		t.Errorf(`Get("b") = %q, want "val:b"`, b)
	}
	c.Get("a") // both cached now
	c.Get("b")
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("loader called %d times, want 2 (one per distinct key)", got)
	}
}

func TestLoaderErrorIsNotCached(t *testing.T) {
	var calls int32
	c := New(time.Minute, func(key string) (int, error) {
		if atomic.AddInt32(&calls, 1) == 1 {
			return 0, errors.New("transient failure")
		}
		return 99, nil
	})
	defer c.Close()

	if _, err := c.Get("k"); err == nil {
		t.Fatal("first Get err = nil, want the loader's error")
	}
	v, err := c.Get("k") // error must not be cached, so this reloads
	if err != nil {
		t.Fatalf("second Get err = %v, want nil", err)
	}
	if v != 99 {
		t.Errorf("second Get = %d, want 99", v)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("loader called %d times, want 2 (a failed load must not be cached)", got)
	}
}

func TestJanitorEvictsExpiredEntries(t *testing.T) {
	c := New(20*time.Millisecond, func(key string) (int, error) { return 1, nil })
	defer c.Close()

	c.Get("a")
	c.Get("b")
	c.Get("c")
	if got := c.Len(); got != 3 {
		t.Fatalf("Len = %d after 3 loads, want 3", got)
	}

	// Wait for the background janitor to evict expired entries — WITHOUT calling
	// Get again (that would prove lazy expiry, not the janitor).
	deadline := time.Now().Add(2 * time.Second)
	for c.Len() != 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if got := c.Len(); got != 0 {
		t.Errorf("Len = %d after expiry, want 0 — the janitor did not sweep expired entries", got)
	}
}

func TestCloseIsIdempotent(t *testing.T) {
	c := New(time.Millisecond, func(key string) (int, error) { return 0, nil })
	c.Close()
	c.Close() // must not panic (no double close of the stop channel)
}

func TestConcurrentUse(t *testing.T) {
	var calls int64
	c := New(5*time.Millisecond, func(key string) (int, error) {
		atomic.AddInt64(&calls, 1)
		return len(key), nil
	})
	defer c.Close()

	keys := []string{"a", "bb", "ccc", "dddd"}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				k := keys[(i+j)%len(keys)]
				v, err := c.Get(k)
				if err != nil {
					t.Errorf("Get(%q) = err %v", k, err)
					return
				}
				if v != len(k) {
					t.Errorf("Get(%q) = %d, want %d", k, v, len(k))
					return
				}
			}
		}(i)
	}
	wg.Wait()
	// No assertion on the exact load count (TTL churn makes it nondeterministic);
	// this test exists to be run under -race.
	_ = atomic.LoadInt64(&calls)
}
