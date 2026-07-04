# 02 — Todo List

The most-asked machine coding question, period. It exercises the full loop: **controlled input → immutable array updates → list rendering with keys → derived filtering**. If you can build this cleanly and talk through it, you clear the bar of most rounds.

## Concepts

- **Controlled input**: the `<input value={text}>` is driven by state; `onChange` is the only way it changes.
- **Immutable updates**: adding is `[...todos, newTodo]`; toggling is `todos.map(...)`; deleting is `todos.filter(...)`. You never call `.push`, `.splice`, or assign to an element.
- **Stable keys**: each todo has an `id`. Use it as the `key` — *never the array index*, because deleting/reordering makes indexes point at the wrong element and React reuses the wrong DOM state.
- **Derived filtering**: the visible list is computed from `todos` + `filter` during render. You don't keep a second "filtered" array in state.

## Your task

Implement `addTodo`, `toggle`, `remove`, and the `visible` derivation in [`TodoList.jsx`](TodoList.jsx).

## Run

```bash
npx vitest run exercises/02-todo-list
```

Run the tests until green, then **try it in the browser**:

```bash
npm run dev
```

Open http://localhost:5173/#02-todo-list (or pick it from the dropdown). Edits hot-reload. Tweak the demo props for this exercise in `dev/main.jsx`.

## Interview follow-ups to expect

- *"Why is the index a bad key here?"* — deleting item 0 shifts every index down; React matches old keys to new positions and can carry over checkbox/focus state to the wrong row.
- *"Add editing / a 'clear completed' button / a remaining-count."* — all are more `map`/`filter`; the state shape shouldn't change.
- *"Persist across reload."* — hook it to `localStorage` (see exercise 06).

## Hints

- `toggle`: `todos.map(t => t.id === id ? { ...t, completed: !t.completed } : t)`.
- `remove`: `todos.filter(t => t.id !== id)`.
- `visible`: `filter === 'active' ? todos.filter(t => !t.completed) : filter === 'completed' ? todos.filter(t => t.completed) : todos`.

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState } from 'react'

let nextId = 1

export default function TodoList() {
  const [todos, setTodos] = useState([])
  const [text, setText] = useState('')
  const [filter, setFilter] = useState('all')

  const addTodo = () => {
    const trimmed = text.trim()
    if (!trimmed) return
    setTodos((ts) => [...ts, { id: nextId++, text: trimmed, completed: false }])
    setText('')
  }

  const toggle = (id) =>
    setTodos((ts) =>
      ts.map((t) => (t.id === id ? { ...t, completed: !t.completed } : t))
    )

  const remove = (id) => setTodos((ts) => ts.filter((t) => t.id !== id))

  const visible =
    filter === 'active'
      ? todos.filter((t) => !t.completed)
      : filter === 'completed'
        ? todos.filter((t) => t.completed)
        : todos

  return (
    <div>
      <input
        aria-label="New todo"
        value={text}
        onChange={(e) => setText(e.target.value)}
      />
      <button onClick={addTodo}>Add</button>

      <div>
        <button onClick={() => setFilter('all')}>All</button>
        <button onClick={() => setFilter('active')}>Active</button>
        <button onClick={() => setFilter('completed')}>Completed</button>
      </div>

      <ul>
        {visible.map((todo) => (
          <li key={todo.id}>
            <label>
              <input
                type="checkbox"
                checked={todo.completed}
                onChange={() => toggle(todo.id)}
              />
              {todo.text}
            </label>
            <button onClick={() => remove(todo.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  )
}
```

</details>
