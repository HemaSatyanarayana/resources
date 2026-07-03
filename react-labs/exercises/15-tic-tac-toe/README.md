# 15 — Tic Tac Toe

The capstone. It's the purest test of the lab's central mantra: **derive, don't store.** The winner, the draw, and whose turn it is are all *functions of the board* (plus one turn boolean). Candidates who store a `winner` state variable end up with sync bugs; candidates who derive it write half the code and never have a stale winner.

## Concepts

- **Minimal state**: `board` (array of 9) and `xIsNext` (boolean). That's it. `winner`, `isDraw`, `current`, and `status` are all derived every render.
- **Immutability**: `const next = [...board]; next[i] = current; setBoard(next)`. Never `board[i] = ...`.
- **Derived winner**: check the 8 winning lines against the board. Pure function, no state.
- **Guard clauses**: a click is a no-op if the game is won or the cell is taken — model illegal moves out of existence.

## Your task

Implement `calculateWinner`, `isDraw`, and `handleClick` in [`TicTacToe.jsx`](TicTacToe.jsx).

## Run

```bash
npx vitest run exercises/15-tic-tac-toe
```

## Interview follow-ups to expect

- *"Why not keep `winner` in state?"* — it's fully determined by `board`; a second copy can disagree with it after an update, and you'd have to remember to recompute it on every move. Deriving makes an inconsistent state unrepresentable.
- *"NxN board / K-in-a-row."* — generalize `LINES` generation instead of hardcoding 8 lines.
- *"Move history / undo (time travel)."* — store an array of board snapshots and an index; this is the canonical React tutorial extension.
- *"Highlight the winning line."* — have `calculateWinner` return the line, not just the mark.

## Hints

```js
function calculateWinner(board) {
  for (const [a, b, c] of LINES) {
    if (board[a] && board[a] === board[b] && board[a] === board[c]) return board[a]
  }
  return null
}
```

- `isDraw = !winner && board.every(Boolean)`.
- `handleClick`: `if (winner || board[i]) return` first, then the immutable update + `setXIsNext(x => !x)`.

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState } from 'react'

const LINES = [
  [0, 1, 2], [3, 4, 5], [6, 7, 8],
  [0, 3, 6], [1, 4, 7], [2, 5, 8],
  [0, 4, 8], [2, 4, 6],
]

function calculateWinner(board) {
  for (const [a, b, c] of LINES) {
    if (board[a] && board[a] === board[b] && board[a] === board[c]) {
      return board[a]
    }
  }
  return null
}

export default function TicTacToe() {
  const [board, setBoard] = useState(() => Array(9).fill(null))
  const [xIsNext, setXIsNext] = useState(true)

  const winner = calculateWinner(board)
  const isDraw = !winner && board.every(Boolean)
  const current = xIsNext ? 'X' : 'O'
  const status = winner
    ? `Winner: ${winner}`
    : isDraw
      ? 'Draw'
      : `Turn: ${current}`

  const handleClick = (i) => {
    if (winner || board[i]) return
    const next = [...board]
    next[i] = current
    setBoard(next)
    setXIsNext((x) => !x)
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
```

</details>
