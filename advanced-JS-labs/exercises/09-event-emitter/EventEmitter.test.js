import { describe, it, expect, vi } from 'vitest'
import { EventEmitter } from './EventEmitter.js'

describe('EventEmitter', () => {
  it('calls subscribers with emitted args', () => {
    const ee = new EventEmitter()
    const a = vi.fn()
    const b = vi.fn()
    ee.on('data', a)
    ee.on('data', b)
    ee.emit('data', 1, 2)
    expect(a).toHaveBeenCalledWith(1, 2)
    expect(b).toHaveBeenCalledWith(1, 2)
  })

  it('does not notify listeners of other events', () => {
    const ee = new EventEmitter()
    const cb = vi.fn()
    ee.on('a', cb)
    ee.emit('b')
    expect(cb).not.toHaveBeenCalled()
  })

  it('off stops future notifications', () => {
    const ee = new EventEmitter()
    const cb = vi.fn()
    ee.on('x', cb)
    ee.emit('x')
    ee.off('x', cb)
    ee.emit('x')
    expect(cb).toHaveBeenCalledTimes(1)
  })

  it('the return value of on() unsubscribes', () => {
    const ee = new EventEmitter()
    const cb = vi.fn()
    const unsubscribe = ee.on('x', cb)
    unsubscribe()
    ee.emit('x')
    expect(cb).not.toHaveBeenCalled()
  })

  it('once fires exactly one time', () => {
    const ee = new EventEmitter()
    const cb = vi.fn()
    ee.once('ping', cb)
    ee.emit('ping', 'a')
    ee.emit('ping', 'b')
    expect(cb).toHaveBeenCalledTimes(1)
    expect(cb).toHaveBeenCalledWith('a')
  })

  it('survives a listener unsubscribing during emit', () => {
    const ee = new EventEmitter()
    const calls = []
    const first = () => {
      calls.push('first')
      ee.off('e', second)
    }
    const second = () => calls.push('second')
    ee.on('e', first)
    ee.on('e', second)
    expect(() => ee.emit('e')).not.toThrow()
  })
})
