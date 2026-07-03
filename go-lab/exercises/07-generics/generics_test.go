package generics

import (
	"reflect"
	"sort"
	"testing"
)

func TestMap(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(x int) int { return x * x })
	if !reflect.DeepEqual(got, []int{1, 4, 9}) {
		t.Errorf("Map squares = %v", got)
	}
	// Type-changing map: int -> string length via a func.
	lengths := Map([]string{"a", "bb", "ccc"}, func(s string) int { return len(s) })
	if !reflect.DeepEqual(lengths, []int{1, 2, 3}) {
		t.Errorf("Map lengths = %v", lengths)
	}
}

func TestFilter(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4, 5, 6}, func(x int) bool { return x%2 == 0 })
	if !reflect.DeepEqual(got, []int{2, 4, 6}) {
		t.Errorf("Filter evens = %v", got)
	}
}

func TestReduce(t *testing.T) {
	sum := Reduce([]int{1, 2, 3, 4}, 0, func(acc, x int) int { return acc + x })
	if sum != 10 {
		t.Errorf("Reduce sum = %d, want 10", sum)
	}
	concat := Reduce([]string{"go", "is", "fun"}, "", func(acc, x string) string {
		if acc == "" {
			return x
		}
		return acc + " " + x
	})
	if concat != "go is fun" {
		t.Errorf("Reduce concat = %q", concat)
	}
}

func TestSum(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Errorf("Sum ints = %d, want 6", got)
	}
	if got := Sum([]float64{1.5, 2.5}); got != 4.0 {
		t.Errorf("Sum floats = %v, want 4.0", got)
	}
	// Named type with float64 underlying — exercises the ~ in the constraint.
	type Celsius float64
	if got := Sum([]Celsius{10, 20, 12}); got != 42 {
		t.Errorf("Sum Celsius = %v, want 42", got)
	}
}

func TestKeys(t *testing.T) {
	got := Keys(map[string]int{"a": 1, "b": 2, "c": 3})
	sort.Strings(got)
	if !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Errorf("Keys = %v", got)
	}
}
