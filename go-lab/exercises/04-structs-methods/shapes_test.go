package shapes

import (
	"math"
	"testing"
)

const eps = 1e-9

func closeEnough(a, b float64) bool { return math.Abs(a-b) < eps }

func TestRectangle(t *testing.T) {
	r := Rectangle{Width: 3, Height: 4}
	if !closeEnough(r.Area(), 12) {
		t.Errorf("Rectangle.Area() = %v, want 12", r.Area())
	}
	if !closeEnough(r.Perimeter(), 14) {
		t.Errorf("Rectangle.Perimeter() = %v, want 14", r.Perimeter())
	}
}

func TestCircle(t *testing.T) {
	c := Circle{Radius: 2}
	if !closeEnough(c.Area(), math.Pi*4) {
		t.Errorf("Circle.Area() = %v, want %v", c.Area(), math.Pi*4)
	}
	if !closeEnough(c.Perimeter(), math.Pi*4) {
		t.Errorf("Circle.Perimeter() = %v, want %v", c.Perimeter(), math.Pi*4)
	}
}

func TestCounter(t *testing.T) {
	var c Counter // zero value must be usable
	if c.Value() != 0 {
		t.Fatalf("zero Counter Value = %d, want 0", c.Value())
	}
	c.Inc()
	c.Inc()
	c.Add(5)
	if c.Value() != 7 {
		t.Errorf("Counter.Value() = %d, want 7", c.Value())
	}
}

func TestBoxPromotion(t *testing.T) {
	b := Box{Label: "crate", Rectangle: Rectangle{Width: 2, Height: 3}}
	if !closeEnough(b.Area(), 6) { // promoted method
		t.Errorf("Box.Area() (promoted) = %v, want 6", b.Area())
	}
	if got, want := b.Describe(), "crate: 6.00"; got != want {
		t.Errorf("Box.Describe() = %q, want %q", got, want)
	}
}
