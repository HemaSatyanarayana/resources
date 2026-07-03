# 04 — Structs & Methods

Structs group data; methods attach behavior. The single most important decision here is **value receiver vs pointer receiver**.

## Concepts

- **Value receiver `func (r Rectangle)`** operates on a *copy*. Great for read-only methods and small structs.
- **Pointer receiver `func (c *Counter)`** operates on the original — required to **mutate** state, and cheaper for large structs.
- **Consistency:** if any method needs a pointer receiver, give *all* methods on that type pointer receivers.
- **The zero value should be useful.** `var c Counter` should already work (count 0).
- **Embedding** (a field with no name) **promotes** the embedded type's fields and methods to the outer type — Go's composition-over-inheritance answer.

## Your task

Implement the methods in [`shapes.go`](shapes.go). You'll need to add `import "math"` (for `Circle`) and `import "fmt"` (for `Box.Describe`).

| Type | Methods | Receiver |
|------|---------|----------|
| `Rectangle` | `Area`, `Perimeter` | value |
| `Circle` | `Area`, `Perimeter` | value |
| `Counter` | `Inc`, `Add`, `Value` | pointer (except `Value`) |
| `Box` | `Describe` (uses *promoted* `Area`) | value |

## Run

```bash
go test -v ./exercises/04-structs-methods/
```

## Hints

- `Counter.Inc` **must** use `*Counter` or the increment is lost on the caller's copy — this is the classic beginner bug the test checks for.
- `fmt.Sprintf("%s: %.2f", b.Label, b.Area())` gives `"crate: 6.00"`.
- Don't declare `Area` on `Box` — the whole point is that it's promoted from the embedded `Rectangle`.

<details>
<summary>Reference solution</summary>

```go
package shapes

import (
	"fmt"
	"math"
)

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

func (c *Counter) Inc()      { c.count++ }
func (c *Counter) Add(n int) { c.count += n }
func (c Counter) Value() int { return c.count }

func (b Box) Describe() string {
	return fmt.Sprintf("%s: %.2f", b.Label, b.Area())
}
```

</details>
