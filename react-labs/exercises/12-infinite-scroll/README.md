# 12 — Infinite Scroll

The modern way to do "load more" — a `IntersectionObserver` watching a **sentinel** element at the bottom of the list. When the sentinel scrolls into view, you fetch the next page and append. Interviewers like it because it forces clean **async guarding** (no double-loads) and **observer cleanup**.

## Concepts

- **`IntersectionObserver`** tells you when an element enters the viewport without wiring up scroll listeners (which fire constantly and are janky). You observe a tiny sentinel `<div>` after the list.
- **Guard against concurrent loads.** The observer can fire repeatedly; bail if `loading` or `done` so you don't fetch page N twice.
- **Append, don't replace.** New items are `[...items, ...next]`.
- **Termination**: an empty page means stop — set `done` and never fetch again.
- **Cleanup**: `observer.disconnect()` on unmount, or you leak the observer.

## Your task

Implement `loadNext` and the observer effect in [`InfiniteScroll.jsx`](InfiniteScroll.jsx). The tests mock `IntersectionObserver` and trigger intersections manually.

> Note: jsdom doesn't implement `IntersectionObserver`, so the test file stubs it and calls the observer callback directly. That's also how you'd unit-test this in a real codebase.

## Run

```bash
npx vitest run exercises/12-infinite-scroll
```

## Interview follow-ups to expect

- *"Scroll listener vs IntersectionObserver?"* — the observer is async, batched, and off the main thread; scroll handlers fire on every pixel and need manual throttling.
- *"Windowing / virtualization."* — for very long lists, only render the visible rows (react-window) so the DOM stays small even as data grows.
- *"Restore scroll position / preserve on back-nav."* — cache items + scroll offset.
- *"What if a fetch fails mid-scroll?"* — surface a retry affordance; don't set `done`.

## Hints

```js
const loadNext = useCallback(async () => {
  if (loading || done) return
  setLoading(true)
  const next = await fetchPage(page)
  if (next.length === 0) setDone(true)
  else {
    setItems((prev) => [...prev, ...next])
    setPage((p) => p + 1)
  }
  setLoading(false)
}, [loading, done, page, fetchPage])

useEffect(() => {
  const el = sentinelRef.current
  if (!el) return
  const observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting) loadNext()
  })
  observer.observe(el)
  return () => observer.disconnect()
}, [loadNext])
```

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState, useRef, useEffect, useCallback } from 'react'

export default function InfiniteScroll({ fetchPage }) {
  const [items, setItems] = useState([])
  const [page, setPage] = useState(0)
  const [loading, setLoading] = useState(false)
  const [done, setDone] = useState(false)
  const sentinelRef = useRef(null)

  const loadNext = useCallback(async () => {
    if (loading || done) return
    setLoading(true)
    const next = await fetchPage(page)
    if (next.length === 0) {
      setDone(true)
    } else {
      setItems((prev) => [...prev, ...next])
      setPage((p) => p + 1)
    }
    setLoading(false)
  }, [loading, done, page, fetchPage])

  useEffect(() => {
    const el = sentinelRef.current
    if (!el) return
    const observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) loadNext()
    })
    observer.observe(el)
    return () => observer.disconnect()
  }, [loadNext])

  return (
    <div>
      <ul>
        {items.map((item, i) => (
          <li key={i}>{item}</li>
        ))}
      </ul>
      {loading && <p>Loading…</p>}
      {done && <p>No more items</p>}
      <div data-testid="sentinel" ref={sentinelRef} />
    </div>
  )
}
```

</details>
