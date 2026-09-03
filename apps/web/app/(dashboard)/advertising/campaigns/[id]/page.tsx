'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { campaignService, metricsService } from '@/lib/api'
import Link from 'next/link'
import { useState } from 'react'
import { MetricCard } from '@/components/dashboard/metric-card'
import { RefreshCw, Pause, Play, DollarSign, Bot, ArrowLeft } from 'lucide-react'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

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

type AdGroup = {
  id: number
  name: string
  status: string
  bid_micros?: number
  metrics?: {
    spend: number
    clicks: number
    conversions: number
    cpa: number
    roas: number
  }
}

type TimeSeriesPoint = { date: string; spend: number; clicks: number }

export default function CampaignDetailPage({ params }: { params: { id: string } }) {
  const queryClient = useQueryClient()
  const [showBudgetDialog, setShowBudgetDialog] = useState(false)
  const [budgetValue, setBudgetValue] = useState('')
  const [errorMsg, setErrorMsg] = useState<string | null>(null)

  const { data: campaignData, isLoading: campaignLoading, error: campaignError } = useQuery({
    queryKey: ['campaign', params.id],
    queryFn: () => campaignService.get(Number(params.id)),
  })

  const { data: metricsData, isLoading: metricsLoading } = useQuery({
    queryKey: ['campaign-metrics', params.id],
    queryFn: () => metricsService.getSummary({ scope_type: 'campaign', scope_id: Number(params.id) }),
  })

  const { data: adGroupsData, isLoading: adGroupsLoading } = useQuery({
    queryKey: ['campaign-adgroups', params.id],
    queryFn: () => campaignService.getAdGroups(Number(params.id)),
  })

  const campaign = campaignData as Campaign | undefined
  const metrics = metricsData || {}
  const adGroups = (Array.isArray(adGroupsData) ? adGroupsData : []) as AdGroup[]

  const pauseMutation = useMutation({
    mutationFn: () => campaignService.pause(Number(params.id)),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaign', params.id] })
      queryClient.invalidateQueries({ queryKey: ['campaigns'] })
    },
  })

  const resumeMutation = useMutation({
    mutationFn: () => campaignService.resume(Number(params.id)),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaign', params.id] })
      queryClient.invalidateQueries({ queryKey: ['campaigns'] })
    },
  })

  const budgetMutation = useMutation({
    mutationFn: (newBudgetCents: number) => campaignService.updateBudget(Number(params.id), newBudgetCents),
    onSuccess: () => {
      setShowBudgetDialog(false)
      setBudgetValue('')
      queryClient.invalidateQueries({ queryKey: ['campaign', params.id] })
      queryClient.invalidateQueries({ queryKey: ['campaigns'] })
    },
    onError: (err: Error) => setErrorMsg(err.message),
  })

  const handleBudgetSave = () => {
    const cents = parseInt(budgetValue, 10) * 100
    if (isNaN(cents) || cents <= 0) {
      setErrorMsg('预算必须为正数')
      return
    }
    budgetMutation.mutate(cents)
  }

  const timeSeries: TimeSeriesPoint[] = []

  const getBudgetBadge = (rate?: number) => {
    if (!rate && rate !== 0) return null
    if (rate > 100) return { class: 'bg-danger-bg text-danger', label: `${rate.toFixed(0)}%` }
    if (rate > 80) return { class: 'bg-warning-bg text-warning', label: `${rate.toFixed(0)}%` }
    return { class: 'bg-success-bg text-success', label: `${rate.toFixed(0)}%` }
  }

  if (campaignError) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/advertising/campaigns" className="btn btn-ghost"><ArrowLeft size={18} /></Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">广告详情</h1>
            <p className="text-sm text-text-muted mt-1">加载失败</p>
          </div>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">{(campaignError as Error).message}</p>
          <button onClick={() => queryClient.invalidateQueries({ queryKey: ['campaign', params.id] })} className="btn btn-primary">重试</button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link href="/advertising/campaigns" className="btn btn-ghost"><ArrowLeft size={18} /></Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">{campaign?.name || '广告详情'}</h1>
            <p className="text-sm text-text-muted mt-1">
              平台: {campaign?.platform || '-'} | 状态: {campaign?.status === 'enabled' ? '启用' : campaign?.status === 'paused' ? '暂停' : campaign?.status || '-'}
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          {campaign?.status === 'enabled' ? (
            <button onClick={() => pauseMutation.mutate()} disabled={pauseMutation.isPending} className="btn btn-secondary">
              <Pause size={16} className="mr-2" />暂停
            </button>
          ) : (
            <button onClick={() => resumeMutation.mutate()} disabled={resumeMutation.isPending} className="btn btn-primary">
              <Play size={16} className="mr-2" />恢复
            </button>
          )}
          <button onClick={() => { setBudgetValue(String((campaign?.daily_budget_cents || 0) / 100)); setShowBudgetDialog(true) }} className="btn btn-secondary">
            <DollarSign size={16} className="mr-2" />改预算
          </button>
          <Link href="/agent" className="btn btn-primary">
            <Bot size={16} className="mr-2" />Agent 优化
          </Link>
        </div>
      </div>

      {campaignLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="card p-4 animate-pulse">
              <div className="h-4 bg-surface-subtle rounded w-1/2 mb-2" />
              <div className="h-8 bg-surface-subtle rounded w-3/4" />
            </div>
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <MetricCard title="花费" value={metrics.spend || 0} format="currency" change={-2.3} />
          <MetricCard title="点击率 CTR" value={metrics.ctr || 0} format="percent" change={1.2} />
          <MetricCard title="转化成本 CPA" value={metrics.cpa || 0} format="currency" change={-5.1} />
          <MetricCard title="广告支出回报率 ROAS" value={metrics.roas || 0} format="ratio" change={3.4} />
        </div>
      )}

      <div className="card p-6">
        <h2 className="font-title-md text-text-primary mb-4">趋势图</h2>
        {timeSeries.length === 0 ? (
          <div className="h-64 flex items-center justify-center text-text-muted">
            暂无趋势数据
          </div>
        ) : (
          <ResponsiveContainer width="100%" height={280}>
            <LineChart data={timeSeries}>
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

      <div className="card overflow-hidden">
        <div className="p-4 border-b border-border-default">
          <h2 className="font-title-md text-text-primary">广告组层级</h2>
        </div>
        {adGroupsLoading ? (
          <div className="p-4 space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-12 bg-surface-subtle rounded animate-pulse" />
            ))}
          </div>
        ) : adGroups.length === 0 ? (
          <div className="p-8 text-center text-text-muted">暂无广告组数据</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border-default">
                  <th className="text-left p-4 font-title-sm text-text-secondary">广告组名称</th>
                  <th className="text-left p-4 font-title-sm text-text-secondary">状态</th>
                  <th className="text-right p-4 font-title-sm text-text-secondary">出价</th>
                  <th className="text-right p-4 font-title-sm text-text-secondary">花费</th>
                  <th className="text-right p-4 font-title-sm text-text-secondary">转化</th>
                  <th className="text-right p-4 font-title-sm text-text-secondary">CPA</th>
                  <th className="text-right p-4 font-title-sm text-text-secondary">ROAS</th>
                </tr>
              </thead>
              <tbody>
                {adGroups.map((ag: AdGroup) => (
                  <tr key={ag.id} className="border-b border-border-default last:border-0 hover:bg-surface-hover">
                    <td className="p-4 font-medium">{ag.name}</td>
                    <td className="p-4">
                      <span className={`badge ${ag.status === 'enabled' ? 'bg-success-bg text-success' : ag.status === 'paused' ? 'bg-warning-bg text-warning' : 'bg-danger-bg text-danger'}`}>
                        {ag.status === 'enabled' ? '启用' : ag.status === 'paused' ? '暂停' : ag.status}
                      </span>
                    </td>
                    <td className="p-4 text-right font-mono">{ag.bid_micros ? `$${(ag.bid_micros / 1_000_000).toFixed(2)}` : '-'}</td>
                    <td className="p-4 text-right font-mono">${ag.metrics?.spend?.toFixed(2) || '-'}</td>
                    <td className="p-4 text-right font-mono">{ag.metrics?.conversions || '-'}</td>
                    <td className="p-4 text-right font-mono">${ag.metrics?.cpa?.toFixed(2) || '-'}</td>
                    <td className="p-4 text-right font-mono">{ag.metrics?.roas?.toFixed(2) || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {showBudgetDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-lg font-semibold mb-4">修改日预算</h2>
            {errorMsg && <p className="text-sm text-danger-500 mb-3">{errorMsg}</p>}
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">新日预算（USD）</label>
                <input
                  className="input"
                  type="number"
                  placeholder="例如：50"
                  value={budgetValue}
                  onChange={(e) => { setBudgetValue(e.target.value); setErrorMsg(null) }}
                />
              </div>
              <div className="flex gap-2 justify-end">
                <button onClick={() => { setShowBudgetDialog(false); setErrorMsg(null) }} className="btn btn-secondary">取消</button>
                <button onClick={handleBudgetSave} disabled={budgetMutation.isPending} className="btn btn-primary">
                  {budgetMutation.isPending ? '保存中...' : '保存'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
