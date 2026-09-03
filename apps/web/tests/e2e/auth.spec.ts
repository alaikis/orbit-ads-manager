import { test, expect } from '@playwright/test'

test.describe('Authentication', () => {
  test('should display login page', async ({ page }) => {
    await page.goto('/login')
    await expect(page.locator('text=Orbit')).toBeVisible()
    await expect(page.locator('text=广告智能中枢')).toBeVisible()
  })

  test('should have login form elements', async ({ page }) => {
    await page.goto('/login')
    await expect(page.locator('input[type="email"]')).toBeVisible()
    await expect(page.locator('input[type="password"]')).toBeVisible()
    await expect(page.locator('button:has-text("登录")')).toBeVisible()
  })

  test('should navigate to register page', async ({ page }) => {
    await page.goto('/login')
    await page.click('text=注册')
    await expect(page).toHaveURL('/register')
    await expect(page.locator('text=创建账户')).toBeVisible()
  })
})
