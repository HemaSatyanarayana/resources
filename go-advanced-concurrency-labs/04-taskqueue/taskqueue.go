// Package taskqueue — Lab 04: a worker pool that retries failing jobs.
//
// Read README.md first; taskqueue_test.go is the spec. Fill in every TODO.
// Run: go test -race -v ./04-taskqueue/
package taskqueue

import "time"

// Job is a unit of work. It returns nil on success, or an error to trigger a
// retry (up to the pool's maxAttempts).
type Job func() error

// Result reports the final outcome of one submitted job.
type Result struct {
	ID       int   // the id you passed to Submit
	Attempts int   // how many times the job ran (1 = succeeded on the first try)
	Err      error // nil on success; the last error if every attempt failed
}

// Pool is a fixed-size worker pool that runs jobs with bounded retries and
// fans results back in on a single channel.
//
// Contract: someone must consume Results() for as long as the pool runs — the
// workers block when the results channel backs up. The usual shape is to drain
// Results() in one goroutine and call Shutdown from another (or after all
// submits), then range over Results() until it closes.
type Pool struct {
	// TODO: fields. You'll need:
	//   - jobs    chan item   (buffered; the work queue workers range over)
	//   - results chan Result (buffered; the fan-in channel Shutdown closes)
	//   - maxAttempts int and backoff time.Duration (retry policy)
	//   - a sync.RWMutex + `closed bool` for close-safe Submit/Shutdown
	//     (same trick as lab 03: Submit takes RLock, Shutdown takes Lock, so no
	//     send can be in flight when Shutdown closes the jobs channel),
	//   - a sync.WaitGroup to know when every worker has exited.
}

// item is the internal envelope carried on the jobs channel.
type item struct {
	id  int
	job Job
}

// NewPool starts a pool of `workers` goroutines. Each job is attempted up to
// `maxAttempts` times, sleeping `backoff` between attempts. Treat workers < 1
// and maxAttempts < 1 as 1.
func NewPool(workers, maxAttempts int, backoff time.Duration) *Pool {
	// TODO:
	//   1. Clamp workers and maxAttempts to a minimum of 1.
	//   2. Build the pool with buffered jobs and results channels.
	//   3. wg.Add(workers), then launch `workers` goroutines each running the
	//      worker loop below.
	panic("TODO: implement NewPool")
}

// worker pulls jobs until the jobs channel is closed and drained, runs each
// with retries, and sends exactly one Result per job.
func (p *Pool) worker() {
	// TODO:
	//   defer wg.Done().
	//   for it := range p.jobs {
	//       attempts := 0
	//       for attempts < maxAttempts:
	//           attempts++; err = it.job()
	//           if err == nil: break
	//           if attempts < maxAttempts: time.Sleep(backoff)
	//       send Result{it.id, attempts, err} on results
	//   }
	panic("TODO: implement worker")
}

// Submit enqueues a job under id. It returns false if the pool has been shut
// down, true once the job is queued.
func (p *Pool) Submit(id int, job Job) bool {
	// TODO:
	//   - RLock (defer RUnlock). If closed: return false.
	//   - Send item{id, job} on the jobs channel, return true.
	//     Holding RLock across the send is what makes it safe: Shutdown can't
	//     close the channel until every RLock is released.
	panic("TODO: implement Submit")
}

// Results returns the fan-in channel of job results. It is closed once Shutdown
// has drained every job, so you can range over it.
func (p *Pool) Results() <-chan Result {
	// TODO: return the results channel.
	panic("TODO: implement Results")
}

// Shutdown stops accepting new jobs, waits for all queued and in-flight jobs to
// finish, then closes the results channel. It is safe to call more than once.
func (p *Pool) Shutdown() {
	// TODO:
	//   - Lock. If already closed: Unlock and return. Set closed = true and
	//     close(jobs), then Unlock.
	//   - wg.Wait() so every worker drains and exits.
	//   - close(results).
	panic("TODO: implement Shutdown")
}
