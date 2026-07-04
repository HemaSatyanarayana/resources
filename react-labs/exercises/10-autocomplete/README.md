# 10 — Autocomplete

The flagship async component. It stacks everything: a **debounced** query (exercise 07), an **async fetch** with **stale-response** handling (exercise 09), **keyboard navigation** with a highlight index (exercise 05), and accessible combobox/listbox semantics. Expect to build a real version of this in a senior round.

## Concepts

- **Debounce the request, not the input.** The input stays instant (controlled); the *fetch* waits for a pause in typing. Fire on every keystroke and you hammer the API and cause flicker.
- **Race conditions are the trap.** If "ap" resolves after "apr", you must not show "ap" results. The same `ignore`/cleanup pattern from `useFetch` cancels the late one.
- **Highlight index is derived UI state.** `active` is an index into `results`; `aria-selected` comes from `i === active`. ArrowDown/Up clamp within bounds.
- **Combobox semantics**: `role="combobox"` + `aria-expanded` on the input, `role="listbox"`/`role="option"` on the menu. This is the accessibility payoff.

## Your task

Implement the debounced-fetch effect, `choose`, and `onKeyDown` in [`Autocomplete.jsx`](Autocomplete.jsx).

## Run

```bash
npx vitest run exercises/10-autocomplete
```

Run the tests until green, then **try it in the browser**:

```bash
npm run dev
```

Open http://localhost:5173/#10-autocomplete (or pick it from the dropdown). Edits hot-reload. Tweak the demo props for this exercise in `dev/main.jsx`.

## Interview follow-ups to expect

- *"A stale request resolves last — what shows?"* — describe the cleanup/`ignore` guard (and `AbortController`).
- *"Highlight follows the mouse too."* — add `onMouseEnter` to set `active`.
- *"Cache results per query."* — memoize a `Map<query, results>` to avoid refetching.
- *"Why debounce and not throttle here?"* — you want the request only after the user *pauses*, i.e. debounce.

## Hints

```js
useEffect(() => {
  if (!query) { setResults([]); setOpen(false); return }
  let ignore = false
  const id = setTimeout(async () => {
    const items = await fetchSuggestions(query)
    if (!ignore) { setResults(items); setOpen(true); setActive(-1) }
  }, debounceMs)
  return () => { ignore = true; clearTimeout(id) }
}, [query, debounceMs, fetchSuggestions])
```

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState, useEffect } from 'react'

export default function Autocomplete({ fetchSuggestions, onSelect, debounceMs = 300 }) {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState([])
  const [open, setOpen] = useState(false)
  const [active, setActive] = useState(-1)

  useEffect(() => {
    if (!query) {
      setResults([])
      setOpen(false)
      return
    }
    let ignore = false
    const id = setTimeout(async () => {
      const items = await fetchSuggestions(query)
      if (!ignore) {
        setResults(items)
        setOpen(true)
        setActive(-1)
      }
    }, debounceMs)
    return () => {
      ignore = true
      clearTimeout(id)
    }
  }, [query, debounceMs, fetchSuggestions])

  const choose = (value) => {
    onSelect?.(value)
    setQuery(value)
    setOpen(false)
    setActive(-1)
  }

  const onKeyDown = (e) => {
    if (!open) return
    if (e.key === 'ArrowDown') {
      setActive((i) => Math.min(i + 1, results.length - 1))
    } else if (e.key === 'ArrowUp') {
      setActive((i) => Math.max(i - 1, 0))
    } else if (e.key === 'Enter') {
      if (active >= 0) choose(results[active])
    } else if (e.key === 'Escape') {
      setOpen(false)
    }
  }

  return (
    <div>
      <input
        aria-label="Search"
        role="combobox"
        aria-expanded={open}
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onKeyDown={onKeyDown}
      />
      {open && results.length > 0 && (
        <ul role="listbox">
          {results.map((item, i) => (
            <li
              key={item}
              role="option"
              aria-selected={i === active}
              onClick={() => choose(item)}
            >
              {item}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
```

</details>
