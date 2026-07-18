// Package pipe — Lab 06: composable pipeline stages with cancellation.
//
// Implement everything below FROM SCRATCH. Delete these comments as you go.
// Read README.md first; the test file is the spec.
//
// Build (every producing stage takes `done` first and must respect it):
//
//	Gen(done <-chan struct{}, nums ...int) <-chan int
//	Map(done <-chan struct{}, in <-chan int, f func(int) int) <-chan int
//	Filter(done <-chan struct{}, in <-chan int, keep func(int) bool) <-chan int
//	Collect(in <-chan int) []int
//
// Run: go test -race -v ./06-pipeline/
package pipe

func Gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int, len(nums))

	go func() {
		defer close(out)
		for _, num := range nums {
			select {
			case out <- num:
			case <-done:
				return
			}
		}
	}()

	return out
}

func Map(done <-chan struct{}, in <-chan int, f func(int) int) <-chan int {
	var out = make(chan int)

	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- f(n):
			case <-done:
				return
			}
		}
	}()

	return out
}

func Filter(done <-chan struct{}, in <-chan int, keep func(int) bool) <-chan int {
	var out = make(chan int)

	go func() {
		defer close(out)
		for i := range in {
			if !keep(i) {
				continue
			}
			select {

			case out <- i:
			case <-done:
				return

			}
		}
	}()

	return out
}

func Collect(in <-chan int) []int {
	var out []int

	for i := range in {
		out = append(out, i)
	}

	return out
}
