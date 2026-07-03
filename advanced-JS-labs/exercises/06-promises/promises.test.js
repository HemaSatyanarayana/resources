import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { promisify, delay, series } from './promises.js'

describe('promisify', () => {
  it('resolves with the callback result', async () => {
    const cbStyle = (a, b, cb) => cb(null, a + b)
    const added = promisify(cbStyle)
    await expect(added(2, 3)).resolves.toBe(5)
  })

  it('rejects with the callback error', async () => {
    const failing = (cb) => cb(new Error('nope'))
    const p = promisify(failing)
    await expect(p()).rejects.toThrow('nope')
  })
})

describe('delay (fake timers)', () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  it('resolves with the value after the delay', async () => {
    const spy = vi.fn()
    delay(50, 'hi').then(spy)
    await vi.advanceTimersByTimeAsync(49)
    expect(spy).not.toHaveBeenCalled()
    await vi.advanceTimersByTimeAsync(1)
    expect(spy).toHaveBeenCalledWith('hi')
  })
})

describe('series', () => {
  it('runs tasks sequentially and collects results in order', async () => {
    const order = []
    const task = (label, ms) => () =>
      new Promise((resolve) =>
        setTimeout(() => {
          order.push(label)
          resolve(label)
        }, ms)
      )

    // Even though "a" is slowest, series must wait for it before starting "b".
    const results = await series([task('a', 20), task('b', 5), task('c', 1)])
    expect(results).toEqual(['a', 'b', 'c'])
    expect(order).toEqual(['a', 'b', 'c'])
  })
})
