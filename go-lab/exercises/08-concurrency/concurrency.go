// Package concurrency drills goroutines, sync.WaitGroup, sync.Mutex, and a
// simple worker pool. Always test this package with -race.
package concurrency

import "sync"

// ParallelSquare squares every element of nums, doing each square in its own
// goroutine, and returns the results IN THE SAME ORDER as the input.
//
// Trick: pre-allocate the result slice and have each goroutine write to its own
// index. Distinct indices don't race, so no mutex is needed — but you DO need a
// WaitGroup to know when all goroutines are done.
func ParallelSquare(nums []int) []int {
	panic("TODO: implement ParallelSquare")
}

// SafeCounter is a counter safe for concurrent use. Its zero value is ready.
type SafeCounter struct {
	mu sync.Mutex
	n  int
}

// Inc increments the counter under the lock.
func (c *SafeCounter) Inc() {
	panic("TODO: implement SafeCounter.Inc")
}

// Value returns the current count under the lock.
func (c *SafeCounter) Value() int {
	panic("TODO: implement SafeCounter.Value")
}

// ConcurrentSum sums nums using `workers` goroutines. Send each number onto a
// jobs channel; each worker adds what it receives to a partial total and, when
// the jobs channel is closed and drained, sends its partial onto a results
// channel. The caller sums the partials.
//
// If workers <= 0, treat it as 1.
func ConcurrentSum(nums []int, workers int) int {
	panic("TODO: implement ConcurrentSum")
}
