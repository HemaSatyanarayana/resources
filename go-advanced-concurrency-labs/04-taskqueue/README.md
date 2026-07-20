# 04 — Retrying Job Pool

A worker pool is the workhorse of concurrent Go: a fixed number of goroutines
chew through a queue of work so you get bounded parallelism instead of one
goroutine per task. Real jobs *fail* — a flaky HTTP call, a locked row — so this
pool adds **retries with backoff**, fans every outcome back in on one channel,
and shuts down **gracefully**: no job is dropped, and nothing leaks.

## The system

```
 Submit(id, job) ─┐                    workers (N goroutines)
 Submit(id, job) ─┼─▶ jobs chan ──▶ ┌─ w1: run→retry→backoff→retry ─┐
 Submit(id, job) ─┘   (buffered)    ├─ w2: ...                      ├─▶ results chan ──▶ you range over it
                                    └─ w3: ...                      ┘
```

Operations:

- **`Submit(id, job)`** → `bool`. Queues `job` (a `func() error`) under `id`.
  Returns `false` if the pool is already shut down, `true` once queued.
- **workers** run each job up to `maxAttempts` times, sleeping `backoff` between
  attempts, and emit exactly **one** `Result{ID, Attempts, Err}` per job.
- **`Results()`** → `<-chan Result`. The fan-in channel. Consume it while the
  pool runs; it's closed when the pool has fully drained.
- **`Shutdown()`** → stop accepting work, wait for every queued and in-flight
  job to finish, then close `Results()`. Idempotent.

### The concurrency design

Three separate problems, three tools:

1. **Bounded fan-out.** `NewPool` launches exactly `workers` goroutines, each
   ranging over a shared `jobs` channel. That's your parallelism cap — no matter
   how many jobs you submit, only `workers` run at once.

2. **Retries with backoff.** Purely local to a worker: a little loop that
   re-runs the job, `time.Sleep(backoff)`-ing between failures, stopping on the
   first success or after `maxAttempts` tries. The `Result` records how many
   attempts it took and the final error (if any).

3. **Graceful drain + close-safety.** This is the hard part, and it's the same
   rule as lab 03:

   > A channel must be closed by exactly one party, and never while a send might
   > be in flight.

   `Submit` sends on `jobs`; `Shutdown` closes `jobs`. If those race, you panic.
   Guard them with a `sync.RWMutex`: **`Submit` takes `RLock` for the whole
   send, `Shutdown` takes the exclusive `Lock` to flip `closed` and `close(jobs)`**.
   `RLock` and `Lock` can't overlap, so no send is ever in flight at close.

   Then draining is automatic: workers `range` over `jobs`, so closing it lets
   them finish the buffered backlog and exit. `Shutdown` does
   `wg.Wait()` to block until all workers are done, *then* `close(results)`.

### Why Shutdown must not close results early

Workers are still sending results as they drain the last jobs. Close `results`
before `wg.Wait()` returns and a worker sends on a closed channel → panic. Order
is non-negotiable: **close jobs → wait for workers → close results.**

## Your task

Implement, in `taskqueue.go` (package `taskqueue`):

```go
type Job func() error
type Result struct { ID, Attempts int; Err error }
type Pool struct { /* your fields */ }

func NewPool(workers, maxAttempts int, backoff time.Duration) *Pool
func (p *Pool) Submit(id int, job Job) bool
func (p *Pool) Results() <-chan Result
func (p *Pool) Shutdown()
```

`Submit` after `Shutdown` returns `false`. `Shutdown` is idempotent and drains
every queued job before closing `Results()`. Exactly one `Result` per accepted
`Submit`.

## Run

```bash
go test -race -v ./04-taskqueue/
```

## Hints

- Wrap the id + job in a small `item` struct to carry on the `jobs` channel.
- Attempt loop: `for attempts < maxAttempts { attempts++; err = job(); if err == nil { break }; if attempts < maxAttempts { time.Sleep(backoff) } }`.
- `Submit`: `RLock`; if `closed` return false; `p.jobs <- item{...}`; return true.
- `Shutdown`: `Lock`; if already `closed` return; set `closed = true`,
  `close(jobs)`, `Unlock`; then `wg.Wait()`; then `close(results)`.
- Buffer both channels so submitters and workers don't lock-step. Remember the
  contract: **something must be draining `Results()`**, or workers wedge and
  `Shutdown` hangs.

<details>
<summary>Reference solution</summary>

```go
package taskqueue

import (
	"sync"
	"time"
)

type Job func() error

type Result struct {
	ID       int
	Attempts int
	Err      error
}

type item struct {
	id  int
	job Job
}

type Pool struct {
	jobs        chan item
	results     chan Result
	maxAttempts int
	backoff     time.Duration

	mu     sync.RWMutex
	closed bool
	wg     sync.WaitGroup
}

func NewPool(workers, maxAttempts int, backoff time.Duration) *Pool {
	if workers < 1 {
		workers = 1
	}
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	p := &Pool{
		jobs:        make(chan item, 128),
		results:     make(chan Result, 128),
		maxAttempts: maxAttempts,
		backoff:     backoff,
	}
	p.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go p.worker()
	}
	return p
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for it := range p.jobs {
		attempts := 0
		var err error
		for attempts < p.maxAttempts {
			attempts++
			err = it.job()
			if err == nil {
				break
			}
			if attempts < p.maxAttempts {
				time.Sleep(p.backoff)
			}
		}
		p.results <- Result{ID: it.id, Attempts: attempts, Err: err}
	}
}

func (p *Pool) Submit(id int, job Job) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.closed {
		return false
	}
	p.jobs <- item{id: id, job: job}
	return true
}

func (p *Pool) Results() <-chan Result { return p.results }

func (p *Pool) Shutdown() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	close(p.jobs)
	p.mu.Unlock()

	p.wg.Wait()
	close(p.results)
}
```

The whole thing hangs on the close-safety pattern from lab 03 (`Submit`/`RLock`
vs. `Shutdown`/`Lock`) plus one ordering rule: **jobs closed → workers joined →
results closed.** Get that order wrong and a draining worker panics on a closed
channel.

</details>
