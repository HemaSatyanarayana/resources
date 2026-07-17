package workerpool

import (
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func TestMapOrder(t *testing.T) {
	inputs := make([]int, 100)
	want := make([]int, 100)
	for i := range inputs {
		inputs[i] = i
		want[i] = i * i
	}
	for _, workers := range []int{1, 2, 5, 8, 0, -3} {
		got := Map(inputs, workers, func(n int) int { return n * n })
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Map(workers=%d) wrong result", workers)
		}
	}
}

func TestMapEmpty(t *testing.T) {
	if got := Map(nil, 4, func(n int) int { return n }); len(got) != 0 {
		t.Errorf("Map(nil) = %v, want empty", got)
	}
}

// TestMapBoundsConcurrency verifies at most `workers` calls to f run at once.
func TestMapBoundsConcurrency(t *testing.T) {
	const workers = 4
	inputs := make([]int, 40)
	for i := range inputs {
		inputs[i] = i
	}

	var inFlight, maxSeen int64
	Map(inputs, workers, func(n int) int {
		cur := atomic.AddInt64(&inFlight, 1)
		for {
			old := atomic.LoadInt64(&maxSeen)
			if cur <= old || atomic.CompareAndSwapInt64(&maxSeen, old, cur) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
		atomic.AddInt64(&inFlight, -1)
		return n
	})

	if maxSeen > workers {
		t.Errorf("observed %d concurrent workers, want <= %d", maxSeen, workers)
	}
	if maxSeen < 2 {
		t.Errorf("observed only %d concurrent worker(s); pool isn't parallel", maxSeen)
	}
}

// TestMapIsConcurrent: 40 tasks of 10ms with 8 workers must finish well under
// the 400ms a sequential run would take.
func TestMapIsConcurrent(t *testing.T) {
	inputs := make([]int, 40)
	for i := range inputs {
		inputs[i] = i
	}
	start := time.Now()
	Map(inputs, 8, func(n int) int {
		time.Sleep(10 * time.Millisecond)
		return n
	})
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("took %v; pool does not look concurrent", elapsed)
	}
}
