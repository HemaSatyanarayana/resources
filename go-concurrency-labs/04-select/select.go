// Package selectlab — Lab 04: the select statement.
//
// Implement everything below FROM SCRATCH. Delete these comments as you go.
// Read README.md first; the test file is the spec.
//
// Build:
//
//	Recv(in <-chan int, timeout time.Duration) (int, bool)  // value or timeout
//	TryRecv(in <-chan int) (int, bool)                      // non-blocking (default case)
//	Merge(a, b <-chan int) <-chan int                       // fan-in via nil-channel trick
//
// Run: go test -race -v ./04-select/
package selectlab

import "time"

func Recv(in <-chan int, timeout time.Duration) (int, bool) {
	select {
	case val := <-in:
		return val, true
	case <-time.After(timeout):
		return 0, false
	}
}

func TryRecv(in <-chan int) (int, bool) {
	select {
	case val := <-in:
		return val, true
	default:
		return 0, false
	}
}

func Merge(a, b <-chan int) <-chan int {
	var out = make(chan int)

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
