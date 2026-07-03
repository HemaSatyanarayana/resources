import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect } from 'vitest'
import NestedCheckboxes from './NestedCheckboxes'

const tree = {
  id: 'all',
  label: 'All',
  children: [
    {
      id: 'fruits',
      label: 'Fruits',
      children: [
        { id: 'apple', label: 'Apple' },
        { id: 'banana', label: 'Banana' },
      ],
    },
    {
      id: 'veg',
      label: 'Vegetables',
      children: [{ id: 'carrot', label: 'Carrot' }],
    },
  ],
}

const box = (name) => screen.getByRole('checkbox', { name })

describe('NestedCheckboxes', () => {
  it('starts fully unchecked', () => {
    render(<NestedCheckboxes tree={tree} />)
    expect(box('All')).not.toBeChecked()
    expect(box('Apple')).not.toBeChecked()
  })

  it('makes the parent indeterminate when only some children are checked', async () => {
    const user = userEvent.setup()
    render(<NestedCheckboxes tree={tree} />)
    await user.click(box('Apple'))
    expect(box('Apple')).toBeChecked()
    expect(box('Fruits')).toBePartiallyChecked()
    expect(box('All')).toBePartiallyChecked()
  })

  it('checks the parent when all children are checked', async () => {
    const user = userEvent.setup()
    render(<NestedCheckboxes tree={tree} />)
    await user.click(box('Apple'))
    await user.click(box('Banana'))
    expect(box('Fruits')).toBeChecked()
    expect(box('Fruits')).not.toBePartiallyChecked()
  })

  it('checking a parent checks every descendant', async () => {
    const user = userEvent.setup()
    render(<NestedCheckboxes tree={tree} />)
    await user.click(box('Fruits'))
    expect(box('Apple')).toBeChecked()
    expect(box('Banana')).toBeChecked()
  })

  it('checking the root checks the whole tree, unchecking clears it', async () => {
    const user = userEvent.setup()
    render(<NestedCheckboxes tree={tree} />)
    await user.click(box('All'))
    expect(box('Apple')).toBeChecked()
    expect(box('Carrot')).toBeChecked()

    await user.click(box('All'))
    expect(box('Apple')).not.toBeChecked()
    expect(box('Carrot')).not.toBeChecked()
  })
})
