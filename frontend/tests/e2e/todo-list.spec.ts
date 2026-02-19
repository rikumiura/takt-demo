import { expect, test } from '@playwright/test'

test('performs TODO CRUD flow', async ({ page }) => {
  const title = `E2E todo ${Date.now()}`

  await page.goto('/')

  await expect(page.getByRole('heading', { name: 'TODO List' })).toBeVisible()

  await page.getByLabel('todo-title-input').fill(title)
  await page.getByRole('button', { name: 'Add' }).click()

  const createdItem = page.getByRole('listitem').filter({ hasText: title })
  await expect(createdItem).toBeVisible()

  const toggle = createdItem.getByRole('checkbox')
  await toggle.click()
  await expect(toggle).toBeChecked()

  await createdItem.getByRole('button', { name: 'Delete' }).click()
  await expect(createdItem).toHaveCount(0)
})
