import { useEffect, useState } from 'react'
import { fetchImages, API_BASE_URL } from '../api/images'
import type { ImageItem } from '../types/image'
import { Modal } from './Modal'

interface FeedProps {
  refreshKey: number
  tags: string[]
}

export function Feed({ refreshKey, tags }: FeedProps) {
  const [images, setImages] = useState<ImageItem[]>([])
  const [nextCursor, setNextCursor] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [selectedImage, setSelectedImage] = useState<ImageItem | null>(null)

  useEffect(() => {
    setLoading(true)
    fetchImages(tags)
      .then((page) => {
        setImages(page.images)
        setNextCursor(page.next_cursor)
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [refreshKey, tags])

  function loadMore() {
    if (!nextCursor || loadingMore) return
    setLoadingMore(true)
    fetchImages(tags, nextCursor)
      .then((page) => {
        setImages((prev) => [...prev, ...page.images])
        setNextCursor(page.next_cursor)
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoadingMore(false))
  }

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
    <div>
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4 p-4">
      {images.map((image) => (
        <div
          key={image.id}
          className="bg-white rounded-lg shadow overflow-hidden cursor-pointer"
          onClick={() => setSelectedImage(image)}
        >
          <img
            src={`${API_BASE_URL}${image.url}`}
            alt={image.title}
            className="w-full h-40 object-cover"
          />
          <div className="p-3">
            <h3 className="font-semibold text-gray-800 truncate">{image.title}</h3>
            <div className="flex flex-wrap gap-1 mt-2">
              {image.tags.map((tag) => (
                <span key={tag} className="text-xs bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full">
                  {tag}
                </span>
              ))}
            </div>
          </div>
        </div>
      ))}
      {selectedImage && (
        <Modal onClose={() => setSelectedImage(null)} panelClassName="max-w-3xl max-h-[90vh] overflow-auto">
          <img
            src={`${API_BASE_URL}${selectedImage.url}`}
            alt={selectedImage.title}
            className="w-full h-auto max-h-[70vh] object-contain rounded-t-lg"
          />
          <div className="p-4">
            <p className="text-gray-800">
              <span className="font-semibold">Title:</span> {selectedImage.title}
            </p>
            <p className="text-gray-800 mt-1">
              <span className="font-semibold">Tags:</span> {selectedImage.tags.join(', ')}
            </p>
          </div>
        </Modal>
      )}
    </div>
    {nextCursor && (
      <p className="text-center pb-8">
        <button
          onClick={loadMore}
          disabled={loadingMore}
          className="px-4 py-2 bg-gray-800 text-white rounded-lg disabled:opacity-50"
        >
          {loadingMore ? 'Loading...' : 'Load more'}
        </button>
      </p>
    )}
    </div>
  )
}