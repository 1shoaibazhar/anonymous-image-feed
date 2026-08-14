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
    <form onSubmit={handleSubmit} className="flex flex-col gap-3">
      <h2 className="text-lg font-semibold text-gray-800">Upload image</h2>
      <div className="flex flex-col gap-1">
        <input
          id="file-input"
          type="file"
          accept="image/jpeg,image/png,image/webp,image/gif"
          onChange={(e) => setFile(e.target.files?.[0] ?? null)}
          className="hidden"
        />
        <label
          htmlFor="file-input"
          className="cursor-pointer inline-block w-fit bg-gray-100 hover:bg-gray-200 border border-gray-300 text-gray-700 px-4 py-1.5 rounded text-sm"
        >
          Choose file
        </label>
        <span className="text-sm text-gray-500">{file ? file.name : 'No file chosen'}</span>
      </div>
      <input
        type="text"
        placeholder="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        className="border rounded px-3 py-1.5 text-sm"
      />
      <input
        type="text"
        placeholder="Tags (comma separated)"
        value={tags}
        onChange={(e) => setTags(e.target.value)}
        className="border rounded px-3 py-1.5 text-sm"
      />
      <button
        type="submit"
        disabled={submitting}
        className="bg-gray-800 text-white px-4 py-1.5 rounded text-sm disabled:opacity-50"
      >
        {submitting ? 'Uploading...' : 'Upload'}
      </button>
      {error && <p className="text-sm text-red-600">{error}</p>}
    </form>
  )
}