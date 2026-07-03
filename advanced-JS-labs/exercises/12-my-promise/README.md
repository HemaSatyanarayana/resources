# 12 — MyPromise

The capstone. Implement a Promise from scratch and you understand promises for good: the **state machine**, **microtask scheduling**, **chaining via a new promise per `.then`**, and **thenable adoption** (why returning a promise from `.then` flattens instead of nesting). This is a genuine senior/staff interview question.

## Concepts

- **State machine.** `pending → fulfilled | rejected`, and it settles exactly **once** — every later `resolve`/`reject` is a no-op.
- **Queue callbacks while pending.** If `.then` is called before the promise settles, stash the handler; run the queue when it settles.
- **Handlers are async.** Even for an already-settled promise, `.then` callbacks must run on the **microtask** queue (`queueMicrotask`), never synchronously — that's what keeps ordering predictable.
- **`.then` returns a NEW promise.** The new one resolves with the handler's return value — or, if that value is a **thenable**, it *adopts* it (this is the flattening that makes `.then(() => fetch(...))` chains work).
- **Pass-through.** A `.then` with no matching handler forwards the value/reason down the chain — this is how a `.catch` at the end catches an earlier rejection.
- **Executor errors reject.** A `throw` inside the executor becomes a rejection.

## Your task

Implement the constructor, `_settle`, and `then` in [`MyPromise.js`](MyPromise.js). Because your `.then` is spec-shaped, `await` will work on your promise — the tests use it.

## Run

```bash
npx vitest run exercises/12-my-promise
```

## Interview follow-ups to expect

- *"Why must `.then` be async even when already resolved?"* — to guarantee consistent ordering (handlers always run after the current synchronous stack) and avoid Zalgo — a function that's sometimes sync, sometimes async.
- *"Microtask vs macrotask?"* — promise callbacks are microtasks; they drain fully before the next `setTimeout` (macrotask). That's why `Promise.resolve().then(f)` beats `setTimeout(f, 0)`.
- *"What does 'adopting a thenable' mean?"* — when you resolve with something that has a `.then`, you wait for *it* and take on its state, rather than fulfilling with the thenable itself.
- *"Implement `finally` / `Promise.all` on top."* — good extensions once the core works.

## Hints

- `_settle` guards `if (this.state !== 'pending') return`, handles the thenable case, then stores state/value and flushes `callbacks` via `queueMicrotask`.
- In `then`, build the `handle` closure and either `queueMicrotask(handle)` (already settled) or `this.callbacks.push(handle)` (pending).

<details>
<summary>Reference solution (try first!)</summary>

```js
export class MyPromise {
  constructor(executor) {
    this.state = 'pending'
    this.value = undefined
    this.callbacks = []

    const resolve = (value) => this._settle('fulfilled', value)
    const reject = (reason) => this._settle('rejected', reason)

    try {
      executor(resolve, reject)
    } catch (err) {
      reject(err)
    }
  }

  _settle(state, value) {
    if (this.state !== 'pending') return

    if (
      state === 'fulfilled' &&
      value != null &&
      (typeof value === 'object' || typeof value === 'function')
    ) {
      let then
      try {
        then = value.then
      } catch (err) {
        return this._settle('rejected', err)
      }
      if (typeof then === 'function') {
        let called = false
        try {
          then.call(
            value,
            (v) => {
              if (!called) {
                called = true
                this._settle('fulfilled', v)
              }
            },
            (r) => {
              if (!called) {
                called = true
                this._settle('rejected', r)
              }
            }
          )
        } catch (err) {
          if (!called) {
            called = true
            this._settle('rejected', err)
          }
        }
        return
      }
    }

    this.state = state
    this.value = value
    this.callbacks.forEach((cb) => queueMicrotask(cb))
    this.callbacks = []
  }

  then(onFulfilled, onRejected) {
    return new MyPromise((resolve, reject) => {
      const handle = () => {
        const cb = this.state === 'fulfilled' ? onFulfilled : onRejected
        if (typeof cb !== 'function') {
          this.state === 'fulfilled' ? resolve(this.value) : reject(this.value)
          return
        }
        try {
          resolve(cb(this.value))
        } catch (err) {
          reject(err)
        }
      }

      if (this.state === 'pending') this.callbacks.push(handle)
      else queueMicrotask(handle)
    })
  }

  catch(onRejected) {
    return this.then(null, onRejected)
  }

  static resolve(value) {
    return new MyPromise((resolve) => resolve(value))
  }

  static reject(reason) {
    return new MyPromise((_, reject) => reject(reason))
  }
}
```

</details>
