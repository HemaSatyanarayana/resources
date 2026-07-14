# 02 — Mutex & RWMutex

When goroutines must share mutable state, a **mutex** (mutual exclusion lock)
serializes access so only one touches it at a time.

## Concepts

### The problem: data races

```go
var n int
for i := 0; i < 1000; i++ {
    go func() { n++ }()   // RACE: read-modify-write from many goroutines
}
```

`n++` is really *load n, add 1, store n*. Two goroutines interleaving those
steps lose updates. Run with `-race` and Go will point right at it. The fix is a
lock.

### `sync.Mutex`

```go
var mu sync.Mutex
mu.Lock()
// ... critical section: exactly one goroutine here at a time ...
mu.Unlock()
```

Idioms & rules:
- Pair every `Lock` with an `Unlock`. Use `defer mu.Unlock()` right after
  `Lock()` so it releases even on early return or panic.
- **Reads need the lock too.** Reading an `int` while another goroutine writes
  it is still a race. `Value()` must lock.
- The zero value is an unlocked mutex — ready to use.
- Keep the critical section small. Don't do slow I/O while holding a lock.
- Embed the mutex next to the data it guards, and never copy the struct after
  first use (that copies the lock). Pass pointers.

### `sync.RWMutex` — many readers, one writer

When reads vastly outnumber writes, `RWMutex` lets readers run in parallel:

```go
var mu sync.RWMutex
mu.RLock(); v := data[k]; mu.RUnlock()   // many readers concurrently
mu.Lock();  data[k] = v; mu.Unlock()     // writers get exclusive access
```

- `RLock`/`RUnlock` for read-only critical sections.
- `Lock`/`Unlock` for anything that mutates.
- A write waits for in-flight reads to finish, and new reads wait for the write.

> **Maps are not concurrency-safe.** Concurrent writes to a built-in `map`
> crash the program with a "concurrent map writes" fatal error — not even
> `-race` needed. Always guard a shared map with a lock (or use `sync.Map`).

## Your task

Create `store.go` (package `mutexlab`) and implement, **from scratch**:

```go
// Counter is a concurrency-safe counter. Its zero value is ready to use.
type Counter struct { /* ... */ }
func (c *Counter) Inc()            // +1
func (c *Counter) Add(delta int)   // +delta
func (c *Counter) Value() int      // current value (also under the lock!)

// Store is a concurrency-safe string->int map, optimized for many readers.
type Store struct { /* ... */ }
func NewStore() *Store                       // returns a ready Store
func (s *Store) Set(key string, val int)     // insert/overwrite
func (s *Store) Get(key string) (int, bool)  // value + found flag
func (s *Store) Len() int                     // number of keys
func (s *Store) Snapshot() map[string]int     // a COPY, safe to keep/mutate
```

Match those names exactly — the test file depends on them.

## Run

```bash
go test -race -v ./02-mutex/
```

## Hints

- `Counter`: embed `mu sync.Mutex` and `n int`. Lock/`defer` Unlock in all three
  methods, including `Value`.
- `Store`: use `sync.RWMutex` + `map[string]int`. `NewStore` must initialize the
  map (a nil map panics on write). `Get`/`Len`/`Snapshot` take `RLock`;
  `Set` takes `Lock`.
- `Snapshot` must build and return a **new** map by copying under `RLock` — never
  hand back the internal map, or callers could read it while you write it.

<details>
<summary>Reference solution</summary>

```go
package mutexlab

import "sync"

type Counter struct {
	mu sync.Mutex
	n  int
}

func (c *Counter) Inc()          { c.Add(1) }
func (c *Counter) Add(delta int) { c.mu.Lock(); c.n += delta; c.mu.Unlock() }
func (c *Counter) Value() int    { c.mu.Lock(); defer c.mu.Unlock(); return c.n }

type Store struct {
	mu sync.RWMutex
	m  map[string]int
}

func NewStore() *Store { return &Store{m: make(map[string]int)} }

func (s *Store) Set(key string, val int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = val
}

func (s *Store) Get(key string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}

func (s *Store) Snapshot() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]int, len(s.m))
	for k, v := range s.m {
		out[k] = v
	}
	return out
}
```

</details>
