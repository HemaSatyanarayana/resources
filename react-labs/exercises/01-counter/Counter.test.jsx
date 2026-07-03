import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect } from 'vitest'
import Counter from './Counter'

const value = () => screen.getByTestId('count').textContent

describe('Counter', () => {
  it('renders the initial value', () => {
    render(<Counter initial={5} />)
    expect(value()).toBe('5')
  })

  it('increments and decrements by step', async () => {
    const user = userEvent.setup()
    render(<Counter initial={0} step={2} />)

    await user.click(screen.getByRole('button', { name: 'Increment' }))
    expect(value()).toBe('2')

    await user.click(screen.getByRole('button', { name: 'Decrement' }))
    expect(value()).toBe('0')
  })

  it('clamps at max and disables Increment there', async () => {
    const user = userEvent.setup()
    render(<Counter initial={9} max={10} />)

    await user.click(screen.getByRole('button', { name: 'Increment' }))
    expect(value()).toBe('10')
    expect(screen.getByRole('button', { name: 'Increment' })).toBeDisabled()

    // Further clicks (if somehow fired) must not exceed max.
    await user.click(screen.getByRole('button', { name: 'Increment' }))
    expect(value()).toBe('10')
  })

  it('clamps at min and disables Decrement there', async () => {
    const user = userEvent.setup()
    render(<Counter initial={1} min={0} />)

    await user.click(screen.getByRole('button', { name: 'Decrement' }))
    expect(value()).toBe('0')
    expect(screen.getByRole('button', { name: 'Decrement' })).toBeDisabled()
  })

  it('resets to the initial value', async () => {
    const user = userEvent.setup()
    render(<Counter initial={3} />)

    await user.click(screen.getByRole('button', { name: 'Increment' }))
    await user.click(screen.getByRole('button', { name: 'Increment' }))
    expect(value()).toBe('5')

    await user.click(screen.getByRole('button', { name: 'Reset' }))
    expect(value()).toBe('3')
  })
})
