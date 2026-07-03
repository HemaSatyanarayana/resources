# 13 — Modal

A modal exposes whether you understand **portals** (rendering outside the parent's DOM so `overflow:hidden`/`z-index` ancestors can't clip it) and **document-level event handling with proper cleanup**. The backdrop-vs-content click distinction is the detail interviewers poke at.

## Concepts

- **`createPortal(node, container)`** renders `node` into `container` (here `document.body`) while keeping it in the React tree — context and events still flow through the component hierarchy, but the DOM lives elsewhere. This escapes clipping/stacking-context bugs.
- **Global listeners belong in an effect, added only while open**, and removed in cleanup. A leftover `keydown` listener firing after unmount is a classic bug.
- **`e.target === e.currentTarget`** is the clean way to detect a click on the backdrop *itself*, not on bubbled clicks from the dialog content. (Alternative: `stopPropagation` on the dialog.)
- **Accessibility**: `role="dialog"` + `aria-modal="true"`. A production version also traps focus and restores it on close.

## Your task

Implement the Escape effect and `onBackdropClick` in [`Modal.jsx`](Modal.jsx).

## Run

```bash
npx vitest run exercises/13-modal
```

## Interview follow-ups to expect

- *"Why a portal?"* — a modal rendered deep in the tree can be clipped by an ancestor's `overflow:hidden` or trapped under its `z-index`/stacking context; a portal to `body` sidesteps both.
- *"Trap focus inside the modal."* — on open, focus the dialog; intercept Tab/Shift+Tab to cycle within focusable children; restore focus to the trigger on close.
- *"Lock body scroll while open."* — toggle `document.body.style.overflow` in the same effect.
- *"Why add the listener in an effect, not the render body?"* — the body is pure and runs repeatedly; effects run after commit and give you a cleanup hook.

## Hints

```js
useEffect(() => {
  if (!isOpen) return
  const onKey = (e) => { if (e.key === 'Escape') onClose?.() }
  document.addEventListener('keydown', onKey)
  return () => document.removeEventListener('keydown', onKey)
}, [isOpen, onClose])

const onBackdropClick = (e) => {
  if (e.target === e.currentTarget) onClose?.()
}
```

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useEffect } from 'react'
import { createPortal } from 'react-dom'

export default function Modal({ isOpen, onClose, children }) {
  useEffect(() => {
    if (!isOpen) return
    const onKey = (e) => {
      if (e.key === 'Escape') onClose?.()
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [isOpen, onClose])

  if (!isOpen) return null

  const onBackdropClick = (e) => {
    if (e.target === e.currentTarget) onClose?.()
  }

  return createPortal(
    <div data-testid="backdrop" onClick={onBackdropClick}>
      <div role="dialog" aria-modal="true">
        {children}
        <button onClick={onClose}>Close</button>
      </div>
    </div>,
    document.body
  )
}
```

</details>
