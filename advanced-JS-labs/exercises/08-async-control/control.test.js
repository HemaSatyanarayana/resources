import { describe, it, expect, vi } from 'vitest'
import { retry, withTimeout, mapLimit } from './control.js'

const resolveIn = (ms, v) => new Promise((r) => setTimeout(() => r(v), ms))
const rejectIn = (ms, e) => new Promise((_, r) => setTimeout(() => r(e), ms))

describe('retry', () => {
  it('resolves once fn succeeds', async () => {
    let attempts = 0
    const fn = vi.fn(async () => {
      attempts += 1
      if (attempts < 3) throw new Error('fail')
      return 'ok'
    })
    await expect(retry(fn, { retries: 5 })).resolves.toBe('ok')
    expect(fn).toHaveBeenCalledTimes(3)
  })

  it('rejects with the last error after exhausting retries', async () => {
    const fn = vi.fn(async () => {
      throw new Error('always')
    })
    await expect(retry(fn, { retries: 2 })).rejects.toThrow('always')
    expect(fn).toHaveBeenCalledTimes(3) // 1 initial + 2 retries
  })
})

describe('withTimeout', () => {
  it('passes through a fast result', async () => {
    await expect(withTimeout(resolveIn(5, 'done'), 50)).resolves.toBe('done')
  })

  it('rejects when the promise is too slow', async () => {
    await expect(withTimeout(resolveIn(50, 'late'), 10)).rejects.toThrow('Timed out')
  })

  it('propagates the underlying rejection', async () => {
    await expect(withTimeout(rejectIn(5, new Error('inner')), 50)).rejects.toThrow('inner')
  })
})

describe('mapLimit', () => {
  it('preserves order and respects the concurrency limit', async () => {
    let active = 0
    let maxActive = 0
    const iteratee = (x) =>
      new Promise((resolve) => {
        active += 1
        maxActive = Math.max(maxActive, active)
        setTimeout(() => {
          active -= 1
          resolve(x * 2)
        }, 10)
      })

    const result = await mapLimit([1, 2, 3, 4, 5], 2, iteratee)
    expect(result).toEqual([2, 4, 6, 8, 10])
    expect(maxActive).toBeLessThanOrEqual(2)
    expect(maxActive).toBeGreaterThan(1)
  })

  it('resolves to [] for empty input', async () => {
    await expect(mapLimit([], 3, async (x) => x)).resolves.toEqual([])
  })
})
