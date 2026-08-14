import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Modal } from './Modal'

describe('Modal', () => {
  it('renders its children', () => {
    render(
      <Modal onClose={vi.fn()}>
        <p>panel content</p>
      </Modal>,
    )

    expect(screen.getByText('panel content')).toBeInTheDocument()
  })

  it('calls onClose when the backdrop is clicked', async () => {
    const onClose = vi.fn()
    const user = userEvent.setup()

    const { container } = render(
      <Modal onClose={onClose}>
        <p>panel content</p>
      </Modal>,
    )

    await user.click(container.firstChild as Element)

    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it('does not call onClose when the panel itself is clicked', async () => {
    const onClose = vi.fn()
    const user = userEvent.setup()

    render(
      <Modal onClose={onClose}>
        <p>panel content</p>
      </Modal>,
    )

    await user.click(screen.getByText('panel content'))

    expect(onClose).not.toHaveBeenCalled()
  })

  it('calls onClose when the close button is clicked', async () => {
    const onClose = vi.fn()
    const user = userEvent.setup()

    render(
      <Modal onClose={onClose}>
        <p>panel content</p>
      </Modal>,
    )

    await user.click(screen.getByRole('button', { name: 'Close' }))

    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it('calls onClose when the Escape key is pressed', async () => {
    const onClose = vi.fn()
    const user = userEvent.setup()

    render(
      <Modal onClose={onClose}>
        <p>panel content</p>
      </Modal>,
    )

    await user.keyboard('{Escape}')

    expect(onClose).toHaveBeenCalledTimes(1)
  })
})
