/**
 * identity(x) -> x
 */
export const identity = (x) => x

/**
 * pipe(...fns)
 * Returns a function that runs the fns LEFT to RIGHT, threading the result:
 *   pipe(f, g, h)(x) === h(g(f(x)))
 * pipe() with no fns behaves like identity.
 */
export function pipe(...fns) {
  // TODO: reduce over fns, applying each to the accumulator.
  throw new Error('TODO: implement pipe')
}

/**
 * compose(...fns)
 * Like pipe but RIGHT to LEFT (the mathematical order):
 *   compose(f, g, h)(x) === f(g(h(x)))
 */
export function compose(...fns) {
  // TODO: reduceRight over fns.
  throw new Error('TODO: implement compose')
}
