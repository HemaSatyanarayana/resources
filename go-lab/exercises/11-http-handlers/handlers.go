// Package httpsvc drills net/http: handlers, ServeMux routing (Go 1.22+
// method+wildcard patterns), query params, and JSON responses.
package httpsvc

import "net/http"

// HealthHandler writes status 200 and the plain-text body "ok".
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement HealthHandler")
}

// GreetHandler reads the "name" query parameter and writes "Hello, <name>!".
// If "name" is absent or empty, use "World". Always status 200, plain text.
func GreetHandler(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement GreetHandler")
}

// Item is the JSON payload returned by ItemHandler.
type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price_cents"`
}

// ItemHandler responds with an Item encoded as JSON. It must:
//   - set the "Content-Type" header to "application/json"
//   - encode this exact item: {ID: 7, Name: "widget", Price: 500}
//   - default status is 200 (no need to call WriteHeader)
func ItemHandler(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement ItemHandler")
}

// NewRouter wires the handlers onto an *http.ServeMux using Go 1.22+ method
// patterns and returns it as an http.Handler:
//
//	GET /health        -> HealthHandler
//	GET /greet         -> GreetHandler
//	GET /items/{id}    -> EchoIDHandler (see below)
//	GET /item          -> ItemHandler
func NewRouter() http.Handler {
	panic("TODO: implement NewRouter")
}

// EchoIDHandler reads the {id} path wildcard (r.PathValue("id")) and writes
// "item <id>" as plain text, e.g. "item 42".
func EchoIDHandler(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement EchoIDHandler")
}
