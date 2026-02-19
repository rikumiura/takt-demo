import { defineConfig } from '@playwright/test'

const frontendPort = 5173
const backendPort = 8080

export default defineConfig({
  testDir: './tests/e2e',
  use: {
    baseURL: `http://127.0.0.1:${frontendPort}`,
    headless: true,
  },
  webServer: [
    {
      command: 'go run ./cmd/server -addr 127.0.0.1:8080 -db ./todo.e2e.db',
      cwd: '../backend',
      port: backendPort,
      reuseExistingServer: true,
    },
    {
      command: 'npm run dev -- --host 127.0.0.1 --port 5173',
      cwd: '.',
      port: frontendPort,
      reuseExistingServer: true,
      env: {
        VITE_API_BASE_URL: `http://127.0.0.1:${backendPort}`,
      },
    },
  ],
})
