/**
 * Reimplement call / apply / bind from scratch. These take the function as the
 * first argument (instead of being methods on Function.prototype) so they're
 * easy to test in isolation — but the mechanics are identical.
 */

/**
 * myCall(fn, thisArg, ...args)
 * Invoke `fn` with `this` set to `thisArg` and the given args. When `thisArg`
 * is null/undefined, default to globalThis. Return fn's result.
 *
 * Trick: temporarily attach `fn` as a property of the context object, call it
 * as a method (so `this` binds to the object), then remove it.
 */
export function myCall(fn, thisArg, ...args) {
  // TODO
  throw new Error('TODO: implement myCall')
}

/**
 * myApply(fn, thisArg, argsArray)
 * Like myCall but args come as an array (or undefined).
 */
export function myApply(fn, thisArg, argsArray) {
  // TODO
  throw new Error('TODO: implement myApply')
}

/**
 * myBind(fn, thisArg, ...bound)
 * Return a new function that, when called, invokes `fn` with `this` = thisArg
 * and the bound args followed by any call-time args.
 */
export function myBind(fn, thisArg, ...bound) {
  // TODO
  throw new Error('TODO: implement myBind')
}
