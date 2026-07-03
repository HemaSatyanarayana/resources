/**
 * sleep(ms)
 * Returns a Promise that resolves (with undefined) after `ms` milliseconds.
 * The building block of `await sleep(1000)`.
 */
export function sleep(ms) {
  // TODO: new Promise that resolves via setTimeout.
  throw new Error('TODO: implement sleep')
}

/**
 * repeat(fn, ms)
 * A setInterval built from setTimeout: call `fn` every `ms`, and return a
 * `cancel` function that stops future calls. (Recursive setTimeout is more
 * accurate than setInterval — the next timer is scheduled only after the
 * current callback finishes, so callbacks can't stack up.)
 */
export function repeat(fn, ms) {
  // TODO: schedule fn on a recurring setTimeout; return () => stop.
  throw new Error('TODO: implement repeat')
}

/**
 * pollUntil(predicate, { interval = 50, timeout = 2000 })
 * Repeatedly call `predicate` (which may be sync or return a Promise) every
 * `interval` ms. Resolve with the predicate's value as soon as it's truthy.
 * Reject with an Error if `timeout` ms elapse first, or if predicate throws.
 */
export function pollUntil(predicate, { interval = 50, timeout = 2000 } = {}) {
  // TODO
  throw new Error('TODO: implement pollUntil')
}
