import { describe, it, expect } from 'vitest'
import { all, allSettled, race, any } from './combinators.js'

const resolveIn = (ms, v) => new Promise((r) => setTimeout(() => r(v), ms))
const rejectIn = (ms, e) => new Promise((_, r) => setTimeout(() => r(e), ms))

describe('all', () => {
  it('resolves in input order regardless of timing', async () => {
    const result = await all([resolveIn(20, 'a'), resolveIn(5, 'b'), 'c'])
    expect(result).toEqual(['a', 'b', 'c'])
  })

  it('resolves to [] for an empty array', async () => {
    await expect(all([])).resolves.toEqual([])
  })

  it('rejects if any input rejects', async () => {
    await expect(all([resolveIn(5, 'a'), rejectIn(1, new Error('boom'))])).rejects.toThrow(
      'boom'
    )
  })
})

describe('allSettled', () => {
  it('reports every outcome without rejecting', async () => {
    const result = await allSettled([Promise.resolve(1), Promise.reject('x')])
    expect(result).toEqual([
      { status: 'fulfilled', value: 1 },
      { status: 'rejected', reason: 'x' },
    ])
  })
})

describe('race', () => {
  it('settles with the first result', async () => {
    await expect(race([resolveIn(20, 'slow'), resolveIn(1, 'fast')])).resolves.toBe('fast')
  })

  it('rejects if the first settled input rejects', async () => {
    await expect(race([rejectIn(1, new Error('first')), resolveIn(20, 'ok')])).rejects.toThrow(
      'first'
    )
  })
})

describe('any', () => {
  it('resolves with the first fulfilled value', async () => {
    await expect(any([rejectIn(1, 'e'), resolveIn(10, 'winner')])).resolves.toBe('winner')
  })

  it('rejects with an AggregateError when all reject', async () => {
    await expect(any([Promise.reject('a'), Promise.reject('b')])).rejects.toMatchObject({
      errors: ['a', 'b'],
    })
  })
})
