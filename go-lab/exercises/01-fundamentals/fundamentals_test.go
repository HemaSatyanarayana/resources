package fundamentals

import (
	"errors"
	"reflect"
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	got := FizzBuzz(15)
	want := []string{
		"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8",
		"Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FizzBuzz(15)\n got = %v\nwant = %v", got, want)
	}
	if len(FizzBuzz(0)) != 0 {
		t.Errorf("FizzBuzz(0) should be empty, got %v", FizzBuzz(0))
	}
}

func TestMax(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{
		{[]int{3}, 3},
		{[]int{1, 9, 4, 9, 2}, 9},
		{[]int{-5, -2, -30}, -2},
	}
	for _, c := range cases {
		got, err := Max(c.in...)
		if err != nil {
			t.Errorf("Max(%v) unexpected error: %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("Max(%v) = %d, want %d", c.in, got, c.want)
		}
	}

	if _, err := Max(); !errors.Is(err, ErrEmpty) {
		t.Errorf("Max() error = %v, want ErrEmpty", err)
	}
}

func TestIsPrime(t *testing.T) {
	primes := map[int]bool{
		-7: false, 0: false, 1: false, 2: true, 3: true, 4: false,
		17: true, 18: false, 19: true, 97: true, 100: false,
	}
	for n, want := range primes {
		if got := IsPrime(n); got != want {
			t.Errorf("IsPrime(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestGCD(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{12, 8, 4},
		{54, 24, 6},
		{17, 5, 1},
		{0, 0, 0},
		{0, 9, 9},
	}
	for _, c := range cases {
		if got := GCD(c.a, c.b); got != c.want {
			t.Errorf("GCD(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}
