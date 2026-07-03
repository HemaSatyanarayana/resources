/**
 * debounce(fn, wait)
 * Returns a debounced function: it postpones calling `fn` until `wait` ms have
 * passed since the LAST call. Rapid-fire calls collapse into a single trailing
 * call, using the most recent arguments.
 * Also expose `.cancel()` to drop any pending call.
 *
 * Use: search-as-you-type, resize/scroll handlers, autosave.
 */
export function debounce(fn, wait) {
  let timer
  function debounced(...args) {
    // TODO: clear the pending timer and schedule a fresh one that calls fn(...args).
    throw new Error('TODO: implement debounce')
  }
  debounced.cancel = () => {
    // TODO: clear the pending timer.
  }
  return debounced
}

/**
 * throttle(fn, wait)
 * Returns a throttled function that invokes `fn` at most once per `wait` ms.
 * LEADING edge: the first call fires immediately. Calls during the cooldown are
 * coalesced into a single TRAILING call (with the latest args) at the end of the
 * window.
 * Also expose `.cancel()`.
 *
 * Use: rate-limiting scroll/mousemove, firing analytics.
 */
export function throttle(fn, wait) {
  let lastInvoke = null // timestamp of the last real invocation; null = never
  let timer = null
  let savedArgs = null
  let savedThis = null

  function throttled(...args) {
    // TODO:
    //   - const now = Date.now(); save args + this
    //   - remaining = lastInvoke === null ? 0 : wait - (now - lastInvoke)
    //   - if remaining <= 0: clear any timer, invoke now, set lastInvoke = now
    //   - else if no timer pending: schedule a trailing invoke after `remaining`
    throw new Error('TODO: implement throttle')
  }
  throttled.cancel = () => {
    // TODO: clear the timer and reset lastInvoke.
  }
  return throttled
}
