'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/lib/api'
import { useState } from 'react'

type Action = { id: number; action_type: string; target_type: string; target_id?: string; risk_level: string; status: string; params: any }

export default function AgentCardsPage() {
  const [tab, setTab] = useState<'pending' | 'approved' | 'rejected'>('pending')
  const [actionError, setActionError] = useState<string | null>(null)
  const queryClient = useQueryClient()
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['agent-actions', tab],
    queryFn: () => apiClient.get<Action[] | { items: Action[] } | { data: Action[] }>(`/agent/actions?status=${tab}`),
  })

  const actions: Action[] = Array.isArray(data) ? data : (data as any)?.items || (data as any)?.data || []

  const approveMutation = useMutation({
    mutationFn: (id: number) => apiClient.post(`/agent/actions/${id}/approve`, {}),
    onSuccess: () => {
      setActionError(null)
      queryClient.invalidateQueries({ queryKey: ['agent-actions', tab] })
    },
    onError: (err: unknown) => setActionError(err instanceof Error ? err.message : '操作失败'),
  })

  const rejectMutation = useMutation({
    mutationFn: (id: number) => apiClient.post(`/agent/actions/${id}/reject`, {}),
    onSuccess: () => {
      setActionError(null)
      queryClient.invalidateQueries({ queryKey: ['agent-actions', tab] })
    },
    onError: (err: unknown) => setActionError(err instanceof Error ? err.message : '操作失败'),
  })

  const getRiskBadge = (risk: string) => {
    switch (risk) {
      case 'high': return 'bg-danger-bg text-danger'
      case 'medium': return 'bg-warning-bg text-warning'
      default: return 'bg-info-bg text-info'
    }
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">待确认动作中心</h1>
          <p className="text-sm text-text-muted mt-1">查看与管理所有待确认的 Agent 动作</p>
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
          <h1 className="text-2xl font-semibold text-text-primary">待确认动作中心</h1>
          <p className="text-sm text-text-muted mt-1">查看与管理所有待确认的 Agent 动作</p>
        </div>
        <button onClick={() => refetch()} className="btn btn-secondary">刷新</button>
      </div>
      {actionError && (
        <div className="card p-3 bg-danger-bg text-danger-500 text-sm">{actionError}</div>
      )}
      <div className="flex gap-2">
        {(['pending', 'approved', 'rejected'] as const).map((t) => (
          <button key={t} onClick={() => setTab(t)} className={`btn ${tab === t ? 'btn-primary' : 'btn-secondary'}`}>
            {t === 'pending' ? '待处理' : t === 'approved' ? '已确认' : '已拒绝'}
          </button>
        ))}
      </div>
      {isLoading ? (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (<div key={i} className="card p-4 animate-pulse h-16 bg-surface-subtle rounded" />))}
        </div>
      ) : actions.length === 0 ? (
        <div className="card p-12 text-center"><p className="text-text-muted">暂无动作</p></div>
      ) : (
        <div className="space-y-3">
          {actions.map((action: Action) => (
            <div key={action.id} className="card p-4">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium text-text-primary">{action.action_type}</h3>
                  <p className="text-sm text-text-muted mt-1">目标: {action.target_type} {action.target_id || ''}</p>
                  <pre className="text-xs text-text-muted mt-2 overflow-auto">{JSON.stringify(action.params)}</pre>
                </div>
                <div className="flex items-center gap-3">
                  <span className={`badge ${getRiskBadge(action.risk_level)}`}>{action.risk_level}</span>
                  <span className={`badge ${action.status === 'pending' ? 'bg-warning-bg text-warning' : 'bg-surface-subtle text-text-muted'}`}>
                    {action.status}
                  </span>
                </div>
              </div>
              {action.status === 'pending' && (
                <div className="mt-4 flex gap-2">
                  <button onClick={() => approveMutation.mutate(action.id)} disabled={approveMutation.isPending} className="btn btn-primary text-sm">
                    {approveMutation.isPending ? '处理中...' : '确认'}
                  </button>
                  <button onClick={() => rejectMutation.mutate(action.id)} disabled={rejectMutation.isPending} className="btn btn-secondary text-sm">
                    {rejectMutation.isPending ? '处理中...' : '拒绝'}
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
