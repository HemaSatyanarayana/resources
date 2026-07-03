import { useState } from 'react'

/**
 * Pagination — page windowing + boundary logic.
 *
 * Props:
 *   items     any[]    the full list
 *   pageSize  number   items per page (default 10)
 *
 * Required UI (the tests rely on this):
 *   - a <ul> containing only the CURRENT page's items as <li>
 *   - a button "Previous" (disabled on the first page)
 *   - a button "Next"     (disabled on the last page)
 *   - one button per page numbered 1..totalPages; the current page's button
 *     has aria-current="page"
 *
 * Rules:
 *   - Page is 1-indexed. totalPages = ceil(items.length / pageSize), min 1.
 *   - Previous/Next move by one page and never go out of [1, totalPages].
 *   - Clicking a page number jumps to it.
 */
export default function Pagination({ items = [], pageSize = 10 }) {
  const [page, setPage] = useState(1)

  const totalPages = /* TODO: Math.max(1, Math.ceil(items.length / pageSize)) */ 1
  const pageItems = /* TODO: items.slice((page-1)*pageSize, page*pageSize) */ items

  const goTo = (p) => {
    // TODO: clamp p to [1, totalPages] and setPage.
  }

  return (
    <div>
      <ul>
        {pageItems.map((item, i) => (
          <li key={i}>{item}</li>
        ))}
      </ul>

      <button onClick={() => goTo(page - 1)} disabled={/* TODO */ false}>
        Previous
      </button>

      {Array.from({ length: totalPages }, (_, i) => {
        const p = i + 1
        return (
          <button
            key={p}
            aria-current={page === p ? 'page' : undefined}
            onClick={() => goTo(p)}
          >
            {p}
          </button>
        )
      })}

      <button onClick={() => goTo(page + 1)} disabled={/* TODO */ false}>
        Next
      </button>
    </div>
  )
}
