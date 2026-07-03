# 13 — Async Cache (Integration: Closures + Promises + Timers)

The capstone that **weaves all three pillars together** — the exact shape of a hard senior interview question: *"Wrap an async function so repeated calls are cached, concurrent identical calls share one request, and entries expire after a TTL."*

- **Closures** hold the private `cache` and `inFlight` maps — state that persists across calls but is invisible to the outside.
- **Promises** power the async work and the in-flight de-duplication: overlapping callers get *the same promise object*, so `fn` runs once.
- **Timers** evict each entry `ttl` ms after it resolves, via `setTimeout`.

This is `useFetch`-style deduping, an SWR/react-query cache, or a DataLoader — in ~15 lines.

## Concepts

- **Two maps, two purposes.** `cache` stores *resolved values*; `inFlight` stores *pending promises*. A key lives in `inFlight` only between the call and its settlement, then moves to `cache` (on success).
- **De-dup by returning the same promise.** When a second call arrives mid-flight, hand back the existing promise instead of starting a new request. `p1 === p2`.
- **Don't cache failures.** Only `cache.set` in the success branch; `.finally` clears `inFlight` for *both* outcomes so a rejected call can be retried.
- **TTL via `setTimeout`.** When a value lands, schedule `cache.delete(key)` after `ttl`. (Lazy expiry — comparing a stored timestamp on read — is the alternative; discuss the trade-off.)

## Your task

Implement `memoizeAsync` in [`asyncCache.js`](asyncCache.js).

## Run

```bash
npx vitest run exercises/13-async-cache
```

> The TTL tests use fake timers (`vi.advanceTimersByTimeAsync`); the de-dup tests use the fact that in-flight calls resolve on the microtask queue.

## Interview follow-ups to expect

- *"Two identical calls race — how many times does `fn` run?"* — once; the second call finds the key in `inFlight` and returns the same promise.
- *"Why not cache the rejected promise?"* — a transient failure would poison the cache forever; clearing `inFlight` in `finally` (and only caching on success) makes failures retryable.
- *"`setTimeout` eviction vs lazy timestamp expiry?"* — timers evict proactively but leave dangling handles (use `.unref()` in Node); lazy expiry has no timers but keeps stale memory until the next read.
- *"Add a max size (LRU) / a custom key resolver / `.clear()`."* — natural extensions once the core works.
- *"Does the pending `setTimeout` keep Node alive?"* — yes; `.unref()` it in real code.

## Hints

```js
if (cache.has(key)) return Promise.resolve(cache.get(key))
if (inFlight.has(key)) return inFlight.get(key)

const promise = Promise.resolve()
  .then(() => fn.apply(this, args))
  .then((value) => {
    cache.set(key, value)
    if (ttl !== Infinity) setTimeout(() => cache.delete(key), ttl)
    return value
  })
  .finally(() => inFlight.delete(key))

inFlight.set(key, promise)
return promise
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export function memoizeAsync(fn, { ttl = Infinity } = {}) {
  const cache = new Map()
  const inFlight = new Map()

  return function (...args) {
    const key = JSON.stringify(args)

    if (cache.has(key)) return Promise.resolve(cache.get(key))
    if (inFlight.has(key)) return inFlight.get(key)

    const promise = Promise.resolve()
      .then(() => fn.apply(this, args))
      .then((value) => {
        cache.set(key, value)
        if (ttl !== Infinity) {
          setTimeout(() => cache.delete(key), ttl)
        }
        return value
      })
      .finally(() => {
        inFlight.delete(key)
      })

    inFlight.set(key, promise)
    return promise
  }
}
```

</details>
