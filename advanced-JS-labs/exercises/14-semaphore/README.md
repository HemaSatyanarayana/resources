# 14 — Semaphore / Mutex (Promises + Closures)

One of the most-asked async-primitive questions: *"Implement a semaphore (or mutex) in JavaScript."* It's the clean pairing of **promises** (acquirers wait on a promise) and **closures** (the permit count and the waiter queue live privately in the closure). A **mutex is just `createSemaphore(1)`**.

## Concepts

- **A permit count + a waiter queue.** When permits are available, `acquire()` returns an already-resolved promise. When they're not, it returns a *pending* promise and stashes its `resolve` in a FIFO queue.
- **`release()` wakes the next waiter.** If someone is queued, hand the permit straight to them by calling their stored `resolve` — don't bump the count (the permit never became free, it was transferred). Only increment `available` when nobody is waiting.
- **`use(fn)` is the safe API.** Acquire, run the work, and release in a `finally` so a throwing task can't leak a permit and deadlock everyone behind it.
- **Semaphore vs the `mapLimit` from exercise 08.** `mapLimit` limits one specific batch; a semaphore is a *reusable primitive* you acquire/release from anywhere — the difference between a for-loop and a lock.

## Your task

Implement `acquire`, `release`, and `use` in [`semaphore.js`](semaphore.js).

## Run

```bash
npx vitest run exercises/14-semaphore
```

## Interview follow-ups to expect

- *"How is this a mutex?"* — size 1: exactly one holder at a time, so critical sections serialize.
- *"What if `use`'s callback throws?"* — the `finally` still releases; otherwise one exception permanently shrinks the pool.
- *"FIFO fairness?"* — a queue (shift the oldest) grants permits in request order; a stack would starve early waiters.
- *"Add a `tryAcquire()` / timeout / cancellation."* — non-blocking acquire, or reject a waiter after N ms.

## Hints

```js
function acquire() {
  if (available > 0) {
    available--
    return Promise.resolve()
  }
  return new Promise((resolve) => queue.push(resolve))
}

function release() {
  if (queue.length > 0) queue.shift()() // pass permit to next waiter
  else available = Math.min(max, available + 1)
}
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export function createSemaphore(max) {
  let available = max
  const queue = []

  function acquire() {
    if (available > 0) {
      available--
      return Promise.resolve()
    }
    return new Promise((resolve) => queue.push(resolve))
  }

  function release() {
    if (queue.length > 0) {
      const resolve = queue.shift()
      resolve()
    } else {
      available = Math.min(max, available + 1)
    }
  }

  async function use(fn) {
    await acquire()
    try {
      return await fn()
    } finally {
      release()
    }
  }

  return { acquire, release, use }
}
```

</details>
