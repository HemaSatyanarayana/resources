// Package store drills the repository pattern: a concurrency-safe in-memory
// data store fronted by a REST/JSON HTTP API. This is the shape of a real
// service's persistence + transport layers, minus a database.
package store

import (
	"net/http"
	"sync"
)

// Task is a single to-do item.
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Store is an in-memory, concurrency-safe repository of Tasks. IDs start at 1
// and increase. Protect every access with the mutex.
type Store struct {
	mu     sync.RWMutex
	tasks  map[int]Task
	nextID int
}

// NewStore returns a ready-to-use Store. (Remember: a nil map panics on write,
// so initialize the map here.)
func NewStore() *Store {
	panic("TODO: implement NewStore")
}

// Create inserts a new task with the given title (Done=false), assigns it the
// next ID, and returns the stored task.
func (s *Store) Create(title string) Task {
	panic("TODO: implement Store.Create")
}

// Get returns the task with id, and whether it existed. Use an RLock (read).
func (s *Store) Get(id int) (Task, bool) {
	panic("TODO: implement Store.Get")
}

// List returns all tasks sorted by ascending ID.
func (s *Store) List() []Task {
	panic("TODO: implement Store.List")
}

// Delete removes the task with id and reports whether it existed.
func (s *Store) Delete(id int) bool {
	panic("TODO: implement Store.Delete")
}

// NewRouter builds the REST API over s and returns it as an http.Handler:
//
//	POST   /tasks        body {"title": "..."} -> 201 + created task JSON
//	GET    /tasks                              -> 200 + JSON array
//	GET    /tasks/{id}                         -> 200 + task JSON, or 404
//	DELETE /tasks/{id}                         -> 204, or 404
//
// Malformed IDs or bodies should return 400.
func NewRouter(s *Store) http.Handler {
	panic("TODO: implement NewRouter")
}
