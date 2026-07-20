package pubsub

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// recvWithin returns the next value from ch, or fails if nothing arrives soon.
func recvWithin[T any](t *testing.T, ch <-chan T, d time.Duration) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(d):
		t.Fatal("timed out waiting for a message")
		var zero T
		return zero
	}
}

// assertNothing fails if any value arrives on ch within d.
func assertNothing[T any](t *testing.T, ch <-chan T, d time.Duration) {
	t.Helper()
	select {
	case v, ok := <-ch:
		if ok {
			t.Fatalf("received unexpected message %v", v)
		}
	case <-time.After(d):
	}
}

func TestPublishSubscribe(t *testing.T) {
	b := NewBroker[string]()
	defer b.Close()

	sub, cancel := b.Subscribe("news")
	defer cancel()

	b.Publish("news", "hello")
	b.Publish("news", "world")

	if got := recvWithin(t, sub, time.Second); got != "hello" {
		t.Errorf("msg 1 = %q, want %q", got, "hello")
	}
	if got := recvWithin(t, sub, time.Second); got != "world" {
		t.Errorf("msg 2 = %q, want %q", got, "world")
	}
}

func TestTopicsAreIsolated(t *testing.T) {
	b := NewBroker[int]()
	defer b.Close()

	sub, cancel := b.Subscribe("a")
	defer cancel()

	b.Publish("b", 42) // different topic
	assertNothing(t, sub, 50*time.Millisecond)
}

func TestFanOutToAllSubscribers(t *testing.T) {
	b := NewBroker[int]()
	defer b.Close()

	subs := make([]<-chan int, 3)
	for i := range subs {
		s, cancel := b.Subscribe("t")
		defer cancel()
		subs[i] = s
	}

	b.Publish("t", 7)
	for i, s := range subs {
		if got := recvWithin(t, s, time.Second); got != 7 {
			t.Errorf("subscriber %d got %d, want 7", i, got)
		}
	}
}

func TestUnsubscribe(t *testing.T) {
	b := NewBroker[int]()
	defer b.Close()

	sub, cancel := b.Subscribe("t")
	cancel()

	// After cancel, the channel must be closed...
	if _, ok := <-sub; ok {
		t.Error("channel not closed after cancel")
	}
	// ...and publishing must not panic or deliver.
	b.Publish("t", 1)
}

func TestCancelIsIdempotent(t *testing.T) {
	b := NewBroker[int]()
	defer b.Close()

	_, cancel := b.Subscribe("t")
	cancel()
	cancel() // must not panic (no double close)
}

func TestCloseClosesAllChannels(t *testing.T) {
	b := NewBroker[int]()

	s1, _ := b.Subscribe("a")
	s2, _ := b.Subscribe("b")

	b.Close()

	if _, ok := <-s1; ok {
		t.Error("s1 not closed after Close")
	}
	if _, ok := <-s2; ok {
		t.Error("s2 not closed after Close")
	}

	b.Close() // idempotent
}

func TestSubscribeAfterCloseIsSafe(t *testing.T) {
	b := NewBroker[int]()
	b.Close()

	sub, cancel := b.Subscribe("t") // must not panic
	cancel()                        // must not panic
	if _, ok := <-sub; ok {
		t.Error("subscribing after Close should return an already-closed channel")
	}
	b.Publish("t", 1) // must not panic
}

func TestSlowSubscriberDoesNotBlockPublish(t *testing.T) {
	b := NewBroker[int]()
	defer b.Close()

	// Never read from this subscriber. Publishing far more than any buffer must
	// still return promptly (messages are dropped for the slow subscriber).
	_, cancel := b.Subscribe("t")
	defer cancel()

	done := make(chan struct{})
	go func() {
		for i := 0; i < 100_000; i++ {
			b.Publish("t", i)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Publish blocked on a slow subscriber")
	}
}

func TestConcurrentUse(t *testing.T) {
	b := NewBroker[int]()
	defer b.Close()

	var received int64
	var wg sync.WaitGroup

	// Subscribers that come and go.
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub, cancel := b.Subscribe("t")
			defer cancel()
			deadline := time.After(100 * time.Millisecond)
			for {
				select {
				case _, ok := <-sub:
					if !ok {
						return
					}
					atomic.AddInt64(&received, 1)
				case <-deadline:
					return
				}
			}
		}()
	}

	// Publishers.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				b.Publish("t", j)
			}
		}()
	}

	wg.Wait()
	// No assertion on the exact count (drops are allowed) — this test exists to
	// be run under -race to prove there are no data races or panics.
	_ = received
}
