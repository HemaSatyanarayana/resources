package mutexlab

import (
	"strconv"
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	var c Counter // zero value must be usable
	const goroutines, per = 100, 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < per; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()
	if got, want := c.Value(), goroutines*per; got != want {
		t.Errorf("Counter.Value() = %d, want %d", got, want)
	}
}

func TestCounterAdd(t *testing.T) {
	var c Counter
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			c.Add(5)
		}()
	}
	wg.Wait()
	if got := c.Value(); got != 50 {
		t.Errorf("Counter.Value() after 10*Add(5) = %d, want 50", got)
	}
}

func TestStore(t *testing.T) {
	s := NewStore()

	if _, ok := s.Get("missing"); ok {
		t.Error("Get on empty store returned ok=true")
	}
	if s.Len() != 0 {
		t.Errorf("empty Store.Len() = %d, want 0", s.Len())
	}

	// Concurrent writers to distinct keys.
	var wg sync.WaitGroup
	const n = 200
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			s.Set(key(i), i)
		}()
	}
	wg.Wait()

	if s.Len() != n {
		t.Errorf("Store.Len() = %d, want %d", s.Len(), n)
	}
	if v, ok := s.Get(key(42)); !ok || v != 42 {
		t.Errorf("Get(%q) = (%d,%v), want (42,true)", key(42), v, ok)
	}

	// Overwrite.
	s.Set(key(42), 4242)
	if v, _ := s.Get(key(42)); v != 4242 {
		t.Errorf("after overwrite Get = %d, want 4242", v)
	}
}

// TestStoreSnapshotIsCopy verifies Snapshot returns an independent copy.
func TestStoreSnapshotIsCopy(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("b", 2)

	snap := s.Snapshot()
	if len(snap) != 2 || snap["a"] != 1 || snap["b"] != 2 {
		t.Fatalf("Snapshot = %v, want {a:1 b:2}", snap)
	}

	// Mutating the snapshot must not affect the store.
	snap["a"] = 999
	delete(snap, "b")
	if v, _ := s.Get("a"); v != 1 {
		t.Errorf("store changed after mutating snapshot: a=%d, want 1", v)
	}
	if _, ok := s.Get("b"); !ok {
		t.Error("store lost key b after deleting from snapshot")
	}
}

// TestStoreConcurrentReadWrite is the real race test: readers and writers at
// once. Run with -race.
func TestStoreConcurrentReadWrite(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); s.Set("k", i) }()
		go func() { defer wg.Done(); s.Get("k"); s.Len(); s.Snapshot() }()
	}
	wg.Wait()
}

func key(i int) string {
	return "k" + strconv.Itoa(i)
}
