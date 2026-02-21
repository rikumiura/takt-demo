import { expect, test } from '@playwright/test'

test('performs TODO CRUD flow', async ({ page }) => {
  const baseTitle = `E2E todo ${Date.now()}`
  const firstTitle = `${baseTitle} A`
  const secondTitle = `${baseTitle} B`
  const editedFirstTitle = `${firstTitle} edited`

  await page.goto('/')

  await expect(page.getByRole('heading', { name: 'TODO List' })).toBeVisible()

  await page.getByLabel('todo-title-input').fill(firstTitle)
  await page.getByRole('button', { name: 'Add' }).click()
  const firstItem = page.getByRole('listitem').filter({ hasText: firstTitle })
  await expect(firstItem).toBeVisible()

  await page.getByLabel('todo-title-input').fill(secondTitle)
  await page.getByRole('button', { name: 'Add' }).click()
  const secondItem = page.getByRole('listitem').filter({ hasText: secondTitle })
  await expect(secondItem).toBeVisible()

  await firstItem.getByRole('button', { name: 'Edit' }).click()
  await page.getByLabel(/edit-title-input-/).fill(editedFirstTitle)
  await page.getByRole('button', { name: 'Save' }).click()

  const editedFirstItem = page.getByRole('listitem').filter({ hasText: editedFirstTitle })
  await expect(editedFirstItem).toBeVisible()

  const toggle = editedFirstItem.getByRole('checkbox')
  await toggle.click()
  await expect(toggle).toBeChecked()

  await secondItem.getByRole('button', { name: 'Delete' }).click()
  await expect(secondItem).toHaveCount(0)

  await page.reload()

  const persistedEditedItem = page.getByRole('listitem').filter({ hasText: editedFirstTitle })
  await expect(persistedEditedItem).toBeVisible()
  await expect(persistedEditedItem.getByRole('checkbox')).toBeChecked()
  await expect(page.getByRole('listitem').filter({ hasText: secondTitle })).toHaveCount(0)
})
