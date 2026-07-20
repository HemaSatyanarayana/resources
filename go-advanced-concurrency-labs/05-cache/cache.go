// Package cache — Lab 05: a singleflight, TTL-expiring, generic cache.
//
// Read README.md first; cache_test.go is the spec. Fill in every TODO.
// Run: go test -race -v ./05-cache/
package cache

import "time"

// entry holds one key's cached value. `ready` is closed once the load finishes,
// which is how concurrent callers wait for an in-flight load and how the janitor
// tells "still loading" from "done" without a second flag.
type entry[V any] struct {
	ready   chan struct{}
	val     V
	err     error
	expires time.Time
}

// Cache caches results of `load`, de-duplicating concurrent misses for the same
// key (only one load runs) and expiring entries after `ttl`. A background
// janitor evicts expired entries; Close stops it.
type Cache[V any] struct {
	// TODO: fields. You'll need:
	//   - ttl time.Duration and load func(string)(V, error),
	//   - mu sync.Mutex guarding items map[string]*entry[V],
	//   - stop chan struct{} + sync.Once to stop the janitor once.
}

// New builds a cache with the given ttl and loader and starts its janitor.
func New[V any](ttl time.Duration, load func(key string) (V, error)) *Cache[V] {
	// TODO:
	//   - Build the Cache with an initialised items map and stop channel.
	//   - Launch the janitor goroutine (sweep roughly every ttl; clamp the
	//     ticker interval to a positive minimum so time.NewTicker can't panic).
	panic("TODO: implement New")
}

// Get returns the cached value for key, loading it (exactly once, even under a
// concurrent stampede) on a miss. A fresh, successful entry is served from
// cache; an expired or previously-errored entry is reloaded. Load errors are
// returned but never cached.
func (c *Cache[V]) Get(key string) (V, error) {
	// TODO:
	//   Lock. If items[key] exists:
	//     - non-blocking check `select { case <-e.ready: ...; default: ... }`.
	//     - ready CLOSED  → the load finished. If e.err == nil && now < e.expires,
	//       Unlock and return e.val (cache hit). Otherwise delete(items,key) and
	//       fall through to reload.
	//     - ready OPEN (default) → a load is in flight: Unlock, `<-e.ready`, then
	//       return e.val, e.err (share the in-flight result — this is singleflight).
	//   Miss/stale: create e := &entry{ready: make(chan struct{})}, store it,
	//   Unlock, then load OUTSIDE the lock:
	//       e.val, e.err = c.load(key); e.expires = now+ttl; close(e.ready).
	//   If e.err != nil: Lock, and only if items[key] is still this same e,
	//   delete it (don't cache failures — but don't clobber a newer entry).
	//   Return e.val, e.err.
	panic("TODO: implement Get")
}

// Len reports how many entries the cache currently holds (including any still
// loading or not yet swept). Handy for observing the janitor.
func (c *Cache[V]) Len() int {
	// TODO: Lock, return len(items).
	panic("TODO: implement Len")
}

// Close stops the janitor goroutine. It is safe to call more than once.
func (c *Cache[V]) Close() {
	// TODO: close the stop channel exactly once (sync.Once).
	panic("TODO: implement Close")
}

// janitor periodically evicts expired entries until Close.
func (c *Cache[V]) janitor(every time.Duration) {
	// TODO: time.NewTicker(every); for { select { case <-tick: c.sweep();
	//       case <-c.stop: return } }.
	panic("TODO: implement janitor")
}

// sweep removes every finished entry that has errored or expired.
func (c *Cache[V]) sweep() {
	// TODO: Lock. For each entry, if its ready channel is closed
	//       (select/default) and it errored or now.After(expires), delete it.
	//       Skip entries still loading (the default branch).
	panic("TODO: implement sweep")
}
