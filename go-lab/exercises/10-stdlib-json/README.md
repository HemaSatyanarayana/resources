# 10 — Standard Library & JSON

Go's standard library is batteries-included. This capstone touches the two you'll reach for constantly: `encoding/json` and `sort`. You'll wire up struct tags, marshal/unmarshal, custom marshaling, and slice sorting.

## Concepts

- **Struct tags** are backtick strings after a field that packages read via reflection: `` `json:"name"` ``. `omitempty` drops zero-valued fields; `-` skips a field entirely.
- **Exported fields only.** `json` (and most reflection-based packages) can only see capitalized fields. `name string` won't marshal — `Name string` will.
- **`json.Marshal` / `json.Unmarshal`** convert between Go values and `[]byte`. Unmarshal takes a **pointer** to the destination.
- **Custom (un)marshaling:** implement `json.Marshaler` (`MarshalJSON() ([]byte, error)`) to control a type's wire format. The bytes you return must be *valid JSON* — a string needs its quotes.
- **`sort.Slice(s, less)`** sorts any slice in place given a `less(i, j) bool` closure. Multi-key sorts tie-break inside `less`.

## Your task

Complete [`json.go`](json.go): add the struct tags, then implement the functions. You'll need `import "encoding/json"`, `import "sort"`, and `import "fmt"`.

| Piece | Skill |
|-------|-------|
| `User` tags | `json:"..."`, `omitempty` |
| `ParseUser` | `json.Unmarshal` into a pointer |
| `ToJSON` | `json.Marshal` |
| `SortByAge` | `sort.Slice` with a tie-break |
| `HexColor.MarshalJSON` | `json.Marshaler`, `fmt.Sprintf` |

## Run

```bash
go test -v ./exercises/10-stdlib-json/
```

## Hints

- Tag syntax is exact: `` `json:"email,omitempty"` `` — no space after the comma.
- `json.Unmarshal(data, &u)` — pass the **address** of your `User`.
- Multi-key `less`: `if a.Age != b.Age { return a.Age < b.Age }; return a.Name < b.Name`.
- `fmt.Sprintf("\"#%06x\"", uint32(c))` yields a quoted, zero-padded, 6-digit lower-case hex string — exactly what a JSON string needs.

<details>
<summary>Reference solution</summary>

```go
package configjson

import (
	"encoding/json"
	"fmt"
	"sort"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"`
	Admin bool   `json:"is_admin"`
}

func ParseUser(data []byte) (User, error) {
	var u User
	err := json.Unmarshal(data, &u)
	return u, err
}

func ToJSON(u User) ([]byte, error) {
	return json.Marshal(u)
}

func SortByAge(users []User) {
	sort.Slice(users, func(i, j int) bool {
		if users[i].Age != users[j].Age {
			return users[i].Age < users[j].Age
		}
		return users[i].Name < users[j].Name
	})
}

func (c HexColor) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"#%06x\"", uint32(c))), nil
}
```

</details>
