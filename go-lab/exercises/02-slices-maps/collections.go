// Package collections drills Go's slice and map idioms.
package collections

// Reverse returns a new slice with the elements of s in reverse order.
// The input slice must not be modified.
func Reverse(s []int) []int {
	panic("TODO: implement Reverse")
}

// Dedup returns a new slice containing the elements of s with duplicates
// removed, preserving the order of first appearance. Use a map as a set.
func Dedup(s []int) []int {
	panic("TODO: implement Dedup")
}

// WordCount splits text on whitespace (use strings.Fields) and returns a map
// from each word to the number of times it appears.
func WordCount(text string) map[string]int {
	panic("TODO: implement WordCount")
}

// Chunk splits s into consecutive sub-slices of at most size elements.
// The final chunk may be shorter. If size <= 0, return nil.
// Example: Chunk([1,2,3,4,5], 2) => [[1,2],[3,4],[5]]
func Chunk(s []int, size int) [][]int {
	panic("TODO: implement Chunk")
}
