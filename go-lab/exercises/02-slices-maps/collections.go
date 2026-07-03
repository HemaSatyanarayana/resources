// Package collections drills Go's slice and map idioms.
package collections

import "strings"

// Reverse returns a new slice with the elements of s in reverse order.
// The input slice must not be modified.
func Reverse(s []int) []int {
	// panic("TODO: implement Reverse")
	var ans []int

	for i := len(s) - 1; i >= 0; i-- {
		ans = append(ans, s[i])
	}
	return ans
}

// Dedup returns a new slice containing the elements of s with duplicates
// removed, preserving the order of first appearance. Use a map as a set.
func Dedup(s []int) []int {
	// panic("TODO: implement Dedup")
	var ans []int
	seen := make(map[int]struct{})

	for _, v := range s {
		if _, ok := seen[v]; !ok {
			ans = append(ans, v)
			seen[v] = struct{}{}
		}

	}

	return ans
}

// WordCount splits text on whitespace (use strings.Fields) and returns a map
// from each word to the number of times it appears.
func WordCount(text string) map[string]int {
	// panic("TODO: implement WordCount")
	var res map[string]int = make(map[string]int)

	words := strings.Fields(text)

	for _, word := range words {
		res[word]++
	}

	return res
}

// Chunk splits s into consecutive sub-slices of at most size elements.
// The final chunk may be shorter. If size <= 0, return nil.
// Example: Chunk([1,2,3,4,5], 2) => [[1,2],[3,4],[5]]
func Chunk(s []int, size int) [][]int {
	// panic("TODO: implement Chunk")
	if size <= 0 {
		return nil
	}

	ans := make([][]int, 0)

	for i := 0; i < len(s); i += size {
		ans = append(ans, s[i:min(i+size, len(s))])
	}

	return ans

}
