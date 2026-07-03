import { useState, useEffect } from 'react'

/**
 * Autocomplete — debounced async search + keyboard navigation.
 *
 * Props:
 *   fetchSuggestions  (query) => Promise<string[]>   (injected so it's testable)
 *   onSelect          (value) => void
 *   debounceMs        number                          (default 300)
 *
 * Required UI (the tests rely on this):
 *   - a text input with aria-label "Search"
 *   - when there are results and the list is open: a <ul role="listbox">
 *     with one <li role="option"> per result; the highlighted option has
 *     aria-selected="true"
 *
 * Behavior:
 *   - Typing updates the query (controlled input).
 *   - The query is DEBOUNCED by `debounceMs`; only after the user stops typing
 *     do you call fetchSuggestions and open the list.
 *   - Empty query => no request, list closed.
 *   - ArrowDown / ArrowUp move the highlight; Enter selects the highlighted
 *     option; Escape closes the list.
 *   - Clicking an option or selecting with Enter: call onSelect(value), put the
 *     value in the input, and close the list.
 */
export default function Autocomplete({ fetchSuggestions, onSelect, debounceMs = 300 }) {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState([])
  const [open, setOpen] = useState(false)
  const [active, setActive] = useState(-1)

  useEffect(() => {
    // TODO: debounce + fetch.
    //   - if query is empty: clear results, close, return.
    //   - else start a setTimeout(debounceMs) that awaits fetchSuggestions(query),
    //     then (unless cancelled) sets results, opens the list, resets active to -1.
    //   - return a cleanup that cancels the timeout AND ignores a late response.
  }, [query, debounceMs, fetchSuggestions])

  const choose = (value) => {
    // TODO: onSelect?.(value); setQuery(value); close the list; reset active.
  }

  const onKeyDown = (e) => {
    if (!open) return
    // TODO:
    //   ArrowDown -> active = min(active + 1, results.length - 1)
    //   ArrowUp   -> active = max(active - 1, 0)
    //   Enter     -> if active >= 0, choose(results[active])
    //   Escape    -> close the list
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
