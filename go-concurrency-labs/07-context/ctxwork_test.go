package ctxlab

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestSleepOrCancelCompletes(t *testing.T) {
	start := time.Now()
	err := SleepOrCancel(context.Background(), 30*time.Millisecond)
	if err != nil {
		t.Errorf("SleepOrCancel err = %v, want nil", err)
	}
	if elapsed := time.Since(start); elapsed < 25*time.Millisecond {
		t.Errorf("returned after %v, expected ~30ms", elapsed)
	}
}

func TestSleepOrCancelCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(10 * time.Millisecond); cancel() }()

	start := time.Now()
	err := SleepOrCancel(ctx, time.Hour)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("did not return promptly on cancel: %v", elapsed)
	}
}

func TestSleepOrCancelTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err := SleepOrCancel(ctx, time.Hour)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want context.DeadlineExceeded", err)
	}
}

func TestFetchAllSuccess(t *testing.T) {
	urls := []string{"a", "b", "c"}
	got, err := FetchAll(context.Background(), urls,
		func(ctx context.Context, url string) (string, error) {
			return strings.ToUpper(url), nil
		})
	if err != nil {
		t.Fatalf("FetchAll err = %v, want nil", err)
	}
	want := []string{"A", "B", "C"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FetchAll = %v, want %v (must be in URL order)", got, want)
	}
}

func TestFetchAllFirstError(t *testing.T) {
	boom := errors.New("boom")
	urls := []string{"ok1", "bad", "ok2"}

	got, err := FetchAll(context.Background(), urls,
		func(ctx context.Context, url string) (string, error) {
			if url == "bad" {
				return "", boom
			}
			return url, nil
		})
	if got != nil {
		t.Errorf("FetchAll results on error = %v, want nil", got)
	}
	if !errors.Is(err, boom) {
		t.Errorf("FetchAll err = %v, want boom", err)
	}
}

// TestFetchAllCancelsSiblings: the failing fetch should cause the still-running
// ones to be canceled via ctx.
func TestFetchAllCancelsSiblings(t *testing.T) {
	boom := errors.New("boom")
	var canceledCount int64

	urls := []string{"fail", "slow1", "slow2"}
	_, err := FetchAll(context.Background(), urls,
		func(ctx context.Context, url string) (string, error) {
			if url == "fail" {
				return "", boom
			}
			// A slow fetch that watches for cancellation.
			select {
			case <-time.After(2 * time.Second):
				return url, nil
			case <-ctx.Done():
				atomic.AddInt64(&canceledCount, 1)
				return "", ctx.Err()
			}
		})

	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want boom", err)
	}
	if atomic.LoadInt64(&canceledCount) == 0 {
		t.Error("sibling fetches were not canceled after the first error")
	}
}
