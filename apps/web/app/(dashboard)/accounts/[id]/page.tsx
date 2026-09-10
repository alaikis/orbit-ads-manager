'use client'

import { useQuery } from '@tanstack/react-query'
import { adAccountService } from '@/lib/api'
import Link from 'next/link'

type Account = { id: number; name: string; platform: string; status: string; currency?: string; customer_id?: string; external_id?: string }

const PLATFORM_LABEL: Record<string, string> = {
  google: 'Google Ads',
  meta: 'Meta Ads',
  bing: 'Bing Ads',
  tiktok: 'TikTok Ads',
}

export default function AccountDetailPage({ params }: { params: { id: string } }) {
  const accountId = Number(params.id)
  const { data, isLoading, error } = useQuery({
    queryKey: ['ad-account', accountId],
    queryFn: () => adAccountService.get(accountId),
  })

  const account = data as Account | undefined

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/accounts" className="btn btn-ghost">返回</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">广告账户详情</h1>
            <p className="text-sm text-text-muted mt-1">加载中...</p>
          </div>
        </div>
        <div className="card p-6 animate-pulse space-y-3">
          <div className="h-6 bg-surface-subtle rounded w-1/3" />
          <div className="h-4 bg-surface-subtle rounded w-1/2" />
        </div>
      </div>
    )
  }

  if (error || !account) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/accounts" className="btn btn-ghost">返回</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">广告账户详情</h1>
            <p className="text-sm text-text-muted mt-1">加载失败</p>
          </div>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">{(error as Error)?.message || '未找到该账户'}</p>
          <Link href="/accounts" className="btn btn-primary">返回列表</Link>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link href="/accounts" className="btn btn-ghost">返回</Link>
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">{account.name}</h1>
          <p className="text-sm text-text-muted mt-1">
            {PLATFORM_LABEL[account.platform] || account.platform} | 状态：{account.status}
          </p>
        </div>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <tbody>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted w-40">账户 ID</td>
              <td className="p-4 font-mono">{account.id}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">账户名称</td>
              <td className="p-4">{account.name}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">平台</td>
              <td className="p-4">{PLATFORM_LABEL[account.platform] || account.platform}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">外部 ID</td>
              <td className="p-4 font-mono">{account.external_id || '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">状态</td>
              <td className="p-4">
                <span className={`badge ${account.status === 'active' ? 'bg-success-bg text-success' : account.status === 'disabled' ? 'bg-danger-bg text-danger' : 'bg-warning-bg text-warning'}`}>
                  {account.status}
                </span>
              </td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">币种</td>
              <td className="p-4">{account.currency || '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">客户 ID</td>
              <td className="p-4">{account.customer_id || '-'}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  )
}