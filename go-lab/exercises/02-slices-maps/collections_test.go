package collections

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	in := []int{1, 2, 3, 4}
	got := Reverse(in)
	want := []int{4, 3, 2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Reverse = %v, want %v", got, want)
	}
	// Original must be untouched.
	if !reflect.DeepEqual(in, []int{1, 2, 3, 4}) {
		t.Errorf("Reverse mutated its input: %v", in)
	}
	if got := Reverse([]int{}); len(got) != 0 {
		t.Errorf("Reverse([]) = %v, want empty", got)
	}
}

func TestDedup(t *testing.T) {
	got := Dedup([]int{3, 1, 3, 2, 1, 2, 5})
	want := []int{3, 1, 2, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Dedup = %v, want %v", got, want)
	}
}

func TestWordCount(t *testing.T) {
	got := WordCount("the cat sat on the mat the")
	want := map[string]int{"the": 3, "cat": 1, "sat": 1, "on": 1, "mat": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordCount = %v, want %v", got, want)
	}
	if len(WordCount("   ")) != 0 {
		t.Errorf("WordCount of blanks should be empty")
	}
}

func TestChunk(t *testing.T) {
	got := Chunk([]int{1, 2, 3, 4, 5}, 2)
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Chunk = %v, want %v", got, want)
	}
	if Chunk([]int{1, 2}, 0) != nil {
		t.Errorf("Chunk with size 0 should be nil")
	}
	exact := Chunk([]int{1, 2, 3, 4}, 2)
	if !reflect.DeepEqual(exact, [][]int{{1, 2}, {3, 4}}) {
		t.Errorf("Chunk exact = %v", exact)
	}
}
