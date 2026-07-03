# 08 — Async Control Flow

Real async work needs orchestration: **retry** flaky calls, **time out** slow ones, and **cap concurrency** so you don't fire 10,000 requests at once. These three utilities are staples of production code and system-design-flavored coding rounds.

## Concepts

- **`retry`** is a loop around `try/await/catch`. Track the last error; only sleep *between* attempts, not after the final one. Real versions use exponential backoff (`delay * 2 ** attempt`) plus jitter.
- **`withTimeout`** is `Promise.race([work, timeoutThatRejects])`. Whichever settles first wins. Always `clearTimeout` so a resolved-fast case doesn't leave a dangling timer (and keep the process alive).
- **`mapLimit`** is a tiny scheduler: keep an index pointer and an in-flight counter; whenever a task finishes, start the next until the pointer runs out. Results go to `results[i]` to preserve order.
- **Racing doesn't cancel the loser.** `withTimeout` rejecting doesn't stop the underlying work — mention `AbortController` for true cancellation.

## Your task

Implement `retry`, `withTimeout`, and `mapLimit` in [`control.js`](control.js).

## Run

```bash
npx vitest run exercises/08-async-control
```

## Interview follow-ups to expect

- *"Add exponential backoff with jitter."* — `delay = base * 2 ** attempt + Math.random() * base`.
- *"Does `withTimeout` cancel the slow request?"* — no; `race` just ignores the loser. Wire an `AbortController` to actually abort it.
- *"Why cap concurrency?"* — avoid exhausting sockets/rate limits/memory; it's the client-side analog of a connection pool.
- *"`mapLimit` where one task rejects — what about the in-flight ones?"* — they keep running (not awaited); discuss cancellation if that matters.

## Hints

```js
export async function retry(fn, { retries = 3, delay = 0 } = {}) {
  let lastError
  for (let attempt = 0; attempt <= retries; attempt++) {
    try {
      return await fn(attempt)
    } catch (err) {
      lastError = err
      if (attempt < retries && delay) {
        await new Promise((r) => setTimeout(r, delay))
      }
    }
  }
  throw lastError
}
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export async function retry(fn, { retries = 3, delay = 0 } = {}) {
  let lastError
  for (let attempt = 0; attempt <= retries; attempt++) {
    try {
      return await fn(attempt)
    } catch (err) {
      lastError = err
      if (attempt < retries && delay) {
        await new Promise((r) => setTimeout(r, delay))
      }
    }
  }
  throw lastError
}

export function withTimeout(promise, ms, message = 'Timed out') {
  let id
  const timeout = new Promise((_, reject) => {
    id = setTimeout(() => reject(new Error(message)), ms)
  })
  return Promise.race([promise, timeout]).finally(() => clearTimeout(id))
}

export function mapLimit(items, limit, iteratee) {
  return new Promise((resolve, reject) => {
    const results = new Array(items.length)
    let index = 0
    let inFlight = 0
    let completed = 0
    let failed = false

    if (items.length === 0) return resolve([])

    const launchNext = () => {
      if (failed) return
      while (inFlight < limit && index < items.length) {
        const i = index++
        inFlight++
        Promise.resolve(iteratee(items[i], i)).then(
          (value) => {
            results[i] = value
            inFlight--
            completed++
            if (completed === items.length) resolve(results)
            else launchNext()
          },
          (err) => {
            failed = true
            reject(err)
          }
        )
      }
    }

    launchNext()
  })
}
```

</details>
