import { useEffect, type ReactNode } from 'react'

interface ModalProps {
  onClose: () => void
  children: ReactNode
  panelClassName?: string
}

export function Modal({ onClose, children, panelClassName = 'w-full max-w-md p-6' }: ModalProps) {
  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', handleKeyDown)
    return () => document.removeEventListener('keydown', handleKeyDown)
  }, [onClose])

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" onClick={onClose}>
      <div
        className={`bg-white rounded-lg shadow-lg relative ${panelClassName}`}
        onClick={(e) => e.stopPropagation()}
      >
        <button
          onClick={onClose}
          aria-label="Close"
          className="absolute top-3 right-3 bg-white/90 hover:bg-white text-gray-600 hover:text-gray-800 rounded-full w-8 h-8 flex items-center justify-center shadow text-xl leading-none"
        >
          ×
        </button>
        {children}
      </div>
    </div>
  )
}
