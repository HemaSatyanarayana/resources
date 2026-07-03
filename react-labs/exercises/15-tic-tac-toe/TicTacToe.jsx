import { useState } from 'react'

/**
 * TicTacToe — derived winner logic, immutable board, turn state.
 *
 * Required UI (the tests rely on this):
 *   - 9 cell buttons with data-testid="cell-0" .. "cell-8"; a played cell shows
 *     "X" or "O", an empty cell is blank
 *   - a status element with data-testid="status" showing:
 *       "Winner: X" / "Winner: O"   when there's a winner
 *       "Draw"                       when the board is full with no winner
 *       "Turn: X" / "Turn: O"        otherwise
 *   - a button "Reset" that clears the board
 *
 * Rules:
 *   - X moves first, then players alternate.
 *   - Clicking a filled cell does nothing. Clicking after the game is over does nothing.
 *   - The board must be updated immutably (copy the array, set one index).
 *   - The winner is DERIVED from the board every render — never stored in state.
 */

const LINES = [
  [0, 1, 2], [3, 4, 5], [6, 7, 8], // rows
  [0, 3, 6], [1, 4, 7], [2, 5, 8], // cols
  [0, 4, 8], [2, 4, 6],            // diagonals
]

function calculateWinner(board) {
  // TODO: for each line [a, b, c], if board[a] is truthy and board[a] === board[b]
  //       === board[c], return board[a]. If none match, return null.
  return null
}

export default function TicTacToe() {
  const [board, setBoard] = useState(() => Array(9).fill(null))
  const [xIsNext, setXIsNext] = useState(true)

  const winner = calculateWinner(board)
  const isDraw = /* TODO: no winner AND every cell filled */ false
  const current = xIsNext ? 'X' : 'O'
  const status = winner
    ? `Winner: ${winner}`
    : isDraw
      ? 'Draw'
      : `Turn: ${current}`

  const handleClick = (i) => {
    // TODO:
    //   - if there's a winner already, or board[i] is filled: return
    //   - copy the board, set index i to `current`, setBoard(copy)
    //   - flip xIsNext
  }

  const reset = () => {
    setBoard(Array(9).fill(null))
    setXIsNext(true)
  }

  return (
    <div>
      <div data-testid="status">{status}</div>
      <div>
        {board.map((cell, i) => (
          <button key={i} data-testid={`cell-${i}`} onClick={() => handleClick(i)}>
            {cell}
          </button>
        ))}
      </div>
      <button onClick={reset}>Reset</button>
    </div>
  )
}
