/**
 * createCounter(start = 0)
 * Returns an object with methods that share ONE private `count` variable via
 * closure. No `count` property should be visible on the returned object.
 *
 *   const c = createCounter(10)
 *   c.increment() // 11
 *   c.value()     // 11
 *   c.decrement() // 10
 *   c.reset()     // back to 10
 *
 * Two counters created separately must not share state.
 */
export function createCounter(start = 0) {
  // TODO: keep `count` in this scope; return { increment, decrement, reset, value }.
  throw new Error('TODO: implement createCounter')
}

/**
 * once(fn)
 * Returns a wrapper that invokes `fn` at most once. Subsequent calls return the
 * first call's cached result and do NOT call `fn` again. Preserve `this`.
 */
export function once(fn) {
  // TODO
  throw new Error('TODO: implement once')
}

/**
 * memoize(fn)
 * Returns a wrapper that caches results by its arguments (JSON is fine for the
 * key). A repeated call with the same args returns the cached value without
 * calling `fn` again.
 */
export function memoize(fn) {
  // TODO
  throw new Error('TODO: implement memoize')
}
