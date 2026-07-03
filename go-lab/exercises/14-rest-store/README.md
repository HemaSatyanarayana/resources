# 14 — Concurrency-Safe REST Store

The capstone of the backend track: a **repository** (in-memory data store) behind a **REST/JSON API**. This mirrors a real service — the `Store` is your persistence layer, `NewRouter` is your transport layer — just with a `map` standing in for a database. Because an HTTP server handles requests concurrently, the store **must** be safe for concurrent use.

## Concepts

- **Repository pattern:** encapsulate data access behind methods (`Create`, `Get`, `List`, `Delete`). Handlers never touch the map directly — they call the store. This keeps transport and persistence decoupled and testable.
- **`sync.RWMutex`:** many readers *or* one writer. Use `RLock`/`RUnlock` for reads (`Get`, `List`) and `Lock`/`Unlock` for writes (`Create`, `Delete`). Under contention this beats a plain `Mutex` for read-heavy workloads.
- **The zero-value-map trap:** writing to a `nil` map panics. Initialize it in `NewStore`.
- **REST conventions:** `POST` creates (201 + body), `GET` reads (200, or 404 if missing), `DELETE` removes (204 No Content). Bad input → 400.
- **Decode/encode JSON** straight from the request/response bodies with `json.NewDecoder(r.Body)` / `json.NewEncoder(w)`.

## Your task

Implement everything in [`store.go`](store.go). You'll need `import "encoding/json"`, `"sort"`, and `"strconv"` in addition to `net/http` and `sync`.

| Piece | Skill |
|-------|-------|
| `NewStore`, `Create`, `Get`, `List`, `Delete` | Repository + `RWMutex` |
| `NewRouter` | REST handlers, status codes, JSON, `PathValue` |

## Run

```bash
go test -race -v ./exercises/14-rest-store/
```

The concurrency test fires 500 simultaneous `Create`s and asserts every ID is unique — it fails loudly (or the race detector trips) if your locking is wrong.

## Hints

- `Create`: lock, `s.nextID++`, build the `Task{ID: s.nextID, ...}`, store it, return it.
- `List`: collect map values into a slice, then `sort.Slice(out, func(i,j int) bool { return out[i].ID < out[j].ID })`.
- Parse the path id with `strconv.Atoi(r.PathValue("id"))`; on error, `http.Error(w, "bad id", 400)`.
- For 201: `w.WriteHeader(http.StatusCreated)` **before** encoding the body. For 204: `w.WriteHeader(http.StatusNoContent)` and write nothing.
- Register routes by method: `mux.HandleFunc("POST /tasks", ...)`, `mux.HandleFunc("GET /tasks/{id}", ...)`, etc.

<details>
<summary>Reference solution</summary>

```go
package store

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

func NewStore() *Store {
	return &Store{tasks: make(map[int]Task)}
}

func (s *Store) Create(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	t := Task{ID: s.nextID, Title: title}
	s.tasks[t.ID] = t
	return t
}

func (s *Store) Get(id int) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *Store) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return false
	}
	delete(s.tasks, id)
	return true
}

func NewRouter(s *Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		t := s.Create(body.Title)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	})

	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.List())
	})

	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "bad id", http.StatusBadRequest)
			return
		}
		t, ok := s.Get(id)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
	})

	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "bad id", http.StatusBadRequest)
			return
		}
		if !s.Delete(id) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}
```

</details>
