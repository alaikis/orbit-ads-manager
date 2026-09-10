'use client'

import { useQuery } from '@tanstack/react-query'
import { feedService } from '@/lib/api'
import { RefreshCw } from 'lucide-react'
import Link from 'next/link'

function formatDate(dateString: string) {
  if (!dateString) return '从未'
  return new Date(dateString).toLocaleString('zh-CN')
}

const STATUS_BADGE: Record<string, string> = {
  active: 'bg-success-bg text-success',
  paused: 'bg-warning-bg text-warning',
  error: 'bg-danger-bg text-danger',
}

export default function FeedsPage() {
  const { data, isLoading, error, refetch } = useQuery({ queryKey: ['feeds'], queryFn: () => feedService.list() })
  const feeds: any[] = (data as any)?.items || []

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">Feed 管理</h1>
          <p className="text-sm text-text-muted mt-1">管理商品 Feed 生成与发布</p>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">加载失败：{(error as Error).message}</p>
          <button onClick={() => refetch()} className="btn btn-primary">重试</button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">Feed 管理</h1>
          <p className="text-sm text-text-muted mt-1">管理商品 Feed 生成与发布</p>
        </div>
        <button onClick={() => refetch()} className="btn btn-secondary"><RefreshCw size={16} /></button>
      </div>
      {isLoading ? (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead><tr className="border-b border-border-default">{['名称','格式','状态','商品数量','最近同步'].map(h => (<th key={h} className="text-left p-4 font-title-sm text-text-secondary">{h}</th>))}</tr></thead>
            <tbody>{[1,2,3].map(i => (<tr key={i} className="border-b border-border-default last:border-0">{[1,2,3,4,5].map(j => (<td key={j} className="p-4"><div className="h-4 bg-surface-subtle rounded animate-pulse" style={{width: `${50 + Math.random()*40}%`}} /></td>))}</tr>))}</tbody>
          </table>
        </div>
      ) : feeds.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted mb-2">尚无 Feed</p>
          <p className="text-sm text-text-muted">请先绑定店铺后生成</p>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border-default">
                <th className="text-left p-4 font-title-sm text-text-secondary">名称</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">格式</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">状态</th>
                <th className="text-right p-4 font-title-sm text-text-secondary">商品数量</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">最近同步</th>
              </tr>
            </thead>
            <tbody>
              {feeds.map((feed: any) => (
                <tr key={feed.id} className="border-b border-border-default last:border-0 hover:bg-surface-hover">
                  <td className="p-4"><Link href={`/feeds/${feed.id}`} className="text-primary-600 hover:underline font-medium">{feed.name}</Link></td>
                  <td className="p-4 text-text-secondary">{feed.format}</td>
                  <td className="p-4"><span className={`badge ${STATUS_BADGE[feed.status] || 'bg-surface-subtle text-text-muted'}`}>{feed.status}</span></td>
                  <td className="p-4 text-right">{feed.product_count ?? '-'}</td>
                  <td className="p-4 text-text-muted">{formatDate(feed.last_synced_at)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
