import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { UploadForm } from './UploadForm'
import * as api from '../api/images'

function selectFile(input: HTMLElement, file: File) {
  return userEvent.setup().upload(input, file)
}

describe('UploadForm', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('shows an error and does not upload when no file is chosen', async () => {
    const uploadSpy = vi.spyOn(api, 'uploadImage')
    const user = userEvent.setup()
    const onUploadSuccess = vi.fn()

    render(<UploadForm onUploadSuccess={onUploadSuccess} />)
    await user.type(screen.getByPlaceholderText('Title'), 'My title')
    await user.click(screen.getByRole('button', { name: 'Upload' }))

    expect(await screen.findByText('Please choose an image file')).toBeInTheDocument()
    expect(uploadSpy).not.toHaveBeenCalled()
    expect(onUploadSuccess).not.toHaveBeenCalled()
  })

  it('shows an error and does not upload when title is blank', async () => {
    const uploadSpy = vi.spyOn(api, 'uploadImage')
    const user = userEvent.setup()
    const file = new File(['data'], 'photo.jpg', { type: 'image/jpeg' })

    render(<UploadForm onUploadSuccess={vi.fn()} />)
    const input = document.getElementById('file-input') as HTMLInputElement
    await selectFile(input, file)
    await user.click(screen.getByRole('button', { name: 'Upload' }))

    expect(await screen.findByText('Please enter a title')).toBeInTheDocument()
    expect(uploadSpy).not.toHaveBeenCalled()
  })

  it('uploads the file with a trimmed title and tags, then calls onUploadSuccess', async () => {
    const created = { id: '1', title: 'My title', tags: ['a'], url: '/uploads/1.jpg', created_at: 'now' }
    const uploadSpy = vi.spyOn(api, 'uploadImage').mockResolvedValue(created)
    const onUploadSuccess = vi.fn()
    const user = userEvent.setup()
    const file = new File(['data'], 'photo.jpg', { type: 'image/jpeg' })

    render(<UploadForm onUploadSuccess={onUploadSuccess} />)
    const input = document.getElementById('file-input') as HTMLInputElement
    await selectFile(input, file)
    await user.type(screen.getByPlaceholderText('Title'), '  My title  ')
    await user.type(screen.getByPlaceholderText('Tags (comma separated)'), ' a ')
    await user.click(screen.getByRole('button', { name: 'Upload' }))

    expect(uploadSpy).toHaveBeenCalledWith(file, 'My title', 'a')
    expect(await screen.findByRole('button', { name: 'Upload' })).toBeEnabled()
    expect(onUploadSuccess).toHaveBeenCalledTimes(1)
    expect(screen.getByPlaceholderText('Title')).toHaveValue('')
  })

  it('shows the server error message when the upload fails', async () => {
    vi.spyOn(api, 'uploadImage').mockRejectedValue(new Error('unsupported file type'))
    const user = userEvent.setup()
    const file = new File(['data'], 'photo.jpg', { type: 'image/jpeg' })

    render(<UploadForm onUploadSuccess={vi.fn()} />)
    const input = document.getElementById('file-input') as HTMLInputElement
    await selectFile(input, file)
    await user.type(screen.getByPlaceholderText('Title'), 'My title')
    await user.click(screen.getByRole('button', { name: 'Upload' }))

    expect(await screen.findByText('unsupported file type')).toBeInTheDocument()
  })
})
