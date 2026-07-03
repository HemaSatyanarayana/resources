# 02 — this & binding

`this` is the interview topic people get wrong most. The rule: **`this` is determined by how a function is *called*, not where it's defined.** Reimplementing `call`, `apply`, and `bind` forces you to internalize that — the trick is that calling a function *as a method of an object* is what binds `this` to that object.

## Concepts

- **Call-site binding.** `obj.fn()` sets `this = obj`. So to force `this = X`, you temporarily make the function a property of `X`, call it as `X.fn()`, then clean up.
- **`call` vs `apply`** differ only in how args arrive (list vs array).
- **`bind`** doesn't call — it returns a new function with `this` (and optionally leading args) locked in. Call-time args append to the bound ones (partial application).
- **Default `this`.** In non-strict code, `null`/`undefined` `thisArg` falls back to `globalThis`.

## Your task

Implement `myCall`, `myApply`, and `myBind` in [`binding.js`](binding.js).

## Run

```bash
npx vitest run exercises/02-this-binding
```

## Interview follow-ups to expect

- *"Why does `const f = obj.method; f()` lose `this`?"* — extracting the function discards the call-site; `f()` is a plain call, so `this` is `undefined`/global. Fix with `bind` or an arrow wrapper.
- *"How do arrow functions differ?"* — they have no own `this`; they capture it lexically from the enclosing scope, and `call`/`bind` can't change it.
- *"Use a `Symbol` for the temporary property — why?"* — so you never clobber a real key on `thisArg`.

## Hints

```js
export function myCall(fn, thisArg, ...args) {
  const ctx = thisArg == null ? globalThis : Object(thisArg)
  const key = Symbol('fn')
  ctx[key] = fn
  try {
    return ctx[key](...args)
  } finally {
    delete ctx[key]
  }
}
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export function myCall(fn, thisArg, ...args) {
  const ctx = thisArg == null ? globalThis : Object(thisArg)
  const key = Symbol('fn')
  ctx[key] = fn
  try {
    return ctx[key](...args)
  } finally {
    delete ctx[key]
  }
}

export function myApply(fn, thisArg, argsArray) {
  return myCall(fn, thisArg, ...(argsArray || []))
}

export function myBind(fn, thisArg, ...bound) {
  return function (...args) {
    return myCall(fn, thisArg, ...bound, ...args)
  }
}
```

</details>
