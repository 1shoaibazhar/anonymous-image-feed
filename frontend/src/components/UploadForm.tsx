import { useState, type SubmitEventHandler } from 'react'
import { uploadImage } from '../api/images'

interface UploadFormProps {
  onUploadSuccess: () => void
}

export function UploadForm({ onUploadSuccess }: UploadFormProps) {
  const [file, setFile] = useState<File | null>(null)
  const [title, setTitle] = useState('')
  const [tags, setTags] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit: SubmitEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault()
    setError(null)

    if (!file) {
      setError('Please choose an image file')
      return
    }
    if (!title.trim()) {
      setError('Please enter a title')
      return
    }

    setSubmitting(true)
    try {
      await uploadImage(file, title.trim(), tags.trim())
      setFile(null)
      setTitle('')
      setTags('')
      const fileInput = document.getElementById('file-input') as HTMLInputElement | null
      if (fileInput) fileInput.value = ''
      onUploadSuccess()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="p-4 bg-white border-b flex flex-wrap gap-3 items-start">
      <input
        id="file-input"
        type="file"
        accept="image/jpeg,image/png,image/webp,image/gif"
        onChange={(e) => setFile(e.target.files?.[0] ?? null)}
        className="text-sm"
      />
      <input
        type="text"
        placeholder="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        className="border rounded px-3 py-1.5 text-sm flex-1 min-w-[150px]"
      />
      <input
        type="text"
        placeholder="Tags (comma separated)"
        value={tags}
        onChange={(e) => setTags(e.target.value)}
        className="border rounded px-3 py-1.5 text-sm flex-1 min-w-[150px]"
      />
      <button
        type="submit"
        disabled={submitting}
        className="bg-gray-800 text-white px-4 py-1.5 rounded text-sm disabled:opacity-50"
      >
        {submitting ? 'Uploading...' : 'Upload'}
      </button>
      {error && <p className="w-full text-sm text-red-600">{error}</p>}
    </form>
  )
}