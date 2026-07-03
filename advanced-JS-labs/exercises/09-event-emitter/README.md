# 09 — Event Emitter

Pub/sub is the backbone of Node (`EventEmitter`), the DOM (`addEventListener`), and every state library. Building one tests your grip on **collections of callbacks**, **unsubscribe handles**, and the subtle **"modify while iterating"** hazard.

## Concepts

- **A map of event → set of listeners.** A `Set` gives O(1) add/remove and dedupes identical callbacks.
- **Return an unsubscribe function from `on`.** This is the modern ergonomic (React effects, RxJS): `const off = ee.on(...); off()`. It closes over the exact event+callback so the caller doesn't have to hold both.
- **`once` wraps the callback** in a self-removing function so it runs at most once. Return its unsubscribe too, so a pending `once` can be cancelled.
- **Iterate a copy in `emit`.** A listener may `off` itself (or another) mid-emit; iterating the live set while mutating it can skip or crash. Snapshot with `[...set]` first.

## Your task

Implement `on`, `off`, `once`, and `emit` in [`EventEmitter.js`](EventEmitter.js).

## Run

```bash
npx vitest run exercises/09-event-emitter
```

## Interview follow-ups to expect

- *"Why snapshot the listeners before calling them?"* — so unsubscribing during dispatch doesn't mutate the collection you're iterating.
- *"Memory leaks?"* — forgotten listeners keep their closures (and captured objects) alive; that's why unsubscribe handles and `once` matter. Node warns past 10 listeners for this reason.
- *"Make it async / ordered / support wildcards."* — `emit` could `await` async listeners in series, or support `'*'` catch-all events.

## Hints

```js
on(event, cb) {
  if (!this.listeners.has(event)) this.listeners.set(event, new Set())
  this.listeners.get(event).add(cb)
  return () => this.off(event, cb)
}

emit(event, ...args) {
  const set = this.listeners.get(event)
  if (!set) return
  for (const cb of [...set]) cb(...args)
}
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export class EventEmitter {
  constructor() {
    this.listeners = new Map()
  }

  on(event, cb) {
    if (!this.listeners.has(event)) this.listeners.set(event, new Set())
    this.listeners.get(event).add(cb)
    return () => this.off(event, cb)
  }

  off(event, cb) {
    this.listeners.get(event)?.delete(cb)
  }

  once(event, cb) {
    const wrapper = (...args) => {
      this.off(event, wrapper)
      cb(...args)
    }
    return this.on(event, wrapper)
  }

  emit(event, ...args) {
    const set = this.listeners.get(event)
    if (!set) return
    for (const cb of [...set]) cb(...args)
  }
}
```

</details>
