import { Feed } from './components/Feed'

function App() {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="p-4 border-b bg-white">
        <h1 className="text-2xl font-semibold text-gray-800">Anon Image Feed</h1>
      </header>
      <Feed />
    </div>
  )
}

export default App