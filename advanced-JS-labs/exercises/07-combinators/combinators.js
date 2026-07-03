/**
 * Polyfill the four promise combinators. Each takes an array of values or
 * promises (wrap every input with Promise.resolve so plain values work too).
 */

/**
 * all(items) — resolves to an array of results IN INPUT ORDER once every input
 * has resolved. Rejects immediately if ANY input rejects.
 * Resolves to [] for an empty array.
 */
export function all(items) {
  // TODO
  throw new Error('TODO: implement all')
}

/**
 * allSettled(items) — never rejects. Resolves once every input settles, to an
 * array of { status: 'fulfilled', value } | { status: 'rejected', reason }.
 */
export function allSettled(items) {
  // TODO
  throw new Error('TODO: implement allSettled')
}

/**
 * race(items) — settles as soon as the FIRST input settles, with that same
 * value or reason (fulfilled OR rejected).
 */
export function race(items) {
  // TODO
  throw new Error('TODO: implement race')
}

/**
 * any(items) — resolves with the first FULFILLED value. If ALL inputs reject,
 * rejects with an AggregateError whose `.errors` holds the reasons in order.
 */
export function any(items) {
  // TODO
  throw new Error('TODO: implement any')
}
