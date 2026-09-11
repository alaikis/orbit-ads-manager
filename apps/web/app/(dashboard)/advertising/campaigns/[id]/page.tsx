'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { campaignService } from '@/lib/api'
import Link from 'next/link'
import { useState } from 'react'

type CampaignDetail = { id: number; name: string; status: string; platform: string; daily_budget?: number; account_id: number }

const STATUS_BADGE: Record<string, string> = {
  enabled: 'bg-success-bg text-success',
  paused: 'bg-warning-bg text-warning',
  active: 'bg-success-bg text-success',
  inactive: 'bg-surface-subtle text-text-muted',
}

export default function CampaignDetailPage({ params }: { params: { id: string } }) {
  const campaignId = Number(params.id)
  const queryClient = useQueryClient()

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['campaign', campaignId],
    queryFn: () => campaignService.get(campaignId),
  })
  const campaign = data as CampaignDetail | undefined

  const { data: adGroupsData } = useQuery({
    queryKey: ['campaign', campaignId, 'ad-groups'],
    queryFn: () => campaignService.listAdGroups(campaignId),
    enabled: !!campaign,
  })
  const adGroups = (adGroupsData as any)?.items || []

  const pauseMutation = useMutation({
    mutationFn: () => campaignService.pause(campaignId),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['campaign', campaignId] }); queryClient.invalidateQueries({ queryKey: ['campaigns'] }) },
  })

  const resumeMutation = useMutation({
    mutationFn: () => campaignService.resume(campaignId),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['campaign', campaignId] }); queryClient.invalidateQueries({ queryKey: ['campaigns'] }) },
  })

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/advertising/campaigns" className="btn btn-ghost">返回广告系列</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">广告系列详情</h1>
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
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !campaign) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/advertising/campaigns" className="btn btn-ghost">返回广告系列</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">广告系列详情</h1>
            <p className="text-sm text-text-muted mt-1">加载失败</p>
          </div>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">{(error as Error)?.message || '未找到该广告系列'}</p>
          <button onClick={() => refetch()} className="btn btn-primary">重试</button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link href="/advertising/campaigns" className="btn btn-ghost">返回广告系列</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">{campaign.name}</h1>
            <p className="text-sm text-text-muted mt-1">
              平台：{campaign.platform} | 状态：<span className={`badge ${STATUS_BADGE[campaign.status] || 'bg-surface-subtle text-text-muted'}`}>{campaign.status}</span>
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          {campaign.status === 'enabled' || campaign.status === 'active' ? (
            <button onClick={() => pauseMutation.mutate()} className="btn btn-secondary">暂停</button>
          ) : (
            <button onClick={() => resumeMutation.mutate()} className="btn btn-secondary">启用</button>
          )}
        </div>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <tbody>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted w-40">系列 ID</td>
              <td className="p-4 font-mono">{campaign.id}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">系列名称</td>
              <td className="p-4">{campaign.name}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">平台</td>
              <td className="p-4">{campaign.platform}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">状态</td>
              <td className="p-4">
                <span className={`badge ${STATUS_BADGE[campaign.status] || 'bg-surface-subtle text-text-muted'}`}>
                  {campaign.status}
                </span>
              </td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">日预算</td>
              <td className="p-4 font-mono">{campaign.daily_budget != null ? `$${(campaign.daily_budget / 100).toFixed(2)}` : '-'}</td>
            </tr>
            <tr>
              <td className="p-4 text-text-muted">广告账户 ID</td>
              <td className="p-4 font-mono">{campaign.account_id ?? '-'}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="card p-6">
        <h2 className="font-title-md text-text-primary mb-4">广告组</h2>
        {adGroups.length === 0 ? (
          <p className="text-sm text-text-muted">暂无广告组</p>
        ) : (
          <div className="space-y-2">
            {adGroups.map((g: any) => (
              <div key={g.id} className="flex items-center justify-between p-3 border border-border-default rounded-md">
                <div>
                  <p className="font-medium text-sm">{g.name}</p>
                  <p className="text-xs text-text-muted">ID: {g.id}</p>
                </div>
                <span className={`badge ${STATUS_BADGE[g.status] || 'bg-surface-subtle text-text-muted'}`}>{g.status}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
