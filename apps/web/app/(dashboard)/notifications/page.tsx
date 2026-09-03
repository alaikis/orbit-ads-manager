'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { notificationService } from '@/lib/api'
import { useState } from 'react'

type Notification = { id: number; type: string; title: string; body: string; read_at?: string; created_at: string }

export default function NotificationsPage() {
  const [filter, setFilter] = useState<'all' | 'unread'>('all')
  const queryClient = useQueryClient()
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['notifications', filter],
    queryFn: () => notificationService.list({ unread_only: filter === 'unread' }),
  })

  const notifications: Notification[] = (data as any)?.items || []

  const markAllReadMutation = useMutation({
    mutationFn: () => notificationService.markAllAsRead(),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['notifications'] }),
    onError: (err: unknown) => alert('操作失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">通知中心</h1>
            <p className="text-sm text-text-muted mt-1">系统通知与事件提醒</p>
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
          <h1 className="text-2xl font-semibold text-text-primary">通知中心</h1>
          <p className="text-sm text-text-muted mt-1">系统通知与事件提醒</p>
        </div>
        <div className="flex gap-2">
          <button onClick={() => setFilter('all')} className={`btn ${filter === 'all' ? 'btn-primary' : 'btn-secondary'}`}>全部</button>
          <button onClick={() => setFilter('unread')} className={`btn ${filter === 'unread' ? 'btn-primary' : 'btn-secondary'}`}>未读</button>
          <button onClick={() => markAllReadMutation.mutate()} disabled={markAllReadMutation.isPending} className="btn btn-secondary">
            {markAllReadMutation.isPending ? '处理中...' : '全部已读'}
          </button>
        </div>
      </div>
      {isLoading ? (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (<div key={i} className="card p-4 animate-pulse h-16 bg-surface-subtle rounded" />))}
        </div>
      ) : notifications.length === 0 ? (
        <div className="card p-12 text-center"><p className="text-text-muted">暂无通知</p></div>
      ) : (
        <div className="space-y-3">
          {notifications.map((n) => (
            <div key={n.id} className={`card p-4 ${!n.read_at ? 'border-l-4 border-l-primary-500' : ''}`}>
              <div className="flex items-start justify-between">
                <div>
                  <h3 className="font-medium text-text-primary">{n.title}</h3>
                  <p className="text-sm text-text-secondary mt-1">{n.body}</p>
                  <p className="text-xs text-text-muted mt-2">{n.created_at}</p>
                </div>
                {!n.read_at && <span className="w-2 h-2 rounded-full bg-primary-500 mt-1" />}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
