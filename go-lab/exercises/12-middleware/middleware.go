// Package middleware drills the classic net/http middleware pattern:
// functions that wrap an http.Handler to add cross-cutting behavior
// (auth, logging, panic recovery) and compose into a chain.
package middleware

import "net/http"

// Middleware wraps a handler, returning a new handler. This
// func(http.Handler) http.Handler shape is the Go convention.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares to h so that the FIRST middleware in the list is
// the OUTERMOST (runs first on the way in, last on the way out).
//
//	Chain(h, A, B) behaves like A(B(h)).
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	panic("TODO: implement Chain")
}

// RequireAuth returns a Middleware that rejects requests whose "Authorization"
// header is not exactly "Bearer secret" with status 401 and body "unauthorized".
// Authorized requests are passed through to next.
func RequireAuth(next http.Handler) http.Handler {
	panic("TODO: implement RequireAuth")
}

// Recover returns a Middleware that recovers from any panic in next, writing
// status 500 and body "internal error" instead of crashing the server.
func Recover(next http.Handler) http.Handler {
	panic("TODO: implement Recover")
}

// Count returns a Middleware that increments *hits once per request (before
// calling next). It is used to prove the chain actually runs the wrapped
// handler. (Tests call it single-threaded, so no locking is required.)
func Count(hits *int) Middleware {
	panic("TODO: implement Count")
}
