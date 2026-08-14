import { useEffect, useState } from 'react'
import { fetchImages } from '../api/images'
import type { ImageItem } from '../types/image'

export function Feed() {
  const [images, setImages] = useState<ImageItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchImages()
      .then(setImages)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return <p className="text-center text-gray-500 mt-8">Loading...</p>
  }

  if (error) {
    return <p className="text-center text-red-600 mt-8">Failed to load images: {error}</p>
  }

  if (images.length === 0) {
    return <p className="text-center text-gray-500 mt-8">No images yet. Be the first to upload one.</p>
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4 p-4">
      {images.map((image) => (
        <div key={image.id} className="bg-white rounded-lg shadow p-3">
          <h3 className="font-semibold text-gray-800 truncate">{image.title}</h3>
          <div className="flex flex-wrap gap-1 mt-2">
            {image.tags.map((tag) => (
              <span key={tag} className="text-xs bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full">
                {tag}
              </span>
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}