'use client'

import { useQuery } from '@tanstack/react-query'
import { ruleService } from '@/lib/api'
import { ArrowLeft } from 'lucide-react'
import Link from 'next/link'

const STATUS_BADGE: Record<string, string> = {
  active: 'bg-success-bg text-success',
  inactive: 'bg-surface-subtle text-text-muted',
}

export default function RuleDetailPage({ params }: { params: { id: string } }) {
  const ruleId = Number(params.id)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['rule', ruleId],
    queryFn: () => ruleService.get(ruleId),
  })

  const rule = data as any

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/rules" className="btn btn-ghost"><ArrowLeft size={16} className="mr-2" />返回规则列表</Link>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">{(error as Error).message}</p>
          <button onClick={() => refetch()} className="btn btn-primary">重试</button>
        </div>
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/rules" className="btn btn-ghost"><ArrowLeft size={16} className="mr-2" />返回规则列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">规则详情</h1>
            <p className="text-sm text-text-muted mt-1">加载中...</p>
          </div>
        </div>
        <div className="card overflow-hidden animate-pulse">
          <div className="p-6 space-y-4">
            <div className="h-5 bg-surface-subtle rounded w-1/4" />
            <div className="space-y-3">
              <div className="h-4 bg-surface-subtle rounded" />
              <div className="h-4 bg-surface-subtle rounded w-5/6" />
              <div className="h-4 bg-surface-subtle rounded w-4/6" />
              <div className="h-4 bg-surface-subtle rounded" />
              <div className="h-4 bg-surface-subtle rounded w-3/4" />
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link href="/rules" className="btn btn-ghost"><ArrowLeft size={16} className="mr-2" />返回规则列表</Link>
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">{rule?.name || '规则详情'}</h1>
          <p className="text-sm text-text-muted mt-1">状态: {rule?.status || '-'}</p>
        </div>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <tbody>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted w-40">规则 ID</td>
              <td className="p-4 font-mono">{rule?.id ?? '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">规则名称</td>
              <td className="p-4 font-medium">{rule?.name || '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">状态</td>
              <td className="p-4">
                <span className={`badge ${STATUS_BADGE[rule?.status || ''] || 'bg-surface-subtle text-text-muted'}`}>
                  {rule?.status === 'active' ? '启用' : '禁用'}
                </span>
              </td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">条件</td>
              <td className="p-4">
                <pre className="font-mono text-xs bg-surface-subtle p-3 rounded-md overflow-auto">
                  {rule?.condition ? JSON.stringify(JSON.parse(rule.condition), null, 2) : '-'}
                </pre>
              </td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">动作</td>
              <td className="p-4">
                <pre className="font-mono text-xs bg-surface-subtle p-3 rounded-md overflow-auto">
                  {rule?.action ? JSON.stringify(JSON.parse(rule.action), null, 2) : '-'}
                </pre>
              </td>
            </tr>
            <tr>
              <td className="p-4 text-text-muted">创建时间</td>
              <td className="p-4">{rule?.created_at ? new Date(rule.created_at).toLocaleString('zh-CN') : '-'}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  )
}