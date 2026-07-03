# 01 — Fundamentals

Warm up with Go's control flow. No structs, no interfaces — just loops, conditionals, `switch`, variadic functions, and Go's signature multiple-return-value pattern.

## Concepts

- **`for` is Go's only loop.** No `while`, no `do`. `for cond {}` is a while loop; `for {}` is infinite; `for i := 0; i < n; i++ {}` is the classic.
- **`switch` needs no `break`** and cases can be expressions: `switch { case x%15==0: ... }`.
- **Multiple returns** are idiomatic, especially `(value, error)`.
- **Variadic params** (`nums ...int`) arrive as a slice; call with `Max(a...)` to spread.
- **Sentinel errors** are package-level `var Err... = errors.New(...)` values callers compare with `errors.Is`.

## Your task

Implement the four functions in [`fundamentals.go`](fundamentals.go):

| Function | Skill |
|----------|-------|
| `FizzBuzz(n)` | Loops, `switch`, `strconv.Itoa` |
| `Max(nums...)` | Variadic, ranging, returning an error |
| `IsPrime(n)` | Loop bounds, early return |
| `GCD(a, b)` | Euclid's algorithm (loop or recursion) |

## Run

```bash
go test -v ./exercises/01-fundamentals/
```

## Hints

- For `FizzBuzz`, check divisibility by 15 (or by both 3 and 5) **first**.
- `strconv.Itoa(i)` turns an `int` into its decimal string.
- For `IsPrime`, you only need to test divisors up to `√n`: `for i := 2; i*i <= n; i++`.
- Euclid: `for b != 0 { a, b = b, a%b }; return a`.

<details>
<summary>Reference solution (try first!)</summary>

```go
package fundamentals

import (
	"errors"
	"strconv"
)

var ErrEmpty = errors.New("fundamentals: no values provided")

func FizzBuzz(n int) []string {
	out := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			out = append(out, "FizzBuzz")
		case i%3 == 0:
			out = append(out, "Fizz")
		case i%5 == 0:
			out = append(out, "Buzz")
		default:
			out = append(out, strconv.Itoa(i))
		}
	}
	return out
}

func Max(nums ...int) (int, error) {
	if len(nums) == 0 {
		return 0, ErrEmpty
	}
	m := nums[0]
	for _, n := range nums[1:] {
		if n > m {
			m = n
		}
	}
	return m, nil
}

func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func GCD(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
```

</details>
