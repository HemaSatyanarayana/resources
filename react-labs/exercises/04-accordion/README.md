# 04 — Accordion

Classic "render a config-driven widget" question. The interesting decision is the **shape of the open/closed state** — one boolean per item scales badly; a single collection of open ids handles both single- and multi-open modes with the same code.

## Concepts

- **State shape drives simplicity.** Track `openIds` (an array or `Set`) rather than a boolean field on each item. `isOpen = openIds.includes(id)` is derived per render.
- **Single vs multi** is just a branch in the toggle: replace the collection with `[id]` vs append to it.
- **Conditional rendering** (`{isOpen && <panel/>}`) keeps closed content out of the DOM — cheaper and simpler to test than CSS hiding.
- **`aria-expanded`** on the header button is the accessible contract for a disclosure widget.

## Your task

Implement `toggle` in [`Accordion.jsx`](Accordion.jsx).

## Run

```bash
npx vitest run exercises/04-accordion
```

Run the tests until green, then **try it in the browser**:

```bash
npm run dev
```

Open http://localhost:5173/#04-accordion (or pick it from the dropdown). Edits hot-reload. Tweak the demo props for this exercise in `dev/main.jsx`.

## Interview follow-ups to expect

- *"Now allow multiple open."* — you should only change the toggle branch, nothing else. That's the payoff of the `openIds` shape.
- *"Animate open/close."* — CSS hiding (`max-height`/`hidden`) instead of unmounting, so there's an element to transition. Discuss the trade-off vs testability.
- *"Keyboard support."* — the header is already a `<button>`, so Enter/Space work for free; add arrow-key navigation between headers for full WAI-ARIA compliance.

## Hints

```js
setOpenIds((prev) => {
  const isOpen = prev.includes(id)
  if (allowMultiple) return isOpen ? prev.filter((x) => x !== id) : [...prev, id]
  return isOpen ? [] : [id]
})
```

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState } from 'react'

export default function Accordion({ items = [], allowMultiple = false }) {
  const [openIds, setOpenIds] = useState([])

  const toggle = (id) => {
    setOpenIds((prev) => {
      const isOpen = prev.includes(id)
      if (allowMultiple) {
        return isOpen ? prev.filter((x) => x !== id) : [...prev, id]
      }
      return isOpen ? [] : [id]
    })
  }

  return (
    <div>
      {items.map((item) => {
        const isOpen = openIds.includes(item.id)
        return (
          <div key={item.id}>
            <button aria-expanded={isOpen} onClick={() => toggle(item.id)}>
              {item.title}
            </button>
            {isOpen && <div role="region">{item.content}</div>}
          </div>
        )
      })}
    </div>
  )
}
```

</details>
