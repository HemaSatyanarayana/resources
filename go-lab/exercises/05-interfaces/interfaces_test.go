package shapesx

import (
	"fmt"
	"math"
	"testing"
)

func TestTotalArea(t *testing.T) {
	got := TotalArea(Rect{W: 2, H: 3}, Circ{R: 1})
	want := 6 + math.Pi
	if math.Abs(got-want) > 1e-4 {
		t.Errorf("TotalArea = %v, want ~%v", got, want)
	}
}

func TestTemperatureStringer(t *testing.T) {
	var s fmt.Stringer = Temperature(23.456)
	if got := s.String(); got != "23.5°C" {
		t.Errorf("Temperature.String() = %q, want %q", got, "23.5°C")
	}
	if got := fmt.Sprintf("%v", Temperature(-4)); got != "-4.0°C" {
		t.Errorf("formatted = %q, want %q", got, "-4.0°C")
	}
}

func TestSortByArea(t *testing.T) {
	shapes := []Shape{Circ{R: 3}, Rect{W: 1, H: 1}, Rect{W: 2, H: 2}}
	SortByArea(shapes)
	for i := 1; i < len(shapes); i++ {
		if shapes[i-1].Area() > shapes[i].Area() {
			t.Fatalf("not sorted ascending: %v", shapes)
		}
	}
	if shapes[0].Area() != 1 {
		t.Errorf("smallest area = %v, want 1", shapes[0].Area())
	}
}

func TestDescribeType(t *testing.T) {
	cases := []struct {
		v    any
		want string
	}{
		{nil, "nil"},
		{42, "int"},
		{"hi", "string"},
		{true, "bool"},
		{Rect{1, 2}, "shape"},
		{3.14, "unknown"},
	}
	for _, c := range cases {
		if got := DescribeType(c.v); got != c.want {
			t.Errorf("DescribeType(%#v) = %q, want %q", c.v, got, c.want)
		}
	}
}
