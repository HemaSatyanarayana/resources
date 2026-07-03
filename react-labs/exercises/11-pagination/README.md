# 11 — Pagination

Deceptively arithmetic-heavy. Interviewers use it to check that you keep **minimal state** (just the current page — everything else is `slice` math) and handle **boundaries** cleanly (disable at the ends, clamp jumps, and the ragged last page).

## Concepts

- **One piece of state: `page`.** The visible slice, the total page count, and the disabled flags are all *derived*.
- **The slice**: `items.slice((page - 1) * pageSize, page * pageSize)`. `slice` safely returns a short array on the last page — no manual bounds needed.
- **`totalPages = Math.max(1, Math.ceil(items.length / pageSize))`** — the `max(1, …)` keeps an empty list at one page.
- **Clamp every navigation** through one `goTo(p)` helper so Previous/Next and number clicks share the same guard.

## Your task

Compute `totalPages`, `pageItems`, implement `goTo`, and set the two `disabled` props in [`Pagination.jsx`](Pagination.jsx).

## Run

```bash
npx vitest run exercises/11-pagination
```

## Interview follow-ups to expect

- *"1000 pages — don't render 1000 buttons."* — window the page numbers: show current ±2 with `1 … 7 8 [9] 10 11 … 42` and ellipses. This is the real follow-up; think about the windowing formula.
- *"Server-side pagination."* — you no longer have all items; fetch `?page=&size=` and store `totalCount` from the response to compute `totalPages`.
- *"Reset to page 1 when `items` changes."* — an effect on `[items]`, or key the component.

## Hints

- `const goTo = (p) => setPage(Math.min(totalPages, Math.max(1, p)))`.
- `disabled={page === 1}` on Previous, `disabled={page === totalPages}` on Next.

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState } from 'react'

export default function Pagination({ items = [], pageSize = 10 }) {
  const [page, setPage] = useState(1)

  const totalPages = Math.max(1, Math.ceil(items.length / pageSize))
  const pageItems = items.slice((page - 1) * pageSize, page * pageSize)
  const goTo = (p) => setPage(Math.min(totalPages, Math.max(1, p)))

  return (
    <div>
      <ul>
        {pageItems.map((item, i) => (
          <li key={i}>{item}</li>
        ))}
      </ul>

      <button onClick={() => goTo(page - 1)} disabled={page === 1}>
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

      <button onClick={() => goTo(page + 1)} disabled={page === totalPages}>
        Next
      </button>
    </div>
  )
}
```

</details>
