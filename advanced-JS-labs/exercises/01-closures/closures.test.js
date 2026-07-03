import { describe, it, expect, vi } from 'vitest'
import { createCounter, once, memoize } from './closures.js'

describe('createCounter', () => {
  it('keeps private state across calls', () => {
    const c = createCounter(10)
    expect(c.value()).toBe(10)
    expect(c.increment()).toBe(11)
    expect(c.increment()).toBe(12)
    expect(c.decrement()).toBe(11)
    expect(c.value()).toBe(11)
  })

  it('resets to the starting value', () => {
    const c = createCounter(5)
    c.increment()
    c.increment()
    c.reset()
    expect(c.value()).toBe(5)
  })

  it('does not leak the counter as a property', () => {
    const c = createCounter()
    expect(c.count).toBeUndefined()
  })

  it('gives each counter independent state', () => {
    const a = createCounter()
    const b = createCounter()
    a.increment()
    expect(a.value()).toBe(1)
    expect(b.value()).toBe(0)
  })
})

describe('once', () => {
  it('invokes the function only the first time', () => {
    const fn = vi.fn(() => 42)
    const wrapped = once(fn)
    expect(wrapped()).toBe(42)
    expect(wrapped()).toBe(42)
    expect(wrapped()).toBe(42)
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('preserves the receiver', () => {
    const obj = {
      value: 7,
      getValue: once(function () {
        return this.value
      }),
    }
    expect(obj.getValue()).toBe(7)
  })
})

describe('memoize', () => {
  it('caches by arguments', () => {
    const add = vi.fn((a, b) => a + b)
    const memo = memoize(add)
    expect(memo(1, 2)).toBe(3)
    expect(memo(1, 2)).toBe(3)
    expect(add).toHaveBeenCalledTimes(1)
    expect(memo(2, 2)).toBe(4)
    expect(add).toHaveBeenCalledTimes(2)
  })
})
