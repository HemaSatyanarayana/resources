package pipeline

import (
	"reflect"
	"sort"
	"testing"
)

func TestGen(t *testing.T) {
	got := Collect(Gen(1, 2, 3))
	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("Collect(Gen(1,2,3)) = %v", got)
	}
	if got := Collect(Gen()); len(got) != 0 {
		t.Errorf("Collect(Gen()) = %v, want empty", got)
	}
}

func TestSquarePipeline(t *testing.T) {
	got := Collect(Square(Gen(1, 2, 3, 4)))
	if !reflect.DeepEqual(got, []int{1, 4, 9, 16}) {
		t.Errorf("Square pipeline = %v, want [1 4 9 16]", got)
	}
}

func TestMergeFanIn(t *testing.T) {
	a := Gen(1, 2, 3)
	b := Gen(4, 5)
	c := Gen(6)
	got := Collect(Merge(a, b, c))
	sort.Ints(got)
	want := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Merge fan-in = %v, want %v", got, want)
	}
}

func TestFanInThenSquare(t *testing.T) {
	merged := Merge(Gen(1, 2), Gen(3, 4))
	got := Collect(Square(merged))
	sort.Ints(got)
	if !reflect.DeepEqual(got, []int{1, 4, 9, 16}) {
		t.Errorf("fan-in then square = %v", got)
	}
}

func TestFirstOf(t *testing.T) {
	a := Gen(42) // has a value ready
	b := make(chan int)
	if got := FirstOf(a, b); got != 42 {
		t.Errorf("FirstOf = %d, want 42", got)
	}
}
