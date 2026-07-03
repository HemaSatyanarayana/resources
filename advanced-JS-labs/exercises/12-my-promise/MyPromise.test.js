import { describe, it, expect } from 'vitest'
import { MyPromise } from './MyPromise.js'

describe('MyPromise', () => {
  it('resolves with a value (awaitable)', async () => {
    const value = await new MyPromise((resolve) => resolve(42))
    expect(value).toBe(42)
  })

  it('resolves asynchronously', async () => {
    const value = await new MyPromise((resolve) => setTimeout(() => resolve('later'), 1))
    expect(value).toBe('later')
  })

  it('chains .then transformations', async () => {
    const value = await new MyPromise((resolve) => resolve(1))
      .then((x) => x + 1)
      .then((x) => x * 2)
    expect(value).toBe(4)
  })

  it('rejects and is caught by .catch', async () => {
    let caught
    await new MyPromise((_, reject) => reject(new Error('bad'))).catch((e) => {
      caught = e.message
    })
    expect(caught).toBe('bad')
  })

  it('rejects when the executor throws', async () => {
    let caught
    await new MyPromise(() => {
      throw new Error('exec')
    }).catch((e) => {
      caught = e.message
    })
    expect(caught).toBe('exec')
  })

  it('flattens a thenable returned from .then', async () => {
    const value = await new MyPromise((resolve) => resolve(1)).then(
      (x) => new MyPromise((r) => r(x + 10))
    )
    expect(value).toBe(11)
  })

  it('passes rejections through handlers that omit onRejected', async () => {
    let caught
    await new MyPromise((_, reject) => reject('boom'))
      .then((x) => x * 2) // no onRejected -> pass through
      .catch((e) => {
        caught = e
      })
    expect(caught).toBe('boom')
  })

  it('supports static resolve and reject', async () => {
    expect(await MyPromise.resolve(5)).toBe(5)
    let reason
    await MyPromise.reject('x').catch((e) => {
      reason = e
    })
    expect(reason).toBe('x')
  })

  it('runs handlers as microtasks (after synchronous code)', async () => {
    const order = []
    order.push('sync-start')
    const p = new MyPromise((resolve) => resolve()).then(() => order.push('microtask'))
    order.push('sync-end')
    await p
    expect(order).toEqual(['sync-start', 'sync-end', 'microtask'])
  })
})
