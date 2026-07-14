# 03 — Channels

Channels are typed pipes that let goroutines **communicate** and **synchronize**
by passing values. This is Go's preferred coordination tool.

## Concepts

### Making and using a channel

```go
ch := make(chan int)   // unbuffered
ch <- 42               // send (blocks until someone receives)
v := <-ch              // receive
```

### Unbuffered vs buffered

- **Unbuffered** `make(chan int)`: a send blocks until another goroutine is
  ready to receive — a *rendezvous*. The two goroutines hand off in lockstep.
- **Buffered** `make(chan int, 3)`: holds up to 3 values. Sends block only when
  the buffer is full; receives block only when it's empty.

A common beginner deadlock:

```go
ch := make(chan int)
ch <- 1        // DEADLOCK: no receiver, and we're the only goroutine
```

The send needs *another* goroutine ready to receive. Either buffer the channel
or do the send/receive from different goroutines.

### Closing & ranging

```go
close(ch)                 // signals "no more values will be sent"
for v := range ch { ... } // loops until ch is closed AND drained
v, ok := <-ch             // ok == false once closed and empty; v is the zero value
```

Rules:
- **Only the sender closes.** Closing from the receiver side, or closing twice,
  or sending after close → panic.
- Closing is a *broadcast*: every receiver sees it.
- You don't have to close every channel — only close when receivers use `range`
  or `ok` to detect completion.

### Directional channel types

Restrict direction in function signatures to document intent and catch misuse:

```go
func produce(out chan<- int)  // send-only
func consume(in <-chan int)   // receive-only
```

A bidirectional `chan int` converts to either automatically. Returning a
`<-chan int` from a producer says "you receive; I own closing."

## Your task

Create `channels.go` (package `channels`) and implement, **from scratch**:

```go
// Generate returns a receive-only channel that emits each of nums in order and
// is then closed. The sending happens in its own goroutine.
func Generate(nums ...int) <-chan int

// Drain receives every value from in until it is closed, returning them in
// order as a slice.
func Drain(in <-chan int) []int

// Take receives up to n values from in and returns them. If in is closed with
// fewer than n values, return what arrived (it must not block forever).
func Take(in <-chan int, n int) []int

// Buffered returns a buffered channel of capacity cap(nums) pre-filled with
// nums and already closed, so a caller can range over it without a separate
// producer goroutine.
func Buffered(nums ...int) <-chan int
```

## Run

```bash
go test -race -v ./03-channels/
```

## Hints

- `Generate`: `out := make(chan int); go func(){ defer close(out); for _, n := range nums { out <- n } }(); return out`.
- `Drain`: `for v := range in { out = append(out, v) }`.
- `Take`: loop `n` times using `v, ok := <-in`; **stop early if `!ok`** (closed),
  otherwise you'll read zero values forever or block.
- `Buffered`: `out := make(chan int, len(nums))`, send all, `close(out)`, return.
  Because it's buffered to exactly `len(nums)`, the sends never block even with
  no receiver yet.

<details>
<summary>Reference solution</summary>

```go
package channels

func Generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func Drain(in <-chan int) []int {
	var out []int
	for v := range in {
		out = append(out, v)
	}
	return out
}

func Take(in <-chan int, n int) []int {
	out := make([]int, 0, n)
	for i := 0; i < n; i++ {
		v, ok := <-in
		if !ok {
			break
		}
		out = append(out, v)
	}
	return out
}

func Buffered(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	for _, n := range nums {
		out <- n
	}
	close(out)
	return out
}
```

</details>
