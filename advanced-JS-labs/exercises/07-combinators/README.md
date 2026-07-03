# 07 — Promise Combinators

Reimplementing `Promise.all` / `allSettled` / `race` / `any` is the single most common "from scratch" promise question. The insight is that each is a small state machine over a **countdown of pending promises**, wrapping every input in `Promise.resolve` so plain values just work.

## Concepts

- **Order vs completion.** `all` resolves in *input* order even though promises settle in *timing* order — write results into `results[i]`, not `results.push`. Track a `remaining` counter; resolve when it hits 0.
- **Fail-fast vs collect-all.** `all` rejects on the first rejection; `allSettled` never rejects and records each outcome; `any` rejects only when *all* reject.
- **`race` is the simplest**: attach the outer `resolve`/`reject` to every input; first one wins.
- **`allSettled` in terms of `all`**: map each input to a promise that *always fulfills* with a `{status, ...}` object, then `all` them.
- **Empty-array edge cases**: `all([])`→`[]`, `any([])`→rejects with an empty `AggregateError`.

## Your task

Implement `all`, `allSettled`, `race`, and `any` in [`combinators.js`](combinators.js).

## Run

```bash
npx vitest run exercises/07-combinators
```

## Interview follow-ups to expect

- *"Why write to an index instead of pushing?"* — pushes land in completion order and scramble the result array; the index preserves input order.
- *"`all` vs `allSettled` — when each?"* — `all` when any failure should abort the whole batch; `allSettled` when you want every result regardless.
- *"`race` vs `any`?"* — `race` settles on the first to *settle* (even a rejection); `any` waits for the first to *fulfill*.

## Hints

```js
export function all(items) {
  return new Promise((resolve, reject) => {
    const results = []
    let remaining = items.length
    if (remaining === 0) return resolve([])
    items.forEach((item, i) => {
      Promise.resolve(item).then((value) => {
        results[i] = value
        if (--remaining === 0) resolve(results)
      }, reject)
    })
  })
}
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export function all(items) {
  return new Promise((resolve, reject) => {
    const results = []
    let remaining = items.length
    if (remaining === 0) return resolve([])
    items.forEach((item, i) => {
      Promise.resolve(item).then((value) => {
        results[i] = value
        if (--remaining === 0) resolve(results)
      }, reject)
    })
  })
}

export function allSettled(items) {
  return all(
    items.map((item) =>
      Promise.resolve(item).then(
        (value) => ({ status: 'fulfilled', value }),
        (reason) => ({ status: 'rejected', reason })
      )
    )
  )
}

export function race(items) {
  return new Promise((resolve, reject) => {
    items.forEach((item) => Promise.resolve(item).then(resolve, reject))
  })
}

export function any(items) {
  return new Promise((resolve, reject) => {
    const errors = []
    let remaining = items.length
    if (remaining === 0) {
      return reject(new AggregateError([], 'All promises were rejected'))
    }
    items.forEach((item, i) => {
      Promise.resolve(item).then(resolve, (err) => {
        errors[i] = err
        if (--remaining === 0) {
          reject(new AggregateError(errors, 'All promises were rejected'))
        }
      })
    })
  })
}
```

</details>
