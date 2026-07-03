import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createRateLimiter } from './rateLimiter.js'

describe('createRateLimiter (token bucket)', () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  it('allows a burst up to capacity, then blocks', () => {
    const rl = createRateLimiter(2, 1000)
    expect(rl.tryAcquire()).toBe(true) // 2 -> 1
    expect(rl.tryAcquire()).toBe(true) // 1 -> 0
    expect(rl.tryAcquire()).toBe(false) // empty
    rl.stop()
  })

  it('refills one token per interval', () => {
    const rl = createRateLimiter(2, 1000)
    rl.tryAcquire()
    rl.tryAcquire()
    expect(rl.tryAcquire()).toBe(false)

    vi.advanceTimersByTime(1000) // +1 token
    expect(rl.tryAcquire()).toBe(true)
    expect(rl.tryAcquire()).toBe(false)
    rl.stop()
  })

  it('never refills beyond capacity', () => {
    const rl = createRateLimiter(2, 1000)
    rl.tryAcquire()
    rl.tryAcquire() // drained

    vi.advanceTimersByTime(5000) // 5 refills, but caps at 2
    expect(rl.availableTokens()).toBe(2)
    expect(rl.tryAcquire()).toBe(true)
    expect(rl.tryAcquire()).toBe(true)
    expect(rl.tryAcquire()).toBe(false)
    rl.stop()
  })

  it('stop() halts refilling', () => {
    const rl = createRateLimiter(1, 1000)
    rl.tryAcquire()
    rl.stop()
    vi.advanceTimersByTime(5000)
    expect(rl.tryAcquire()).toBe(false) // no refill after stop
  })
})
