# 03 — Star Rating

A reusable widget with a subtlety interviewers watch for: it has **two** pieces of state — what you've *selected* and what you're *hovering* — and the rendered fill is *derived* from both. Getting the hover/restore behavior right without a tangle of flags is the whole exercise.

## Concepts

- **Two states, one derived value.** `selected` persists; `hovered` is transient (`null` when not hovering). The value you actually render is `hovered ?? selected`.
- **`Array.from({ length: count })`** is the clean way to render N of something without a real array.
- **Derived fill**: star `n` is filled iff `n <= displayValue`. No per-star boolean state.
- **Accessibility**: each star is a real `<button>` with a descriptive `aria-label`, so it's keyboard-focusable and screen-reader-friendly.

## Your task

Fill in `displayValue`, `select`, the `filled` calc, and the two mouse handlers in [`StarRating.jsx`](StarRating.jsx).

## Run

```bash
npx vitest run exercises/03-star-rating
```

Run the tests until green, then **try it in the browser**:

```bash
npm run dev
```

Open http://localhost:5173/#03-star-rating (or pick it from the dropdown). Edits hot-reload. Tweak the demo props for this exercise in `dev/main.jsx`.

## Interview follow-ups to expect

- *"Make it controlled."* — accept a `value` prop and call `onChange`; drop internal `selected` and read from props. Discuss controlled vs uncontrolled.
- *"Support half-stars."* — the value becomes a float; fill logic compares `n - 0.5`.
- *"Keyboard support (arrow keys)."* — arrow left/right adjusts a focused rating.

## Hints

- `const displayValue = hovered ?? selected` — `??` (not `||`) so a hovered `0` still works.
- Render `filled` as `position <= displayValue`.
- `onMouseEnter={() => setHovered(position)}`, `onMouseLeave={() => setHovered(null)}`.

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState } from 'react'

export default function StarRating({ count = 5, defaultValue = 0, onChange }) {
  const [selected, setSelected] = useState(defaultValue)
  const [hovered, setHovered] = useState(null)

  const displayValue = hovered ?? selected

  const select = (n) => {
    setSelected(n)
    onChange?.(n)
  }

  return (
    <div>
      {Array.from({ length: count }, (_, i) => {
        const position = i + 1
        const filled = position <= displayValue
        return (
          <button
            key={position}
            aria-label={`Rate ${position} star${position === 1 ? '' : 's'}`}
            data-filled={filled}
            onClick={() => select(position)}
            onMouseEnter={() => setHovered(position)}
            onMouseLeave={() => setHovered(null)}
          >
            {filled ? '★' : '☆'}
          </button>
        )
      })}
    </div>
  )
}
```

</details>
