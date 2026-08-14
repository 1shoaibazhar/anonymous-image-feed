const API_BASE_URL = 'http://localhost:8080'

import type { ImageItem } from '../types/image.ts'

export async function fetchImages(): Promise<ImageItem[]> {
  const response = await fetch(`${API_BASE_URL}/api/images`)
  if (!response.ok) {
    throw new Error(`Failed to fetch images: ${response.status}`)
  }
  return response.json()
}