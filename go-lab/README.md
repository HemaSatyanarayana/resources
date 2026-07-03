# Go Lab — Mastering Go Through Exercises

A hands-on, test-driven curriculum for mastering Go. Each exercise ships with:

- **A README** explaining the concept, the task, and what "good Go" looks like.
- **Starter code** with function signatures and `panic("TODO")` bodies for you to implement.
- **A test suite** — implement until `go test` is green.
- **A reference solution** hidden in a `<details>` block in each README (peek only after trying).

The whole repo is a single Go module (`go-lab`). Every exercise is its own package, so you can work on them independently.

## How to use this

1. Pick an exercise directory, e.g. `exercises/01-fundamentals`.
2. Read its `README.md`.
3. Open the starter file (e.g. `fizzbuzz.go`) and replace each `panic("TODO")`.
4. Run the tests for that exercise:

   ```bash
   go test ./exercises/01-fundamentals/
   ```

5. Green? Move on. Red? Read the failure, fix, repeat.

Run **everything** at once:

```bash
go test ./...
```

Run with verbose output and the race detector (great habit for the concurrency chapters):

```bash
go test -v -race ./...
```

## Curriculum

### Track A — Language mastery (01–10)

| #  | Topic | What you'll master |
|----|-------|--------------------|
| 01 | [Fundamentals](exercises/01-fundamentals) | Control flow, loops, multiple returns, `switch` |
| 02 | [Slices & Maps](exercises/02-slices-maps) | `append`, slicing tricks, map idioms, the comma-ok pattern |
| 03 | [Strings & Runes](exercises/03-strings-runes) | UTF-8, `rune` vs `byte`, `strings`/`unicode`, builders |
| 04 | [Structs & Methods](exercises/04-structs-methods) | Value vs pointer receivers, embedding, constructors |
| 05 | [Interfaces](exercises/05-interfaces) | Implicit satisfaction, `Stringer`, `sort.Interface`, type switches |
| 06 | [Errors](exercises/06-errors) | Sentinel errors, wrapping, `errors.Is`/`As`, custom error types |
| 07 | [Generics](exercises/07-generics) | Type parameters, constraints, generic containers |
| 08 | [Concurrency](exercises/08-concurrency) | Goroutines, `WaitGroup`, `Mutex`, worker pools |
| 09 | [Channels & Pipelines](exercises/09-channels-pipeline) | Directional channels, `select`, fan-in/fan-out, cancellation |
| 10 | [Stdlib & JSON](exercises/10-stdlib-json) | `encoding/json`, struct tags, `sort`, custom marshaling |

### Track B — Backend mastery (11–15)

Building real HTTP services with only the standard library. Everything is tested with `net/http/httptest` — no external database or framework required.

| #  | Topic | What you'll master |
|----|-------|--------------------|
| 11 | [HTTP Handlers & Routing](exercises/11-http-handlers) | `net/http`, `ServeMux` method patterns, query params, JSON responses |
| 12 | [Middleware](exercises/12-middleware) | `func(Handler) Handler`, chaining, auth, panic recovery |
| 13 | [Context](exercises/13-context) | Cancellation, deadlines, request-scoped values |
| 14 | [REST Store](exercises/14-rest-store) | Repository pattern, `RWMutex`, full CRUD JSON API |
| 15 | [Rate Limiting](exercises/15-rate-limiting) | Token bucket, injectable clocks, HTTP 429 middleware |

## Recommended order

Do them in numeric order — later exercises assume earlier concepts. Track A chapters 08–09 (concurrency) are the payoff of the language track and prerequisites for Track B, which leans on goroutine-safe state throughout. Run the backend and concurrency chapters with `-race`.

## Idioms this lab drills into you

- **Handle every error.** No `_ = err` unless you mean it.
- **Accept interfaces, return structs.**
- **The zero value should be useful.** Design types so `var x T` is ready to go.
- **Don't communicate by sharing memory; share memory by communicating** (channels over locks — when it fits).
- **`gofmt` is not optional.** Run `go fmt ./...` before you commit.

## Toolbelt

```bash
go fmt ./...      # format
go vet ./...      # catch suspicious constructs
go test ./...     # run all tests
go test -race ./... # detect data races
go doc strings.Builder  # read docs from the terminal
```

Happy hacking. 🐹
