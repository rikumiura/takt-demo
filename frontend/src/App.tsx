import { useEffect, useState } from 'react'
import type { TodoItem } from './types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

export function App() {
  const [items, setItems] = useState<TodoItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    void fetch(`${API_BASE_URL}/api/todos`)
      .then(async (res) => {
        if (!res.ok) {
          throw new Error(`request failed: ${res.status}`)
        }
        return (await res.json()) as TodoItem[]
      })
      .then((data) => {
        setItems(data)
      })
      .catch(() => {
        setError('TODO一覧の取得に失敗しました')
      })
      .finally(() => {
        setLoading(false)
      })
  }, [])

  return (
    <main className="container">
      <h1>TODO List</h1>
      {loading && <p>Loading...</p>}
      {error && <p role="alert">{error}</p>}
      {!loading && !error && (
        <ul aria-label="todo-list">
          {items.map((item) => (
            <li key={item.id}>{item.title}</li>
          ))}
        </ul>
      )}
    </main>
  )
}
