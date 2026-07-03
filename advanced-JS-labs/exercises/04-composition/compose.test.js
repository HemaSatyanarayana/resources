import { describe, it, expect } from 'vitest'
import { identity, pipe, compose } from './compose.js'

const inc = (x) => x + 1
const double = (x) => x * 2
const negate = (x) => -x

describe('identity', () => {
  it('returns its argument', () => {
    expect(identity(5)).toBe(5)
  })
})

describe('pipe', () => {
  it('applies left to right', () => {
    // inc then double: (3 + 1) * 2 = 8
    expect(pipe(inc, double)(3)).toBe(8)
  })

  it('with no functions acts like identity', () => {
    expect(pipe()(42)).toBe(42)
  })

  it('threads through many functions', () => {
    expect(pipe(inc, double, negate)(3)).toBe(-8)
  })
})

describe('compose', () => {
  it('applies right to left', () => {
    // double then inc: (3 * 2) + 1 = 7
    expect(compose(inc, double)(3)).toBe(7)
  })

  it('is the mirror of pipe', () => {
    const value = 3
    expect(compose(inc, double, negate)(value)).toBe(pipe(negate, double, inc)(value))
  })
})
