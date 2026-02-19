import { expect, test } from '@playwright/test'

test('shows TODO list items', async ({ page }) => {
  await page.goto('/')

  await expect(page.getByRole('heading', { name: 'TODO List' })).toBeVisible()
  await expect(page.getByText('Buy milk')).toBeVisible()
  await expect(page.getByText('Read Go docs')).toBeVisible()
})
