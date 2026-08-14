import { useState } from 'react'
import { Feed } from './components/Feed'
import { UploadForm } from './components/UploadForm'
import { TagFilter } from './components/TagFilter'
import { Modal } from './components/Modal'
import { useLiveFeed } from './hooks/useLiveFeed'

function App() {
  const [refreshKey, setRefreshKey] = useState(0)
  const [selectedTags, setSelectedTags] = useState<string[]>([])
  const [uploadOpen, setUploadOpen] = useState(false)

  useLiveFeed(() => setRefreshKey((k) => k + 1))

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="p-4 border-b bg-white flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-gray-800">Anon Image Feed</h1>
        <button
          onClick={() => setUploadOpen(true)}
          className="bg-gray-800 text-white px-4 py-1.5 rounded text-sm"
        >
          Upload
        </button>
      </header>
      <TagFilter selectedTags={selectedTags} onChange={setSelectedTags} refreshKey={refreshKey} />
      <Feed refreshKey={refreshKey} tags={selectedTags} />
      {uploadOpen && (
        <Modal onClose={() => setUploadOpen(false)}>
          <UploadForm
            onUploadSuccess={() => {
              setRefreshKey((k) => k + 1)
              setUploadOpen(false)
            }}
          />
        </Modal>
      )}
    </div>
  )
}

export default App