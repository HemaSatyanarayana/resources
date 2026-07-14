// Package channels — Lab 03: channel fundamentals.
//
// Implement everything below FROM SCRATCH. Delete these comments as you go.
// Read README.md first; the test file is the spec.
//
// Build:
//
//	Generate(nums ...int) <-chan int   // emit then close, in a goroutine
//	Drain(in <-chan int) []int         // range to a slice
//	Take(in <-chan int, n int) []int   // up to n values, stop on close
//	Buffered(nums ...int) <-chan int   // pre-filled, already-closed buffered chan
//
// Run: go test -race -v ./03-channels/
package channels

func Generate(nums ...int) <-chan int {
	var numsProducer = make(chan int)

	go func() {
		defer close(numsProducer)
		for _, val := range nums {
			numsProducer <- val
		}
	}()

	return numsProducer
}

func Drain(in <-chan int) []int {
	var nums = make([]int, 0)

	for val := range in {
		nums = append(nums, val)
	}

	return nums
}

func Take(in <-chan int, n int) []int {
	var nums = make([]int, 0, n)

	for range n {
		val, ok := <-in
		if !ok {
			break
		}
		nums = append(nums, val)
	}

	return nums
}

func Buffered(nums ...int) <-chan int {
	var out = make(chan int)
	go func() {
		defer close(out)
		for _, num := range nums {
			out <- num
		}
	}()
	return out
}
