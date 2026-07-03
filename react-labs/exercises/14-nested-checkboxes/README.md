# 14 — Nested Checkboxes

The recursion question. A file-tree / permissions widget where a parent reflects its children: checked when all are, **indeterminate** when some are. The senior move is choosing a state model where parent status is *derived*, not stored — otherwise you drown in sync bugs.

## Concepts

- **Recursive components.** `TreeNode` renders itself for each child. The base case is a leaf (no children).
- **Store leaves, derive parents.** Keep a `Set` of checked *leaf* ids. A parent's checked/indeterminate status is computed from its leaf descendants each render. Storing a checkbox state per node forces you to keep parents and children manually in sync — the source of every bug in this problem.
- **`indeterminate` is a DOM property, not an attribute.** You can't set it in JSX; set it via a `ref` callback: `el.indeterminate = someChecked`.
- **Toggle is all-or-nothing on a subtree.** Clicking a node sets every leaf under it to the same new value.

## Your task

Implement `allChecked` / `someChecked` in `TreeNode` and `onToggle` in [`NestedCheckboxes.jsx`](NestedCheckboxes.jsx). A `leafIds` helper is provided.

## Run

```bash
npx vitest run exercises/14-nested-checkboxes
```

## Interview follow-ups to expect

- *"Why store leaves instead of a checked flag per node?"* — parent state is a pure function of its leaves, so deriving it can never drift; storing both means every toggle must update ancestors and descendants consistently.
- *"Make it controlled."* — lift the `Set` to props with an `onChange`.
- *"Huge tree — performance?"* — memoize `leafIds` per node, or precompute a parent→leaves index so each toggle is O(subtree) not O(tree).

## Hints

- `allChecked = checkedCount === leaves.length && leaves.length > 0`.
- `someChecked = checkedCount > 0 && !allChecked`.
- In `onToggle`: `const next = new Set(checked)`; if all leaves are in it, delete them all, else add them all; `setChecked(next)`.

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState } from 'react'

function leafIds(node) {
  if (!node.children || node.children.length === 0) return [node.id]
  return node.children.flatMap(leafIds)
}

function TreeNode({ node, checked, onToggle }) {
  const leaves = leafIds(node)
  const checkedCount = leaves.filter((id) => checked.has(id)).length
  const allChecked = checkedCount === leaves.length && leaves.length > 0
  const someChecked = checkedCount > 0 && !allChecked

  return (
    <li>
      <label>
        <input
          type="checkbox"
          checked={allChecked}
          ref={(el) => {
            if (el) el.indeterminate = someChecked
          }}
          onChange={() => onToggle(node)}
        />
        {node.label}
      </label>
      {node.children && (
        <ul>
          {node.children.map((child) => (
            <TreeNode key={child.id} node={child} checked={checked} onToggle={onToggle} />
          ))}
        </ul>
      )}
    </li>
  )
}

export default function NestedCheckboxes({ tree }) {
  const [checked, setChecked] = useState(() => new Set())

  const onToggle = (node) => {
    const leaves = leafIds(node)
    const allChecked = leaves.every((id) => checked.has(id))
    const next = new Set(checked)
    leaves.forEach((id) => (allChecked ? next.delete(id) : next.add(id)))
    setChecked(next)
  }

  return (
    <ul>
      <TreeNode node={tree} checked={checked} onToggle={onToggle} />
    </ul>
  )
}
```

</details>
