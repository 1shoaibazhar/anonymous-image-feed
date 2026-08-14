export const API_BASE_URL = 'http://localhost:8080'

import type { ImageItem } from '../types/image'

export async function fetchImages(tags: string[] = []): Promise<ImageItem[]> {
  const query = tags.length > 0 ? `?tags=${tags.map(encodeURIComponent).join(',')}` : ''
  const response = await fetch(`${API_BASE_URL}/api/images${query}`)
  if (!response.ok) {
    throw new Error(`Failed to fetch images: ${response.status}`)
  }
  return response.json()
}

export async function fetchTags(): Promise<string[]> {
  const response = await fetch(`${API_BASE_URL}/api/tags`)
  if (!response.ok) {
    throw new Error(`Failed to fetch tags: ${response.status}`)
  }
  return response.json()
}

export async function uploadImage(file: File, title: string, tags: string): Promise<ImageItem> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('title', title)
  formData.append('tags', tags)

  const response = await fetch(`${API_BASE_URL}/api/upload`, {
    method: 'POST',
    body: formData,
  })

  if (!response.ok) {
    const message = await response.text()
    throw new Error(message || `Upload failed: ${response.status}`)
  }

  return response.json()
}