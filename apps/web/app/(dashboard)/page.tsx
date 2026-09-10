'use client'

import { useQuery } from '@tanstack/react-query'
import { metricsService, workspaceService } from '@/lib/api'
import { MetricCard } from '@/components/dashboard/metric-card'
import { QuickActions } from '@/components/dashboard/quick-actions'

export default function DashboardPage() {
  const { data: workspaceData } = useQuery({
    queryKey: ['workspace', 'current'],
    queryFn: () => workspaceService.getCurrent(),
  })

  const workspaceId = workspaceData?.id

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['dashboard', workspaceId],
    queryFn: () => metricsService.getSummary(workspaceId ? { workspace_id: workspaceId } : undefined),
  })

  const metrics = data || {}

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">工作台</h1>
          <p className="text-sm text-text-muted mt-1">欢迎回来，查看你的广告投放概况</p>
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
      <div>
        <h1 className="text-2xl font-semibold text-text-primary">工作台</h1>
        <p className="text-sm text-text-muted mt-1">欢迎回来，查看你的广告投放概况</p>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="card p-6 animate-pulse">
              <div className="h-4 bg-surface-subtle rounded w-1/2 mb-2" />
              <div className="h-8 bg-surface-subtle rounded w-3/4" />
            </div>
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
          <MetricCard title="广告花费" value={metrics.spend || 0} format="currency" change={-2.3} />
          <MetricCard title="点击率 CTR" value={metrics.ctr || 0} format="percent" change={1.2} />
          <MetricCard title="转化成本 CPA" value={metrics.cpa || 0} format="currency" change={-5.1} />
          <MetricCard title="广告支出回报率 ROAS" value={metrics.roas || 0} format="ratio" change={3.4} />
        </div>
      )}

      <QuickActions />

      <div className="card p-6">
        <h2 className="font-title-md text-text-primary mb-4">同步状态</h2>
        <div className="flex items-center gap-2">
          <span className="w-2.5 h-2.5 rounded-full bg-success-500" />
          <span className="text-sm text-text-secondary">全部同步正常</span>
        </div>
      </div>
    </div>
  )
}
