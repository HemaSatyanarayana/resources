// Package ratelimiter — Lab 02: a concurrency-safe token-bucket rate limiter.
//
// Read README.md first; ratelimiter_test.go is the spec. Fill in every TODO.
// Run: go test -race -v ./02-ratelimiter/
package ratelimiter

import (
	"context"
	"time"
)

// Limiter hands out tokens from a bucket that refills over time.
type Limiter struct {
	// TODO: fields. A buffered `chan struct{}` makes a great bucket — its
	// length IS the number of available tokens, and channel ops are already
	// safe for concurrent use. You'll also want a `done chan struct{}` to stop
	// the refill goroutine, and a `sync.Once` so Stop is idempotent.
}

// NewTokenBucket returns a limiter that holds up to `burst` tokens and adds
// `refill` tokens every `interval`, never exceeding `burst`. It starts full.
// A background goroutine performs the refills; call Stop to end it.
// Treat burst < 1 and refill < 1 as 1.
func NewTokenBucket(burst, refill int, interval time.Duration) *Limiter {
	// TODO:
	//   1. Clamp burst/refill to a minimum of 1.
	//   2. Make the bucket channel with capacity `burst` and pre-fill it full.
	//   3. Launch a goroutine with a time.Ticker: on each tick, add up to
	//      `refill` tokens using a NON-blocking send (select/default) so it
	//      never blocks when the bucket is full; return when `done` closes.
	panic("TODO: implement NewTokenBucket")
}

// Allow consumes one token and returns true if one was available right now,
// otherwise returns false immediately (never blocks).
func (l *Limiter) Allow() bool {
	// TODO: non-blocking receive from the bucket (select with a default).
	panic("TODO: implement Allow")
}

// Wait blocks until a token is available or ctx is done. It returns nil once it
// consumes a token, or ctx.Err() if the context is cancelled/expired first.
func (l *Limiter) Wait(ctx context.Context) error {
	// TODO: select over (receive from bucket) and (<-ctx.Done()).
	panic("TODO: implement Wait")
}

// Stop halts the background refill goroutine. It is safe to call more than once.
func (l *Limiter) Stop() {
	// TODO: close `done` exactly once (sync.Once).
	panic("TODO: implement Stop")
}
