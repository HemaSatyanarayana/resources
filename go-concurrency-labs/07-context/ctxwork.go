// Package ctxlab — Lab 07: context cancellation & deadlines.
//
// Implement everything below FROM SCRATCH. Delete these comments as you go.
// Read README.md first; the test file is the spec.
//
// Build:
//
//	SleepOrCancel(ctx context.Context, d time.Duration) error
//	    Sleep d, or return ctx.Err() early if ctx is canceled.
//
//	FetchAll(ctx, urls, fetch) ([]string, error)
//	    Run fetch on every url concurrently (results in URL order). On the first
//	    error, cancel the rest and return that error.
//
// Run: go test -race -v ./07-context/
package ctxlab

import (
	"context"
	"sync"
	"time"
)

func SleepOrCancel(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}

func FetchAll(ctx context.Context, urls []string, fetch func(context.Context, string) (string, error)) ([]string, error) {
	var out = make([]string, len(urls))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var firstErr error
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(len(urls))
	for i, url := range urls {
		go func() {
			defer wg.Done()
			resp, err := fetch(ctx, url)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
					cancel()
				}
				mu.Unlock()
			}
			out[i] = resp
		}()
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	return out, nil

}
