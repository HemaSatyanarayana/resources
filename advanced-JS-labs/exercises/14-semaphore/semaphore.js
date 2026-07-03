/**
 * createSemaphore(max)  — combines PROMISES + CLOSURES.
 *
 * A semaphore lets at most `max` holders proceed at once; extra acquirers WAIT
 * (as promises) until a permit is released. A semaphore of size 1 is a MUTEX.
 *
 * Returns { acquire, release, use }:
 *   - acquire()  -> Promise that resolves when a permit is granted
 *   - release()  -> give a permit back (waking the next waiter, if any)
 *   - use(fn)    -> acquire, run `fn` (await it), release in a finally
 *
 * State lives in the closure: `available` permits + a `queue` of waiting
 * resolvers (FIFO).
 */
export function createSemaphore(max) {
  let available = max
  const queue = [] // resolvers waiting for a permit

  function acquire() {
    // TODO:
    //   - if available > 0: available--, return Promise.resolve()
    //   - else: return a new Promise and push its `resolve` onto `queue`
    throw new Error('TODO: implement acquire')
  }

  function release() {
    // TODO:
    //   - if there's a waiter in `queue`: shift it and call resolve() (the permit
    //     passes straight to them; `available` stays the same)
    //   - else: available = Math.min(max, available + 1)
    throw new Error('TODO: implement release')
  }

  async function use(fn) {
    // TODO: await acquire(); try { return await fn() } finally { release() }
    throw new Error('TODO: implement use')
  }

  return { acquire, release, use }
}
