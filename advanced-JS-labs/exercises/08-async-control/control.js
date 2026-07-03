/**
 * retry(fn, { retries = 3, delay = 0 })
 * Call `fn` (which returns a Promise). If it rejects, try again — up to
 * `retries` additional times — waiting `delay` ms between attempts. Resolve
 * with the first success; if every attempt fails, reject with the LAST error.
 * (So total attempts = retries + 1.)
 */
export async function retry(fn, { retries = 3, delay = 0 } = {}) {
  // TODO: loop attempt = 0..retries; try/await fn(); on error remember it and,
  // if attempts remain, await a delay before retrying.
  throw new Error('TODO: implement retry')
}

/**
 * withTimeout(promise, ms, message = 'Timed out')
 * Return a Promise that settles like `promise`, but rejects with an Error(message)
 * if `promise` hasn't settled within `ms`. Clean up the timer either way.
 */
export function withTimeout(promise, ms, message = 'Timed out') {
  // TODO: Promise.race the input against a timeout promise; clearTimeout in finally.
  throw new Error('TODO: implement withTimeout')
}

/**
 * mapLimit(items, limit, iteratee)
 * Like items.map(iteratee) but runs at most `limit` iteratee calls concurrently.
 * `iteratee(item, index)` returns a Promise. Resolve to the results array IN
 * INPUT ORDER. Reject as soon as any iteratee rejects.
 */
export function mapLimit(items, limit, iteratee) {
  // TODO: keep a running pointer + in-flight counter; launch up to `limit` at a
  // time, and launch the next when one finishes.
  throw new Error('TODO: implement mapLimit')
}
