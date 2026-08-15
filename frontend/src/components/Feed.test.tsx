import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Feed } from './Feed'
import * as api from '../api/images'
import type { ImageItem } from '../types/image'

const sampleImages: ImageItem[] = [
  { id: '1', title: 'First', tags: ['cats'], url: '/uploads/1.jpg', created_at: 'now' },
  { id: '2', title: 'Second', tags: [], url: '/uploads/2.jpg', created_at: 'now' },
]

describe('Feed', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('shows a loading state before images resolve', () => {
    vi.spyOn(api, 'fetchImages').mockReturnValue(new Promise(() => {}))

    render(<Feed refreshKey={0} tags={[]} />)

    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('shows an empty state when there are no images', async () => {
    vi.spyOn(api, 'fetchImages').mockResolvedValue({ images: [], next_cursor: null })

    render(<Feed refreshKey={0} tags={[]} />)

    expect(await screen.findByText('No images yet. Be the first to upload one.')).toBeInTheDocument()
  })

  it('shows an error state when fetching fails', async () => {
    vi.spyOn(api, 'fetchImages').mockRejectedValue(new Error('boom'))

    render(<Feed refreshKey={0} tags={[]} />)

    expect(await screen.findByText('Failed to load images: boom')).toBeInTheDocument()
  })

  it('renders the fetched images with their titles and tags', async () => {
    vi.spyOn(api, 'fetchImages').mockResolvedValue({ images: sampleImages, next_cursor: null })

    render(<Feed refreshKey={0} tags={[]} />)

    expect(await screen.findByText('First')).toBeInTheDocument()
    expect(screen.getByText('Second')).toBeInTheDocument()
    expect(screen.getByText('cats')).toBeInTheDocument()
  })

  it('opens a modal with the full title and tags when an image is clicked', async () => {
    vi.spyOn(api, 'fetchImages').mockResolvedValue({ images: sampleImages, next_cursor: null })
    const user = userEvent.setup()

    render(<Feed refreshKey={0} tags={[]} />)
    const firstTitle = await screen.findByText('First')
    await user.click(firstTitle)

    expect(await screen.findByText(/Title:/)).toBeInTheDocument()
    expect(screen.getByText(/Tags:/)).toBeInTheDocument()
  })

  it('refetches when refreshKey or tags change', async () => {
    const fetchSpy = vi.spyOn(api, 'fetchImages').mockResolvedValue({ images: [], next_cursor: null })

    const { rerender } = render(<Feed refreshKey={0} tags={[]} />)
    await waitFor(() => expect(fetchSpy).toHaveBeenCalledTimes(1))

    rerender(<Feed refreshKey={1} tags={[]} />)
    await waitFor(() => expect(fetchSpy).toHaveBeenCalledTimes(2))

    rerender(<Feed refreshKey={1} tags={['cats']} />)
    await waitFor(() => expect(fetchSpy).toHaveBeenCalledTimes(3))
    expect(fetchSpy).toHaveBeenLastCalledWith(['cats'])
  })

  it('shows a load more button when there is a next cursor, and fetches the next page on click', async () => {
    const fetchSpy = vi
      .spyOn(api, 'fetchImages')
      .mockResolvedValueOnce({ images: [sampleImages[0]], next_cursor: 'cursor1' })
      .mockResolvedValueOnce({ images: [sampleImages[1]], next_cursor: null })
    const user = userEvent.setup()

    render(<Feed refreshKey={0} tags={[]} />)
    await screen.findByText('First')

    const loadMoreButton = screen.getByRole('button', { name: 'Load more' })
    await user.click(loadMoreButton)

    expect(await screen.findByText('Second')).toBeInTheDocument()
    expect(fetchSpy).toHaveBeenLastCalledWith([], 'cursor1')
    expect(screen.queryByRole('button', { name: 'Load more' })).not.toBeInTheDocument()
  })
})
