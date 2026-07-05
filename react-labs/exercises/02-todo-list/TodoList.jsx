import { useState } from "react";

/**
 * TodoList — controlled input + immutable list updates + filtering.
 *
 * Required UI (the tests rely on this):
 *   - a text input with aria-label "New todo"
 *   - a button labelled "Add"           -> adds the trimmed input as a todo
 *   - a <ul>; each todo is an <li> containing:
 *       * a checkbox whose accessible name is the todo text (toggles completed)
 *       * a button labelled "Delete"    -> removes that todo
 *   - three filter buttons: "All", "Active", "Completed"
 *
 * Rules:
 *   - Adding trims whitespace; empty/whitespace-only input adds nothing.
 *   - "Active" shows only not-completed todos; "Completed" shows only completed.
 *   - Completed todos render a checked checkbox.
 *   - Never mutate state — always build new arrays/objects.
 */
let nextId = 1;

export default function TodoList() {
  const [todos, setTodos] = useState([]); // [{ id, text, completed }]
  const [text, setText] = useState("");
  const [filter, setFilter] = useState("all"); // 'all' | 'active' | 'completed'

  const addTodo = () => {
    // TODO: trim `text`; if empty, do nothing. Otherwise append
    // { id: nextId++, text: trimmed, completed: false } and clear the input.

    if (text.trim() === "") {
      return;
    }

    setTodos((todos) => [
      ...todos,
      { id: nextId++, text: text.trim(), completed: false },
    ]);

    setText("");
    nextId++;
  };

  const toggle = (id) => {
    // TODO: return a NEW array where the matching todo's `completed` is flipped.
    if (todos.length === 0) {
      return;
    }

    setTodos((todos) =>
      todos.map((todo) =>
        todo.id === id ? { ...todo, completed: !todo.completed } : todo
      )
    );
  };

  const remove = (id) => {
    // TODO: filter out the todo with this id.

    if (todos.length === 0) {
      return;
    }

    setTodos((todos) => todos.filter((todo) => todo.id !== id));
  };

  const visible = todos.filter((todo) => {
    if (filter === "all") {
      return true;
    } else if (filter === "active") {
      return !todo.completed;
    } else if (filter === "completed") {
      return todo.completed;
    }
  }); // TODO: derive the visible list from `filter`.

  return (
    <div>
      <input
        aria-label="New todo"
        value={text}
        onChange={(e) => setText(e.target.value)}
      />
      <button onClick={addTodo}>Add</button>

      <div>
        <button onClick={() => setFilter("all")}>All</button>
        <button onClick={() => setFilter("active")}>Active</button>
        <button onClick={() => setFilter("completed")}>Completed</button>
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
  );
}
