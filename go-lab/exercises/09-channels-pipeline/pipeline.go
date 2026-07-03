// Package pipeline drills channels: directional types, closing, ranging,
// fan-in with select, and building composable stages.
package pipeline

import "sync"

// Gen returns a receive-only channel that emits each of nums in order, then
// closes. Do the sending in a goroutine so Gen returns immediately.
func Gen(nums ...int) <-chan int {
	panic("TODO: implement Gen")
}

// Square is a pipeline STAGE: it reads every value from in, emits its square on
// the returned channel, and closes the returned channel when in is drained.
func Square(in <-chan int) <-chan int {
	panic("TODO: implement Square")
}

// Merge fans several input channels into one. It forwards every value from
// every input to the returned channel, and closes it only after ALL inputs are
// drained. Use a sync.WaitGroup to know when to close.
func Merge(ins ...<-chan int) <-chan int {
	panic("TODO: implement Merge")
}

// Collect drains in completely and returns all received values as a slice.
// (This is how a pipeline's final stage is consumed.)
func Collect(in <-chan int) []int {
	panic("TODO: implement Collect")
}

// FirstOf returns the first value to arrive on either a or b using select.
// Exactly one receive should happen.
func FirstOf(a, b <-chan int) int {
	panic("TODO: implement FirstOf")
}

// ensure the sync import is referenced by the package even before you implement
// Merge; delete this once you use sync.WaitGroup.
var _ = sync.WaitGroup{}
