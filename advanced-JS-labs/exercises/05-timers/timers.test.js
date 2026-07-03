import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { sleep, repeat, pollUntil } from './timers.js'

describe('sleep (fake timers)', () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  it('resolves after the delay', async () => {
    let done = false
    sleep(100).then(() => {
      done = true
    })
    expect(done).toBe(false)
    await vi.advanceTimersByTimeAsync(100)
    expect(done).toBe(true)
  })
})

describe('repeat (fake timers)', () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  it('calls fn on each interval and stops when cancelled', async () => {
    const fn = vi.fn()
    const cancel = repeat(fn, 100)
    await vi.advanceTimersByTimeAsync(350)
    expect(fn).toHaveBeenCalledTimes(3)

    cancel()
    await vi.advanceTimersByTimeAsync(500)
    expect(fn).toHaveBeenCalledTimes(3)
  })
})

describe('pollUntil (real timers)', () => {
  it('resolves with the predicate value once it is truthy', async () => {
    let n = 0
    const result = await pollUntil(
      () => {
        n += 1
        return n >= 3 ? 'ready' : false
      },
      { interval: 5, timeout: 500 }
    )
    expect(result).toBe('ready')
  })

  it('supports async predicates', async () => {
    let n = 0
    const result = await pollUntil(async () => (++n >= 2 ? n : 0), {
      interval: 5,
      timeout: 500,
    })
    expect(result).toBe(2)
  })

  it('rejects when it times out', async () => {
    await expect(
      pollUntil(() => false, { interval: 5, timeout: 30 })
    ).rejects.toThrow(/timed out/i)
  })
})
