'use client'

import { useQuery } from '@tanstack/react-query'
import { api, reportService } from '@/lib/api'
import { MetricCard } from '@/components/dashboard/metric-card'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

type ReportPoint = { date: string; spend: number; clicks: number; impressions: number }

export default function ReportsPage() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['reports'],
    queryFn: () => api.get('/reports/summary'),
  })

  const summary = data || {}
  const points: ReportPoint[] = []

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">数据报表</h1>
          <p className="text-sm text-text-muted mt-1">跨平台聚合报表与趋势分析</p>
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
        <h1 className="text-2xl font-semibold text-text-primary">数据报表</h1>
        <p className="text-sm text-text-muted mt-1">跨平台聚合报表与趋势分析</p>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="card p-4 animate-pulse">
              <div className="h-4 bg-surface-subtle rounded w-1/2 mb-2" />
              <div className="h-8 bg-surface-subtle rounded w-3/4" />
            </div>
          ))}
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
            <MetricCard title="总花费" value={summary.spend || 0} format="currency" />
            <MetricCard title="展示数" value={summary.impressions || 0} format="number" />
            <MetricCard title="点击数" value={summary.clicks || 0} format="number" />
            <MetricCard title="转化数" value={summary.conversions || 0} format="number" />
          </div>

          <div className="card p-6">
            <h2 className="font-title-md text-text-primary mb-4">趋势图</h2>
            {points.length === 0 ? (
              <div className="h-64 flex items-center justify-center text-text-muted">
                暂无趋势数据
              </div>
            ) : (
              <ResponsiveContainer width="100%" height={280}>
                <LineChart data={points}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
                  <XAxis dataKey="date" tick={{ fontSize: 12 }} stroke="#9CA3AF" />
                  <YAxis tick={{ fontSize: 12 }} stroke="#9CA3AF" />
                  <Tooltip
                    contentStyle={{ background: 'white', border: '1px solid #E5E7EB', borderRadius: 6, fontSize: 12 }}
                  />
                  <Line type="monotone" dataKey="spend" stroke="#6366F1" strokeWidth={2} dot={false} />
                  <Line type="monotone" dataKey="clicks" stroke="#10B981" strokeWidth={2} dot={false} />
                </LineChart>
              </ResponsiveContainer>
            )}
          </div>
        </>
      )}
    </div>
  )
}


