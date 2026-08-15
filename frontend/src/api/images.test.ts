import { describe, it, expect, vi, beforeEach } from 'vitest'
import { fetchImages, fetchTags, uploadImage, API_BASE_URL } from './images'

function mockFetchResponse(body: unknown, ok = true, status = 200) {
  return {
    ok,
    status,
    json: () => Promise.resolve(body),
    text: () => Promise.resolve(typeof body === 'string' ? body : JSON.stringify(body)),
  } as Response
}

describe('fetchImages', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
  })

  it('requests images without a query when no tags or cursor are given', async () => {
    const page = {
      images: [{ id: '1', title: 'a', tags: [], url: '/uploads/a.jpg', created_at: 'now' }],
      next_cursor: null,
    }
    vi.mocked(fetch).mockResolvedValueOnce(mockFetchResponse(page))

    const result = await fetchImages()

    expect(fetch).toHaveBeenCalledWith(`${API_BASE_URL}/api/images`)
    expect(result).toEqual(page)
  })

  it('appends an encoded, comma-joined tags query when tags are given', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(mockFetchResponse({ images: [], next_cursor: null }))

    await fetchImages(['cats', 'funny animals'])

    expect(fetch).toHaveBeenCalledWith(`${API_BASE_URL}/api/images?tags=cats,funny%20animals`)
  })

  it('appends an encoded cursor query when a cursor is given', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(mockFetchResponse({ images: [], next_cursor: null }))

    await fetchImages([], 'abc+def')

    expect(fetch).toHaveBeenCalledWith(`${API_BASE_URL}/api/images?cursor=abc%2Bdef`)
  })

  it('combines tags and cursor queries when both are given', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(mockFetchResponse({ images: [], next_cursor: null }))

    await fetchImages(['cats'], 'abc')

    expect(fetch).toHaveBeenCalledWith(`${API_BASE_URL}/api/images?tags=cats&cursor=abc`)
  })

  it('throws when the response is not ok', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(mockFetchResponse(null, false, 500))

    await expect(fetchImages()).rejects.toThrow('Failed to fetch images: 500')
  })
})

describe('fetchTags', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
  })

  it('returns the parsed tag list', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(mockFetchResponse(['cats', 'dogs']))

    const result = await fetchTags()

    expect(fetch).toHaveBeenCalledWith(`${API_BASE_URL}/api/tags`)
    expect(result).toEqual(['cats', 'dogs'])
  })

  it('throws when the response is not ok', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(mockFetchResponse(null, false, 404))

    await expect(fetchTags()).rejects.toThrow('Failed to fetch tags: 404')
  })
})

describe('uploadImage', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
  })

  it('posts multipart form data with the file, title, and tags', async () => {
    const created = { id: '1', title: 'My title', tags: ['a'], url: '/uploads/1.jpg', created_at: 'now' }
    vi.mocked(fetch).mockResolvedValueOnce(mockFetchResponse(created, true, 201))
    const file = new File(['data'], 'photo.jpg', { type: 'image/jpeg' })

    const result = await uploadImage(file, 'My title', 'a,b')

    expect(fetch).toHaveBeenCalledTimes(1)
    const [url, init] = vi.mocked(fetch).mock.calls[0]
    expect(url).toBe(`${API_BASE_URL}/api/upload`)
    expect(init?.method).toBe('POST')
    const body = init?.body as FormData
    expect(body.get('file')).toBe(file)
    expect(body.get('title')).toBe('My title')
    expect(body.get('tags')).toBe('a,b')
    expect(result).toEqual(created)
  })

  it('throws the response body text when the upload fails', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(mockFetchResponse('title is required', false, 400))
    const file = new File(['data'], 'photo.jpg', { type: 'image/jpeg' })

    await expect(uploadImage(file, '', '')).rejects.toThrow('title is required')
  })
})
