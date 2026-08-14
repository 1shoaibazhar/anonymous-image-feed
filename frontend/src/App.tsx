import { useState } from 'react'
import { Feed } from './components/Feed'
import { UploadForm } from './components/UploadForm'
import { useLiveFeed } from './hooks/useLiveFeed'

function App() {
  const [refreshKey, setRefreshKey] = useState(0)

  useLiveFeed(() => setRefreshKey((k) => k + 1))

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="p-4 border-b bg-white">
        <h1 className="text-2xl font-semibold text-gray-800">Anon Image Feed</h1>
      </header>
      <UploadForm onUploadSuccess={() => setRefreshKey((k) => k + 1)} />
      <Feed refreshKey={refreshKey} />
    </div>
  )
}

export default App