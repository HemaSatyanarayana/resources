# 12 — Middleware

Middleware is how backends add cross-cutting concerns — auth, logging, metrics, panic recovery, CORS — without touching every handler. In Go the pattern is a plain function that **wraps one handler and returns another**.

## Concepts

- **The shape:** `func(http.Handler) http.Handler`. The returned handler does some work, then calls `next.ServeHTTP(w, r)` (or short-circuits with an error response).
- **Composition:** wrapping `A(B(C(handler)))` runs `A` first on the way in and last on the way out — like nested layers of an onion.
- **Short-circuiting:** auth middleware that writes a 401 and simply *doesn't* call `next` stops the chain.
- **Panic recovery:** a `defer`/`recover()` in middleware turns a handler panic into a clean 500 instead of killing the process.
- **Order matters:** put `Recover` outermost so it catches panics from everything inside it.

## Your task

Implement everything in [`middleware.go`](middleware.go):

| Function | Skill |
|----------|-------|
| `Chain` | Compose middlewares, first = outermost |
| `RequireAuth` | Inspect a header, short-circuit with 401 |
| `Recover` | `defer` + `recover()` → 500 |
| `Count` | A middleware *factory* (closure over state) |

## Run

```bash
go test -v ./exercises/12-middleware/
```

## Hints

- Wrap with `http.HandlerFunc(func(w, r){ ... })` inside each middleware.
- `Chain` applies middlewares in **reverse** so the first one ends up outermost:
  `for i := len(mws) - 1; i >= 0; i-- { h = mws[i](h) }`.
- `Recover`: `defer func(){ if rec := recover(); rec != nil { http.Error(w, "internal error", 500) } }()` **before** calling `next`.
- `http.Error(w, msg, code)` sets the status and writes the message in one call (it appends a newline — the tests trim it).
- `Count` returns a `Middleware`: `return func(next http.Handler) http.Handler { ... }`.

<details>
<summary>Reference solution</summary>

```go
package middleware

import "net/http"

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Count(hits *int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*hits++
			next.ServeHTTP(w, r)
		})
	}
}
```

</details>
