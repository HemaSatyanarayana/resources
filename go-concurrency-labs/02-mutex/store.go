// Package mutexlab — Lab 02: sync.Mutex & sync.RWMutex.
//
// Implement everything below FROM SCRATCH. Delete these comments as you go.
// Read README.md first; the test file is the spec.
//
// Build:
//
//	type Counter struct{...}          // zero value ready
//	  (c *Counter) Inc()
//	  (c *Counter) Add(delta int)
//	  (c *Counter) Value() int
//
//	type Store struct{...}            // string->int map, RWMutex-guarded
//	  NewStore() *Store
//	  (s *Store) Set(key string, val int)
//	  (s *Store) Get(key string) (int, bool)
//	  (s *Store) Len() int
//	  (s *Store) Snapshot() map[string]int   // a copy
//
// Run: go test -race -v ./02-mutex/
package mutexlab

import "sync"

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++

}

func (c *Counter) Add(delta int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

type Store struct {
	mu    sync.RWMutex
	store map[string]int
}

func NewStore() *Store {
	return &Store{store: make(map[string]int)}
}

func (s *Store) Set(key string, val int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = val
}

func (s *Store) Get(key string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.store[key]
	return val, ok
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.store)
}

func (s *Store) Snapshot() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var copy = make(map[string]int)

	for key, val := range s.store {
		copy[key] = val
	}

	return copy
}
