# 07 — useDebounce

Debouncing is the single most common async sub-task in machine coding — search boxes, resize handlers, autosave. The trap interviewers watch for is the **missing cleanup**: an effect that starts a `setTimeout` and never clears it leaks timers and publishes stale values. Get the cleanup reflex burned in here; you'll reuse this hook in exercise 10.

## Concepts

- **Effect cleanup runs before the next effect and on unmount.** When `value` changes, React runs your cleanup (clearing the old timer) *before* running the effect again (starting a new timer). That cancellation is what makes debouncing work.
- **One timer at a time.** Because each change clears the previous timeout, only the final value in a burst ever calls `setDebounced`.
- **Dependencies matter**: `[value, delay]` — the effect re-runs whenever either changes.

## Your task

Implement the effect body in [`useDebounce.js`](useDebounce.js): start a timeout, return a cleanup that clears it.

## Run

```bash
npx vitest run exercises/07-use-debounce
```

**Manual testing:** this is a hook-only exercise — there's no UI to render, so drive it through the test file above (or a scratch component).

## Interview follow-ups to expect

- *"What happens without the cleanup?"* — every keystroke stacks a new timer; they all eventually fire, and the value flickers through the whole history instead of settling on the last one.
- *"Debounce vs throttle?"* — debounce waits for silence (fires once, after the last event); throttle fires at most once per interval during activity.
- *"Debounce a callback instead of a value."* — return a memoized function that internally manages a timer ref (`useDebouncedCallback`).

## Hints

```js
useEffect(() => {
  const id = setTimeout(() => setDebounced(value), delay)
  return () => clearTimeout(id)
}, [value, delay])
```

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState, useEffect } from 'react'

export function useDebounce(value, delay = 500) {
  const [debounced, setDebounced] = useState(value)

  useEffect(() => {
    const id = setTimeout(() => setDebounced(value), delay)
    return () => clearTimeout(id)
  }, [value, delay])

  return debounced
}
```

</details>
