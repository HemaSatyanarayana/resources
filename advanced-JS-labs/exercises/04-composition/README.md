# 04 — Composition

`compose` and `pipe` glue small functions into a pipeline. They're two lines of `reduce`, but they cement how data flows through a chain and are the backbone of functional libraries (Redux middleware, RxJS operators, Ramda).

## Concepts

- **`pipe`** reads in execution order (left → right): `pipe(f, g)(x)` runs `f` then `g`. Most people find this the natural one.
- **`compose`** is the math convention (right → left): `compose(f, g)(x) = f(g(x))`.
- **`reduce` vs `reduceRight`** is the *only* difference between them.
- **Point-free style**: you build new functions by combining existing ones without mentioning the data.

## Your task

Implement `pipe` and `compose` in [`compose.js`](compose.js).

## Run

```bash
npx vitest run exercises/04-composition
```

## Interview follow-ups to expect

- *"Make it work with multiple initial arguments."* — let the first function be variadic: `(...args) => rest.reduce((acc, fn) => fn(acc), first(...args))`.
- *"Support async (a promise pipeline)."* — `reduce` starting from `Promise.resolve(x)` and `await`ing each step.
- *"Where have you seen this?"* — Redux's `compose` for `applyMiddleware`, RxJS `pipe`.

## Hints

```js
export const pipe = (...fns) => (x) => fns.reduce((acc, fn) => fn(acc), x)
export const compose = (...fns) => (x) => fns.reduceRight((acc, fn) => fn(acc), x)
```

<details>
<summary>Reference solution (try first!)</summary>

```js
export const identity = (x) => x

export function pipe(...fns) {
  return (x) => fns.reduce((acc, fn) => fn(acc), x)
}

export function compose(...fns) {
  return (x) => fns.reduceRight((acc, fn) => fn(acc), x)
}
```

</details>
