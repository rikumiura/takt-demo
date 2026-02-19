import { render, screen, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { App } from './App'

const originalFetch = global.fetch

afterEach(() => {
  vi.restoreAllMocks()
  global.fetch = originalFetch
})

describe('App', () => {
  it('shows todos after successful fetch', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => [
        { id: 1, title: 'Buy milk' },
        { id: 2, title: 'Read Go docs' },
      ],
    } as Response)

    render(<App />)

    expect(screen.getByText('Loading...')).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.getByText('Buy milk')).toBeInTheDocument()
      expect(screen.getByText('Read Go docs')).toBeInTheDocument()
    })
  })

  it('shows error when fetch fails', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 500,
    } as Response)

    render(<App />)

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(
        'TODO一覧の取得に失敗しました',
      )
    })
  })
})
