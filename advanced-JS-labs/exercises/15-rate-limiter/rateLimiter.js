/**
 * createRateLimiter(capacity, refillIntervalMs)  — combines TIMERS + CLOSURES.
 *
 * The token-bucket algorithm. A bucket holds up to `capacity` tokens and refills
 * one token every `refillIntervalMs` (via a timer). Each request spends a token.
 * This permits bursts (up to `capacity`) while capping the long-run rate — the
 * key difference from throttle.
 *
 * Returns { tryAcquire, availableTokens, stop }:
 *   - tryAcquire()     -> true and consume a token if any remain, else false
 *   - availableTokens()-> current token count (for tests/introspection)
 *   - stop()           -> clear the refill timer
 *
 * State (token count + timer id) lives in the closure.
 */
export function createRateLimiter(capacity, refillIntervalMs) {
  let tokens = capacity

  // TODO: start a setInterval that refills one token every refillIntervalMs,
  //       clamped so it never exceeds `capacity`. Keep its id for stop().
  const id = null

  function tryAcquire() {
    // TODO: if tokens > 0, tokens-- and return true; else return false.
    throw new Error('TODO: implement tryAcquire')
  }

  function availableTokens() {
    return tokens
  }

  function stop() {
    // TODO: clearInterval(id)
    throw new Error('TODO: implement stop')
  }

  return { tryAcquire, availableTokens, stop }
}
