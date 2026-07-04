import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

// Grab every exercise component (default export). Hook-only exercises
// (06, 07, 09) have no default component and are skipped automatically.
const modules = import.meta.glob(
  ['../exercises/*/*.jsx', '!../exercises/*/*.test.jsx'],
  { eager: true },
)

// Demo props per exercise slug — tweak these to exercise different behaviour.
const demoProps = {
  '01-counter': { initial: 0, min: 0, max: 10, step: 1 },
  '03-star-rating': { max: 5 },
  '04-accordion': {
    items: [
      { id: 'a', title: 'Section A', content: 'Content A' },
      { id: 'b', title: 'Section B', content: 'Content B' },
    ],
  },
  '05-tabs': {
    tabs: [
      { id: 't1', label: 'One', content: 'Panel one' },
      { id: 't2', label: 'Two', content: 'Panel two' },
    ],
  },
}

// Build slug -> component map from the file paths.
const exercises = {}
for (const [path, mod] of Object.entries(modules)) {
  const slug = path.match(/exercises\/([^/]+)\//)[1]
  // Skip test files (they don't have a default component export anyway).
  if (path.endsWith('.test.jsx')) continue
  if (mod.default) exercises[slug] = mod.default
}

const slugs = Object.keys(exercises).sort()
const current = window.location.hash.slice(1) || slugs[0]

// Render the picker.
const controls = document.getElementById('controls')
controls.innerHTML =
  '<label>Exercise: <select id="picker">' +
  slugs
    .map(
      (s) =>
        `<option value="${s}" ${s === current ? 'selected' : ''}>${s}</option>`,
    )
    .join('') +
  '</select></label>'
document.getElementById('picker').addEventListener('change', (e) => {
  window.location.hash = e.target.value
  window.location.reload()
})

const Component = exercises[current]
const root = createRoot(document.getElementById('root'))
if (Component) {
  root.render(
    <StrictMode>
      <Component {...(demoProps[current] || {})} />
    </StrictMode>,
  )
} else {
  document.getElementById('root').textContent =
    `No visual component for "${current}" (hook-only exercise — use the test file).`
}
