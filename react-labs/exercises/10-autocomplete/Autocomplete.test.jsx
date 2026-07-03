import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import Autocomplete from './Autocomplete'

const FRUITS = ['apple', 'apricot', 'banana', 'cherry']
const makeFetch = () =>
  vi.fn(async (q) => FRUITS.filter((f) => f.startsWith(q.toLowerCase())))

describe('Autocomplete', () => {
  it('shows debounced suggestions as you type', async () => {
    const user = userEvent.setup()
    render(<Autocomplete fetchSuggestions={makeFetch()} debounceMs={100} />)

    await user.type(screen.getByLabelText('Search'), 'ap')
    const list = await screen.findByRole('listbox')
    expect(within(list).getByText('apple')).toBeInTheDocument()
    expect(within(list).getByText('apricot')).toBeInTheDocument()
    expect(within(list).queryByText('banana')).not.toBeInTheDocument()
  })

  it('does not open for an empty query', async () => {
    const user = userEvent.setup()
    const fetchSuggestions = makeFetch()
    render(<Autocomplete fetchSuggestions={fetchSuggestions} debounceMs={50} />)

    await user.type(screen.getByLabelText('Search'), 'a')
    await screen.findByRole('listbox')
    await user.clear(screen.getByLabelText('Search'))
    await waitFor(() =>
      expect(screen.queryByRole('listbox')).not.toBeInTheDocument()
    )
  })

  it('selects an option by clicking', async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()
    render(
      <Autocomplete fetchSuggestions={makeFetch()} onSelect={onSelect} debounceMs={50} />
    )

    await user.type(screen.getByLabelText('Search'), 'ba')
    const option = await screen.findByText('banana')
    await user.click(option)

    expect(onSelect).toHaveBeenCalledWith('banana')
    expect(screen.getByLabelText('Search')).toHaveValue('banana')
    expect(screen.queryByRole('listbox')).not.toBeInTheDocument()
  })

  it('navigates with the keyboard and selects with Enter', async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()
    render(
      <Autocomplete fetchSuggestions={makeFetch()} onSelect={onSelect} debounceMs={50} />
    )

    await user.type(screen.getByLabelText('Search'), 'ap')
    await screen.findByRole('listbox')

    await user.keyboard('{ArrowDown}') // highlight "apple"
    expect(screen.getByText('apple')).toHaveAttribute('aria-selected', 'true')

    await user.keyboard('{ArrowDown}') // highlight "apricot"
    expect(screen.getByText('apricot')).toHaveAttribute('aria-selected', 'true')

    await user.keyboard('{Enter}')
    expect(onSelect).toHaveBeenCalledWith('apricot')
  })

  it('closes on Escape', async () => {
    const user = userEvent.setup()
    render(<Autocomplete fetchSuggestions={makeFetch()} debounceMs={50} />)

    await user.type(screen.getByLabelText('Search'), 'ch')
    await screen.findByRole('listbox')
    await user.keyboard('{Escape}')
    expect(screen.queryByRole('listbox')).not.toBeInTheDocument()
  })
})
