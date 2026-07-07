// Package generics drills type parameters, constraints, and generic helpers.
package generics

// Number is a type constraint: any of these underlying numeric kinds. The ~
// means "any type whose underlying type is this", so `type Celsius float64`
// also satisfies Number.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Map applies f to every element of s and returns a new slice of results.
// T is the input element type, U the output element type.
func Map[T, U any](s []T, f func(T) U) []U {
	// panic("TODO: implement Map")
	var result = make([]U, 0, len(s))
	for _, v := range s {
		result = append(result, f(v))
	}
	return result
}

// Filter returns a new slice containing only the elements for which pred is true.
func Filter[T any](s []T, pred func(T) bool) []T {
	// panic("TODO: implement Filter")
	var result = make([]T, 0, len(s))

	for _, v := range s {
		if pred(v) {
			result = append(result, v)
		}
	}

	return result
}

// Reduce folds s into a single value, starting from init and combining with f.
// Example: Reduce([1,2,3], 0, func(acc, x int) int { return acc + x }) == 6
func Reduce[T, U any](s []T, init U, f func(U, T) U) U {
	// panic("TODO: implement Reduce")
	var ans = init
	for _, v := range s {
		ans = f(ans, v)
	}
	return ans
}

// Sum adds up all elements. The Number constraint lets it work for ints,
// floats, and any type whose underlying type is numeric.
func Sum[T Number](s []T) T {
	// panic("TODO: implement Sum")
	var result T = 0
	for _, v := range s {
		result = result + v
	}

	return T(result)
}

// Keys returns the keys of m in an unspecified order. K must be comparable
// (a requirement for map keys); V can be anything.
func Keys[K comparable, V any](m map[K]V) []K {
	// panic("TODO: implement Keys")
	var result = make([]K, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}
