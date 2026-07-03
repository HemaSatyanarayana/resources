# 11 — Deep Clone & Equal

Recursion over arbitrary data, with the two edge cases that trip people up: **circular references** (needs a `WeakMap`) and **special objects** (`Date`, `RegExp`). `deepEqual` adds the `Object.is` subtlety around `NaN` and `-0`.

## Concepts

- **Recurse on structure, copy leaves.** Primitives return as-is; arrays/objects are rebuilt key by key.
- **Circular refs** would infinite-loop a naive clone. Keep a `WeakMap` from *original → clone*; before recursing into an object, record its clone, and return the existing clone if you meet it again. `WeakMap` because keys are objects and shouldn't be pinned in memory.
- **Special objects** don't survive a plain key copy — `new Date(value)`, `new RegExp(value)` reconstruct them.
- **`Object.is` vs `===`.** `Object.is(NaN, NaN)` is `true` and `Object.is(0, -0)` is `false` — the right primitive comparison for equality.
- **`structuredClone`** is the built-in now; this exercise is about understanding what it does.

## Your task

Implement `deepClone` and `deepEqual` in [`deep.js`](deep.js).

## Run

```bash
npx vitest run exercises/11-deep
```

## Interview follow-ups to expect

- *"Why `WeakMap` and not `Map`?"* — it holds keys weakly so cloned sources can be garbage-collected; also keys must be objects, which is exactly what we're tracking.
- *"`JSON.parse(JSON.stringify(x))` — why not?"* — it drops `undefined`/functions/`Symbol`, mangles `Date` to a string, throws on cycles, and loses `Map`/`Set`.
- *"Extend to `Map`/`Set`."* — branch on `instanceof` and rebuild with cloned entries.

## Hints

```js
export function deepClone(value, seen = new WeakMap()) {
  if (value === null || typeof value !== 'object') return value
  if (value instanceof Date) return new Date(value)
  if (value instanceof RegExp) return new RegExp(value.source, value.flags)
  if (seen.has(value)) return seen.get(value)
  const clone = Array.isArray(value) ? [] : {}
  seen.set(value, clone)
  for (const key of Object.keys(value)) clone[key] = deepClone(value[key], seen)
  return clone
}
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export function deepClone(value, seen = new WeakMap()) {
  if (value === null || typeof value !== 'object') return value
  if (value instanceof Date) return new Date(value)
  if (value instanceof RegExp) return new RegExp(value.source, value.flags)
  if (seen.has(value)) return seen.get(value)

  const clone = Array.isArray(value) ? [] : {}
  seen.set(value, clone)
  for (const key of Object.keys(value)) {
    clone[key] = deepClone(value[key], seen)
  }
  return clone
}

export function deepEqual(a, b) {
  if (Object.is(a, b)) return true
  if (typeof a !== 'object' || a === null) return false
  if (typeof b !== 'object' || b === null) return false
  if (Array.isArray(a) !== Array.isArray(b)) return false

  const keysA = Object.keys(a)
  const keysB = Object.keys(b)
  if (keysA.length !== keysB.length) return false

  return keysA.every(
    (key) =>
      Object.prototype.hasOwnProperty.call(b, key) && deepEqual(a[key], b[key])
  )
}
```

</details>
