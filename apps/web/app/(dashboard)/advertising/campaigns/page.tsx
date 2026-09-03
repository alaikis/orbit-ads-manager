'use client'

import { useQuery } from '@tanstack/react-query'
import { campaignService } from '@/lib/api'
import Link from 'next/link'
import { useState, useMemo } from 'react'
import { MetricCard } from '@/components/dashboard/metric-card'
import { RefreshCw, Download, Search, Filter } from 'lucide-react'

type Campaign = {
  id: number
  name: string
  platform: string
  status: string
  daily_budget_cents?: number
  metrics?: {
    spend: number
    impressions: number
    clicks: number
    ctr: number
    cpc: number
    conversions: number
    cpa: number
    roas: number
    budget_rate: number
  }
}

type FilterState = {
  search: string
  platform: string
  status: string
}

export default function CampaignsPage() {
  const [filters, setFilters] = useState<FilterState>({
    search: '',
    platform: '',
    status: '',
  })
  const [showFilters, setShowFilters] = useState(false)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['campaigns', filters],
    queryFn: () => campaignService.list(),
  })

  const campaigns = ((data as any)?.items || []) as Campaign[]

  const filtered = useMemo(() => {
    return campaigns.filter((c) => {
      if (filters.search && !c.name.toLowerCase().includes(filters.search.toLowerCase())) return false
      if (filters.platform && c.platform !== filters.platform) return false
      if (filters.status && c.status !== filters.status) return false
      return true
    })
  }, [campaigns, filters])

  const summary = useMemo(() => {
    const spend = filtered.reduce((sum, c) => sum + (c.metrics?.spend || 0), 0)
    const clicks = filtered.reduce((sum, c) => sum + (c.metrics?.clicks || 0), 0)
    const conversions = filtered.reduce((sum, c) => sum + (c.metrics?.conversions || 0), 0)
    const ctr = clicks > 0 ? (clicks / (filtered.reduce((sum, c) => sum + (c.metrics?.impressions || 0), 0) || 1)) * 100 : 0
    const cpa = conversions > 0 ? spend / conversions : 0
    const roas = spend > 0 ? (filtered.reduce((sum, c) => sum + (c.metrics?.roas || 0), 0) / filtered.length) : 0
    return { spend, ctr, cpa, roas, count: filtered.length }
  }, [filtered])

  const getBudgetBadge = (rate?: number) => {
    if (!rate && rate !== 0) return null
    if (rate > 100) return { class: 'bg-danger-bg text-danger', label: `${rate.toFixed(0)}%` }
    if (rate > 80) return { class: 'bg-warning-bg text-warning', label: `${rate.toFixed(0)}%` }
    return { class: 'bg-success-bg text-success', label: `${rate.toFixed(0)}%` }
  }

  const platforms = useMemo(() => {
    const set = new Set(campaigns.map((c) => c.platform).filter(Boolean))
    return Array.from(set)
  }, [campaigns])

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">广告管理</h1>
          <p className="text-sm text-text-muted mt-1">跨平台广告系列统一视图</p>
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
          <h1 className="text-2xl font-semibold text-text-primary">广告管理</h1>
          <p className="text-sm text-text-muted mt-1">跨平台广告系列统一视图</p>
        </div>
        <div className="flex gap-2">
          <button onClick={() => setShowFilters(!showFilters)} className={`btn ${showFilters ? 'btn-primary' : 'btn-secondary'}`}>
            <Filter size={16} className="mr-2" />筛选
          </button>
          <button onClick={() => refetch()} className="btn btn-secondary">
            <RefreshCw size={16} className="mr-2" />刷新
          </button>
          <button className="btn btn-secondary">
            <Download size={16} className="mr-2" />导出
          </button>
        </div>
      </div>

      {showFilters && (
        <div className="card p-4">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="relative">
              <Search size={16} className="absolute left-3 top-2.5 text-text-muted" />
              <input
                className="input pl-9"
                placeholder="搜索广告系列..."
                value={filters.search}
                onChange={(e) => setFilters((f) => ({ ...f, search: e.target.value }))}
              />
            </div>
            <select
              className="input"
              value={filters.platform}
              onChange={(e) => setFilters((f) => ({ ...f, platform: e.target.value }))}
            >
              <option value="">全部平台</option>
              {platforms.map((p) => (
                <option key={p} value={p}>{p}</option>
              ))}
            </select>
            <select
              className="input"
              value={filters.status}
              onChange={(e) => setFilters((f) => ({ ...f, status: e.target.value }))}
            >
              <option value="">全部状态</option>
              <option value="enabled">启用</option>
              <option value="paused">暂停</option>
              <option value="removed">移除</option>
            </select>
            <button onClick={() => setFilters({ search: '', platform: '', status: '' })} className="btn btn-secondary">清除筛选</button>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
        <MetricCard title="广告系列数" value={summary.count} format="number" />
        <MetricCard title="总花费" value={summary.spend} format="currency" />
        <MetricCard title="平均 CTR" value={summary.ctr} format="percent" />
        <MetricCard title="平均 ROAS" value={summary.roas} format="ratio" />
      </div>

      {isLoading ? (
        <div className="card p-4 animate-pulse space-y-3">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="h-12 bg-surface-subtle rounded" />
          ))}
        </div>
      ) : filtered.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted mb-4">{campaigns.length === 0 ? '还没有广告系列，请先授权广告账户' : '没有匹配的广告系列，请清除筛选条件'}</p>
          {campaigns.length === 0 && (
            <Link href="/accounts" className="btn btn-primary">前往授权账户</Link>
          )}
        </div>
      ) : (
        <div className="card overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border-default">
                  <th className="text-left p-4 font-title-sm text-text-secondary min-w-[200px]">广告系列</th>
                  <th className="text-left p-4 font-title-sm text-text-secondary">状态</th>
                  <th className="text-right p-4 font-title-sm text-text-secondary">日预算</th>
                  <th className="text-right p-4 font-title-sm text-text-secondary">花费</th>
                  <th className="text-right p-4 font-title-sm text-text-secondary">CTR</th>
                  <th className="text-right p-4 font-title-sm text-text-secondary">CPA</th>
                  <th className="text-right p-4 font-title-sm text-text-secondary">ROAS</th>
                  <th className="text-center p-4 font-title-sm text-text-secondary">预算消耗</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((c: Campaign) => {
                  const badge = getBudgetBadge(c.metrics?.budget_rate)
                  return (
                    <tr key={c.id} className="border-b border-border-default last:border-0 hover:bg-surface-hover">
                      <td className="p-4">
                        <Link href={`/advertising/campaigns/${c.id}`} className="text-primary-600 hover:underline font-medium">
                          {c.name}
                        </Link>
                        <div className="text-xs text-text-muted mt-0.5">{c.platform}</div>
                      </td>
                      <td className="p-4">
                        <span className={`badge ${c.status === 'enabled' ? 'bg-success-bg text-success' : c.status === 'paused' ? 'bg-warning-bg text-warning' : 'bg-danger-bg text-danger'}`}>
                          {c.status === 'enabled' ? '启用' : c.status === 'paused' ? '暂停' : '移除'}
                        </span>
                      </td>
                      <td className="p-4 text-right font-mono">${((c.daily_budget_cents || 0) / 100).toFixed(2)}</td>
                      <td className="p-4 text-right font-mono">${c.metrics?.spend?.toFixed(2) || '-'}</td>
                      <td className="p-4 text-right font-mono">{c.metrics?.ctr?.toFixed(2) || '-'}%</td>
                      <td className="p-4 text-right font-mono">${c.metrics?.cpa?.toFixed(2) || '-'}</td>
                      <td className="p-4 text-right font-mono">{c.metrics?.roas?.toFixed(2) || '-'}</td>
                      <td className="p-4 text-center">
                        {badge ? <span className={`badge ${badge.class}`}>{badge.label}</span> : '-'}
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  )
}
