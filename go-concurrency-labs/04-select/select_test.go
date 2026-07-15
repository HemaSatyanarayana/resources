package selectlab

import (
	"sort"
	"testing"
	"time"
)

func TestRecv(t *testing.T) {
	// Value available in time.
	ch := make(chan int, 1)
	ch <- 99
	if v, ok := Recv(ch, 100*time.Millisecond); !ok || v != 99 {
		t.Errorf("Recv ready = (%d,%v), want (99,true)", v, ok)
	}

	// Timeout: nothing ever sent.
	empty := make(chan int)
	start := time.Now()
	if v, ok := Recv(empty, 30*time.Millisecond); ok || v != 0 {
		t.Errorf("Recv timeout = (%d,%v), want (0,false)", v, ok)
	}
	if elapsed := time.Since(start); elapsed < 25*time.Millisecond {
		t.Errorf("Recv returned after %v, expected to wait ~30ms", elapsed)
	}

	// Value arrives shortly before the timeout.
	slow := make(chan int)
	go func() { time.Sleep(10 * time.Millisecond); slow <- 7 }()
	if v, ok := Recv(slow, 200*time.Millisecond); !ok || v != 7 {
		t.Errorf("Recv delayed = (%d,%v), want (7,true)", v, ok)
	}
}

func TestTryRecv(t *testing.T) {
	empty := make(chan int)
	if v, ok := TryRecv(empty); ok || v != 0 {
		t.Errorf("TryRecv(empty) = (%d,%v), want (0,false)", v, ok)
	}

	ready := make(chan int, 1)
	ready <- 5
	if v, ok := TryRecv(ready); !ok || v != 5 {
		t.Errorf("TryRecv(ready) = (%d,%v), want (5,true)", v, ok)
	}
}

func TestMerge(t *testing.T) {
	a := gen(1, 3, 5)
	b := gen(2, 4, 6)

	var got []int
	for v := range Merge(a, b) {
		got = append(got, v)
	}
	sort.Ints(got)
	want := []int{1, 2, 3, 4, 5, 6}
	if len(got) != len(want) {
		t.Fatalf("Merge produced %v, want the 6 values %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Merge sorted = %v, want %v", got, want)
		}
	}
}

// TestMergeCloses ensures the merged channel actually closes (both inputs
// drained). If Merge doesn't handle closed inputs via the nil-channel trick,
// this hangs.
func TestMergeCloses(t *testing.T) {
	done := make(chan int)
	go func() {
		count := 0
		for range Merge(gen(1, 2), gen(3)) {
			count++
		}
		done <- count
	}()
	select {
	case c := <-done:
		if c != 3 {
			t.Errorf("Merge emitted %d values, want 3", c)
		}
	case <-time.After(time.Second):
		t.Fatal("Merge never closed its output channel")
	}
}

func gen(nums ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, n := range nums {
			ch <- n
		}
	}()
	return ch
}
