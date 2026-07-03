# 03 — Currying

Currying turns `f(a, b, c)` into `f(a)(b)(c)` — a function that collects arguments one (or several) at a time until it has enough, then fires. It's a favorite because it combines **closures** (accumulating args) with **`fn.length`** (arity introspection), and the "infinite sum" variant tests whether you know about `valueOf`.

## Concepts

- **`fn.length`** is the number of declared parameters — how `curry` knows when it has "enough."
- **Accumulate via closure + recursion.** Each partial call returns a new function that remembers the args so far and waits for more.
- **`partial`** is the simpler cousin: fix some leading args now, supply the rest later. (This is what `bind` does for args.)
- **Infinite currying** has no fixed arity, so it can't know when to stop — instead it returns a *callable that is also a number* by overriding `valueOf`, which the engine calls during numeric coercion (`+x`, `Number(x)`).

## Your task

Implement `curry`, `partial`, and `add` in [`curry.js`](curry.js).

## Run

```bash
npx vitest run exercises/03-currying
```

## Interview follow-ups to expect

- *"How does `curry` know the function is 'complete'?"* — it compares collected args against `fn.length`.
- *"Why does `+add(1)(2)` work?"* — `add` returns a function whose `valueOf` returns the running total; the unary `+` coerces it via `valueOf`.
- *"Curry vs partial application?"* — currying always takes one logical step at a time and is fully generic; partial just pre-binds some arguments.

## Hints

```js
export function curry(fn) {
  return function curried(...args) {
    if (args.length >= fn.length) return fn.apply(this, args)
    return (...next) => curried.apply(this, [...args, ...next])
  }
}
```

For `add`: return `fn` where `fn = (b) => add(a + b)` and `fn.valueOf = () => a`.

<details>
<summary>Reference solution (try first!)</summary>

```js
export function curry(fn) {
  return function curried(...args) {
    if (args.length >= fn.length) return fn.apply(this, args)
    return (...next) => curried.apply(this, [...args, ...next])
  }
}

export function partial(fn, ...preset) {
  return (...args) => fn(...preset, ...args)
}

export function add(a) {
  const fn = (b) => add(a + b)
  fn.valueOf = () => a
  return fn
}
```

</details>
