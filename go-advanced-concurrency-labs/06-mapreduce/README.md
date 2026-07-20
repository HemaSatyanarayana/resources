# 06 — MapReduce Word Count

MapReduce is the canonical big-data pattern: **map** a function over a huge input
in parallel, then **reduce** the partial results into one answer. In Go it's a
three-stage pipeline wired with channels — the culmination of everything the
first course taught (goroutines, channels, `WaitGroup`, `select`, context) and
the composition skills from labs 01–05. Your job: count every word across a set
of documents, in parallel, with a result that's **byte-for-byte identical no
matter how many workers you use.**

## The system

```
        source                 map (fan-out, N mappers)          reduce (fan-in)
 docs ─▶ chan string ─┬─▶ mapper: tokenize + count locally ─┐
                      ├─▶ mapper: ...            (local map) ┼─▶ chan map ─▶ sum into result
                      └─▶ mapper: ...                        ┘
```

- **`WordCount(ctx, docs, workers)`** → `map[string]int`. Splits the work across
  `workers` mapper goroutines and returns the merged counts.
- Tokenizing is done for you (`tokenize`): lowercase, split on non-alphanumeric.
  So `"Hello, world!"` → `["hello", "world"]`. Focus on the concurrency.

### The three stages

1. **Source.** One goroutine streams `docs` onto an unbuffered `chan string` and
   closes it when done. It `select`s the send against `<-ctx.Done()` so a
   cancelled context stops the feed immediately instead of pushing every doc.

2. **Map (fan-out).** `workers` goroutines each `range` over the source channel.
   Because a `range` over a shared channel hands each value to exactly one
   receiver, the documents are load-balanced across mappers automatically — no
   coordination needed. Each mapper accumulates into its **own local
   `map[string]int`**, so mappers never share memory and there's nothing to lock.
   When the source is drained, each mapper emits its local map on a `partials`
   channel.

3. **Reduce (fan-in).** The main goroutine `range`s over `partials` and sums each
   local map into the final `result`. A separate `go func(){ wg.Wait();
   close(partials) }()` closes the fan-in channel once every mapper has emitted,
   which ends the reduce loop.

### Why the result is deterministic

Every word's total is a **sum** of per-mapper subtotals, and integer addition is
associative and commutative — so it doesn't matter which mapper saw which
document, or in what order the reducer merges the partials. `workers=1` and
`workers=64` produce the exact same map. (Contrast with anything order-dependent,
like building a list, which would *not* be deterministic here.)

### The "local map" trick

The reason this needs no mutex at all: each mapper writes only to its own map,
and the reducer reads a mapper's map only *after* receiving it on a channel —
which is a happens-before edge. Shared-nothing mapping + channel handoff = a
race-free pipeline with zero locks. That's the payoff of the whole course.

### Cancellation

Both the source send and the mapper's partial-emit `select` on `<-ctx.Done()`.
On cancellation the source stops feeding, mappers finish their (now short) range
and exit, the fan-in closes, and `WordCount` returns a possibly-partial map
promptly — without leaking a goroutine or deadlocking.

## Your task

Implement, in `mapreduce.go` (package `mapreduce`):

```go
func WordCount(ctx context.Context, docs []string, workers int) map[string]int
```

Deterministic across worker counts. `workers < 1` behaves as `1`. Empty input
returns an empty (non-nil) map. Cancellation returns promptly. Use `tokenize`.

## Run

```bash
go test -race -v ./06-mapreduce/
```

## Hints

- Source: `go func(){ defer close(source); for _, d := range docs { select { case source <- d: case <-ctx.Done(): return } } }()`.
- Mapper: local map, `for d := range source { for _, w := range tokenize(d) { local[w]++ } }`, then `select { case partials <- local: case <-ctx.Done(): }`.
- Closer: `go func(){ wg.Wait(); close(partials) }()`.
- Reducer: `for local := range partials { for w, n := range local { result[w] += n } }`.
- Don't share a single result map across mappers — that's a data race and needs
  a lock. Local maps + fan-in is the point.

<details>
<summary>Reference solution</summary>

```go
package mapreduce

import (
	"context"
	"strings"
	"sync"
	"unicode"
)

func tokenize(s string) []string {
	return strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

func WordCount(ctx context.Context, docs []string, workers int) map[string]int {
	if workers < 1 {
		workers = 1
	}

	// Stage 1: source.
	source := make(chan string)
	go func() {
		defer close(source)
		for _, d := range docs {
			select {
			case source <- d:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Stage 2: map (fan-out).
	partials := make(chan map[string]int)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			local := make(map[string]int)
			for d := range source {
				for _, w := range tokenize(d) {
					local[w]++
				}
			}
			select {
			case partials <- local:
			case <-ctx.Done():
			}
		}()
	}
	go func() {
		wg.Wait()
		close(partials)
	}()

	// Stage 3: reduce (fan-in).
	result := make(map[string]int)
	for local := range partials {
		for w, n := range local {
			result[w] += n
		}
	}
	return result
}
```

Notice there is not a single mutex. The concurrency is entirely expressed as
channel ownership: the source owns `source` (and closes it), each mapper owns its
own `local` map until it hands it off, and the closer owns `partials`. Ownership
handoffs *are* the synchronization — the lesson the whole course was building to.

</details>
