import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import StarRating from './StarRating'

const stars = () => screen.getAllByRole('button')
const filledCount = () =>
  stars().filter((b) => b.getAttribute('data-filled') === 'true').length

describe('StarRating', () => {
  it('renders `count` stars', () => {
    render(<StarRating count={5} />)
    expect(stars()).toHaveLength(5)
  })

  it('shows the default value as filled', () => {
    render(<StarRating count={5} defaultValue={3} />)
    expect(filledCount()).toBe(3)
  })

  it('selects a rating on click and calls onChange', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()
    render(<StarRating count={5} onChange={onChange} />)

    await user.click(screen.getByRole('button', { name: 'Rate 4 stars' }))
    expect(onChange).toHaveBeenCalledWith(4)
    expect(filledCount()).toBe(4)
  })

  it('previews on hover and restores the selection on mouse leave', async () => {
    const user = userEvent.setup()
    render(<StarRating count={5} defaultValue={2} />)

    await user.hover(screen.getByRole('button', { name: 'Rate 5 stars' }))
    expect(filledCount()).toBe(5)

    await user.unhover(screen.getByRole('button', { name: 'Rate 5 stars' }))
    expect(filledCount()).toBe(2)
  })

  it('uses the singular label for the first star', () => {
    render(<StarRating count={3} />)
    expect(screen.getByRole('button', { name: 'Rate 1 star' })).toBeInTheDocument()
  })
})
