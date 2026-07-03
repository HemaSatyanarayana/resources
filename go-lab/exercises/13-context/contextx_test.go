package contextx

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSumWithContextCompletes(t *testing.T) {
	sum, err := SumWithContext(context.Background(), []int{1, 2, 3, 4}, time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum != 10 {
		t.Errorf("sum = %d, want 10", sum)
	}
}

func TestSumWithContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled
	sum, err := SumWithContext(ctx, []int{1, 2, 3}, time.Millisecond)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
	if sum != 0 {
		t.Errorf("sum = %d, want 0 (nothing processed)", sum)
	}
}

func TestSumWithContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	// 100 items * 5ms = 500ms of work, but the deadline is 25ms.
	nums := make([]int, 100)
	for i := range nums {
		nums[i] = 1
	}
	sum, err := SumWithContext(ctx, nums, 5*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want context.DeadlineExceeded", err)
	}
	if sum <= 0 || sum >= 100 {
		t.Errorf("partial sum = %d, expected some progress but not all", sum)
	}
}

func TestRequestID(t *testing.T) {
	ctx := WithRequestID(context.Background(), "abc-123")
	id, ok := RequestID(ctx)
	if !ok || id != "abc-123" {
		t.Errorf("RequestID = (%q, %v), want (abc-123, true)", id, ok)
	}
	if _, ok := RequestID(context.Background()); ok {
		t.Error("RequestID on bare context should be (\"\", false)")
	}
}

func TestRaceWorkWins(t *testing.T) {
	got, err := Race(context.Background(), func() int { return 99 })
	if err != nil || got != 99 {
		t.Errorf("Race = (%d, %v), want (99, nil)", got, err)
	}
}

func TestRaceContextWins(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := Race(ctx, func() int {
		time.Sleep(200 * time.Millisecond)
		return 1
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want context.DeadlineExceeded", err)
	}
}
