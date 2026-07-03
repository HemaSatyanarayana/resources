import { describe, it, expect } from 'vitest'
import { curry, partial, add } from './curry.js'

const add3 = (a, b, c) => a + b + c

describe('curry', () => {
  it('supports every argument grouping', () => {
    const c = curry(add3)
    expect(c(1)(2)(3)).toBe(6)
    expect(c(1, 2)(3)).toBe(6)
    expect(c(1)(2, 3)).toBe(6)
    expect(c(1, 2, 3)).toBe(6)
  })

  it('is reusable from a partially applied stage', () => {
    const c = curry(add3)
    const addTo10 = c(4)(6) // needs one more
    expect(addTo10(1)).toBe(11)
    expect(addTo10(2)).toBe(12)
  })
})

describe('partial', () => {
  it('pre-fills leading arguments', () => {
    expect(partial(add3, 1, 2)(3)).toBe(6)
    expect(partial(add3, 1)(2, 3)).toBe(6)
  })
})

describe('add (infinite currying)', () => {
  it('accumulates and coerces to a number', () => {
    expect(Number(add(1)(2)(3))).toBe(6)
    expect(+add(5)(10)).toBe(15)
    expect(Number(add(1)(2)(3)(4)(5))).toBe(15)
  })
})
