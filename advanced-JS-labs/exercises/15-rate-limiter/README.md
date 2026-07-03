# 15 — Rate Limiter / Token Bucket (Timers + Closures)

*"Implement a rate limiter"* is a staple that sits between coding and system design. The **token-bucket** algorithm is the canonical answer, and it's a clean pairing of **timers** (a background refill) and **closures** (the private token count). It's distinct from throttle: a bucket **permits bursts** up to `capacity`, then enforces the steady refill rate.

## Concepts

- **Tokens = permission to act.** Each request spends one. Empty bucket → request denied (`false`). This models "N requests per interval" precisely.
- **Refill on a timer.** `setInterval` adds one token every `refillIntervalMs`, **clamped to `capacity`** so an idle bucket doesn't accumulate infinite tokens (that would defeat the limit — the classic bug).
- **Burst vs steady rate.** Capacity sets how big a burst you tolerate; the refill interval sets the sustained rate. Throttle can't express "allow 10 quickly, then 1/sec."
- **`stop()` clears the timer.** A dangling `setInterval` leaks and (in Node) keeps the process alive — hence a teardown method (real code would also `.unref()` the timer).

## Your task

Implement the refill interval, `tryAcquire`, and `stop` in [`rateLimiter.js`](rateLimiter.js).

## Run

```bash
npx vitest run exercises/15-rate-limiter
```

> Tests use fake timers, so `vi.advanceTimersByTime(n)` drives the refills deterministically.

## Interview follow-ups to expect

- *"Token bucket vs leaky bucket vs fixed/sliding window?"* — token bucket allows bursts; leaky bucket smooths output to a constant rate; windows count events per time slice (and have boundary spikes).
- *"Why clamp at capacity?"* — without it, an idle limiter banks unlimited tokens and the next burst is unbounded.
- *"Make it async — `await limiter.acquire()` that waits for a token."* — combine this with the semaphore idea: queue waiters and wake them on refill.
- *"Lazy refill instead of a timer?"* — compute `tokens += elapsed / interval` on each call using `Date.now()`; no `setInterval`, no dangling handle.
- *"Distributed rate limiting?"* — move the bucket to Redis (`INCR` + `EXPIRE`) so it's shared across servers.

## Hints

```js
const id = setInterval(() => {
  tokens = Math.min(capacity, tokens + 1)
}, refillIntervalMs)

function tryAcquire() {
  if (tokens > 0) {
    tokens--
    return true
  }
  return false
}
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export function createRateLimiter(capacity, refillIntervalMs) {
  let tokens = capacity

  const id = setInterval(() => {
    tokens = Math.min(capacity, tokens + 1)
  }, refillIntervalMs)

  function tryAcquire() {
    if (tokens > 0) {
      tokens--
      return true
    }
    return false
  }

  function availableTokens() {
    return tokens
  }

  function stop() {
    clearInterval(id)
  }

  return { tryAcquire, availableTokens, stop }
}
```

</details>
