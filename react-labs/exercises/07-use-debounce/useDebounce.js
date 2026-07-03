import { useState, useEffect } from 'react'

/**
 * useDebounce — returns a copy of `value` that only updates after `value` has
 * stopped changing for `delay` milliseconds.
 *
 * Behavior:
 *   - Returns the current `value` immediately on first render.
 *   - Each time `value` changes, start a timer; when it fires, publish the value.
 *   - If `value` changes again before the timer fires, the pending update is
 *     cancelled (only the latest value ever lands).
 *
 * The cleanup that cancels the pending timer is the entire point of this
 * exercise — an effect that sets a timeout MUST clear it on cleanup.
 */
export function useDebounce(value, delay = 500) {
  const [debounced, setDebounced] = useState(value)

  useEffect(() => {
    // TODO:
    //   1. start a setTimeout that calls setDebounced(value) after `delay`
    //   2. return a cleanup function that clears that timeout
  }, [value, delay])

  return debounced
}
