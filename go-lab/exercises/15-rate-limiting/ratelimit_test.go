package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeClock lets us control time deterministically.
func newFakeBucket(capacity int, ratePerSec float64) (*TokenBucket, *time.Time) {
	b := NewTokenBucket(capacity, ratePerSec)
	base := time.Unix(0, 0)
	clock := base
	b.now = func() time.Time { return clock }
	b.last = clock
	return b, &clock
}

func TestBucketStartsFull(t *testing.T) {
	b, _ := newFakeBucket(3, 1)
	for i := 0; i < 3; i++ {
		if !b.Allow() {
			t.Fatalf("request %d should be allowed (bucket starts full)", i+1)
		}
	}
	if b.Allow() {
		t.Error("4th request should be denied — bucket empty")
	}
}

func TestBucketRefills(t *testing.T) {
	b, clock := newFakeBucket(3, 1) // 1 token/sec
	for i := 0; i < 3; i++ {
		b.Allow() // drain
	}
	if b.Allow() {
		t.Fatal("bucket should be empty")
	}

	// Advance 2 seconds -> +2 tokens.
	*clock = clock.Add(2 * time.Second)
	if !b.Allow() {
		t.Error("first request after 2s refill should be allowed")
	}
	if !b.Allow() {
		t.Error("second request after 2s refill should be allowed")
	}
	if b.Allow() {
		t.Error("third request after 2s refill should be denied")
	}
}

func TestBucketCapsAtCapacity(t *testing.T) {
	b, clock := newFakeBucket(2, 1)
	*clock = clock.Add(100 * time.Second) // would add 100 tokens, but cap is 2
	if !b.Allow() {
		t.Error("first request should be allowed")
	}
	if !b.Allow() {
		t.Error("second request should be allowed")
	}
	if b.Allow() {
		t.Error("refill must cap at capacity (2), not accumulate")
	}
}

func TestBucketConcurrent(t *testing.T) {
	// Bucket of exactly 100 tokens, no refill during the test window.
	b, _ := newFakeBucket(100, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex
	allowed := 0
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if b.Allow() {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if allowed != 100 {
		t.Errorf("allowed = %d, want exactly 100 (no over-issuing under concurrency)", allowed)
	}
}

func TestMiddleware(t *testing.T) {
	b, _ := newFakeBucket(1, 0) // one token, never refills
	h := Middleware(b, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	rr1 := httptest.NewRecorder()
	h.ServeHTTP(rr1, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr1.Code != http.StatusOK {
		t.Errorf("first request = %d, want 200", rr1.Code)
	}

	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("second request = %d, want 429", rr2.Code)
	}
	if got := strings.TrimSpace(rr2.Body.String()); got != "rate limited" {
		t.Errorf("body = %q, want %q", got, "rate limited")
	}
}
