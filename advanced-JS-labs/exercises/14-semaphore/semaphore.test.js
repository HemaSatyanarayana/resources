import { describe, it, expect } from 'vitest'
import { createSemaphore } from './semaphore.js'

const delay = (ms) => new Promise((r) => setTimeout(r, ms))

describe('createSemaphore', () => {
  it('lets up to `max` holders run at once', async () => {
    const sem = createSemaphore(2)
    let active = 0
    let maxActive = 0

    const task = async () => {
      await sem.acquire()
      active++
      maxActive = Math.max(maxActive, active)
      await delay(10)
      active--
      sem.release()
    }

    await Promise.all(Array.from({ length: 6 }, task))
    expect(maxActive).toBe(2)
    expect(active).toBe(0)
  })

  it('acquire resolves immediately while permits are free, then queues', async () => {
    const sem = createSemaphore(1)
    await sem.acquire() // takes the only permit

    let secondAcquired = false
    sem.acquire().then(() => {
      secondAcquired = true
    })
    await Promise.resolve()
    expect(secondAcquired).toBe(false) // still blocked

    sem.release() // hand the permit to the waiter
    await Promise.resolve()
    await Promise.resolve()
    expect(secondAcquired).toBe(true)
  })

  it('use() runs the fn and always releases, even on throw', async () => {
    const sem = createSemaphore(1)
    await expect(sem.use(async () => Promise.reject(new Error('fail')))).rejects.toThrow(
      'fail'
    )
    // If the permit leaked, this acquire would hang forever.
    await expect(
      Promise.race([sem.use(async () => 'ok'), delay(50).then(() => 'timeout')])
    ).resolves.toBe('ok')
  })

  it('a mutex is just createSemaphore(1) — serializes access', async () => {
    const mutex = createSemaphore(1)
    const order = []
    const critical = (label, ms) =>
      mutex.use(async () => {
        order.push(`${label}-start`)
        await delay(ms)
        order.push(`${label}-end`)
      })

    await Promise.all([critical('a', 20), critical('b', 5)])
    // b must not start until a finished.
    expect(order).toEqual(['a-start', 'a-end', 'b-start', 'b-end'])
  })
})
