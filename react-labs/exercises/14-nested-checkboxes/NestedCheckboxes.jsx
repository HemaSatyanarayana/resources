import { useState } from 'react'

/**
 * NestedCheckboxes — a recursive checkbox tree with parent/child sync.
 *
 * Props:
 *   tree  a node: { id, label, children?: node[] }
 *
 * Behavior:
 *   - Checking a parent checks every descendant; unchecking unchecks them all.
 *   - A parent is CHECKED when all of its leaf descendants are checked,
 *     INDETERMINATE when only some are, and unchecked when none are.
 *   - Each checkbox's accessible name is its node label.
 *
 * State model: keep a Set of checked LEAF ids. Every parent's checked/
 * indeterminate status is DERIVED from its leaves — don't store it.
 */

// Returns the ids of all leaf descendants of `node` (or [node.id] if it's a leaf).
function leafIds(node) {
  if (!node.children || node.children.length === 0) return [node.id]
  return node.children.flatMap(leafIds)
}

function TreeNode({ node, checked, onToggle }) {
  const leaves = leafIds(node)
  const checkedCount = leaves.filter((id) => checked.has(id)).length

  const allChecked = /* TODO: every leaf checked */ false
  const someChecked = /* TODO: at least one, but not all */ false

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
    // TODO:
    //   - compute the node's leaf ids
    //   - if they're ALL currently checked, remove them all; otherwise add them all
    //   - setChecked to a NEW Set (never mutate the existing one in place)
  }

  return (
    <ul>
      <TreeNode node={tree} checked={checked} onToggle={onToggle} />
    </ul>
  )
}
