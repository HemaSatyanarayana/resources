package pipe

import (
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestPipelineCompose(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	nums := Gen(done, 1, 2, 3, 4, 5, 6)
	evens := Filter(done, nums, func(n int) bool { return n%2 == 0 })
	squared := Map(done, evens, func(n int) int { return n * n })

	got := Collect(squared)
	want := []int{4, 16, 36}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("pipeline = %v, want %v", got, want)
	}
}

func TestGen(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	if got := Collect(Gen(done, 7, 8, 9)); !reflect.DeepEqual(got, []int{7, 8, 9}) {
		t.Errorf("Gen = %v, want [7 8 9]", got)
	}
	if got := Collect(Gen(done)); len(got) != 0 {
		t.Errorf("Gen() = %v, want empty", got)
	}
}

func TestMapAndFilter(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	doubled := Collect(Map(done, Gen(done, 1, 2, 3), func(n int) int { return n * 2 }))
	if !reflect.DeepEqual(doubled, []int{2, 4, 6}) {
		t.Errorf("Map = %v, want [2 4 6]", doubled)
	}

	big := Collect(Filter(done, Gen(done, 5, 15, 25, 3), func(n int) bool { return n > 10 }))
	if !reflect.DeepEqual(big, []int{15, 25}) {
		t.Errorf("Filter = %v, want [15 25]", big)
	}
}

// TestEarlyCancelNoLeak closes done after taking a single value and verifies the
// stage goroutines actually exit (no goroutine leak).
func TestEarlyCancelNoLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	done := make(chan struct{})
	// An "infinite" source: keep feeding the same numbers by chaining a big
	// slice; the consumer will bail after one value.
	nums := make([]int, 10000)
	for i := range nums {
		nums[i] = i
	}
	out := Map(done, Gen(done, nums...), func(n int) int { return n + 1 })

	// Take exactly one value, then cancel.
	v := <-out
	if v != 1 {
		t.Errorf("first value = %d, want 1", v)
	}
	close(done)

	// Give the stages a moment to notice done and return.
	time.Sleep(50 * time.Millisecond)
	runtime.GC()

	after := runtime.NumGoroutine()
	if after > before+2 { // small slack for scheduler/GC goroutines
		t.Errorf("goroutines leaked: before=%d after=%d", before, after)
	}
}
