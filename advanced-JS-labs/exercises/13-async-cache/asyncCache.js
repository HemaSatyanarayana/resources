/**
 * memoizeAsync(fn, { ttl = Infinity })  — the integration exercise.
 * Combines all three pillars at once:
 *   - CLOSURES: two private Maps held in the closure (a value cache + an
 *     in-flight registry).
 *   - PROMISES: `fn` is async; concurrent callers share the SAME pending promise.
 *   - TIMERS:   a cached entry is evicted `ttl` ms after it resolves (setTimeout).
 *
 * Returns a wrapped async function. Behavior:
 *   - Key each call by its arguments (JSON is fine).
 *   - If a resolved value is cached for that key, return it (wrapped in a promise).
 *   - If a call with that key is already IN FLIGHT, return the exact same promise
 *     (de-duplication — don't call `fn` twice for overlapping calls).
 *   - Otherwise call `fn`, cache the resolved value, and schedule its eviction
 *     after `ttl` ms (skip scheduling when ttl is Infinity).
 *   - Do NOT cache rejections — a failed call should be retryable.
 */
export function memoizeAsync(fn, { ttl = Infinity } = {}) {
  const cache = new Map() // key -> resolved value
  const inFlight = new Map() // key -> pending promise

  return function (...args) {
    const key = JSON.stringify(args)

    // TODO:
    //   1. if `cache` has `key`, return Promise.resolve(cachedValue)
    //   2. if `inFlight` has `key`, return the existing promise
    //   3. otherwise build the promise chain:
    //        Promise.resolve()
    //          .then(() => fn.apply(this, args))
    //          .then((value) => { cache.set(key, value);
    //                             if (ttl !== Infinity) setTimeout(() => cache.delete(key), ttl);
    //                             return value })
    //          .finally(() => inFlight.delete(key))
    //      store it in `inFlight` under `key`, and return it.
    throw new Error('TODO: implement memoizeAsync')
  }
}
