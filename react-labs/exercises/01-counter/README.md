# 01 — Counter

The interview warm-up. It looks trivial, but interviewers use it to check three habits immediately: **functional state updates**, **derived disabled state**, and **clamping**. Nail the fundamentals here fast so you have time for the hard questions.

## Concepts

- **`useState`** holds the count. That is the *only* state you need — `disabled` is derived from it, not stored.
- **Functional updates** (`setCount(c => c + step)`) read the latest value and are safe under batching. Prefer them whenever the next value depends on the previous one.
- **Clamping** keeps the value inside `[min, max]`: `Math.min(max, Math.max(min, next))`.
- **Derived UI**: `disabled={count >= max}` — compute it during render, never keep it in state.

## Your task

Implement the three handlers and the two `disabled` props in [`Counter.jsx`](Counter.jsx):

| Piece | Skill |
|-------|-------|
| `increment` / `decrement` | Functional updates + clamping |
| `reset` | Reset to `initial` |
| `disabled` on each button | Deriving UI from state |

## Run

```bash
npx vitest run exercises/01-counter
```

## Interview follow-ups to expect

- *"Why `setCount(c => c + 1)` instead of `setCount(count + 1)`?"* — because `count` is a stale closure value; the functional form always sees the latest committed state, which matters under React's batching and in async callbacks.
- *"Add a `step` prop / long-press to auto-increment."* — the clamping should still hold.
- *"Store `disabled` in state?"* — no; it's derivable, and duplicated state drifts out of sync.

## Hints

- Do the clamp in one place. A tiny helper `clamp(n) => Math.min(max, Math.max(min, n))` keeps both handlers clean.
- `disabled={count >= max}` and `disabled={count <= min}`.

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState } from 'react'

export default function Counter({
  initial = 0,
  min = -Infinity,
  max = Infinity,
  step = 1,
}) {
  const [count, setCount] = useState(initial)
  const clamp = (n) => Math.min(max, Math.max(min, n))

  const increment = () => setCount((c) => clamp(c + step))
  const decrement = () => setCount((c) => clamp(c - step))
  const reset = () => setCount(initial)

  return (
    <div>
      <output data-testid="count">{count}</output>
      <button onClick={decrement} disabled={count <= min}>
        Decrement
      </button>
      <button onClick={increment} disabled={count >= max}>
        Increment
      </button>
      <button onClick={reset}>Reset</button>
    </div>
  )
}
```

</details>
