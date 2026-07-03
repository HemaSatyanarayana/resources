import { renderHook, act } from '@testing-library/react'
import { describe, it, expect, beforeEach } from 'vitest'
import { useToggle, usePrevious, useLocalStorage } from './hooks'

describe('useToggle', () => {
  it('flips the value', () => {
    const { result } = renderHook(() => useToggle(false))
    expect(result.current[0]).toBe(false)
    act(() => result.current[1]())
    expect(result.current[0]).toBe(true)
  })

  it('forces a value when passed a boolean', () => {
    const { result } = renderHook(() => useToggle(true))
    act(() => result.current[1](true))
    expect(result.current[0]).toBe(true)
    act(() => result.current[1](false))
    expect(result.current[0]).toBe(false)
  })

  it('keeps a stable toggle identity across renders', () => {
    const { result, rerender } = renderHook(() => useToggle())
    const first = result.current[1]
    act(() => result.current[1]())
    rerender()
    expect(result.current[1]).toBe(first)
  })
})

describe('usePrevious', () => {
  it('returns undefined on first render, then the prior value', () => {
    const { result, rerender } = renderHook(({ v }) => usePrevious(v), {
      initialProps: { v: 1 },
    })
    expect(result.current).toBeUndefined()
    rerender({ v: 2 })
    expect(result.current).toBe(1)
    rerender({ v: 3 })
    expect(result.current).toBe(2)
  })
})

describe('useLocalStorage', () => {
  beforeEach(() => localStorage.clear())

  it('initializes from the provided default and persists changes', () => {
    const { result } = renderHook(() => useLocalStorage('theme', 'light'))
    expect(result.current[0]).toBe('light')
    act(() => result.current[1]('dark'))
    expect(result.current[0]).toBe('dark')
    expect(JSON.parse(localStorage.getItem('theme'))).toBe('dark')
  })

  it('reads an existing stored value on mount', () => {
    localStorage.setItem('count', JSON.stringify(42))
    const { result } = renderHook(() => useLocalStorage('count', 0))
    expect(result.current[0]).toBe(42)
  })

  it('supports the functional-updater form', () => {
    const { result } = renderHook(() => useLocalStorage('n', 1))
    act(() => result.current[1]((n) => n + 4))
    expect(result.current[0]).toBe(5)
  })
})
