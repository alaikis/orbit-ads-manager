'use client'

import { useQuery } from '@tanstack/react-query'
import { workspaceService } from '@/lib/api'
import { ArrowLeft } from 'lucide-react'
import Link from 'next/link'

type Workspace = {
  id: number
  name: string
  plan: string
  status: string
  member_count?: number
  created_at?: string
  updated_at?: string
}

const PLAN_BADGE: Record<string, string> = {
  beta: 'badge-info',
  pro: 'badge-primary',
  enterprise: 'badge-warning',
}

const STATUS_BADGE: Record<string, string> = {
  active: 'bg-success-bg text-success',
  inactive: 'bg-danger-bg text-danger',
}

export default function WorkspaceDetailPage({ params }: { params: { id: string } }) {
  const workspaceId = Number(params.id)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['workspace', workspaceId],
    queryFn: () => workspaceService.get(workspaceId),
  })

  const workspace = data as Workspace | undefined

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/workspaces" className="btn btn-ghost inline-flex items-center gap-2"><ArrowLeft size={18} />返回工作区列表</Link>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">{(error as Error).message}</p>
          <button onClick={() => refetch()} className="btn btn-primary">重试</button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link href="/workspaces" className="btn btn-ghost inline-flex items-center gap-2"><ArrowLeft size={18} />返回工作区列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">{workspace?.name || '工作区详情'}</h1>
            <p className="text-sm text-text-muted mt-1">ID: {workspace?.id || workspaceId}</p>
          </div>
        </div>
        <button onClick={() => refetch()} className="btn btn-secondary">刷新</button>
      </div>

      {isLoading ? (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <tbody>
              {[1, 2, 3, 4, 5, 6, 7].map((i) => (
                <tr key={i} className="animate-pulse">
                  <td className="px-6 py-4">
                    <div className="h-4 bg-surface-subtle rounded w-24" />
                  </td>
                  <td className="px-6 py-4">
                    <div className="h-4 bg-surface-subtle rounded w-48" />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <tbody className="divide-y divide-border-default">
              <tr className="hover:bg-surface-hover">
                <td className="px-6 py-3 text-text-muted font-medium w-40">工作区 ID</td>
                <td className="px-6 py-3 font-mono">{workspace?.id || '-'}</td>
              </tr>
              <tr className="hover:bg-surface-hover">
                <td className="px-6 py-3 text-text-muted font-medium w-40">工作区名称</td>
                <td className="px-6 py-3 font-medium text-text-primary">{workspace?.name || '-'}</td>
              </tr>
              <tr className="hover:bg-surface-hover">
                <td className="px-6 py-3 text-text-muted font-medium w-40">计划</td>
                <td className="px-6 py-3">
                  <span className={`badge ${PLAN_BADGE[workspace?.plan || ''] || 'bg-surface-subtle text-text-secondary'}`}>
                    {workspace?.plan || '-'}
                  </span>
                </td>
              </tr>
              <tr className="hover:bg-surface-hover">
                <td className="px-6 py-3 text-text-muted font-medium w-40">状态</td>
                <td className="px-6 py-3">
                  <span className={`badge ${STATUS_BADGE[workspace?.status || ''] || 'bg-surface-subtle text-text-secondary'}`}>
                    {workspace?.status || '-'}
                  </span>
                </td>
              </tr>
              <tr className="hover:bg-surface-hover">
                <td className="px-6 py-3 text-text-muted font-medium w-40">成员数量</td>
                <td className="px-6 py-3">{workspace?.member_count ?? '-'}</td>
              </tr>
              <tr className="hover:bg-surface-hover">
                <td className="px-6 py-3 text-text-muted font-medium w-40">创建时间</td>
                <td className="px-6 py-3">{workspace?.created_at ? new Date(workspace.created_at).toLocaleString('zh-CN') : '-'}</td>
              </tr>
              <tr className="hover:bg-surface-hover">
                <td className="px-6 py-3 text-text-muted font-medium w-40">更新时间</td>
                <td className="px-6 py-3">{workspace?.updated_at ? new Date(workspace.updated_at).toLocaleString('zh-CN') : '-'}</td>
              </tr>
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}