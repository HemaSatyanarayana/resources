# 05 — Tabs

Tabs look like Accordion's cousin, but interviewers use them to probe **accessibility depth**. Anyone can toggle a panel; the signal is whether you reach for the real ARIA tab pattern — `role="tablist"`/`role="tab"`/`role="tabpanel"`, `aria-selected`, and **roving tabindex** with arrow-key navigation.

## Concepts

- **Single source of truth**: `activeId`. Everything else (which panel, which `aria-selected`) is derived.
- **Roving tabindex**: exactly one tab is in the tab order (`tabIndex={0}`); the rest are `-1`. Arrow keys move between them. This is how real tab widgets avoid trapping keyboard users in a long list of tabs.
- **Wrapping navigation** with modulo: `(index + 1) % length` and `(index - 1 + length) % length`.
- **Semantics over `<div>`s**: `role="tab"` on a `<button>` gives you focus + Enter/Space for free.

## Your task

Implement `onKeyDown` and wire up `aria-selected` / `tabIndex` in [`Tabs.jsx`](Tabs.jsx).

## Run

```bash
npx vitest run exercises/05-tabs
```

## Interview follow-ups to expect

- *"What's roving tabindex and why?"* — keeping every tab at `tabIndex={0}` forces a keyboard user to Tab through all of them; the ARIA pattern puts arrow keys in charge of moving *within* the widget and Tab moves *past* it.
- *"Lazy-render panels."* — only mount the active panel (as here) vs. mount all and hide — trade DOM size for state preservation.
- *"Home/End keys."* — jump to first/last tab.

## Hints

```js
const idx = tabs.findIndex((t) => t.id === activeId)
if (e.key === 'ArrowRight') setActiveId(tabs[(idx + 1) % tabs.length].id)
if (e.key === 'ArrowLeft') setActiveId(tabs[(idx - 1 + tabs.length) % tabs.length].id)
```

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState } from 'react'

export default function Tabs({ tabs = [], defaultActiveId }) {
  const [activeId, setActiveId] = useState(defaultActiveId ?? tabs[0]?.id)
  const activeTab = tabs.find((t) => t.id === activeId)

  const onKeyDown = (e) => {
    const idx = tabs.findIndex((t) => t.id === activeId)
    if (e.key === 'ArrowRight') {
      setActiveId(tabs[(idx + 1) % tabs.length].id)
    } else if (e.key === 'ArrowLeft') {
      setActiveId(tabs[(idx - 1 + tabs.length) % tabs.length].id)
    }
  }

  return (
    <div>
      <div role="tablist" onKeyDown={onKeyDown}>
        {tabs.map((tab) => {
          const isActive = tab.id === activeId
          return (
            <button
              key={tab.id}
              role="tab"
              aria-selected={isActive}
              tabIndex={isActive ? 0 : -1}
              onClick={() => setActiveId(tab.id)}
            >
              {tab.label}
            </button>
          )
        })}
      </div>
      <div role="tabpanel">{activeTab?.content}</div>
    </div>
  )
}
```

</details>
