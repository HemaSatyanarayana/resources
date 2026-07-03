# 07 — Generics

Generics (type parameters, added in Go 1.18) let you write one function that works for many types — without `interface{}` and runtime type assertions. This exercise builds the classic `Map`/`Filter`/`Reduce` toolkit.

## Concepts

- **Type parameters** go in square brackets: `func Map[T, U any](...)`.
- **Constraints** restrict what a type parameter can be. `any` allows anything; `comparable` allows `==`/map keys; a **constraint interface** (a set of types) allows operators.
- **The `~` token** ("approximation") means "any type whose *underlying* type is this" — so `type Celsius float64` still satisfies a `~float64` constraint.
- **Type inference** usually lets you call `Map(xs, f)` without spelling out `[int, string]`.
- **Don't over-genericize.** Reach for generics when you're truly writing the same logic for many types (containers, algorithms) — not by default.

## Your task

Implement the five generic functions in [`generics.go`](generics.go):

| Function | Constraint(s) |
|----------|---------------|
| `Map[T, U any]` | `any`, `any` |
| `Filter[T any]` | `any` |
| `Reduce[T, U any]` | `any`, `any` |
| `Sum[T Number]` | custom `Number` union |
| `Keys[K comparable, V any]` | `comparable`, `any` |

## Run

```bash
go test -v ./exercises/07-generics/
```

## Hints

- Pre-size result slices with `make([]U, 0, len(s))` — you know the max length.
- `Reduce`'s accumulator has type `U`; thread it through the loop: `acc = f(acc, x)`.
- `Sum` starts from the zero value `var total T` — which is `0` for every numeric type.
- Map iteration order is randomized in Go; that's why the test sorts `Keys` before comparing.

<details>
<summary>Reference solution</summary>

```go
package generics

func Map[T, U any](s []T, f func(T) U) []U {
	out := make([]U, 0, len(s))
	for _, v := range s {
		out = append(out, f(v))
	}
	return out
}

func Filter[T any](s []T, pred func(T) bool) []T {
	out := make([]T, 0, len(s))
	for _, v := range s {
		if pred(v) {
			out = append(out, v)
		}
	}
	return out
}

func Reduce[T, U any](s []T, init U, f func(U, T) U) U {
	acc := init
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

func Sum[T Number](s []T) T {
	var total T
	for _, v := range s {
		total += v
	}
	return total
}

func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
```

</details>
