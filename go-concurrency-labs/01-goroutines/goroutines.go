// Package goroutines — Lab 01: goroutines & sync.WaitGroup.
//
// Implement everything below FROM SCRATCH. Delete these comments as you go.
// Read README.md first; the test file is the spec.
//
// Build these two functions:
//
//	func ParallelMap(nums []int, f func(int) int) []int
//	    Apply f to each element in its own goroutine; return results in order.
//
//	func WaitAll(fns []func())
//	    Run each function in its own goroutine; return after all finish.
//
// Run: go test -race -v ./01-goroutines/
package goroutines

import "sync"

func ParallelMap(nums []int, f func(int) int) []int {
	var wg sync.WaitGroup

	var ans = make([]int, len(nums))

	for i, num := range nums {
		wg.Add(1)

		go func() {
			defer wg.Done()
			ans[i] = f(num)
		}()

	}
	wg.Wait()

	return ans
}

func WaitAll(fns []func()) {
	var wg sync.WaitGroup

	for _, f := range fns {
		wg.Add(1)

		go func() {
			defer wg.Done()
			f()
		}()
	}

	wg.Wait()
}
