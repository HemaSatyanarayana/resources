# 05 — Timers

`setTimeout`/`setInterval` are where the **event loop** becomes tangible. This exercise turns them into the promise-based primitives you actually reach for — `sleep`, a safe `setInterval`, and `pollUntil` — and highlights why recursive `setTimeout` beats `setInterval`.

## Concepts

- **`setTimeout` schedules a macrotask.** It runs *after* the current synchronous code and all pending microtasks (promises). The delay is a *minimum*, not a guarantee.
- **Wrap a timer in a Promise** to get `sleep`: `new Promise(res => setTimeout(res, ms))`.
- **Recursive `setTimeout` > `setInterval`** when the work might take a while: `setInterval` fires on a fixed schedule and can queue callbacks back-to-back if one runs long; recursive `setTimeout` schedules the next tick only *after* the current one finishes.
- **Polling** needs both a repeating timer and a deadline — track `Date.now()` against a start time.

## Your task

Implement `sleep`, `repeat`, and `pollUntil` in [`timers.js`](timers.js).

## Run

```bash
npx vitest run exercises/05-timers
```

> The `sleep`/`repeat` tests use Vitest **fake timers** (`vi.advanceTimersByTimeAsync`) — the standard way to test time-based code without actually waiting.

## Interview follow-ups to expect

- *"`setTimeout(fn, 0)` vs `Promise.resolve().then(fn)` — which runs first?"* — the promise. Microtasks drain completely before the next macrotask (the timer).
- *"Why can `setInterval` drift or overlap?"* — it ignores how long the callback takes; recursive `setTimeout` self-corrects.
- *"Cancel an in-flight `sleep`."* — accept an `AbortSignal`; reject on `abort` and `clearTimeout`.

## Hints

```js
export const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export function repeat(fn, ms) {
  let id
  let stopped = false
  const tick = () => {
    if (stopped) return
    fn()
    id = setTimeout(tick, ms)
  }
  id = setTimeout(tick, ms)
  return () => {
    stopped = true
    clearTimeout(id)
  }
}
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export function repeat(fn, ms) {
  let id
  let stopped = false
  const tick = () => {
    if (stopped) return
    fn()
    id = setTimeout(tick, ms)
  }
  id = setTimeout(tick, ms)
  return () => {
    stopped = true
    clearTimeout(id)
  }
}

export function pollUntil(predicate, { interval = 50, timeout = 2000 } = {}) {
  return new Promise((resolve, reject) => {
    const start = Date.now()
    const check = async () => {
      let result
      try {
        result = await predicate()
      } catch (err) {
        return reject(err)
      }
      if (result) return resolve(result)
      if (Date.now() - start >= timeout) {
        return reject(new Error('pollUntil: timed out'))
      }
      setTimeout(check, interval)
    }
    check()
  })
}
```

</details>
