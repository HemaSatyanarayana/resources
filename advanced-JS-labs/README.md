# Advanced JS Labs — Master Closures, Timers & Promises

A hands-on, **test-driven** curriculum for the hard parts of JavaScript — the concepts that show up in senior interviews and "implement this from scratch" rounds: **closures, `this`/binding, currying, timers (`setTimeout`/`setInterval`), promises and combinators, debounce/throttle, and a hand-rolled Promises/A+ implementation.**

Each exercise ships with:

- **A README** explaining the concept, the task, the interview follow-ups, and what "good JS" looks like.
- **Starter code** with function/class signatures and `TODO` bodies for you to implement.
- **A test suite** — implement until `npm test` is green.
- **A reference solution** hidden in a `<details>` block in each README (peek only after trying).

No framework, no DOM — just the language. Every exercise is its own module.

## Setup

```bash
cd advanced-JS-labs
npm install
```

## How to use this

1. Pick an exercise directory, e.g. `exercises/01-closures`.
2. Read its `README.md`.
3. Open the starter file (e.g. `closures.js`) and replace each `TODO`.
4. Run the tests for that exercise:

   ```bash
   npx vitest run exercises/01-closures
   ```

5. Green? Move on. Red? Read the failure, fix, repeat.

Run **everything**:

```bash
npm test
```

Watch mode (the way to actually work through these):

```bash
npm run test:watch
```

## Curriculum

### Track A — Functions, scope & closures (01–04)

The mental model everything else is built on: what a closure captures, how `this` is decided, and how functions compose.

| #  | Exercise | What you'll master |
|----|----------|--------------------|
| 01 | [Closures](exercises/01-closures) | Private state, `once`, `memoize` — what a closure captures and when |
| 02 | [this & binding](exercises/02-this-binding) | Implement `call` / `apply` / `bind` from scratch |
| 03 | [Currying](exercises/03-currying) | `curry`, `partial`, infinite currying via `valueOf` |
| 04 | [Composition](exercises/04-composition) | `compose`, `pipe`, point-free style with `reduce` |

### Track B — Timers & async foundations (05–08)

The event loop made concrete: timers, promisification, the combinators, and real async control flow.

| #  | Exercise | What you'll master |
|----|----------|--------------------|
| 05 | [Timers](exercises/05-timers) | `sleep`, `pollUntil`, `setInterval` built from `setTimeout` |
| 06 | [Promises](exercises/06-promises) | `promisify`, `delay`, sequential chains with `reduce` |
| 07 | [Promise Combinators](exercises/07-combinators) | Polyfill `all` / `allSettled` / `race` / `any` |
| 08 | [Async Control Flow](exercises/08-async-control) | `retry` w/ backoff, `withTimeout`, concurrency-limited `mapLimit` |

### Track C — Patterns & internals (09–12)

The classic "build this utility" questions, ending with a from-scratch Promise.

| #  | Exercise | What you'll master |
|----|----------|--------------------|
| 09 | [Event Emitter](exercises/09-event-emitter) | Pub/sub, `on`/`off`/`once`/`emit`, unsubscribe handles |
| 10 | [Debounce & Throttle](exercises/10-debounce-throttle) | Trailing/leading edges, `cancel`, timer bookkeeping |
| 11 | [Deep Clone & Equal](exercises/11-deep) | Recursion, circular refs with `WeakMap`, `Object.is` |
| 12 | [MyPromise](exercises/12-my-promise) | A Promises/A+ subset: states, microtask scheduling, thenable flattening |

### Track D — Combination primitives (13–15)

Where these concepts show up **together** — the hardest interview flavor. Famous async primitives, each labeled by the pillars it fuses.

| #  | Exercise | Pillars | What you'll master |
|----|----------|---------|--------------------|
| 13 | [Async Cache](exercises/13-async-cache) | timers + promises + closures | `memoizeAsync`: closure caches + in-flight promise de-dup + `setTimeout` TTL eviction |
| 14 | [Semaphore / Mutex](exercises/14-semaphore) | promises + closures | permit count + FIFO waiter queue handing out promises; mutex = size 1 |
| 15 | [Rate Limiter](exercises/15-rate-limiter) | timers + closures | token bucket: closure token count refilled by a `setInterval`, capped at capacity |

## Recommended order

Numeric order. Track A builds the closure/`this` intuition that the timer and promise code leans on constantly. Track B turns "I've used `async/await`" into "I understand what it does." Track C's `MyPromise` (12) is a capstone for promise internals. **Track D (13) is the true finale**: it combines closures + promises + timers in one utility — do it only after Track B, since it reuses every idea from `memoize` (01), `pollUntil` (05), and the combinators (07).

## Where the pillars combine — the drill set

Want to specifically practice combining these concepts? Here's the map, by pairing:

**All three (timers + promises + closures)**
- **[13 Async Cache](exercises/13-async-cache)** — all three, by design.
- **[05 Timers](exercises/05-timers)** — `pollUntil` closes over a deadline, loops with `setTimeout`, returns a promise.
- **[08 Async Control](exercises/08-async-control)** — `retry`/`withTimeout`/`mapLimit` thread closured state through promises and timers.

**Promises + closures**
- **[14 Semaphore / Mutex](exercises/14-semaphore)** — a waiter queue in a closure handing out promises.

**Timers + closures**
- **[15 Rate Limiter](exercises/15-rate-limiter)** — a closure token count refilled by a timer.
- **[10 Debounce & Throttle](exercises/10-debounce-throttle)** — timer bookkeeping in a closure.

**Timers + promises**
- **[05 Timers](exercises/05-timers)** (`sleep`) and **[08 Async Control](exercises/08-async-control)** (`withTimeout`) — promisified timers.

## The mental models this lab drills into you

- **A closure is a function plus the variables it captured** — captured by *reference*, not value. The `for (var i …)` bug and the counter factory are two sides of this coin.
- **`this` is set at call time**, by *how* a function is called (`obj.fn()`, `fn.call(x)`, `new`), not where it's defined. Arrow functions are the exception: they capture `this` lexically.
- **The event loop runs the stack to empty, then drains all microtasks (promises), then takes one macrotask (`setTimeout`).** This is why a resolved promise's `.then` runs before a `setTimeout(…, 0)`.
- **A promise is a state machine**: `pending → fulfilled | rejected`, settled exactly once, with callbacks scheduled as microtasks.
- **Debounce waits for silence; throttle rate-limits.** Know which one a problem wants.

## Toolbelt

```bash
npm test                                  # run everything once
npm run test:watch                        # re-run on save
npx vitest run exercises/07-combinators   # one exercise
npx vitest run -t "rejects"               # tests matching a name
```

Happy hacking. 🟨
