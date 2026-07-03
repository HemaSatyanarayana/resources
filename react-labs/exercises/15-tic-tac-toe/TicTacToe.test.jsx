import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect } from 'vitest'
import TicTacToe from './TicTacToe'

const status = () => screen.getByTestId('status').textContent
const cell = (i) => screen.getByTestId(`cell-${i}`)

async function play(user, order) {
  for (const i of order) await user.click(cell(i))
}

describe('TicTacToe', () => {
  it('starts on X and alternates turns', async () => {
    const user = userEvent.setup()
    render(<TicTacToe />)
    expect(status()).toBe('Turn: X')

    await user.click(cell(0))
    expect(cell(0)).toHaveTextContent('X')
    expect(status()).toBe('Turn: O')

    await user.click(cell(1))
    expect(cell(1)).toHaveTextContent('O')
    expect(status()).toBe('Turn: X')
  })

  it('ignores clicks on a filled cell', async () => {
    const user = userEvent.setup()
    render(<TicTacToe />)
    await user.click(cell(0)) // X
    await user.click(cell(0)) // no-op
    expect(cell(0)).toHaveTextContent('X')
    expect(status()).toBe('Turn: O')
  })

  it('detects a winner and freezes the game', async () => {
    const user = userEvent.setup()
    render(<TicTacToe />)
    // X: 0,1,2 (top row) — O: 3,4
    await play(user, [0, 3, 1, 4, 2])
    expect(status()).toBe('Winner: X')

    await user.click(cell(5)) // game over -> ignored
    expect(cell(5)).toHaveTextContent('')
  })

  it('detects a draw', async () => {
    const user = userEvent.setup()
    render(<TicTacToe />)
    await play(user, [0, 1, 2, 4, 3, 5, 7, 6, 8])
    expect(status()).toBe('Draw')
  })

  it('resets the board', async () => {
    const user = userEvent.setup()
    render(<TicTacToe />)
    await play(user, [0, 1, 2])
    await user.click(screen.getByRole('button', { name: 'Reset' }))
    expect(cell(0)).toHaveTextContent('')
    expect(status()).toBe('Turn: X')
  })
})
