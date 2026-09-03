'use client'

import { useAuthStore } from '@/lib/auth-store'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { AppShell } from '@/components/layout/app-shell'

export default function ProtectedLayout({ children }: { children: React.ReactNode }) {
  const user = useAuthStore((s: any) => s.user)
  const router = useRouter()
  const [hydrated, setHydrated] = useState(false)

  useEffect(() => {
    // Wait for client-side hydration to read localStorage
    setHydrated(true)
  }, [])

  useEffect(() => {
    if (hydrated && !user) {
      router.push('/login')
    }
  }, [hydrated, user, router])

  if (!hydrated) return null
  if (!user) return null

  return <AppShell>{children}</AppShell>
}
