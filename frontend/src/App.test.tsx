import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { App } from './App'

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

describe('App', () => {
  it('shows todos after successful fetch', async () => {
    globalThis.fetch = vi
      .fn()
      .mockResolvedValueOnce(
        createJSONResponse(200, [
          { id: 1, title: 'Buy milk', completed: false },
          { id: 2, title: 'Read Go docs', completed: true },
        ]),
      )

    render(<App />)

    expect(screen.getByText('Loading...')).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.getByText('Buy milk')).toBeInTheDocument()
      expect(screen.getByText('Read Go docs')).toBeInTheDocument()
    })
  })

  it('shows error when initial fetch fails', async () => {
    globalThis.fetch = vi.fn().mockResolvedValueOnce(createStatusResponse(500))

    render(<App />)

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(
        'TODO一覧の取得に失敗しました',
      )
    })
  })

  it('creates a todo', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(createJSONResponse(200, []))
      .mockResolvedValueOnce(
        createJSONResponse(201, {
          id: 3,
          title: 'Write tests',
          completed: false,
        }),
      )
    globalThis.fetch = fetchMock

    render(<App />)

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(1)
    })

    fireEvent.change(screen.getByLabelText('todo-title-input'), {
      target: { value: 'Write tests' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Add' }))

    await waitFor(() => {
      expect(screen.getByText('Write tests')).toBeInTheDocument()
    })

    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      'http://localhost:8080/api/todos',
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: 'Write tests' }),
      },
    )
  })

  it.each([
    { caseName: 'empty input', title: '' },
    { caseName: 'whitespace input', title: '   ' },
  ])(
    'shows validation error and does not send create request for $caseName',
    async ({ title }) => {
      const fetchMock = vi.fn().mockResolvedValueOnce(createJSONResponse(200, []))
      globalThis.fetch = fetchMock

      render(<App />)

      await waitFor(() => {
        expect(fetchMock).toHaveBeenCalledTimes(1)
      })

      fireEvent.change(screen.getByLabelText('todo-title-input'), {
        target: { value: title },
      })
      fireEvent.click(screen.getByRole('button', { name: 'Add' }))

      await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent(
          'TODOタイトルを入力してください',
        )
      })
      expect(fetchMock).toHaveBeenCalledTimes(1)
    },
  )

  it('toggles todo completion', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(
        createJSONResponse(200, [
          {
            id: 1,
            title: 'Write tests',
            completed: false,
          },
        ]),
      )
      .mockResolvedValueOnce(
        createJSONResponse(200, {
          id: 1,
          title: 'Write tests',
          completed: true,
        }),
      )
    globalThis.fetch = fetchMock

    render(<App />)

    await screen.findByText('Write tests')

    fireEvent.click(screen.getByLabelText('toggle-1'))

    await waitFor(() => {
      expect(screen.getByLabelText('toggle-1')).toBeChecked()
    })

    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      'http://localhost:8080/api/todos/1',
      {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ completed: true }),
      },
    )
  })

  it('deletes a todo', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(
        createJSONResponse(200, [{ id: 1, title: 'Write tests', completed: false }]),
      )
      .mockResolvedValueOnce(createStatusResponse(204))
    globalThis.fetch = fetchMock

    render(<App />)

    await screen.findByText('Write tests')

    fireEvent.click(screen.getByRole('button', { name: 'Delete' }))

    await waitFor(() => {
      expect(screen.queryByText('Write tests')).not.toBeInTheDocument()
    })

    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      'http://localhost:8080/api/todos/1',
      {
        method: 'DELETE',
      },
    )
  })

  it('shows error when create request fails', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(createJSONResponse(200, []))
      .mockResolvedValueOnce(createStatusResponse(500))
    globalThis.fetch = fetchMock

    render(<App />)

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(1)
    })

    fireEvent.change(screen.getByLabelText('todo-title-input'), {
      target: { value: 'Write tests' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Add' }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(
        'TODOの追加に失敗しました',
      )
    })
  })

  it('shows error when update request fails', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(
        createJSONResponse(200, [{ id: 1, title: 'Write tests', completed: false }]),
      )
      .mockResolvedValueOnce(createStatusResponse(500))
    globalThis.fetch = fetchMock

    render(<App />)

    await screen.findByText('Write tests')

    fireEvent.click(screen.getByLabelText('toggle-1'))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(
        'TODOの更新に失敗しました',
      )
    })

    expect(screen.getByLabelText('todo-list')).toBeInTheDocument()
    expect(screen.getByLabelText('toggle-1')).toBeInTheDocument()
    expect(screen.getByLabelText('toggle-1')).not.toBeChecked()
    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      'http://localhost:8080/api/todos/1',
      {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ completed: true }),
      },
    )
  })

  it('shows error when delete request fails', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(
        createJSONResponse(200, [{ id: 1, title: 'Write tests', completed: false }]),
      )
      .mockResolvedValueOnce(createStatusResponse(500))
    globalThis.fetch = fetchMock

    render(<App />)

    await screen.findByText('Write tests')

    fireEvent.click(screen.getByRole('button', { name: 'Delete' }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(
        'TODOの削除に失敗しました',
      )
    })

    expect(screen.getByLabelText('todo-list')).toBeInTheDocument()
    expect(screen.getByText('Write tests')).toBeInTheDocument()
    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      'http://localhost:8080/api/todos/1',
      {
        method: 'DELETE',
      },
    )
  })
})
