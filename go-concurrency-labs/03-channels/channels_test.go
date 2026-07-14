package channels

import (
	"reflect"
	"testing"
	"time"
)

func TestGenerateAndDrain(t *testing.T) {
	got := Drain(Generate(1, 2, 3, 4))
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Drain(Generate(...)) = %v, want %v", got, want)
	}

	if got := Drain(Generate()); len(got) != 0 {
		t.Errorf("Drain(Generate()) = %v, want empty", got)
	}
}

// TestGenerateCloses makes sure Generate closes its channel (otherwise Drain
// would block forever). If this test hangs, you forgot to close.
func TestGenerateCloses(t *testing.T) {
	done := make(chan []int)
	go func() { done <- Drain(Generate(5, 6)) }()
	select {
	case got := <-done:
		if !reflect.DeepEqual(got, []int{5, 6}) {
			t.Errorf("got %v, want [5 6]", got)
		}
	case <-time.After(time.Second):
		t.Fatal("Drain never returned — Generate probably didn't close its channel")
	}
}

func TestTake(t *testing.T) {
	// Take fewer than available.
	if got := Take(Generate(1, 2, 3, 4, 5), 3); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("Take(...,3) = %v, want [1 2 3]", got)
	}
	// Take more than available — must stop at close, not block.
	done := make(chan []int)
	go func() { done <- Take(Generate(1, 2), 10) }()
	select {
	case got := <-done:
		if !reflect.DeepEqual(got, []int{1, 2}) {
			t.Errorf("Take(2 avail, 10) = %v, want [1 2]", got)
		}
	case <-time.After(time.Second):
		t.Fatal("Take blocked past channel close")
	}
	// Take zero.
	if got := Take(Generate(1, 2), 0); len(got) != 0 {
		t.Errorf("Take(...,0) = %v, want empty", got)
	}
}

func TestBuffered(t *testing.T) {
	ch := Buffered(7, 8, 9)
	// Should be rangeable without any producer goroutine.
	got := Drain(ch)
	if !reflect.DeepEqual(got, []int{7, 8, 9}) {
		t.Errorf("Drain(Buffered(7,8,9)) = %v, want [7 8 9]", got)
	}
}
