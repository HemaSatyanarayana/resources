// Package shapes drills structs, value vs pointer receivers, and embedding.
package shapes

// NOTE: Circle.Area/Perimeter need "math" (math.Pi) and Box.Describe needs
// "fmt". Add those imports yourself when you implement them.

// Rectangle is defined by its width and height.
type Rectangle struct {
	Width  float64
	Height float64
}

// Area returns width * height. Use a VALUE receiver — Area doesn't mutate.
func (r Rectangle) Area() float64 {
	panic("TODO: implement Rectangle.Area")
}

// Perimeter returns 2 * (width + height).
func (r Rectangle) Perimeter() float64 {
	panic("TODO: implement Rectangle.Perimeter")
}

// Circle is defined by its radius.
type Circle struct {
	Radius float64
}

// Area returns π r². Use math.Pi.
func (c Circle) Area() float64 {
	panic("TODO: implement Circle.Area")
}

// Perimeter returns the circumference, 2πr.
func (c Circle) Perimeter() float64 {
	panic("TODO: implement Circle.Perimeter")
}

// Counter accumulates a running total. Its zero value (count 0) is ready to use.
type Counter struct {
	count int
}

// Inc increases the counter by 1. It MUST use a pointer receiver so the change
// is visible to the caller.
func (c *Counter) Inc() {
	panic("TODO: implement Counter.Inc")
}

// Add increases the counter by n (pointer receiver).
func (c *Counter) Add(n int) {
	panic("TODO: implement Counter.Add")
}

// Value returns the current count (value receiver is fine — read only).
func (c Counter) Value() int {
	panic("TODO: implement Counter.Value")
}

// Box embeds a Rectangle and adds a Label. Because Rectangle is embedded
// (no field name), its Area/Perimeter methods are promoted to Box.
type Box struct {
	Label string
	Rectangle
}

// Describe returns a string like "crate: 6.00" where the number is the Box's
// area, formatted with two decimals. Demonstrate that b.Area() works via
// promotion — do NOT redeclare Area on Box.
func (b Box) Describe() string {
	panic("TODO: implement Box.Describe")
}
