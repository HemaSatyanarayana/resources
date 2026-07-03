import { describe, it, expect } from 'vitest'
import { myCall, myApply, myBind } from './binding.js'

function greet(greeting, punctuation) {
  return `${greeting}, ${this.name}${punctuation}`
}

describe('myCall', () => {
  it('sets `this` and passes args', () => {
    expect(myCall(greet, { name: 'Ada' }, 'Hi', '!')).toBe('Hi, Ada!')
  })

  it('borrows a method for another object', () => {
    const obj = { name: 'Bob' }
    function whoAmI() {
      return this.name
    }
    expect(myCall(whoAmI, obj)).toBe('Bob')
  })

  it('defaults to globalThis when thisArg is null', () => {
    globalThis.name = undefined
    function readGlobal() {
      return this === globalThis
    }
    expect(myCall(readGlobal, null)).toBe(true)
  })
})

describe('myApply', () => {
  it('accepts args as an array', () => {
    expect(myApply(greet, { name: 'Cy' }, ['Hey', '?'])).toBe('Hey, Cy?')
  })

  it('works with no args array', () => {
    function count() {
      return arguments.length
    }
    expect(myApply(count, {})).toBe(0)
  })
})

describe('myBind', () => {
  it('binds `this` permanently', () => {
    const bound = myBind(greet, { name: 'Di' }, 'Yo')
    expect(bound('!')).toBe('Yo, Di!')
  })

  it('supports partial application', () => {
    function add(a, b, c) {
      return a + b + c
    }
    const add5 = myBind(add, null, 5)
    expect(add5(10, 20)).toBe(35)
  })
})
