import { test, expect } from '@playwright/test'

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[type="email"]', 'test@example.com')
    await page.fill('input[type="password"]', 'password123')
    await page.click('button:has-text("登录")')
    await page.waitForURL('/')
  })

  test('should display dashboard', async ({ page }) => {
    await expect(page.locator('text=工作台')).toBeVisible()
  })

  test('should show metric cards', async ({ page }) => {
    await expect(page.locator('text=广告花费')).toBeVisible()
    await expect(page.locator('text=点击率 CTR')).toBeVisible()
    await expect(page.locator('text=转化成本 CPA')).toBeVisible()
    await expect(page.locator('text=广告支出回报率 ROAS')).toBeVisible()
  })

  test('should navigate to stores page', async ({ page }) => {
    await page.click('text=店铺')
    await expect(page).toHaveURL('/stores')
    await expect(page.locator('text=店铺管理')).toBeVisible()
  })
})
