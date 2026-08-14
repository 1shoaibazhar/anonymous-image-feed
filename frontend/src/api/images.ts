export const API_BASE_URL = 'http://localhost:8080'

import type { ImageItem } from '../types/image'

export async function fetchImages(): Promise<ImageItem[]> {
  const response = await fetch(`${API_BASE_URL}/api/images`)
  if (!response.ok) {
    throw new Error(`Failed to fetch images: ${response.status}`)
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