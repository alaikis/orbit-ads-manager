'use client'

import { useQuery } from '@tanstack/react-query'
import { feedService } from '@/lib/api'
import { RefreshCw } from 'lucide-react'

export default function FeedsPage() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['feeds'],
    queryFn: () => feedService.list(),
  })

  const feeds = (data as any)?.items || []

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">Feed 管理</h1>
            <p className="text-sm text-text-muted mt-1">管理商品 Feed 生成与发布</p>
          </div>
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
        <div className="space-y-4">
          {[1, 2].map((i) => (<div key={i} className="card p-4 animate-pulse h-16 bg-surface-subtle rounded" />))}
        </div>
      ) : feeds.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted">尚无 Feed，请先绑定店铺后生成</p>
        </div>
      ) : (
        <div className="space-y-4">
          {feeds.map((feed: any) => (
            <div key={feed.id} className="card p-4">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium text-text-primary">{feed.name}</h3>
                  <p className="text-sm text-text-muted mt-1">格式: {feed.format} | 状态: {feed.status}</p>
                </div>
                <button className="btn btn-secondary text-sm">立即生成</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
