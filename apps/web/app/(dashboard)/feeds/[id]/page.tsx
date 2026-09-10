'use client'

import { useQuery } from '@tanstack/react-query'
import { feedService } from '@/lib/api'
import { ArrowLeft } from 'lucide-react'
import Link from 'next/link'

const FORMAT_LABEL: Record<string, string> = {
  google_shopping_xml: 'Google Shopping XML',
  google_shopping_tsv: 'Google Shopping TSV',
  meta_catalog_csv: 'Meta Catalog CSV',
}

const STATUS_BADGE: Record<string, string> = {
  active: 'bg-success-bg text-success',
  paused: 'bg-warning-bg text-warning',
  error: 'bg-danger-bg text-danger',
}

export default function FeedDetailPage({ params }: { params: { id: string } }) {
  const feedId = Number(params.id)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['feed', feedId],
    queryFn: () => feedService.get(feedId),
  })

  const feed = data as any

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/feeds" className="btn btn-ghost"><ArrowLeft size={16} className="mr-2" />返回 Feed 列表</Link>
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
          <Link href="/feeds" className="btn btn-ghost"><ArrowLeft size={16} className="mr-2" />返回 Feed 列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">Feed 详情</h1>
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
        <Link href="/feeds" className="btn btn-ghost"><ArrowLeft size={16} className="mr-2" />返回 Feed 列表</Link>
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">{feed?.name || 'Feed 详情'}</h1>
          <p className="text-sm text-text-muted mt-1">格式: {FORMAT_LABEL[feed?.format || ''] || feed?.format || '-'} | 状态: {feed?.status || '-'}</p>
        </div>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <tbody>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted w-40">Feed ID</td>
              <td className="p-4 font-mono">{feed?.id ?? '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">Feed 名称</td>
              <td className="p-4">{feed?.name || '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">格式</td>
              <td className="p-4">{FORMAT_LABEL[feed?.format || ''] || feed?.format || '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">状态</td>
              <td className="p-4">
                <span className={`badge ${STATUS_BADGE[feed?.status || ''] || 'bg-surface-subtle text-text-muted'}`}>
                  {feed?.status || '-'}
                </span>
              </td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">商品数量</td>
              <td className="p-4">{feed?.product_count ?? '-'}</td>
            </tr>
            <tr>
              <td className="p-4 text-text-muted">最近同步时间</td>
              <td className="p-4">{feed?.last_synced_at || '-'}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  )
}