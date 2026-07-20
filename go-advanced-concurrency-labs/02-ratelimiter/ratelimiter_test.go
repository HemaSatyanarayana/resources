package ratelimiter

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAllowBurst(t *testing.T) {
	// Big interval so no refill happens during the test.
	l := NewTokenBucket(3, 1, time.Hour)
	defer l.Stop()

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("Allow #%d = false, want true (bucket starts full with burst=3)", i+1)
		}
	}
	if l.Allow() {
		t.Errorf("Allow #4 = true, want false (bucket should be empty)")
	}
}

func TestRefill(t *testing.T) {
	l := NewTokenBucket(1, 1, 25*time.Millisecond)
	defer l.Stop()

	if !l.Allow() {
		t.Fatal("first Allow = false, want true")
	}
	if l.Allow() {
		t.Fatal("second Allow = true, want false (token not yet refilled)")
	}

	time.Sleep(40 * time.Millisecond) // let one refill tick happen
	if !l.Allow() {
		t.Error("Allow after refill = false, want true")
	}
}

func TestWaitBlocksUntilRefill(t *testing.T) {
	l := NewTokenBucket(1, 1, 30*time.Millisecond)
	defer l.Stop()

	if !l.Allow() {
		t.Fatal("setup Allow = false, want true")
	}

	start := time.Now()
	if err := l.Wait(context.Background()); err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
	if elapsed := time.Since(start); elapsed < 15*time.Millisecond {
		t.Errorf("Wait returned after %v — it did not block for the refill", elapsed)
	}
}

func TestWaitRespectsContext(t *testing.T) {
	l := NewTokenBucket(1, 1, time.Hour) // no refill will come in time
	defer l.Stop()

	if !l.Allow() {
		t.Fatal("setup Allow = false, want true")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := l.Wait(ctx)
	if err == nil {
		t.Fatal("Wait = nil, want context error (no token was ever available)")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("Wait err = %v, want context.DeadlineExceeded", err)
	}
}

func TestConcurrentAllowGivesExactlyBurst(t *testing.T) {
	const burst = 100
	l := NewTokenBucket(burst, 1, time.Hour) // no refill during test
	defer l.Stop()

	var granted int64
	var wg sync.WaitGroup
	for i := 0; i < 4*burst; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.Allow() {
				atomic.AddInt64(&granted, 1)
			}
		}()
	}
	wg.Wait()

	if granted != burst {
		t.Errorf("granted %d tokens, want exactly %d — bucket is not concurrency-safe", granted, burst)
	}
}

func TestStopIsIdempotent(t *testing.T) {
	l := NewTokenBucket(2, 1, 10*time.Millisecond)
	l.Stop()
	l.Stop() // must not panic (no double close)
}

func TestStopEndsRefill(t *testing.T) {
	// After Stop, no further refills should occur. Drain, stop, wait, still empty.
	l := NewTokenBucket(1, 1, 10*time.Millisecond)
	if !l.Allow() {
		t.Fatal("setup Allow = false, want true")
	}
	l.Stop()
	time.Sleep(40 * time.Millisecond)
	if l.Allow() {
		t.Error("Allow after Stop = true, want false (refill goroutine should be dead)")
	}
}
