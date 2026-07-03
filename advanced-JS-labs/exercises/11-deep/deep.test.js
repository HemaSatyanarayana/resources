import { describe, it, expect } from 'vitest'
import { deepClone, deepEqual } from './deep.js'

describe('deepClone', () => {
  it('returns primitives unchanged', () => {
    expect(deepClone(5)).toBe(5)
    expect(deepClone('hi')).toBe('hi')
    expect(deepClone(null)).toBe(null)
  })

  it('deep-copies nested structures', () => {
    const src = { a: 1, b: { c: [1, 2, { d: 3 }] } }
    const clone = deepClone(src)
    expect(clone).toEqual(src)
    expect(clone).not.toBe(src)
    expect(clone.b).not.toBe(src.b)
    expect(clone.b.c).not.toBe(src.b.c)

    clone.b.c[2].d = 999
    expect(src.b.c[2].d).toBe(3) // original untouched
  })

  it('reconstructs Date and RegExp', () => {
    const src = { when: new Date(0), re: /abc/gi }
    const clone = deepClone(src)
    expect(clone.when).toEqual(src.when)
    expect(clone.when).not.toBe(src.when)
    expect(clone.re.source).toBe('abc')
    expect(clone.re.flags).toBe('gi')
  })

  it('handles circular references', () => {
    const src = { name: 'loop' }
    src.self = src
    const clone = deepClone(src)
    expect(clone.name).toBe('loop')
    expect(clone.self).toBe(clone) // cycle preserved, not the original
    expect(clone.self).not.toBe(src)
  })
})

describe('deepEqual', () => {
  it('compares nested structures', () => {
    expect(deepEqual({ a: [1, 2], b: { c: 3 } }, { a: [1, 2], b: { c: 3 } })).toBe(true)
    expect(deepEqual({ a: 1 }, { a: 2 })).toBe(false)
    expect(deepEqual({ a: 1 }, { a: 1, b: 2 })).toBe(false)
  })

  it('uses Object.is semantics for primitives', () => {
    expect(deepEqual(NaN, NaN)).toBe(true)
    expect(deepEqual(0, -0)).toBe(false)
  })

  it('distinguishes arrays from objects', () => {
    expect(deepEqual([1, 2], { 0: 1, 1: 2 })).toBe(false)
  })
})
