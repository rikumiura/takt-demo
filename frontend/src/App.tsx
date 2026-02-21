import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { createTodo, deleteTodo, listTodos, updateTodo } from './api/todos'
import type { TodoItem } from './types'

export function App() {
  const [items, setItems] = useState<TodoItem[]>([])
  const [loading, setLoading] = useState(true)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [mutationError, setMutationError] = useState<string | null>(null)
  const [newTitle, setNewTitle] = useState('')
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editingTitle, setEditingTitle] = useState('')

  useEffect(() => {
    void loadTodos()
  }, [])

  async function loadTodos() {
    setLoading(true)
    try {
      const data = await listTodos()
      setItems(data)
      setLoadError(null)
    } catch {
      setLoadError('TODO一覧の取得に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  async function handleCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const title = newTitle.trim()
    if (title.length === 0) {
      setMutationError('TODOタイトルを入力してください')
      return
    }

    try {
      const created = await createTodo({ title })
      setItems((prevItems) => [...prevItems, created])
      setNewTitle('')
      setMutationError(null)
    } catch {
      setMutationError('TODOの追加に失敗しました')
    }
  }

  async function handleToggle(item: TodoItem) {
    try {
      const updated = await updateTodo(item.id, { completed: !item.completed })
      setItems((prevItems) =>
        prevItems.map((currentItem) =>
          currentItem.id === item.id ? updated : currentItem,
        ),
      )
      setMutationError(null)
    } catch {
      setMutationError('TODOの更新に失敗しました')
    }
  }

  function startEditing(item: TodoItem) {
    setEditingId(item.id)
    setEditingTitle(item.title)
    setMutationError(null)
  }

  function cancelEditing() {
    setEditingId(null)
    setEditingTitle('')
  }

  async function handleSaveEdit(id: number) {
    const title = editingTitle.trim()
    if (title.length === 0) {
      setMutationError('TODOタイトルを入力してください')
      return
    }

    try {
      const updated = await updateTodo(id, { title })
      setItems((prevItems) =>
        prevItems.map((currentItem) =>
          currentItem.id === id ? updated : currentItem,
        ),
      )
      setEditingId(null)
      setEditingTitle('')
      setMutationError(null)
    } catch {
      setMutationError('TODOの更新に失敗しました')
    }
  }

  async function handleDelete(id: number) {
    try {
      await deleteTodo(id)
      setItems((prevItems) =>
        prevItems.filter((currentItem) => currentItem.id !== id),
      )
      if (editingId === id) {
        setEditingId(null)
        setEditingTitle('')
      }
      setMutationError(null)
    } catch {
      setMutationError('TODOの削除に失敗しました')
    }
  }

  return (
    <main className="container">
      <h1>TODO List</h1>

      <form aria-label="todo-create-form" onSubmit={(event) => void handleCreate(event)}>
        <input
          aria-label="todo-title-input"
          value={newTitle}
          onChange={(event) => setNewTitle(event.target.value)}
          placeholder="Add a task"
        />
        <button type="submit">Add</button>
      </form>

      {loading && <p>Loading...</p>}
      {loadError && <p role="alert">{loadError}</p>}
      {mutationError && <p role="alert">{mutationError}</p>}
      {!loading && !loadError && (
        <ul aria-label="todo-list">
          {items.map((item) => {
            const isEditing = editingId === item.id
            return (
              <li key={item.id} className="todo-row">
                {isEditing ? (
                  <input
                    aria-label={`edit-title-input-${item.id}`}
                    value={editingTitle}
                    onChange={(event) => setEditingTitle(event.target.value)}
                  />
                ) : (
                  <label>
                    <input
                      aria-label={`toggle-${item.id}`}
                      type="checkbox"
                      checked={item.completed}
                      onChange={() => void handleToggle(item)}
                    />
                    <span className={item.completed ? 'todo-completed' : ''}>
                      {item.title}
                    </span>
                  </label>
                )}
                {isEditing ? (
                  <div>
                    <button
                      type="button"
                      onClick={() => void handleSaveEdit(item.id)}
                    >
                      Save
                    </button>
                    <button type="button" onClick={cancelEditing}>
                      Cancel
                    </button>
                  </div>
                ) : (
                  <div>
                    <button type="button" onClick={() => startEditing(item)}>
                      Edit
                    </button>
                    <button
                      type="button"
                      onClick={() => void handleDelete(item.id)}
                    >
                      Delete
                    </button>
                  </div>
                )}
              </li>
            )
          })}
        </ul>
      )}
    </main>
  )
}
