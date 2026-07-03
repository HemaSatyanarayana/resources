# 01 — Closures

A closure is **a function bundled with the variables it captured from the scope where it was defined**. That single idea powers private state, `once`, memoization, and the whole module pattern. Interviewers open with these to see if you *really* understand capture-by-reference.

## Concepts

- **Capture by reference, not value.** The inner function sees the live variable, so `createCounter`'s methods all mutate the same `count`.
- **Private state.** Nothing outside the factory can touch `count` — there's no property to reach. This is encapsulation without classes.
- **Independent instances.** Each call to `createCounter` makes a *new* scope, so two counters don't share state.
- **`once` / `memoize`** stash their bookkeeping (a `called` flag, a cache `Map`) in the closure so it survives between calls but stays hidden.

## Your task

Implement `createCounter`, `once`, and `memoize` in [`closures.js`](closures.js).

## Run

```bash
npx vitest run exercises/01-closures
```

## Interview follow-ups to expect

- *"Fix the `for (var i = 0; i < 3; i++) setTimeout(() => log(i))` bug."* — all three log `3` because they close over the *same* `i`. Fix with `let` (per-iteration binding) or an IIFE that captures the value.
- *"Do closures cause memory leaks?"* — they can: a closure keeps its captured scope alive. A forgotten listener holding a big object won't be collected.
- *"Make `memoize` cache by a custom key / bound size (LRU)."* — accept a `resolver` and evict oldest entries.

## Hints

- `once`: a boolean flag + a stored `result`; `if (!called) { called = true; result = fn.apply(this, args) }`.
- `memoize`: `const key = JSON.stringify(args)`; check a `Map` before calling.

<details>
<summary>Reference solution (try first!)</summary>

```js
export function createCounter(start = 0) {
  let count = start
  return {
    increment: () => ++count,
    decrement: () => --count,
    reset: () => {
      count = start
    },
    value: () => count,
  }
}

export function once(fn) {
  let called = false
  let result
  return function (...args) {
    if (!called) {
      called = true
      result = fn.apply(this, args)
    }
    return result
  }
}

export function memoize(fn) {
  const cache = new Map()
  return function (...args) {
    const key = JSON.stringify(args)
    if (cache.has(key)) return cache.get(key)
    const result = fn.apply(this, args)
    cache.set(key, result)
    return result
  }
}
```

</details>
