// Package workerpool — Lab 05: bounded concurrency with a worker pool.
//
// Implement everything below FROM SCRATCH. Delete these comments as you go.
// Read README.md first; the test file is the spec.
//
// Build:
//
//	Map(inputs []int, workers int, f func(int) int) []int
//	    Process inputs with exactly `workers` goroutines; results in INPUT ORDER.
//	    Treat workers <= 0 as 1.
//
// Run: go test -race -v ./05-worker-pool/
package workerpool

import "sync"

type Job struct{ idx, val int }
type Result struct{ idx, val int }

func Map(inputs []int, workers int, f func(int) int) []int {
	if workers <= 0 {
		workers = 1
	}

	var jobs = make(chan Job, len(inputs))
	go func() {
		defer close(jobs)
		for i, val := range inputs {
			jobs <- Job{idx: i, val: val}
		}
	}()

	var results = make(chan Result, len(inputs))
	var wg sync.WaitGroup

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- Result{idx: job.idx, val: f(job.val)}
			}
		}()
	}

	go func() {
		wg.Wait()

		close(results)
	}()

	var ans = make([]int, len(inputs))
	for result := range results {
		ans[result.idx] = result.val
	}

	return ans
}
