import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createTodo, deleteTodo, listTodos, updateTodo } from './todos'

const originalFetch = globalThis.fetch

function createJSONResponse(status: number, body: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response
}

function createStatusResponse(status: number): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
  } as Response
}

beforeEach(() => {
  vi.stubEnv('VITE_API_BASE_URL', 'http://localhost:8080')
})

afterEach(() => {
  vi.restoreAllMocks()
  vi.unstubAllEnvs()
  globalThis.fetch = originalFetch
})

describe('todo api', () => {
  it('lists todos', async () => {
    globalThis.fetch = vi
      .fn()
      .mockResolvedValueOnce(
        createJSONResponse(200, [{ id: 1, title: 'Buy milk', completed: false }]),
      )

    const todos = await listTodos()

    expect(todos).toEqual([{ id: 1, title: 'Buy milk', completed: false }])
    expect(globalThis.fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/todos',
      { method: 'GET' },
    )
  })

  it('creates todo', async () => {
    globalThis.fetch = vi
      .fn()
      .mockResolvedValueOnce(
        createJSONResponse(201, { id: 2, title: 'Read docs', completed: false }),
      )

    const todo = await createTodo({ title: 'Read docs' })

    expect(todo).toEqual({ id: 2, title: 'Read docs', completed: false })
    expect(globalThis.fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/todos',
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: 'Read docs' }),
      },
    )
  })

  it('updates todo completed', async () => {
    globalThis.fetch = vi
      .fn()
      .mockResolvedValueOnce(
        createJSONResponse(200, { id: 2, title: 'Read docs', completed: true }),
      )

    const todo = await updateTodo(2, { completed: true })

    expect(todo).toEqual({ id: 2, title: 'Read docs', completed: true })
    expect(globalThis.fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/todos/2',
      {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ completed: true }),
      },
    )
  })

  it('updates todo title', async () => {
    globalThis.fetch = vi
      .fn()
      .mockResolvedValueOnce(
        createJSONResponse(200, { id: 2, title: 'Updated docs', completed: false }),
      )

    const todo = await updateTodo(2, { title: 'Updated docs' })

    expect(todo).toEqual({ id: 2, title: 'Updated docs', completed: false })
    expect(globalThis.fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/todos/2',
      {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: 'Updated docs' }),
      },
    )
  })

  it('deletes todo', async () => {
    globalThis.fetch = vi.fn().mockResolvedValueOnce(createStatusResponse(204))

    await expect(deleteTodo(2)).resolves.toBeUndefined()
    expect(globalThis.fetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/todos/2',
      { method: 'DELETE' },
    )
  })

  it('fails fast when base URL is not configured', async () => {
    vi.unstubAllEnvs()

    await expect(listTodos()).rejects.toThrow('VITE_API_BASE_URL is required')
  })
})
