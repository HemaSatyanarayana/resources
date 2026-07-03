import { useState, useRef, useEffect, useCallback } from 'react'

/**
 * InfiniteScroll — load more pages when a sentinel scrolls into view.
 *
 * Props:
 *   fetchPage  (page) => Promise<string[]>   page is 0-indexed; an empty array
 *                                            means there are no more pages.
 *
 * Required UI (the tests rely on this):
 *   - a <ul> of all loaded items as <li>
 *   - a sentinel element with data-testid="sentinel" that you observe with an
 *     IntersectionObserver; when it intersects, load the next page and append.
 *   - text "Loading…" while a page request is in flight
 *   - text "No more items" once fetchPage returns an empty page
 *
 * Rules:
 *   - Never fire two loads at once (guard on `loading`) and never load after done.
 *   - Disconnect the observer on cleanup.
 */
export default function InfiniteScroll({ fetchPage }) {
  const [items, setItems] = useState([])
  const [page, setPage] = useState(0)
  const [loading, setLoading] = useState(false)
  const [done, setDone] = useState(false)
  const sentinelRef = useRef(null)

  const loadNext = useCallback(async () => {
    // TODO:
    //   - bail out if `loading` or `done`
    //   - setLoading(true); const next = await fetchPage(page)
    //   - if next is empty -> setDone(true)
    //     else -> append to items and increment page
    //   - setLoading(false) at the end
  }, [loading, done, page, fetchPage])

  useEffect(() => {
    // TODO:
    //   - grab sentinelRef.current; if missing, return
    //   - create an IntersectionObserver whose callback calls loadNext() when
    //     entries[0].isIntersecting is true
    //   - observe the sentinel; return a cleanup that disconnects the observer
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
