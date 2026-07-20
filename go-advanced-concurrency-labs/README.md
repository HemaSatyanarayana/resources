# Go Advanced Concurrency Labs 🚀🐹

The sequel to [`go-concurrency-labs`](../go-concurrency-labs). That course drilled
the **primitives** one at a time (goroutines, mutexes, channels, select, pools,
pipelines, context, atomics). This one makes you **compose** them into small,
realistic, end-to-end projects — the way you'd actually use concurrency in
production.

Each lab is a self-contained mini-project: a concurrent web crawler, a rate
limiter, a pub/sub broker, a retrying job pool, a singleflight cache, and a
MapReduce word counter. You build the whole thing from a guided scaffold.

## How this works

Same loop as the first course — **the test file is the spec**:

```
NN-project/
  README.md         ← read first: the design, the exact API to build, hints
  <name>.go         ← a guided scaffold: signatures + TODOs. You fill the bodies.
  <name>_test.go    ← the spec. Don't edit it. Make it pass — with -race.
```

The scaffold compiles but every function `panic("TODO")`s, so tests fail loudly
until you implement them. The test file tells you the exact names/signatures the
graders expect; the README explains the design and *why* it's shaped that way.

### The loop

1. `cd` into a lab, read `README.md` — understand the **whole system** before coding.
2. Open the test file — treat it as the specification.
3. Fill in the scaffold's TODOs.
4. Run with the race detector (non-negotiable — these are real concurrent systems):

   ```bash
   go test -race -v ./01-crawler/
   ```

5. Green *and* race-clean? Move on. Otherwise read the output and fix.

Run one lab, or everything:

```bash
go test -race ./03-pubsub/     # one lab
go test -race ./...            # all labs (unfinished ones will fail/panic — expected)
```

> Each README hides a **reference solution** in a `<details>` block. Peek only
> after a genuine attempt — the point is to fight the deadlock yourself.

## Curriculum

| #  | Project | Concurrency skills you'll master |
|----|---------|----------------------------------|
| 01 | [Concurrent Web Crawler](01-crawler) | Bounded fan-out, shared visited-set dedup, **termination detection** (`WaitGroup` over a growing frontier), context cancellation |
| 02 | [Token-Bucket Rate Limiter](02-ratelimiter) | Channel-as-bucket, a background refill ticker, `Allow`/blocking `Wait`, `sync.Once` shutdown |
| 03 | [In-Memory Pub/Sub Broker](03-pubsub) | Generics, `RWMutex` fan-out, non-blocking (drop-on-full) delivery, safe channel close under concurrency |
| 04 | [Retrying Job Pool](04-taskqueue) | Worker pool + **retries with backoff**, graceful drain on `Shutdown`, results fan-in, close-safety |
| 05 | [Singleflight TTL Cache](05-cache) | Request **de-duplication** (one loader call for N concurrent misses), TTL expiry, a janitor goroutine |
| 06 | [MapReduce Word Count](06-mapreduce) | Three-stage fan-out/fan-in pipeline, deterministic results across worker counts, cancellation |

Do them in order — each assumes you're fluent in everything the first course
covered.

## Ground rules

- **Always** run with `-race`. A green run without `-race` proves nothing here.
- No goroutine leaks: every goroutine you start must have a guaranteed exit.
- No `time.Sleep` to "fix" a race — if a test only passes with a sleep, the
  synchronization is wrong.
