// Package shapesx drills interfaces: implicit satisfaction, fmt.Stringer,
// sort.Interface, and type switches.
package shapesx

import (
	"fmt"
	"math"
	"sort"
)

// Shape is anything with an area. Note: types satisfy this IMPLICITLY —
// there is no "implements" keyword.
type Shape interface {
	Area() float64
}

// Rect and Circ both satisfy Shape.
type Rect struct{ W, H float64 }
type Circ struct{ R float64 }

// Area implementations — fill these in (use 3.14159 or math.Pi for Circ).
func (r Rect) Area() float64 {
	// panic("TODO: implement Rect.Area")
	return r.H * r.W
}

func (c Circ) Area() float64 {
	// panic("TODO: implement Circ.Area")
	return math.Pi * c.R * c.R
}

// TotalArea sums the areas of all shapes. It accepts the Shape interface, so it
// works for any current or future shape type.
func TotalArea(shapes ...Shape) float64 {
	// panic("TODO: implement TotalArea")
	var total float64
	for _, v := range shapes {
		total += v.Area()
	}
	return total
}

// Temperature in Celsius. Give it a String() method so it satisfies
// fmt.Stringer and prints like "23.5°C" (one decimal place).
type Temperature float64

// String makes Temperature a fmt.Stringer.
func (t Temperature) String() string {
	// panic("TODO: implement Temperature.String")
	return fmt.Sprintf("%.1f°C", t)
}

// ByArea attaches sort.Interface to a slice of Shape, sorting ascending by area.
type ByArea []Shape

func (a ByArea) Len() int {
	// panic("TODO: implement ByArea.Len")
	return len(a)
}

func (a ByArea) Less(i, j int) bool {
	// panic("TODO: implement ByArea.Less")
	return a[i].Area() < a[j].Area()
}

func (a ByArea) Swap(i, j int) {
	// panic("TODO: implement ByArea.Swap")
	a[i], a[j] = a[j], a[i]
}

// SortByArea sorts shapes in place, ascending by area, using sort.Sort and the
// ByArea type above.
func SortByArea(shapes []Shape) {
	// panic("TODO: implement SortByArea")

	sort.Sort(ByArea(shapes))
}

// DescribeType returns a label for the dynamic type of v using a type switch:
//
//	nil            -> "nil"
//	int            -> "int"
//	string         -> "string"
//	bool           -> "bool"
//	implements Shape -> "shape"
//	anything else  -> "unknown"
func DescribeType(v any) string {
	// panic("TODO: implement DescribeType")
	switch v.(type) {
	case nil:
		return "nil"
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	case Shape:
		return "shape"
	default:
		return "unknown"
	}
}
