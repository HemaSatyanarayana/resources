package taskqueue

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// collect drains every Result the pool emits into a slice, in the background.
// The pool's contract is that someone must consume Results() while it runs, so
// every test starts a collector, then calls Shutdown, then wait()s for the
// closed Results channel to hand back the full slice.
func collect(p *Pool) func() []Result {
	var mu sync.Mutex
	var out []Result
	done := make(chan struct{})
	go func() {
		defer close(done)
		for r := range p.Results() {
			mu.Lock()
			out = append(out, r)
			mu.Unlock()
		}
	}()
	return func() []Result {
		<-done
		return out
	}
}

func TestAllSucceed(t *testing.T) {
	p := NewPool(4, 3, time.Millisecond)
	wait := collect(p)

	const n = 50
	for i := 0; i < n; i++ {
		if !p.Submit(i, func() error { return nil }) {
			t.Fatalf("Submit(%d) = false, want true", i)
		}
	}
	p.Shutdown()
	res := wait()

	if len(res) != n {
		t.Fatalf("got %d results, want %d", len(res), n)
	}
	for _, r := range res {
		if r.Err != nil {
			t.Errorf("job %d err = %v, want nil", r.ID, r.Err)
		}
		if r.Attempts != 1 {
			t.Errorf("job %d attempts = %d, want 1 (succeeded first try)", r.ID, r.Attempts)
		}
	}
}

func TestRetryThenSucceed(t *testing.T) {
	p := NewPool(2, 5, time.Millisecond)
	wait := collect(p)

	var calls int32
	p.Submit(1, func() error {
		if atomic.AddInt32(&calls, 1) < 3 {
			return errors.New("boom")
		}
		return nil
	})
	p.Shutdown()
	res := wait()

	if len(res) != 1 {
		t.Fatalf("got %d results, want 1", len(res))
	}
	if res[0].Err != nil {
		t.Errorf("err = %v, want nil (job succeeds on its 3rd attempt)", res[0].Err)
	}
	if res[0].Attempts != 3 {
		t.Errorf("attempts = %d, want 3", res[0].Attempts)
	}
}

func TestRetriesExhausted(t *testing.T) {
	p := NewPool(2, 3, time.Millisecond)
	wait := collect(p)

	p.Submit(7, func() error { return errors.New("always fails") })
	p.Shutdown()
	res := wait()

	if len(res) != 1 {
		t.Fatalf("got %d results, want 1", len(res))
	}
	if res[0].Attempts != 3 {
		t.Errorf("attempts = %d, want 3 (should try up to maxAttempts)", res[0].Attempts)
	}
	if res[0].Err == nil {
		t.Error("err = nil, want the job's error after exhausting retries")
	}
}

func TestBackoffDelaysRetries(t *testing.T) {
	const backoff = 20 * time.Millisecond
	p := NewPool(1, 3, backoff)
	wait := collect(p)

	var calls int32
	start := time.Now()
	p.Submit(1, func() error {
		if atomic.AddInt32(&calls, 1) < 3 {
			return errors.New("retry me")
		}
		return nil
	})
	p.Shutdown()
	res := wait()
	elapsed := time.Since(start)

	if res[0].Attempts != 3 {
		t.Fatalf("attempts = %d, want 3", res[0].Attempts)
	}
	// Three attempts means two failures, so two backoff gaps must have elapsed.
	if elapsed < 2*backoff {
		t.Errorf("elapsed %v, want >= %v — backoff not applied between retries", elapsed, 2*backoff)
	}
}

func TestSubmitAfterShutdown(t *testing.T) {
	p := NewPool(2, 1, time.Millisecond)
	wait := collect(p)
	p.Shutdown()
	_ = wait()

	if p.Submit(1, func() error { return nil }) {
		t.Error("Submit after Shutdown = true, want false")
	}
}

func TestShutdownIsIdempotent(t *testing.T) {
	p := NewPool(2, 1, time.Millisecond)
	wait := collect(p)
	p.Submit(1, func() error { return nil })
	p.Shutdown()
	p.Shutdown() // must not panic (no double close of the results channel)
	_ = wait()
}

func TestGracefulDrain(t *testing.T) {
	p := NewPool(3, 1, time.Millisecond)
	wait := collect(p)

	var ran int32
	const n = 200
	for i := 0; i < n; i++ {
		p.Submit(i, func() error {
			atomic.AddInt32(&ran, 1)
			return nil
		})
	}
	p.Shutdown() // must block until every queued job has finished
	res := wait()

	if got := atomic.LoadInt32(&ran); got != n {
		t.Errorf("%d jobs ran, want %d — Shutdown did not drain the queue", got, n)
	}
	if len(res) != n {
		t.Errorf("got %d results, want %d", len(res), n)
	}
}

func TestConcurrentSubmitters(t *testing.T) {
	p := NewPool(8, 2, time.Millisecond)
	wait := collect(p)

	var submitted int64
	var wg sync.WaitGroup
	for g := 0; g < 10; g++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				if p.Submit(base*100+i, func() error { return nil }) {
					atomic.AddInt64(&submitted, 1)
				}
			}
		}(g)
	}
	wg.Wait()
	p.Shutdown()
	res := wait()

	if int64(len(res)) != submitted {
		t.Errorf("got %d results, want %d (exactly one per accepted submit)", len(res), submitted)
	}
}
