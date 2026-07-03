import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect } from 'vitest'
import TodoList from './TodoList'

async function addTodo(user, text) {
  await user.clear(screen.getByLabelText('New todo'))
  await user.type(screen.getByLabelText('New todo'), text)
  await user.click(screen.getByRole('button', { name: 'Add' }))
}

describe('TodoList', () => {
  it('adds a todo', async () => {
    const user = userEvent.setup()
    render(<TodoList />)
    await addTodo(user, 'Buy milk')
    expect(screen.getByText('Buy milk')).toBeInTheDocument()
  })

  it('trims input and ignores empty todos', async () => {
    const user = userEvent.setup()
    render(<TodoList />)
    await addTodo(user, '   ')
    expect(screen.queryAllByRole('listitem')).toHaveLength(0)

    await addTodo(user, '  Walk dog  ')
    expect(screen.getByRole('checkbox', { name: 'Walk dog' })).toBeInTheDocument()
  })

  it('clears the input after adding', async () => {
    const user = userEvent.setup()
    render(<TodoList />)
    await addTodo(user, 'Something')
    expect(screen.getByLabelText('New todo')).toHaveValue('')
  })

  it('toggles completion', async () => {
    const user = userEvent.setup()
    render(<TodoList />)
    await addTodo(user, 'Task A')
    const box = screen.getByRole('checkbox', { name: 'Task A' })
    expect(box).not.toBeChecked()
    await user.click(box)
    expect(box).toBeChecked()
  })

  it('deletes a todo', async () => {
    const user = userEvent.setup()
    render(<TodoList />)
    await addTodo(user, 'Delete me')
    const item = screen.getByText('Delete me').closest('li')
    await user.click(within(item).getByRole('button', { name: 'Delete' }))
    expect(screen.queryByText('Delete me')).not.toBeInTheDocument()
  })

  it('filters by active and completed', async () => {
    const user = userEvent.setup()
    render(<TodoList />)
    await addTodo(user, 'Active one')
    await addTodo(user, 'Done one')
    await user.click(screen.getByRole('checkbox', { name: 'Done one' }))

    await user.click(screen.getByRole('button', { name: 'Active' }))
    expect(screen.getByText('Active one')).toBeInTheDocument()
    expect(screen.queryByText('Done one')).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'Completed' }))
    expect(screen.queryByText('Active one')).not.toBeInTheDocument()
    expect(screen.getByText('Done one')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'All' }))
    expect(screen.getAllByRole('listitem')).toHaveLength(2)
  })
})
