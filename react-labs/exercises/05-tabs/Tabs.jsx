import { useState } from 'react'

/**
 * Tabs — active-index state, keyboard navigation, ARIA tab semantics.
 *
 * Props:
 *   tabs            [{ id, label, content }]
 *   defaultActiveId string  which tab starts active (defaults to the first tab)
 *
 * Required UI (the tests rely on this):
 *   - a container with role="tablist"
 *   - one <button role="tab"> per tab, text = label, with aria-selected="true"|"false"
 *   - a single element with role="tabpanel" rendering the ACTIVE tab's content
 *   - the active tab has tabIndex={0}, inactive tabs tabIndex={-1} (roving tabindex)
 *
 * Rules:
 *   - Clicking a tab activates it.
 *   - With focus inside the tablist: ArrowRight moves to the next tab (wrapping),
 *     ArrowLeft to the previous (wrapping), and activates it.
 */
export default function Tabs({ tabs = [], defaultActiveId }) {
  const [activeId, setActiveId] = useState(defaultActiveId ?? tabs[0]?.id)

  const activeTab = tabs.find((t) => t.id === activeId)

  const onKeyDown = (e) => {
    // TODO: on ArrowRight / ArrowLeft, compute the next index (wrapping with
    // modulo) and setActiveId to that tab's id.
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
              aria-selected={/* TODO: isActive */ false}
              tabIndex={/* TODO: isActive ? 0 : -1 */ 0}
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
