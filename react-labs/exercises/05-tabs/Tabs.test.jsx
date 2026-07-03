import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect } from 'vitest'
import Tabs from './Tabs'

const tabs = [
  { id: 'one', label: 'One', content: 'Content one' },
  { id: 'two', label: 'Two', content: 'Content two' },
  { id: 'three', label: 'Three', content: 'Content three' },
]

const panel = () => screen.getByRole('tabpanel')

describe('Tabs', () => {
  it('activates the first tab by default', () => {
    render(<Tabs tabs={tabs} />)
    expect(screen.getByRole('tab', { name: 'One' })).toHaveAttribute(
      'aria-selected',
      'true'
    )
    expect(panel()).toHaveTextContent('Content one')
  })

  it('honors defaultActiveId', () => {
    render(<Tabs tabs={tabs} defaultActiveId="two" />)
    expect(panel()).toHaveTextContent('Content two')
  })

  it('switches on click', async () => {
    const user = userEvent.setup()
    render(<Tabs tabs={tabs} />)
    await user.click(screen.getByRole('tab', { name: 'Three' }))
    expect(panel()).toHaveTextContent('Content three')
    expect(screen.getByRole('tab', { name: 'One' })).toHaveAttribute(
      'aria-selected',
      'false'
    )
  })

  it('sets a roving tabindex', () => {
    render(<Tabs tabs={tabs} />)
    expect(screen.getByRole('tab', { name: 'One' })).toHaveAttribute('tabindex', '0')
    expect(screen.getByRole('tab', { name: 'Two' })).toHaveAttribute('tabindex', '-1')
  })

  it('navigates with arrow keys, wrapping around', async () => {
    const user = userEvent.setup()
    render(<Tabs tabs={tabs} />)

    await user.click(screen.getByRole('tab', { name: 'One' }))
    await user.keyboard('{ArrowRight}')
    expect(panel()).toHaveTextContent('Content two')

    await user.keyboard('{ArrowLeft}{ArrowLeft}')
    expect(panel()).toHaveTextContent('Content three')
  })
})
