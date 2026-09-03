'use client'

import { useQuery } from '@tanstack/react-query'
import { workspaceService } from '@/lib/api'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { Package } from 'lucide-react'
import Link from 'next/link'

interface WorkspaceGuardProps {
  children: React.ReactNode
}

export function WorkspaceGuard({ children }: WorkspaceGuardProps) {
  const router = useRouter()
  const { data, isLoading, error } = useQuery({
    queryKey: ['workspace', 'current'],
    queryFn: () => workspaceService.getCurrent(),
    retry: false,
  })

  useEffect(() => {
    if (error) {
      // If 401 or 403, the layout will handle redirect to login
    }
  }, [error])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-text-muted text-sm">检查工作空间配置...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="card p-8 text-center">
        <p className="text-danger-500">加载工作空间失败</p>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="card p-8 max-w-2xl mx-auto mt-12">
        <div className="flex flex-col items-center text-center">
          <div className="w-16 h-16 rounded-full bg-primary-50 flex items-center justify-center mb-4">
            <Package size={28} className="text-primary-600" />
          </div>
          <h2 className="text-xl font-semibold text-text-primary mb-2">需要先配置工作空间</h2>
          <p className="text-text-muted mb-6 max-w-md">
            在管理店铺、广告账户和商品之前，请先创建一个工作空间来组织你的资源。
          </p>
          <Link
            href="/workspaces"
            className="btn btn-primary inline-flex items-center gap-2"
          >
            前往创建工作空间 →
          </Link>
        </div>
      </div>
    )
  }

  return <>{children}</>
}
