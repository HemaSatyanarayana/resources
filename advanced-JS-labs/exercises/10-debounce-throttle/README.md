# 10 — Debounce & Throttle

The two most-confused utilities in front-end interviews. Both limit how often `fn` runs, but with opposite intent: **debounce waits for silence; throttle enforces a steady rate.** You'll implement both with `.cancel()`, which forces you to manage the timer bookkeeping precisely.

## Concepts

- **Debounce** = "wait until the user stops." Every call clears the previous timer and starts a new one, so `fn` runs only once, `wait` ms after the *last* call. Great for search-as-you-type and autosave.
- **Throttle** = "at most once per interval." The leading call fires immediately; further calls in the window are coalesced into one trailing call with the latest args. Great for scroll/mousemove/resize.
- **Trailing vs leading edge.** Debounce here is trailing. Throttle here is leading + trailing — track the timestamp of the last *invocation* to decide whether a call fires now or gets queued.
- **`.cancel()`** clears the pending timer (and resets throttle's bookkeeping) so nothing fires later.

## Your task

Implement `debounce` and `throttle` in [`rateLimit.js`](rateLimit.js).

## Run

```bash
npx vitest run exercises/10-debounce-throttle
```

> Tests use fake timers, so leading-edge throttle relies on `Date.now()` (which fake timers advance). Track the last-invocation timestamp rather than assuming the clock starts anywhere specific.

## Interview follow-ups to expect

- *"Which for a search box? For infinite-scroll?"* — debounce the search (fire after typing stops); throttle the scroll (fire at a steady rate while scrolling).
- *"Add a `{ leading, trailing }` options object."* — this is what lodash does; the flags gate the immediate call and the tail call.
- *"Add `.flush()`."* — invoke any pending trailing call immediately.

## Hints

```js
export function debounce(fn, wait) {
  let timer
  function debounced(...args) {
    clearTimeout(timer)
    timer = setTimeout(() => fn.apply(this, args), wait)
  }
  debounced.cancel = () => clearTimeout(timer)
  return debounced
}
```

For throttle, compute `remaining = lastInvoke === null ? 0 : wait - (Date.now() - lastInvoke)`; invoke now if `remaining <= 0`, else schedule the trailing call once.

<details>
<summary>Reference solution (try first!)</summary>

```js
export function debounce(fn, wait) {
  let timer
  function debounced(...args) {
    clearTimeout(timer)
    timer = setTimeout(() => fn.apply(this, args), wait)
  }
  debounced.cancel = () => clearTimeout(timer)
  return debounced
}

export function throttle(fn, wait) {
  let lastInvoke = null
  let timer = null
  let savedArgs = null
  let savedThis = null

  const invoke = (time) => {
    lastInvoke = time
    fn.apply(savedThis, savedArgs)
  }

  function throttled(...args) {
    const now = Date.now()
    savedArgs = args
    savedThis = this
    const remaining = lastInvoke === null ? 0 : wait - (now - lastInvoke)

    if (remaining <= 0) {
      if (timer) {
        clearTimeout(timer)
        timer = null
      }
      invoke(now)
    } else if (!timer) {
      timer = setTimeout(() => {
        timer = null
        invoke(Date.now())
      }, remaining)
    }
  }

  throttled.cancel = () => {
    if (timer) clearTimeout(timer)
    timer = null
    lastInvoke = null
  }

  return throttled
}
```

</details>
