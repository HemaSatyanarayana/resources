# 05 — Interfaces

Interfaces are Go's abstraction mechanism, and they're **structural**: a type satisfies an interface simply by having the right methods. No `implements` keyword, no explicit declaration. This is the heart of idiomatic Go design.

## Concepts

- **Implicit satisfaction:** `Rect` satisfies `Shape` just by having `Area() float64`.
- **Accept interfaces, return structs:** functions take the smallest interface they need (`TotalArea(...Shape)`), so they work with types that don't exist yet.
- **`fmt.Stringer`** (`String() string`) customizes how a value prints with `%v`/`%s`.
- **`sort.Interface`** (`Len`, `Less`, `Swap`) lets `sort.Sort` order any collection.
- **Type switch** (`switch v := x.(type)`) branches on an `any`'s dynamic type.

## Your task

Implement everything in [`interfaces.go`](interfaces.go):

| Piece | Skill |
|-------|-------|
| `Rect.Area`, `Circ.Area` | Satisfy `Shape` implicitly |
| `TotalArea` | Program to the interface |
| `Temperature.String` | Implement `fmt.Stringer` |
| `ByArea` + `SortByArea` | Implement `sort.Interface` |
| `DescribeType` | Type switch over `any` |

## Run

```bash
go test -v ./exercises/05-interfaces/
```

## Hints

- `fmt.Sprintf("%.1f°C", float64(t))` for the `Temperature` — note the explicit conversion, or you'll recurse into `String()` forever if you use `%v` on `t` itself!
- `sort.Sort(ByArea(shapes))` — the conversion is how you attach the sort methods.
- In the type switch, put the `nil` case first (`case nil:`) and the `Shape` case before `default`.

<details>
<summary>Reference solution</summary>

```go
package shapesx

import (
	"fmt"
	"math"
	"sort"
)

func (r Rect) Area() float64 { return r.W * r.H }
func (c Circ) Area() float64 { return math.Pi * c.R * c.R }

func TotalArea(shapes ...Shape) float64 {
	var sum float64
	for _, s := range shapes {
		sum += s.Area()
	}
	return sum
}

func (t Temperature) String() string {
	return fmt.Sprintf("%.1f°C", float64(t))
}

func (a ByArea) Len() int           { return len(a) }
func (a ByArea) Less(i, j int) bool { return a[i].Area() < a[j].Area() }
func (a ByArea) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func SortByArea(shapes []Shape) { sort.Sort(ByArea(shapes)) }

func DescribeType(v any) string {
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
```

</details>
