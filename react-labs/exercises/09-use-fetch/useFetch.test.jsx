import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, afterEach } from 'vitest'
import { useFetch } from './useFetch'

afterEach(() => vi.restoreAllMocks())

describe('useFetch', () => {
  it('starts in a loading state', () => {
    vi.stubGlobal('fetch', vi.fn(() => new Promise(() => {}))) // never resolves
    const { result } = renderHook(() => useFetch('/api/thing'))
    expect(result.current.loading).toBe(true)
    expect(result.current.data).toBeNull()
    expect(result.current.error).toBeNull()
  })

  it('resolves with parsed data', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(async () => ({ ok: true, json: async () => ({ id: 1, name: 'Ada' }) }))
    )
    const { result } = renderHook(() => useFetch('/api/user'))
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(result.current.data).toEqual({ id: 1, name: 'Ada' })
    expect(result.current.error).toBeNull()
  })

  it('reports an error on a non-ok response', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => ({ ok: false, status: 500 })))
    const { result } = renderHook(() => useFetch('/api/boom'))
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(result.current.error).toBeInstanceOf(Error)
    expect(result.current.data).toBeNull()
  })

  it('reports network rejections', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => { throw new Error('offline') }))
    const { result } = renderHook(() => useFetch('/api/x'))
    await waitFor(() => expect(result.current.error).toBeInstanceOf(Error))
    expect(result.current.error.message).toBe('offline')
  })

  it('refetches when the url changes', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(async (url) => ({ ok: true, json: async () => ({ url }) }))
    )
    const { result, rerender } = renderHook(({ url }) => useFetch(url), {
      initialProps: { url: '/a' },
    })
    await waitFor(() => expect(result.current.data).toEqual({ url: '/a' }))
    rerender({ url: '/b' })
    await waitFor(() => expect(result.current.data).toEqual({ url: '/b' }))
  })
})
