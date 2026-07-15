# 04 — Select

`select` lets one goroutine wait on **multiple** channel operations at once. It's
how you build timeouts, cancellation, and fan-in.

## Concepts

### The basic form

```go
select {
case v := <-a:
    // a was ready first
case b <- x:
    // sending on b became possible
case <-time.After(time.Second):
    // nothing else was ready within 1s
}
```

- `select` blocks until **one** case can proceed, then runs exactly that one.
- If several are ready at once, it picks one **at random** (prevents starvation).
- An empty `select {}` blocks forever.

### Non-blocking with `default`

```go
select {
case v := <-ch:
    use(v)
default:
    // ch had nothing ready right now — don't block
}
```

`default` runs immediately if no other case is ready. Great for "try to receive,
but don't wait."

### Timeouts

`time.After(d)` returns a channel that delivers a value after `d`. Put it in a
`select` to bound how long you'll wait:

```go
select {
case v := <-ch:
    return v, true
case <-time.After(d):
    return 0, false   // timed out
}
```

### Cancellation with a `done` channel

A closed channel is *always* ready to receive (yielding the zero value). That
makes `<-done` a perfect "stop now" signal that every goroutine can watch:

```go
for {
    select {
    case <-done:
        return               // someone closed done → exit
    case out <- next():
        // produced a value
    }
}
```

### The nil-channel trick (fan-in that ends cleanly)

Receiving from a `nil` channel **blocks forever**. That's not a bug — it's a
tool: set a channel variable to `nil` to *disable* its `select` case. To merge
two channels until **both** close:

```go
for a != nil || b != nil {
    select {
    case v, ok := <-a:
        if !ok { a = nil; continue }  // a done → disable this case
        out <- v
    case v, ok := <-b:
        if !ok { b = nil; continue }
        out <- v
    }
}
```

Without this, a closed channel stays "ready" forever and floods you with zero
values.

## Your task

Create `select.go` (package `selectlab`) and implement, **from scratch**:

```go
// Recv waits up to timeout for a value from in. It returns (value, true) if one
// arrives in time, or (0, false) on timeout.
func Recv(in <-chan int, timeout time.Duration) (int, bool)

// TryRecv attempts a non-blocking receive from in. Returns (value, true) if a
// value was immediately available, else (0, false).
func TryRecv(in <-chan int) (int, bool)

// Merge fans two channels into one output channel that is closed only after
// BOTH inputs are closed and fully drained. Use the nil-channel trick.
func Merge(a, b <-chan int) <-chan int
```

## Run

```bash
go test -race -v ./04-select/
```

## Hints

- `Recv`: two-case `select` — `case v := <-in` and `case <-time.After(timeout)`.
- `TryRecv`: two-case `select` with a `default`.
- `Merge`: start a goroutine that runs the nil-channel loop above, sending to
  `out`, and `close(out)` (via `defer`) once both inputs are nil. Return `out`.

<details>
<summary>Reference solution</summary>

```go
package selectlab

import "time"

func Recv(in <-chan int, timeout time.Duration) (int, bool) {
	select {
	case v := <-in:
		return v, true
	case <-time.After(timeout):
		return 0, false
	}
}

func TryRecv(in <-chan int) (int, bool) {
	select {
	case v := <-in:
		return v, true
	default:
		return 0, false
	}
}

func Merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}
```

</details>
