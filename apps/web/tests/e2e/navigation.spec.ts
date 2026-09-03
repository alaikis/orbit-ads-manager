import { test, expect } from '@playwright/test'

test.describe('Navigation', () => {
  test('should navigate to advertising campaigns', async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[type="email"]', 'test@example.com')
    await page.fill('input[type="password"]', 'password123')
    await page.click('button:has-text("登录")')
    await page.waitForURL('/')

    await page.click('text=广告管理')
    await expect(page).toHaveURL('/advertising/campaigns')
    await expect(page.locator('text=广告管理')).toBeVisible()
  })

  test('should navigate to agent page', async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[type="email"]', 'test@example.com')
    await page.fill('input[type="password"]', 'password123')
    await page.click('button:has-text("登录")')
    await page.waitForURL('/')

    await page.click('text=Agent 对话')
    await expect(page).toHaveURL('/agent')
  })

  test('should navigate to rules page', async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[type="email"]', 'test@example.com')
    await page.fill('input[type="password"]', 'password123')
    await page.click('button:has-text("登录")')
    await page.waitForURL('/')

    await page.click('text=规则管理')
    await expect(page).toHaveURL('/rules')
    await expect(page.locator('text=规则管理')).toBeVisible()
  })
})
