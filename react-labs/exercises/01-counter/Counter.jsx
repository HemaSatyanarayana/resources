import { useState } from 'react'

/**
 * Counter — the classic warm-up.
 *
 * Props:
 *   initial  number  starting value (default 0)
 *   min      number  lower bound, inclusive (default -Infinity)
 *   max      number  upper bound, inclusive (default Infinity)
 *   step     number  amount to add/subtract (default 1)
 *
 * Required UI (the tests rely on this):
 *   - an element with data-testid="count" showing the current value
 *   - a button labelled "Increment"
 *   - a button labelled "Decrement"
 *   - a button labelled "Reset"
 *
 * Rules:
 *   - Increment/Decrement move by `step`, but the value must stay within [min, max].
 *   - The Increment button is disabled when the value is at max.
 *   - The Decrement button is disabled when the value is at min.
 *   - Reset returns to `initial`.
 */
export default function Counter({
  initial = 0,
  min = -Infinity,
  max = Infinity,
  step = 1,
}) {
  const [count, setCount] = useState(initial)

  const increment = () => {
    // TODO: increase by `step`, clamped to `max`. Use a functional update.
  }

  const decrement = () => {
    // TODO: decrease by `step`, clamped to `min`.
  }

  const reset = () => {
    // TODO: return to `initial`.
  }

  return (
    <div>
      <output data-testid="count">{count}</output>
      <button onClick={decrement} disabled={/* TODO: at min? */ false}>
        Decrement
      </button>
      <button onClick={increment} disabled={/* TODO: at max? */ false}>
        Increment
      </button>
      <button onClick={reset}>Reset</button>
    </div>
  )
}
