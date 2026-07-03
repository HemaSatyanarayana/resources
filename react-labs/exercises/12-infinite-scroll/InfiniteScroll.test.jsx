import { render, screen, act, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import InfiniteScroll from './InfiniteScroll'

// jsdom has no IntersectionObserver — mock a controllable one.
const observers = new Set()
class MockIntersectionObserver {
  constructor(cb) {
    this.cb = cb
    observers.add(this)
  }
  observe() {}
  unobserve() {}
  disconnect() {
    observers.delete(this)
  }
}

async function scrollIntoView() {
  await act(async () => {
    for (const o of observers) o.cb([{ isIntersecting: true }])
  })
}

beforeEach(() => {
  observers.clear()
  vi.stubGlobal('IntersectionObserver', MockIntersectionObserver)
})

// 2 pages of 3 items, then empty.
const makeFetch = () =>
  vi.fn(async (page) => {
    if (page >= 2) return []
    return [0, 1, 2].map((i) => `Item ${page * 3 + i + 1}`)
  })

describe('InfiniteScroll', () => {
  it('starts empty until the sentinel is seen', () => {
    render(<InfiniteScroll fetchPage={makeFetch()} />)
    expect(screen.queryAllByRole('listitem')).toHaveLength(0)
  })

  it('loads the first page when the sentinel intersects', async () => {
    render(<InfiniteScroll fetchPage={makeFetch()} />)
    await scrollIntoView()
    await waitFor(() => expect(screen.getAllByRole('listitem')).toHaveLength(3))
    expect(screen.getByText('Item 1')).toBeInTheDocument()
  })

  it('appends subsequent pages', async () => {
    render(<InfiniteScroll fetchPage={makeFetch()} />)
    await scrollIntoView()
    await waitFor(() => expect(screen.getAllByRole('listitem')).toHaveLength(3))
    await scrollIntoView()
    await waitFor(() => expect(screen.getAllByRole('listitem')).toHaveLength(6))
    expect(screen.getByText('Item 6')).toBeInTheDocument()
  })

  it('stops and shows a done message, without extra fetches', async () => {
    const fetchPage = makeFetch()
    render(<InfiniteScroll fetchPage={fetchPage} />)
    await scrollIntoView() // page 0
    await scrollIntoView() // page 1
    await scrollIntoView() // page 2 -> empty -> done
    await waitFor(() => expect(screen.getByText('No more items')).toBeInTheDocument())

    const callsSoFar = fetchPage.mock.calls.length
    await scrollIntoView()
    expect(fetchPage.mock.calls.length).toBe(callsSoFar)
  })
})
