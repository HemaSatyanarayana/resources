package concurrency

import (
	"reflect"
	"sync"
	"testing"
)

func TestParallelSquare(t *testing.T) {
	got := ParallelSquare([]int{1, 2, 3, 4, 5})
	want := []int{1, 4, 9, 16, 25}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParallelSquare = %v, want %v", got, want)
	}
	if got := ParallelSquare(nil); len(got) != 0 {
		t.Errorf("ParallelSquare(nil) = %v, want empty", got)
	}
}

func TestSafeCounter(t *testing.T) {
	var c SafeCounter
	const goroutines, per = 50, 200
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
		t.Errorf("SafeCounter.Value() = %d, want %d", got, want)
	}
}

func TestConcurrentSum(t *testing.T) {
	nums := make([]int, 1000)
	want := 0
	for i := range nums {
		nums[i] = i + 1
		want += i + 1
	}
	for _, workers := range []int{1, 3, 8, 0} {
		if got := ConcurrentSum(nums, workers); got != want {
			t.Errorf("ConcurrentSum(workers=%d) = %d, want %d", workers, got, want)
		}
	}
	if got := ConcurrentSum(nil, 4); got != 0 {
		t.Errorf("ConcurrentSum(nil) = %d, want 0", got)
	}
}
