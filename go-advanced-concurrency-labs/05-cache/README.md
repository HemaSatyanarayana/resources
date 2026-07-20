# 05 — Singleflight TTL Cache

Caches sit in front of expensive things — a database row, an HTTP call, a
rendered template. The naive concurrent cache has a nasty failure mode: a
**cache stampede**. A hot key expires, 500 requests miss simultaneously, and all
500 fire the expensive load at once — often knocking over the very thing the
cache was protecting. This lab builds a cache that collapses those 500 misses
into **one** load (that's *singleflight*), expires entries after a TTL, and runs
a background **janitor** to evict stale entries so the map doesn't grow forever.

## The system

```
 Get("k") ─┐
 Get("k") ─┼─▶ miss? ─▶ ONE load("k") ───▶ value stored with expires = now+ttl
 Get("k") ─┘            (the other callers                    │
                         wait on it, share result)            ▼
                                              janitor ticks ──▶ evict expired
```

Operations (generic over the value type `V`):

- **`New(ttl, load)`** builds the cache and starts the janitor. `load` is
  `func(key string) (V, error)`.
- **`Get(key)`** → `(V, error)`. Serves a fresh cached value, or loads on a
  miss. Concurrent misses for the same key trigger **exactly one** `load`; the
  rest wait and share its result.
- **`Len()`** → number of entries currently held (so tests — and you — can watch
  the janitor work).
- **`Close()`** stops the janitor. Idempotent.

### The key idea: one channel does two jobs

Each entry carries a `ready chan struct{}` that starts open and is **closed once
the load finishes**. That single channel gives you both things you need:

1. **Singleflight.** The first `Get` to miss inserts an entry (with an open
   `ready`) *under the lock*, then loads outside the lock. Any other `Get` that
   arrives finds that entry, sees `ready` still open, drops the lock, and blocks
   on `<-e.ready`. When the loader closes `ready`, every waiter wakes and reads
   the same `e.val, e.err`. N misses, one load.

2. **"Done vs. still loading."** The janitor (and `Get`) must not touch an
   entry's `val`/`expires` while the loader is still writing them. A
   *non-blocking* receive tells them apart safely:

   ```go
   select {
   case <-e.ready: // load finished — val/err/expires are set, safe to read
   default:        // still loading — leave it alone
   }
   ```

### Why the load happens outside the lock

If you held the mutex during `load`, a slow load of key `"a"` would block every
`Get` of every *other* key. So: insert the placeholder entry under the lock,
release the lock, then load. The lock protects the **map**, never the slow I/O.

### The two rules that keep it correct

- **Never cache a failure.** After a load errors, delete the entry — but only if
  it's still *your* entry (`items[key] == e`). A later `Get` may have already
  replaced it; don't clobber the newcomer. (This is the classic compare-then-act
  under the lock.)
- **Expiry is checked on read *and* swept in the background.** `Get` treats an
  expired entry as a miss and reloads (lazy). The janitor deletes expired
  entries on a timer (eager) so keys that are never requested again don't leak.

## Your task

Implement, in `cache.go` (package `cache`):

```go
type Cache[V any] struct { /* your fields */ }

func New[V any](ttl time.Duration, load func(key string) (V, error)) *Cache[V]
func (c *Cache[V]) Get(key string) (V, error)
func (c *Cache[V]) Len() int
func (c *Cache[V]) Close()
```

One `load` per key per TTL window. Concurrent misses de-duplicate. Errors are
returned but not cached. The janitor evicts expired entries and stops on `Close`.

## Run

```bash
go test -race -v ./05-cache/
```

## Hints

- Entry: `struct { ready chan struct{}; val V; err error; expires time.Time }`.
- `Get` cache-hit test: `e.err == nil && time.Now().Before(e.expires)`.
- In-flight wait: `c.mu.Unlock(); <-e.ready; return e.val, e.err`.
- Miss: create the entry, store it, unlock, load, set `expires`, `close(e.ready)`.
- Error cleanup: `if cur, ok := c.items[key]; ok && cur == e { delete(c.items, key) }`.
- Janitor: `time.NewTicker(ttl)` in a `for { select { <-tick / <-stop } }` loop;
  clamp the interval to a positive minimum so the ticker can't panic on `ttl<=0`.
- `sweep`: under the lock, for each *finished* entry (`ready` closed) that
  errored or is past `expires`, delete it.

<details>
<summary>Reference solution</summary>

```go
package cache

import (
	"sync"
	"time"
)

type entry[V any] struct {
	ready   chan struct{}
	val     V
	err     error
	expires time.Time
}

type Cache[V any] struct {
	ttl      time.Duration
	load     func(key string) (V, error)
	mu       sync.Mutex
	items    map[string]*entry[V]
	stop     chan struct{}
	stopOnce sync.Once
}

func New[V any](ttl time.Duration, load func(key string) (V, error)) *Cache[V] {
	c := &Cache[V]{
		ttl:   ttl,
		load:  load,
		items: make(map[string]*entry[V]),
		stop:  make(chan struct{}),
	}
	sweepEvery := ttl
	if sweepEvery <= 0 {
		sweepEvery = time.Millisecond
	}
	go c.janitor(sweepEvery)
	return c
}

func (c *Cache[V]) Get(key string) (V, error) {
	c.mu.Lock()
	if e, ok := c.items[key]; ok {
		select {
		case <-e.ready:
			if e.err == nil && time.Now().Before(e.expires) {
				c.mu.Unlock()
				return e.val, nil
			}
			delete(c.items, key) // stale or errored — drop and reload
		default:
			c.mu.Unlock()
			<-e.ready
			return e.val, e.err
		}
	}
	e := &entry[V]{ready: make(chan struct{})}
	c.items[key] = e
	c.mu.Unlock()

	e.val, e.err = c.load(key)
	e.expires = time.Now().Add(c.ttl)
	close(e.ready)

	if e.err != nil {
		c.mu.Lock()
		if cur, ok := c.items[key]; ok && cur == e {
			delete(c.items, key)
		}
		c.mu.Unlock()
	}
	return e.val, e.err
}

func (c *Cache[V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

func (c *Cache[V]) Close() {
	c.stopOnce.Do(func() { close(c.stop) })
}

func (c *Cache[V]) janitor(every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			c.sweep()
		case <-c.stop:
			return
		}
	}
}

func (c *Cache[V]) sweep() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, e := range c.items {
		select {
		case <-e.ready:
			if e.err != nil || now.After(e.expires) {
				delete(c.items, k)
			}
		default:
		}
	}
}
```

The elegance is that one `ready` channel is the whole synchronization story:
closing it publishes the load's result to every waiter (a happens-before edge,
so no data race on `val`/`err`/`expires`), *and* its open/closed state is how
`Get` and the janitor distinguish an in-flight load from a finished one.

</details>
