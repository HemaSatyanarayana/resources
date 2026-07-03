import { useState, useRef, useEffect, useCallback } from 'react'

/**
 * useToggle — boolean state with a memoized toggle.
 *
 * Returns [value, toggle, setValue] where:
 *   - toggle()      flips the value
 *   - toggle(next)  sets it to `next` when `next` is a boolean
 *   - setValue      the raw setter
 *
 * `toggle` must be stable across renders (wrap it in useCallback).
 */
export function useToggle(initial = false) {
  const [value, setValue] = useState(initial)

  const toggle = useCallback((next) => {
    // TODO: if `next` is a boolean, set to it; otherwise flip the current value.
    // Use a functional update so it doesn't depend on `value`.
  }, [])

  return [value, toggle, setValue]
}

/**
 * usePrevious — returns the value from the PREVIOUS render.
 * On the first render it returns undefined.
 *
 * Hint: a ref survives renders but doesn't trigger them. Read it during render,
 * then update it in an effect (which runs *after* render commits).
 */
export function usePrevious(value) {
  const ref = useRef(undefined)

  // TODO: after each commit, store the current value in ref.current
  //       so the NEXT render sees it as "previous".

  return ref.current
}

/**
 * useLocalStorage — state that persists to window.localStorage under `key`.
 *
 * Returns [value, setValue]. setValue works like useState's setter, including
 * the functional-updater form. Reads the existing stored value on mount (lazy
 * initializer) and writes JSON on every change.
 */
export function useLocalStorage(key, initial) {
  const [value, setValue] = useState(() => {
    // TODO: read localStorage[key]; if present, JSON.parse it; else return `initial`.
    return initial
  })

  useEffect(() => {
    // TODO: write JSON.stringify(value) to localStorage[key].
  }, [key, value])

  return [value, setValue]
}
