import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { memoizeAsync } from './asyncCache.js'

describe('memoizeAsync — in-flight de-duplication (closures + promises)', () => {
  it('shares one promise for concurrent calls with the same args', async () => {
    const fn = vi.fn(async (x) => x + 1)
    const m = memoizeAsync(fn)

    const p1 = m(5)
    const p2 = m(5) // same args, still in flight
    expect(p1).toBe(p2) // literally the same promise

    const [a, b] = await Promise.all([p1, p2])
    expect(a).toBe(6)
    expect(b).toBe(6)
    expect(fn).toHaveBeenCalledTimes(1) // computed once
  })

  it('caches a resolved value for later calls', async () => {
    const fn = vi.fn(async (x) => x * 10)
    const m = memoizeAsync(fn)

    expect(await m(2)).toBe(20)
    expect(await m(2)).toBe(20)
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('does NOT cache rejections (so a later call can retry)', async () => {
    let fail = true
    const fn = vi.fn(async () => {
      if (fail) throw new Error('boom')
      return 'ok'
    })
    const m = memoizeAsync(fn)

    await expect(m()).rejects.toThrow('boom')
    fail = false
    await expect(m()).resolves.toBe('ok')
    expect(fn).toHaveBeenCalledTimes(2)
  })
})

describe('memoizeAsync — TTL eviction (timers)', () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  it('recomputes after the entry expires', async () => {
    const fn = vi.fn(async (x) => x * 2)
    const m = memoizeAsync(fn, { ttl: 1000 })

    expect(await m(3)).toBe(6)
    expect(await m(3)).toBe(6) // still cached
    expect(fn).toHaveBeenCalledTimes(1)

    await vi.advanceTimersByTimeAsync(1000) // TTL elapses -> eviction

    expect(await m(3)).toBe(6)
    expect(fn).toHaveBeenCalledTimes(2) // recomputed
  })

  it('keeps the entry cached before the TTL elapses', async () => {
    const fn = vi.fn(async (x) => x * 2)
    const m = memoizeAsync(fn, { ttl: 1000 })

    expect(await m(3)).toBe(6)
    await vi.advanceTimersByTimeAsync(999)
    expect(await m(3)).toBe(6)
    expect(fn).toHaveBeenCalledTimes(1)
  })
})
