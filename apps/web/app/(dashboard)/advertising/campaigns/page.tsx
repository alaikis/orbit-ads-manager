'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { campaignService } from '@/lib/api'
import { useState } from 'react'
import Link from 'next/link'

type Campaign = { id: number; name: string; status: string; platform: string; account_id: number; daily_budget?: number }

const STATUS_BADGE: Record<string, string> = {
  enabled: 'bg-success-bg text-success',
  paused: 'bg-warning-bg text-warning',
  active: 'bg-success-bg text-success',
  inactive: 'bg-surface-subtle text-text-muted',
}

export default function CampaignsPage() {
  const queryClient = useQueryClient()
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['campaigns'],
    queryFn: () => campaignService.list(),
  })
  const campaigns: Campaign[] = (data as any)?.items || []

  const pauseMutation = useMutation({
    mutationFn: (id: number) => campaignService.pause(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['campaigns'] }),
    onError: (err: unknown) => alert('暂停失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const resumeMutation = useMutation({
    mutationFn: (id: number) => campaignService.resume(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['campaigns'] }),
    onError: (err: unknown) => alert('启用失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">广告系列</h1>
          <p className="text-sm text-text-muted mt-1">管理各平台的广告系列</p>
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
          <h1 className="text-2xl font-semibold text-text-primary">广告系列</h1>
          <p className="text-sm text-text-muted mt-1">管理各平台的广告系列</p>
        </div>
        <button onClick={() => refetch()} className="btn btn-secondary">刷新</button>
      </div>

      {isLoading ? (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead><tr className="border-b border-border-default">{['名称','平台','状态','日预算','操作'].map(h => (<th key={h} className="text-left p-4 font-title-sm text-text-secondary">{h}</th>))}</tr></thead>
            <tbody>{[1,2,3].map(i => (<tr key={i} className="border-b border-border-default last:border-0">{[1,2,3,4,5].map(j => (<td key={j} className="p-4"><div className="h-4 bg-surface-subtle rounded animate-pulse" style={{width: `${50 + Math.random()*40}%`}} /></td>))}</tr>))}</tbody>
          </table>
        </div>
      ) : campaigns.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted mb-2">暂无广告系列</p>
          <p className="text-sm text-text-muted">请先连接广告账户以同步系列数据</p>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border-default">
                <th className="text-left p-4 font-title-sm text-text-secondary">名称</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">平台</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">状态</th>
                <th className="text-right p-4 font-title-sm text-text-secondary">日预算</th>
                <th className="text-right p-4 font-title-sm text-text-secondary">操作</th>
              </tr>
            </thead>
            <tbody>
              {campaigns.map((campaign) => (
                <tr key={campaign.id} className="border-b border-border-default last:border-0 hover:bg-surface-hover">
                  <td className="p-4">
                    <Link href={`/advertising/campaigns/${campaign.id}`} className="font-medium text-primary-600 hover:underline">
                      {campaign.name}
                    </Link>
                  </td>
                  <td className="p-4 text-text-secondary">{campaign.platform}</td>
                  <td className="p-4">
                    <span className={`badge ${STATUS_BADGE[campaign.status] || 'bg-surface-subtle text-text-muted'}`}>
                      {campaign.status}
                    </span>
                  </td>
                  <td className="p-4 text-right font-mono">
                    {campaign.daily_budget != null ? `$${(campaign.daily_budget / 100).toFixed(2)}` : '-'}
                  </td>
                  <td className="p-4 text-right space-x-2">
                    {campaign.status === 'enabled' || campaign.status === 'active' ? (
                      <button onClick={() => pauseMutation.mutate(campaign.id)} className="btn btn-ghost text-sm text-warning">暂停</button>
                    ) : (
                      <button onClick={() => resumeMutation.mutate(campaign.id)} className="btn btn-ghost text-sm text-success">启用</button>
                    )}
                    <Link href={`/advertising/campaigns/${campaign.id}`} className="btn btn-ghost text-sm">详情</Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
