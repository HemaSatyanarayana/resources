# 08 — Stopwatch

A timer looks easy until the two classic bugs bite: **stale closures** (the tick reads an old `ms`) and **leaked intervals** (Start twice, or forget to clear on unmount). This exercise trains the `useRef`-for-mutable-ids + functional-update pattern that fixes both.

## Concepts

- **`useRef` for the interval id.** The id isn't render data — it shouldn't cause re-renders and must survive them so Stop can clear it. That's exactly what a ref is for.
- **Functional updates dodge stale closures.** `setMs(m => m + 100)` always adds to the latest value. If you wrote `setMs(ms + 100)`, the interval callback would forever see the `ms` captured when it was created (0), and the clock would freeze at 0.1s.
- **Every interval needs a clear.** On Stop, on Reset, and in an unmount cleanup effect. A running interval after unmount is a memory leak and a "setState on unmounted component" warning.
- **Guard double-start** so two intervals can't run at once (the `disabled` prop helps, but guard in code too).

## Your task

Implement `start`, `stop`, `reset`, and the cleanup effect in [`Stopwatch.jsx`](Stopwatch.jsx).

## Run

```bash
npx vitest run exercises/08-stopwatch
```

## Interview follow-ups to expect

- *"Why a ref and not state for the interval id?"* — you need to read/overwrite it without re-rendering, and it must persist across renders.
- *"Your clock drifts."* — `setInterval` isn't precise. For accuracy, store a start timestamp (`Date.now()`) and compute elapsed as `now - start`, using the interval only to trigger re-renders.
- *"Add lap times."* — push the current elapsed onto a `laps` array.

## Hints

```js
const start = () => {
  if (running) return
  setRunning(true)
  intervalRef.current = setInterval(() => setMs((m) => m + 100), 100)
}
const stop = () => {
  setRunning(false)
  clearInterval(intervalRef.current)
  intervalRef.current = null
}
```

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState, useRef, useEffect } from 'react'

export default function Stopwatch() {
  const [ms, setMs] = useState(0)
  const [running, setRunning] = useState(false)
  const intervalRef = useRef(null)

  const start = () => {
    if (running) return
    setRunning(true)
    intervalRef.current = setInterval(() => setMs((m) => m + 100), 100)
  }

  const stop = () => {
    setRunning(false)
    clearInterval(intervalRef.current)
    intervalRef.current = null
  }

  const reset = () => {
    stop()
    setMs(0)
  }

  useEffect(() => () => clearInterval(intervalRef.current), [])

  return (
    <div>
      <output data-testid="elapsed">{(ms / 1000).toFixed(1)}</output>
      <button onClick={start} disabled={running}>
        Start
      </button>
      <button onClick={stop} disabled={!running}>
        Stop
      </button>
      <button onClick={reset}>Reset</button>
    </div>
  )
}
```

</details>
