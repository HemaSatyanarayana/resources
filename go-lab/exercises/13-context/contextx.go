// Package contextx drills context.Context: cancellation, deadlines, and
// request-scoped values — the backbone of well-behaved Go servers.
package contextx

import (
	"context"
	"time"
)

// SumWithContext adds up nums, simulating perItem of work before counting each
// one. Between items it must respect ctx: if ctx is cancelled or its deadline
// passes, stop immediately and return the sum-so-far plus ctx.Err().
// If it finishes all items, return the full sum and nil.
//
// Use a select over <-ctx.Done() and <-time.After(perItem).
func SumWithContext(ctx context.Context, nums []int, perItem time.Duration) (int, error) {
	panic("TODO: implement SumWithContext")
}

// ctxKey is an UNEXPORTED type used for context value keys. Using a private
// type prevents key collisions with other packages — never use a bare string.
type ctxKey int

const requestIDKey ctxKey = 0

// WithRequestID returns a copy of ctx carrying the given request id.
func WithRequestID(ctx context.Context, id string) context.Context {
	panic("TODO: implement WithRequestID")
}

// RequestID extracts the request id previously stored with WithRequestID.
// If none is present, it returns ("", false).
func RequestID(ctx context.Context) (string, bool) {
	panic("TODO: implement RequestID")
}

// Race runs work in a goroutine and returns whichever happens first:
//   - work finishes -> (its result, nil)
//   - ctx is done   -> (0, ctx.Err())
//
// Use a result channel and a select. (The goroutine may keep running after a
// ctx cancellation; that's acceptable for this exercise — make the channel
// buffered so the goroutine never blocks forever.)
func Race(ctx context.Context, work func() int) (int, error) {
	panic("TODO: implement Race")
}
