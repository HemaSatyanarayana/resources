# 06 — Pipelines

A pipeline is a series of stages connected by channels. Each stage is a goroutine
that **receives** from an inbound channel, does some work, and **sends** to an
outbound one. Stages compose like Lego bricks — and the hard part is shutting
them all down cleanly.

## Concepts

### The stage shape

Every stage looks the same:

```go
func stage(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)          // always close what you own
        for v := range in {       // ends when `in` closes
            out <- transform(v)
        }
    }()
    return out
}
```

Compose them by feeding one into the next:

```go
nums := Gen(done, 1, 2, 3, 4)
evens := Filter(done, nums, isEven)
squared := Map(done, evens, square)
result := Collect(squared)
```

### The leak problem

What if the consumer stops early — takes 2 values and walks away? The upstream
stages are still blocked on `out <- v` with nobody receiving. Those goroutines
**leak**: they live forever, holding memory. In a long-running server this is a
slow death.

### The fix: a `done` channel

Pass a `done` channel into every stage. Each stage's send becomes a `select`
that *also* watches `done`. When the consumer closes `done`, every stage's send
unblocks and the goroutine returns:

```go
func stage(done <-chan struct{}, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for v := range in {
            select {
            case out <- transform(v):
            case <-done:
                return            // consumer bailed → stop, no leak
            }
        }
    }()
    return out
}
```

Closing `done` once unblocks the whole chain. This is the classic Go pipeline
cancellation pattern (and the direct ancestor of `context.Context`, Lab 07).

> Convention: `done` is `chan struct{}` because it carries no data — only the
> signal of being closed. `struct{}` uses zero bytes.

## Your task

Create `pipeline.go` (package `pipe`) and implement, **from scratch**, every
stage taking a `done <-chan struct{}` first parameter:

```go
func Gen(done <-chan struct{}, nums ...int) <-chan int
func Map(done <-chan struct{}, in <-chan int, f func(int) int) <-chan int
func Filter(done <-chan struct{}, in <-chan int, keep func(int) bool) <-chan int
func Collect(in <-chan int) []int
```

`Gen` emits `nums`; `Map` transforms; `Filter` passes through only values where
`keep` returns true; `Collect` drains to a slice. Every producing stage must
respect `done`.

## Run

```bash
go test -race -v ./06-pipeline/
```

## Hints

- In `Gen`, the send is `select { case out <- n: case <-done: return }`.
- In `Filter`, only send when `keep(v)` is true; skip otherwise.
- `Collect` just `range`s — it's the consumer; it doesn't need `done`.
- Test the leak fix by closing `done` early and confirming stages return.

<details>
<summary>Reference solution</summary>

```go
package pipe

func Gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func Map(done <-chan struct{}, in <-chan int, f func(int) int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case out <- f(v):
			case <-done:
				return
			}
		}
	}()
	return out
}

func Filter(done <-chan struct{}, in <-chan int, keep func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			if !keep(v) {
				continue
			}
			select {
			case out <- v:
			case <-done:
				return
			}
		}
	}()
	return out
}

func Collect(in <-chan int) []int {
	var out []int
	for v := range in {
		out = append(out, v)
	}
	return out
}
```

</details>
