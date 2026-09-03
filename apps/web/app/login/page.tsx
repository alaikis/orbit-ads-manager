'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { authService } from '@/lib/api'
import { useAuthStore } from '@/lib/auth-store'

export default function LoginPage() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      const res = await authService.login({ email, password })
      useAuthStore.getState().setAuth(res.user, res.access_token, res.refresh_token)
      // Wait a tick to ensure localStorage is written
      await new Promise(resolve => setTimeout(resolve, 50))
      router.push('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen">
      <div className="hidden lg:flex lg:w-1/2 bg-primary-600 items-center justify-center">
        <div className="text-white text-center px-12">
          <h1 className="text-4xl font-bold mb-4">Orbit</h1>
          <p className="text-xl opacity-90">广告智能中枢</p>
          <p className="mt-4 text-sm opacity-75">AI Agent 驱动的多平台广告自动化管理平台</p>
        </div>
      </div>
      <div className="flex w-full lg:w-1/2 items-center justify-center bg-surface-page px-6">
        <div className="w-full max-w-[400px]">
          <div className="lg:hidden text-center mb-8">
            <h1 className="text-3xl font-bold text-text-primary">Orbit</h1>
            <p className="text-text-secondary mt-2">广告智能中枢</p>
          </div>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-text-secondary mb-1.5">
                邮箱
              </label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input"
                placeholder="you@example.com"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-text-secondary mb-1.5">
                密码
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="input"
                placeholder="••••••••"
                required
              />
            </div>
            {error && (
              <div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">
                {error}
              </div>
            )}
            <button type="submit" className="btn btn-primary w-full" disabled={loading}>
              {loading ? '登录中...' : '登录'}
            </button>
            <p className="text-center text-sm text-text-muted">
              还没有账户？{' '}
              <Link href="/register" className="text-primary-600 hover:underline">
                注册
              </Link>
            </p>
          </form>
        </div>
      </div>
    </div>
  )
}
