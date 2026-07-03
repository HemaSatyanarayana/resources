/**
 * deepClone(value)
 * Recursively copy `value` so the result shares NO references with the original
 * (mutating the clone never affects the source). Handle:
 *   - primitives (return as-is)
 *   - arrays and plain objects (recurse)
 *   - Date and RegExp (reconstruct)
 *   - circular references (an object that contains itself) — use a WeakMap of
 *     already-cloned sources so you don't recurse forever.
 */
export function deepClone(value, seen = new WeakMap()) {
  // TODO
  throw new Error('TODO: implement deepClone')
}

/**
 * deepEqual(a, b)
 * Structural equality. True when a and b have the same shape and equal leaves.
 *   - primitives compare with Object.is (so NaN === NaN, and +0 !== -0)
 *   - arrays/objects: same keys, recursively equal values
 *   - an array and a non-array are never equal
 */
export function deepEqual(a, b) {
  // TODO
  throw new Error('TODO: implement deepEqual')
}
