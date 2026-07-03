import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect } from 'vitest'
import Accordion from './Accordion'

const items = [
  { id: 'a', title: 'Section A', content: 'Body A' },
  { id: 'b', title: 'Section B', content: 'Body B' },
  { id: 'c', title: 'Section C', content: 'Body C' },
]

describe('Accordion', () => {
  it('starts fully collapsed', () => {
    render(<Accordion items={items} />)
    expect(screen.queryByText('Body A')).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Section A' })).toHaveAttribute(
      'aria-expanded',
      'false'
    )
  })

  it('expands a section on click', async () => {
    const user = userEvent.setup()
    render(<Accordion items={items} />)
    await user.click(screen.getByRole('button', { name: 'Section A' }))
    expect(screen.getByText('Body A')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Section A' })).toHaveAttribute(
      'aria-expanded',
      'true'
    )
  })

  it('collapses again on second click', async () => {
    const user = userEvent.setup()
    render(<Accordion items={items} />)
    await user.click(screen.getByRole('button', { name: 'Section A' }))
    await user.click(screen.getByRole('button', { name: 'Section A' }))
    expect(screen.queryByText('Body A')).not.toBeInTheDocument()
  })

  it('single-open mode closes the previous section', async () => {
    const user = userEvent.setup()
    render(<Accordion items={items} />)
    await user.click(screen.getByRole('button', { name: 'Section A' }))
    await user.click(screen.getByRole('button', { name: 'Section B' }))
    expect(screen.queryByText('Body A')).not.toBeInTheDocument()
    expect(screen.getByText('Body B')).toBeInTheDocument()
  })

  it('allowMultiple keeps sections independent', async () => {
    const user = userEvent.setup()
    render(<Accordion items={items} allowMultiple />)
    await user.click(screen.getByRole('button', { name: 'Section A' }))
    await user.click(screen.getByRole('button', { name: 'Section B' }))
    expect(screen.getByText('Body A')).toBeInTheDocument()
    expect(screen.getByText('Body B')).toBeInTheDocument()
  })
})
