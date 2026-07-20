// Package mapreduce — Lab 06: a concurrent MapReduce word counter.
//
// Read README.md first; mapreduce_test.go is the spec. Fill in every TODO.
// Run: go test -race -v ./06-mapreduce/
package mapreduce

import (
	"context"
	"strings"
	"unicode"
)

// tokenize lowercases s and splits it into words on any non-letter/non-digit
// rune. This is provided so you can focus on the concurrency, not the parsing.
//
//	tokenize("Hello, world!") == []string{"hello", "world"}
func tokenize(s string) []string {
	return strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

// WordCount counts every word across docs using a three-stage, `workers`-wide
// pipeline: a source that streams docs, a fan-out of mapper goroutines that
// tokenize and count locally, and a single reducer that fans the partial counts
// back in. The result is identical no matter how many workers are used (integer
// addition is associative and commutative). If ctx is cancelled it returns
// promptly with a possibly-partial result. Treat workers < 1 as 1.
func WordCount(ctx context.Context, docs []string, workers int) map[string]int {
	// TODO: build the three stages.
	//
	// 0. Clamp workers < 1 -> 1.
	//
	// 1. SOURCE. Make `source := make(chan string)`. Launch one goroutine that
	//    ranges docs and, for each, `select { case source <- d: case <-ctx.Done(): return }`,
	//    then `close(source)` (defer it). The ctx case is what makes cancellation prompt.
	//
	// 2. MAP (fan-out). Make `partials := make(chan map[string]int)` and a
	//    sync.WaitGroup. Add(workers); launch `workers` goroutines. Each builds a
	//    LOCAL map, `for d := range source { for _, w := range tokenize(d) { local[w]++ } }`,
	//    then emits its local map: `select { case partials <- local: case <-ctx.Done(): }`.
	//    Separately: `go func(){ wg.Wait(); close(partials) }()` so the reducer
	//    knows when every mapper is done.
	//
	// 3. REDUCE (fan-in). `result := make(map[string]int)`; `for local := range
	//    partials { for w, n := range local { result[w] += n } }`; return result.
	//
	// Why it's deterministic: each word's total is a sum of per-mapper counts,
	// and addition doesn't care about order — so the map is identical for any
	// worker count. Why local maps: mappers never touch shared state, so there's
	// nothing to lock; the only synchronization is the channels.
	panic("TODO: implement WordCount")
}
