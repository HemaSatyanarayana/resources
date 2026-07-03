import { render, screen, act, fireEvent } from '@testing-library/react'
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import Stopwatch from './Stopwatch'

const elapsed = () => screen.getByTestId('elapsed').textContent
const click = (name) => fireEvent.click(screen.getByRole('button', { name }))

// Fake timers + userEvent don't mix well, so drive clicks with fireEvent
// (synchronous) and advance the clock with vi.advanceTimersByTime inside act().
describe('Stopwatch', () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  it('starts at 0.0', () => {
    render(<Stopwatch />)
    expect(elapsed()).toBe('0.0')
  })

  it('counts up while running', () => {
    render(<Stopwatch />)
    act(() => click('Start'))
    act(() => vi.advanceTimersByTime(1000))
    expect(elapsed()).toBe('1.0')
  })

  it('freezes on Stop', () => {
    render(<Stopwatch />)
    act(() => click('Start'))
    act(() => vi.advanceTimersByTime(1000))
    act(() => click('Stop'))
    act(() => vi.advanceTimersByTime(2000))
    expect(elapsed()).toBe('1.0')
  })

  it('resets to 0.0', () => {
    render(<Stopwatch />)
    act(() => click('Start'))
    act(() => vi.advanceTimersByTime(1500))
    act(() => click('Reset'))
    expect(elapsed()).toBe('0.0')
  })

  it('disables Start while running and Stop while stopped', () => {
    render(<Stopwatch />)
    expect(screen.getByRole('button', { name: 'Stop' })).toBeDisabled()
    act(() => click('Start'))
    expect(screen.getByRole('button', { name: 'Start' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'Stop' })).toBeEnabled()
  })
})
