import { useEffect, useRef } from 'react'

const WS_URL = 'ws://localhost:8080/ws'

interface WSMessage {
  type: string
  payload: unknown
}

export function useLiveFeed(onImageCreated: () => void) {
  const callbackRef = useRef(onImageCreated)
  callbackRef.current = onImageCreated

  useEffect(() => {
    let cancelled = false
    const socket = new WebSocket(WS_URL)

    socket.onopen = () => {
      if (cancelled) socket.close()
    }

    socket.onmessage = (event) => {
      try {
        const message: WSMessage = JSON.parse(event.data)
        if (message.type === 'image.created') {
          callbackRef.current()
        }
      } catch {
        // ignore malformed messages
      }
    }

    socket.onerror = (event) => {
      console.error('WebSocket error', event)
    }

    return () => {
      cancelled = true
      if (socket.readyState === WebSocket.OPEN) {
        socket.close()
      }
    }
  }, [])
}