'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'

export default function AuthCallbackPage() {
  const router = useRouter()
  useEffect(() => {
    const platform = new URLSearchParams(window.location.search).get('platform')
    if (platform) {
      router.push(`/stores?connected=${platform}`)
    } else {
      router.push('/stores')
    }
  }, [router])
  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <div className="text-lg font-medium text-text-primary">正在连接...</div>
        <p className="text-sm text-text-muted mt-2">请稍候，正在完成授权绑定</p>
      </div>
    </div>
  )
}
