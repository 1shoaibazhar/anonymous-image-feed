import { useEffect, useState } from 'react'
import { fetchTags } from '../api/images'

interface TagFilterProps {
  selectedTags: string[]
  onChange: (tags: string[]) => void
  refreshKey: number
}

export function TagFilter({ selectedTags, onChange, refreshKey }: TagFilterProps) {
  const [allTags, setAllTags] = useState<string[]>([])

  useEffect(() => {
    fetchTags().then(setAllTags).catch(() => setAllTags([]))
  }, [refreshKey])

  if (allTags.length === 0) {
    return null
  }

  function toggle(tag: string) {
    if (selectedTags.includes(tag)) {
      onChange(selectedTags.filter((t) => t !== tag))
    } else {
      onChange([...selectedTags, tag])
    }
  }

  return (
    <div className="flex flex-wrap gap-2 p-4 pb-0">
      {allTags.map((tag) => {
        const active = selectedTags.includes(tag)
        return (
          <button
            key={tag}
            onClick={() => toggle(tag)}
            className={`text-xs px-3 py-1 rounded-full border ${
              active
                ? 'bg-gray-800 text-white border-gray-800'
                : 'bg-white text-gray-600 border-gray-300'
            }`}
          >
            {tag}
          </button>
        )
      })}
    </div>
  )
}
