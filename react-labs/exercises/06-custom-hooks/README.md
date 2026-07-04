# 06 — Custom Hooks

Custom hooks are where interviewers check that you *understand* React rather than just use it: the rules of hooks, `useRef` as a mutable box, lazy initializers, and effect cleanup timing. You'll build three hooks that show up constantly in real codebases and follow-up questions.

## Concepts

- **A hook is just a function that calls other hooks.** Extracting logic into `useX` is how you share stateful behavior without HOCs or render props.
- **`useRef` is a mutable box that survives renders but never triggers one.** `usePrevious` leans on the fact that the effect updating the ref runs *after* render, so during render `ref.current` still holds the previous value.
- **Lazy initializer**: `useState(() => expensive())` runs the function only on the first render. Perfect for reading `localStorage` once.
- **Stable identity with `useCallback`**: returning a memoized `toggle` means consumers can pass it to `memo`'d children or effect deps without causing re-runs.
- **Functional updates** (`setValue(v => ...)`) let the setter work without closing over the current value.

## Your task

Implement `useToggle`, `usePrevious`, and `useLocalStorage` in [`hooks.js`](hooks.js). Tests use `@testing-library/react`'s `renderHook`.

## Run

```bash
npx vitest run exercises/06-custom-hooks
```

**Manual testing:** this is a hook-only exercise — there's no UI to render, so drive it through the test file above (or a scratch component).

## Interview follow-ups to expect

- *"Why does `usePrevious` need a ref instead of state?"* — storing the previous value in state would trigger an extra render and create a loop; a ref updates silently after commit.
- *"Make `useLocalStorage` sync across tabs."* — subscribe to the `storage` window event and update state when another tab writes.
- *"Why `useCallback` on `toggle`?"* — referential stability so it's safe in dependency arrays and `React.memo` children.

## Hints

- `useToggle`: `setValue(v => typeof next === 'boolean' ? next : !v)`.
- `usePrevious`: `useEffect(() => { ref.current = value }, [value])`; return `ref.current`.
- `useLocalStorage`: read with `const raw = localStorage.getItem(key); return raw != null ? JSON.parse(raw) : initial`, then persist in an effect.

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState, useRef, useEffect, useCallback } from 'react'

export function useToggle(initial = false) {
  const [value, setValue] = useState(initial)
  const toggle = useCallback((next) => {
    setValue((v) => (typeof next === 'boolean' ? next : !v))
  }, [])
  return [value, toggle, setValue]
}

export function usePrevious(value) {
  const ref = useRef(undefined)
  useEffect(() => {
    ref.current = value
  }, [value])
  return ref.current
}

export function useLocalStorage(key, initial) {
  const [value, setValue] = useState(() => {
    try {
      const raw = localStorage.getItem(key)
      return raw != null ? JSON.parse(raw) : initial
    } catch {
      return initial
    }
  })

  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(value))
  }, [key, value])

  return [value, setValue]
}
```

</details>
