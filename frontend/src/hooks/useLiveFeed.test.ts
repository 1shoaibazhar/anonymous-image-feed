import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook } from '@testing-library/react'
import { useLiveFeed } from './useLiveFeed'

class FakeWebSocket {
  static instances: FakeWebSocket[] = []
  static OPEN = 1
  static CONNECTING = 0

  url: string
  readyState = FakeWebSocket.CONNECTING
  onopen: (() => void) | null = null
  onmessage: ((event: { data: string }) => void) | null = null
  onerror: ((event: unknown) => void) | null = null
  close = vi.fn(() => {
    this.readyState = 3
  })

  constructor(url: string) {
    this.url = url
    FakeWebSocket.instances.push(this)
  }

  emitMessage(data: unknown) {
    this.onmessage?.({ data: JSON.stringify(data) })
  }
}

describe('useLiveFeed', () => {
  beforeEach(() => {
    FakeWebSocket.instances = []
    vi.stubGlobal('WebSocket', FakeWebSocket)
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('connects to the websocket URL', () => {
    renderHook(() => useLiveFeed(vi.fn()))

    expect(FakeWebSocket.instances).toHaveLength(1)
    expect(FakeWebSocket.instances[0].url).toBe('ws://localhost:8080/ws')
  })

  it('calls the callback when an image.created message arrives', () => {
    const onImageCreated = vi.fn()
    renderHook(() => useLiveFeed(onImageCreated))

    FakeWebSocket.instances[0].emitMessage({ type: 'image.created', payload: {} })

    expect(onImageCreated).toHaveBeenCalledTimes(1)
  })

  it('ignores messages of other types', () => {
    const onImageCreated = vi.fn()
    renderHook(() => useLiveFeed(onImageCreated))

    FakeWebSocket.instances[0].emitMessage({ type: 'something.else', payload: {} })

    expect(onImageCreated).not.toHaveBeenCalled()
  })

  it('ignores malformed JSON without throwing', () => {
    const onImageCreated = vi.fn()
    renderHook(() => useLiveFeed(onImageCreated))
    const socket = FakeWebSocket.instances[0]

    expect(() => socket.onmessage?.({ data: 'not json' })).not.toThrow()
    expect(onImageCreated).not.toHaveBeenCalled()
  })

  it('always calls the latest callback even if it changes between renders', () => {
    const firstCallback = vi.fn()
    const secondCallback = vi.fn()
    const { rerender } = renderHook(({ cb }) => useLiveFeed(cb), {
      initialProps: { cb: firstCallback },
    })

    rerender({ cb: secondCallback })
    FakeWebSocket.instances[0].emitMessage({ type: 'image.created', payload: {} })

    expect(firstCallback).not.toHaveBeenCalled()
    expect(secondCallback).toHaveBeenCalledTimes(1)
  })

  it('closes the socket on unmount if it is open', () => {
    const { unmount } = renderHook(() => useLiveFeed(vi.fn()))
    const socket = FakeWebSocket.instances[0]
    socket.readyState = FakeWebSocket.OPEN

    unmount()

    expect(socket.close).toHaveBeenCalledTimes(1)
  })
})
