import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { TagFilter } from './TagFilter'
import * as api from '../api/images'

describe('TagFilter', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('renders nothing while there are no tags', async () => {
    vi.spyOn(api, 'fetchTags').mockResolvedValue([])

    const { container } = render(<TagFilter selectedTags={[]} onChange={vi.fn()} />)

    await waitFor(() => expect(api.fetchTags).toHaveBeenCalled())
    expect(container).toBeEmptyDOMElement()
  })

  it('renders nothing when fetchTags rejects', async () => {
    vi.spyOn(api, 'fetchTags').mockRejectedValue(new Error('network error'))

    const { container } = render(<TagFilter selectedTags={[]} onChange={vi.fn()} />)

    await waitFor(() => expect(api.fetchTags).toHaveBeenCalled())
    expect(container).toBeEmptyDOMElement()
  })

  it('renders fetched tags and adds a tag to the selection when clicked', async () => {
    vi.spyOn(api, 'fetchTags').mockResolvedValue(['cats', 'dogs'])
    const onChange = vi.fn()
    const user = userEvent.setup()

    render(<TagFilter selectedTags={[]} onChange={onChange} />)

    const catsButton = await screen.findByRole('button', { name: 'cats' })
    await user.click(catsButton)

    expect(onChange).toHaveBeenCalledWith(['cats'])
  })

  it('removes an already-selected tag when clicked again', async () => {
    vi.spyOn(api, 'fetchTags').mockResolvedValue(['cats', 'dogs'])
    const onChange = vi.fn()
    const user = userEvent.setup()

    render(<TagFilter selectedTags={['cats']} onChange={onChange} />)

    const catsButton = await screen.findByRole('button', { name: 'cats' })
    await user.click(catsButton)

    expect(onChange).toHaveBeenCalledWith([])
  })
})
