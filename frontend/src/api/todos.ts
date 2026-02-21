import type { TodoItem } from '../types'

type TodoCreateInput = {
  title: string
}

type TodoUpdateInput = {
  title: string
  completed?: boolean
} | {
  completed: boolean
  title?: string
}

function getApiBaseURL(): string {
  const baseURL = import.meta.env.VITE_API_BASE_URL
  if (!baseURL) {
    throw new Error('VITE_API_BASE_URL is required')
  }
  return baseURL
}

async function handleJSONResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    throw new Error(`request failed: ${response.status}`)
  }
  return (await response.json()) as T
}

export async function listTodos(): Promise<TodoItem[]> {
  const response = await fetch(`${getApiBaseURL()}/api/todos`, {
    method: 'GET',
  })
  return await handleJSONResponse<TodoItem[]>(response)
}

export async function createTodo(input: TodoCreateInput): Promise<TodoItem> {
  const response = await fetch(`${getApiBaseURL()}/api/todos`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })
  return await handleJSONResponse<TodoItem>(response)
}

export async function updateTodo(
  id: number,
  input: TodoUpdateInput,
): Promise<TodoItem> {
  const response = await fetch(`${getApiBaseURL()}/api/todos/${id}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  })
  return await handleJSONResponse<TodoItem>(response)
}

export async function deleteTodo(id: number): Promise<void> {
  const response = await fetch(`${getApiBaseURL()}/api/todos/${id}`, {
    method: 'DELETE',
  })

  if (!response.ok) {
    throw new Error(`request failed: ${response.status}`)
  }
}
