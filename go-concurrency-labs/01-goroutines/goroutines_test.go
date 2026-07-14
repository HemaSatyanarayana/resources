package goroutines

import (
	"reflect"
	"testing"
	"time"
)

func TestParallelMap(t *testing.T) {
	got := ParallelMap([]int{1, 2, 3, 4, 5}, func(n int) int { return n * n })
	want := []int{1, 4, 9, 16, 25}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParallelMap squares = %v, want %v", got, want)
	}

	// Order must be preserved even though work happens concurrently.
	got = ParallelMap([]int{10, 20, 30}, func(n int) int { return n + 1 })
	want = []int{11, 21, 31}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParallelMap +1 = %v, want %v", got, want)
	}

	if got := ParallelMap(nil, func(n int) int { return n }); len(got) != 0 {
		t.Errorf("ParallelMap(nil) = %v, want empty", got)
	}
}

// TestParallelMapIsConcurrent fails if the work is done sequentially: 100 tasks
// that each sleep 20ms must finish in well under 100*20ms if run concurrently.
func TestParallelMapIsConcurrent(t *testing.T) {
	const n = 100
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i
	}
	start := time.Now()
	got := ParallelMap(nums, func(x int) int {
		time.Sleep(20 * time.Millisecond)
		return x * 2
	})
	if elapsed := time.Since(start); elapsed > 1*time.Second {
		t.Fatalf("took %v; ParallelMap does not look concurrent", elapsed)
	}
	for i := range got {
		if got[i] != i*2 {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], i*2)
		}
	}
}

func TestWaitAll(t *testing.T) {
	const n = 50
	// Each fn writes its own slot — race-free. WaitAll must not return until
	// every slot is written.
	done := make([]bool, n)
	fns := make([]func(), n)
	for i := range fns {
		fns[i] = func() {
			time.Sleep(time.Millisecond)
			done[i] = true
		}
	}

	WaitAll(fns)

	for i, d := range done {
		if !d {
			t.Fatalf("fn %d did not finish before WaitAll returned", i)
		}
	}
}

// TestWaitAllEmpty must return promptly and not panic.
func TestWaitAllEmpty(t *testing.T) {
	doneCh := make(chan struct{})
	go func() { WaitAll(nil); close(doneCh) }()
	select {
	case <-doneCh:
	case <-time.After(time.Second):
		t.Fatal("WaitAll(nil) did not return")
	}
}
