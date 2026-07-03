import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect } from 'vitest'
import Pagination from './Pagination'

const items = Array.from({ length: 23 }, (_, i) => `Item ${i + 1}`)

describe('Pagination', () => {
  it('shows only the first page initially', () => {
    render(<Pagination items={items} pageSize={10} />)
    expect(screen.getAllByRole('listitem')).toHaveLength(10)
    expect(screen.getByText('Item 1')).toBeInTheDocument()
    expect(screen.queryByText('Item 11')).not.toBeInTheDocument()
  })

  it('renders one page button per page', () => {
    render(<Pagination items={items} pageSize={10} />)
    // 23 items / 10 => 3 pages (+ Previous + Next = 5 buttons total)
    expect(screen.getByRole('button', { name: '1' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '3' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '4' })).not.toBeInTheDocument()
  })

  it('advances with Next', async () => {
    const user = userEvent.setup()
    render(<Pagination items={items} pageSize={10} />)
    await user.click(screen.getByRole('button', { name: 'Next' }))
    expect(screen.getByText('Item 11')).toBeInTheDocument()
    expect(screen.queryByText('Item 1')).not.toBeInTheDocument()
  })

  it('disables Previous on the first page and Next on the last', async () => {
    const user = userEvent.setup()
    render(<Pagination items={items} pageSize={10} />)
    expect(screen.getByRole('button', { name: 'Previous' })).toBeDisabled()

    await user.click(screen.getByRole('button', { name: '3' }))
    expect(screen.getByRole('button', { name: 'Next' })).toBeDisabled()
    expect(screen.getAllByRole('listitem')).toHaveLength(3) // 23 - 20
  })

  it('marks the current page with aria-current', async () => {
    const user = userEvent.setup()
    render(<Pagination items={items} pageSize={10} />)
    expect(screen.getByRole('button', { name: '1' })).toHaveAttribute(
      'aria-current',
      'page'
    )
    await user.click(screen.getByRole('button', { name: '2' }))
    expect(screen.getByRole('button', { name: '2' })).toHaveAttribute(
      'aria-current',
      'page'
    )
    expect(screen.getByRole('button', { name: '1' })).not.toHaveAttribute(
      'aria-current'
    )
  })
})
