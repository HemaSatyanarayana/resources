# 11 — HTTP Handlers & Routing

Welcome to the **backend track**. Go's `net/http` is a production-grade HTTP stack in the standard library — no framework required. This exercise covers handlers, the modern `ServeMux` router, query params, path wildcards, and JSON responses.

## Concepts

- **A handler** is anything with `ServeHTTP(w http.ResponseWriter, r *http.Request)`. `http.HandlerFunc` adapts a plain function to that interface.
- **`http.ResponseWriter`** is how you reply: `w.Write([]byte(...))`, `w.WriteHeader(code)`, `w.Header().Set(k, v)`. **Set headers before writing the body**, and `WriteHeader` before `Write`.
- **`*http.Request`** carries the input: `r.URL.Query().Get("name")`, `r.Method`, `r.Body`, and (Go 1.22+) `r.PathValue("id")`.
- **`http.ServeMux`** (Go 1.22+) supports **method + wildcard patterns**: `mux.HandleFunc("GET /items/{id}", h)`. Unmatched routes get an automatic 404.
- **Test without a network** using `net/http/httptest`: `httptest.NewRecorder()` captures a response; `httptest.NewServer(h)` spins up a real local server.

## Your task

Implement everything in [`handlers.go`](handlers.go). You'll need `import "net/http"`, `"encoding/json"`, and `"fmt"`.

| Function | Skill |
|----------|-------|
| `HealthHandler` | Write a plain-text body |
| `GreetHandler` | Read a query param, default value |
| `ItemHandler` | JSON encoding + `Content-Type` header |
| `EchoIDHandler` | `r.PathValue("id")` |
| `NewRouter` | `ServeMux` method patterns |

## Run

```bash
go test -v ./exercises/11-http-handlers/
```

## Hints

- `fmt.Fprintf(w, "Hello, %s!", name)` writes directly to the response.
- Set the header *before* encoding: `w.Header().Set("Content-Type", "application/json")` then `json.NewEncoder(w).Encode(item)`.
- Register routes like `mux.HandleFunc("GET /items/{id}", EchoIDHandler)`.
- `httptest.NewRequest` + `httptest.NewRecorder` let you call a handler directly in a unit test.

<details>
<summary>Reference solution</summary>

```go
package httpsvc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func ItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Item{ID: 7, Name: "widget", Price: 500})
}

func EchoIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "item %s", r.PathValue("id"))
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", HealthHandler)
	mux.HandleFunc("GET /greet", GreetHandler)
	mux.HandleFunc("GET /items/{id}", EchoIDHandler)
	mux.HandleFunc("GET /item", ItemHandler)
	return mux
}
```

</details>
