/**
 * curry(fn)
 * Returns a curried version of `fn`. It collects arguments until it has at
 * least `fn.length` of them, then invokes `fn`. All of these are equivalent:
 *
 *   const add = (a, b, c) => a + b + c
 *   const c = curry(add)
 *   c(1)(2)(3)     // 6
 *   c(1, 2)(3)     // 6
 *   c(1)(2, 3)     // 6
 *   c(1, 2, 3)     // 6
 */
export function curry(fn) {
  // TODO: return a `curried` function that, if it has enough args, calls fn;
  // otherwise returns a function collecting more.
  throw new Error('TODO: implement curry')
}

/**
 * partial(fn, ...preset)
 * Returns a function with the `preset` args pre-filled; call-time args follow.
 *
 *   const add = (a, b, c) => a + b + c
 *   partial(add, 1, 2)(3) // 6
 */
export function partial(fn, ...preset) {
  // TODO
  throw new Error('TODO: implement partial')
}

/**
 * add(a) — infinite currying that coerces to a number.
 * add(1)(2)(3) can be chained any number of times; converting the result to a
 * number yields the running sum:
 *
 *   Number(add(1)(2)(3)) // 6
 *   +add(5)(10)          // 15
 *
 * Hint: return a function that also has a custom `valueOf`.
 */
export function add(a) {
  // TODO
  throw new Error('TODO: implement add')
}
