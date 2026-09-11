'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient, reportService, workspaceService } from '@/lib/api'
import { MetricCard } from '@/components/dashboard/metric-card'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { useState } from 'react'
import Link from 'next/link'

type ReportPoint = { date: string; spend: number; clicks: number; impressions: number; conversions?: number }

export default function ReportsPage() {
  const queryClient = useQueryClient()
  const { data: workspaceData } = useQuery({
    queryKey: ['workspace', 'current'],
    queryFn: () => workspaceService.getCurrent(),
  })

  const workspaceId = workspaceData?.id
  const [dateRange, setDateRange] = useState('last_30d')
  const [scopeType, setScopeType] = useState('account')
  const [scopeId, setScopeId] = useState('')

  const summaryQuery = useQuery({
    queryKey: ['reports', 'summary', workspaceId, dateRange, scopeType, scopeId],
    queryFn: () => apiClient.get<any>('/reports/summary', {
      workspace_id: workspaceId ? String(workspaceId) : undefined,
      date_range: dateRange,
      scope_type: scopeType,
      scope_id: scopeId || undefined,
    }),
  })

  const timeseriesQuery = useQuery({
    queryKey: ['reports', 'timeseries', workspaceId, dateRange, scopeType, scopeId],
    queryFn: () => apiClient.get<any>('/reports/timeseries', {
      workspace_id: workspaceId ? String(workspaceId) : undefined,
      scope_type: scopeType,
      scope_id: scopeId || undefined,
      step: 'day',
    }),
  })

  const summary = summaryQuery.data || {}
  const points: ReportPoint[] = (timeseriesQuery.data?.items || []).map((item: any) => {
    const m = item.metrics || {}
    return {
      date: item.date,
      spend: Number(m.spend || 0),
      clicks: Number(m.clicks || 0),
      impressions: Number(m.impressions || 0),
      conversions: m.conversions !== undefined ? Number(m.conversions) : undefined,
    }
  })

  const isLoading = summaryQuery.isLoading || timeseriesQuery.isLoading
  const error = summaryQuery.error || timeseriesQuery.error

  const exportMutation = useMutation({
    mutationFn: () => reportService.download(0),
    onError: (err: unknown) => alert('导出失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const handleExport = () => {
    const params = new URLSearchParams()
    if (workspaceId) params.set('workspace_id', String(workspaceId))
    if (scopeType) params.set('scope_type', scopeType)
    if (scopeId) params.set('scope_id', scopeId)
    if (dateRange) params.set('date_range', dateRange)
    const token = localStorage.getItem('orbit-auth') ? JSON.parse(localStorage.getItem('orbit-auth')!).state.accessToken : ''
    fetch(`/api/v1/reports/export?${params.toString()}`, {
      headers: { 'Authorization': `Bearer ${token}` },
    })
      .then((res) => res.text())
      .then((csv) => {
        const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
        const url = URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = `report-${dateRange}-${Date.now()}.csv`
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        URL.revokeObjectURL(url)
      })
      .catch(() => alert('导出失败'))
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">数据报表</h1>
          <p className="text-sm text-text-muted mt-1">跨平台聚合报表与趋势分析</p>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">加载失败：{(error as Error).message}</p>
          <button onClick={() => { summaryQuery.refetch(); timeseriesQuery.refetch() }} className="btn btn-primary">重试</button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">数据报表</h1>
          <p className="text-sm text-text-muted mt-1">跨平台聚合报表与趋势分析</p>
        </div>
        <div className="flex gap-2">
          <select className="input" value={dateRange} onChange={(e) => setDateRange(e.target.value)}>
            <option value="last_7d">最近 7 天</option>
            <option value="last_30d">最近 30 天</option>
            <option value="last_90d">最近 90 天</option>
          </select>
          <button onClick={handleExport} disabled={exportMutation.isPending} className="btn btn-secondary">导出 CSV</button>
          <button onClick={() => { summaryQuery.refetch(); timeseriesQuery.refetch() }} className="btn btn-secondary">刷新</button>
        </div>
      </div>

      <div className="card p-4 flex items-center gap-4">
        <label className="text-sm text-text-secondary">范围类型</label>
        <select className="input" value={scopeType} onChange={(e) => setScopeType(e.target.value)}>
          <option value="account">广告账户</option>
          <option value="campaign">广告系列</option>
          <option value="adgroup">广告组</option>
        </select>
        <input className="input max-w-xs" placeholder="范围 ID（可选）" value={scopeId} onChange={(e) => setScopeId(e.target.value)} />
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
                  <Line type="monotone" dataKey="impressions" stroke="#F59E0B" strokeWidth={2} dot={false} />
                </LineChart>
              </ResponsiveContainer>
            )}
          </div>
        </>
      )}
    </div>
  )
}


