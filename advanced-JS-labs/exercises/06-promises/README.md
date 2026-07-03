# 06 — Promises

Before the combinators, get fluent with the promise fundamentals: **wrapping callbacks** (`promisify`), **timing** (`delay`), and **sequencing** (`series`). The `series`-vs-`Promise.all` distinction — sequential vs concurrent — is a question you'll get verbatim.

## Concepts

- **The `new Promise(executor)` bridge.** To adapt any callback API, call the async thing inside the executor and resolve/reject from its callback. This is the *only* place you should hand-construct a Promise.
- **Node callback convention**: `(err, result)` — error first. `promisify` rejects on `err`, resolves on `result`.
- **Sequential chaining with `reduce`.** Start from `Promise.resolve([])` and `.then` each task onto the chain so the next starts only after the previous resolves.
- **Sequential ≠ concurrent.** `series` runs one at a time (total time = sum); `Promise.all` runs all at once (total time = max). Choose based on whether tasks are independent and whether you must limit load.

## Your task

Implement `promisify`, `delay`, and `series` in [`promises.js`](promises.js).

## Run

```bash
npx vitest run exercises/06-promises
```

## Interview follow-ups to expect

- *"Convert `series` to run concurrently."* — `Promise.all(tasks.map(t => t()))`. Discuss the trade-off.
- *"What's the difference between returning a value and returning a promise from `.then`?"* — a returned promise is *adopted*: the chain waits for it (flattening). A plain value passes straight through.
- *"`Promise.resolve(x)` where `x` is a thenable?"* — it adopts `x`'s eventual state rather than wrapping it.

## Hints

```js
export const promisify = (fn) => (...args) =>
  new Promise((resolve, reject) => {
    fn(...args, (err, result) => (err ? reject(err) : resolve(result)))
  })

export const series = (tasks) =>
  tasks.reduce(
    (chain, task) => chain.then((acc) => task().then((r) => [...acc, r])),
    Promise.resolve([])
  )
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export function promisify(fn) {
  return (...args) =>
    new Promise((resolve, reject) => {
      fn(...args, (err, result) => (err ? reject(err) : resolve(result)))
    })
}

export function delay(ms, value) {
  return new Promise((resolve) => setTimeout(() => resolve(value), ms))
}

export function series(tasks) {
  return tasks.reduce(
    (chain, task) => chain.then((acc) => task().then((r) => [...acc, r])),
    Promise.resolve([])
  )
}
```

</details>
