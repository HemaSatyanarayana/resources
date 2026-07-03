// Package fundamentals drills Go's control flow, loops, and multiple returns.
package fundamentals

import (
	"errors"
)

// ErrEmpty is returned by Max when called with no arguments.
var ErrEmpty = errors.New("fundamentals: no values provided")

// FizzBuzz returns a slice of length n where index i (0-based) holds the
// FizzBuzz string for the number i+1:
//   - "Fizz"     if divisible by 3
//   - "Buzz"     if divisible by 5
//   - "FizzBuzz" if divisible by both
//   - the decimal number otherwise (e.g. "7")
func FizzBuzz(n int) []string {
	panic("TODO: implement FizzBuzz")

}

// Max returns the largest of the provided values. If no values are given it
// returns ErrEmpty. Use the variadic signature and a range loop.
func Max(nums ...int) (int, error) {
	panic("TODO: implement Max")

}

// IsPrime reports whether n is a prime number. Numbers < 2 are not prime.
// Aim for an O(sqrt(n)) trial-division loop.
func IsPrime(n int) bool {
	panic("TODO: implement IsPrime")

}

// GCD returns the greatest common divisor of a and b using Euclid's algorithm.
// GCD(0, 0) is defined here as 0.
func GCD(a, b int) int {
	panic("TODO: implement GCD")

}
