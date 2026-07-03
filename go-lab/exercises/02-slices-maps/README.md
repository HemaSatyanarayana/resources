# 02 — Slices & Maps

Slices and maps are Go's workhorse data structures. This is where most Go bugs are born, so learn the idioms cold.

## Concepts

- **A slice is a view** (pointer, length, capacity) over a backing array. Copying a slice header is cheap but two slices can share the same array — mutating one can surprise the other.
- **`append` may reallocate.** `s = append(s, x)` — always reassign the result.
- **Pre-size when you can:** `make([]T, 0, n)` avoids repeated growth.
- **Maps are references.** A `nil` map reads fine but panics on write; `make(map[K]V)` first.
- **Comma-ok:** `v, ok := m[k]` distinguishes "absent" from "zero value".
- **A `map[T]struct{}`** (or `map[T]bool`) is Go's idiomatic **set**.

## Your task

Implement the four functions in [`collections.go`](collections.go):

| Function | Skill |
|----------|-------|
| `Reverse` | Build a *new* slice; don't mutate the input |
| `Dedup` | Map-as-set, preserve first-seen order |
| `WordCount` | `strings.Fields`, map increment |
| `Chunk` | Slicing with `s[i:j]`, edge cases |

## Run

```bash
go test -v ./exercises/02-slices-maps/
```

## Hints

- To avoid aliasing bugs in `Chunk`, slice the original directly — the tests compare values, not backing arrays. But beware: appending into a sub-slice of the original could clobber neighbors. Slicing with three indices `s[i:j:j]` caps capacity and is the safe pro move.
- `strings.Fields` handles arbitrary runs of whitespace and returns no empty strings.
- Incrementing a missing map key just works: `m[w]++` starts from 0.
- `min(a, b)` is a builtin since Go 1.21 — handy for the last chunk bound.

<details>
<summary>Reference solution</summary>

```go
package collections

import "strings"

func Reverse(s []int) []int {
	out := make([]int, len(s))
	for i, v := range s {
		out[len(s)-1-i] = v
	}
	return out
}

func Dedup(s []int) []int {
	seen := make(map[int]struct{}, len(s))
	out := make([]int, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func WordCount(text string) map[string]int {
	counts := make(map[string]int)
	for _, w := range strings.Fields(text) {
		counts[w]++
	}
	return counts
}

func Chunk(s []int, size int) [][]int {
	if size <= 0 {
		return nil
	}
	var out [][]int
	for i := 0; i < len(s); i += size {
		end := min(i+size, len(s))
		out = append(out, s[i:end:end])
	}
	return out
}
```

</details>
