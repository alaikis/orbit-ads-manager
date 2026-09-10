'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { adAccountService } from '@/lib/api'
import { apiClient } from '@/lib/api/client'
import { useState } from 'react'

type Account = { id: number; name: string; platform: string; status: string; currency?: string; customer_id?: string; external_id?: string }

const OAUTH_PLATFORMS = new Set(['google', 'meta', 'bing'])

export default function AccountsPage() {
  const queryClient = useQueryClient()
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editingAccount, setEditingAccount] = useState<Account | null>(null)
  const [editName, setEditName] = useState('')
  const [editCurrency, setEditCurrency] = useState('')
  const [editCustomerId, setEditCustomerId] = useState('')
  const [name, setName] = useState('')
  const [platform, setPlatform] = useState('google')
  const [externalId, setExternalId] = useState('')
  const [currency, setCurrency] = useState('USD')
  const [customerId, setCustomerId] = useState('')
  const [formError, setFormError] = useState('')
  const [actionError, setActionError] = useState<string | null>(null)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['ad-accounts'],
    queryFn: () => adAccountService.list(),
  })

  const accounts: Account[] = (data as any)?.items || []

  const createMutation = useMutation({
    mutationFn: (payload: { name: string; platform: string; external_id: string; currency: string; customer_id?: string }) =>
      adAccountService.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ad-accounts'] })
      setShowCreateForm(false)
      setName(''); setExternalId(''); setCustomerId('')
      setFormError('')
    },
    onError: (err: unknown) => setFormError(err instanceof Error ? err.message : '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: (payload: { id: number; data: { name?: string; currency?: string; customer_id?: string; status?: string } }) =>
      adAccountService.update(payload.id, payload.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ad-accounts'] })
      setEditingAccount(null)
    },
    onError: (err: unknown) => setActionError(err instanceof Error ? err.message : '更新失败'),
  })

  const handleCreate = () => {
    setFormError('')
    if (!name.trim() || !externalId.trim()) {
      setFormError('请填写账户名称和外部 ID')
      return
    }

    if (OAUTH_PLATFORMS.has(platform)) {
      adAccountService.authorize(platform)
        .then((res) => {
          const params = new URLSearchParams({
            state: res.state,
            platform,
            name: name.trim(),
            external_id: externalId.trim(),
            currency: currency.trim().toUpperCase() || 'USD',
            ...(customerId.trim() ? { customer_id: customerId.trim() } : {}),
          } as Record<string, string>)
          window.location.href = `/api/v1/oauth/connect/${platform}?${params.toString()}`
        })
        .catch((err: unknown) => setFormError(err instanceof Error ? err.message : '授权失败'))
      return
    }

    createMutation.mutate({
      name: name.trim(),
      platform,
      external_id: externalId.trim(),
      currency: currency.trim().toUpperCase() || 'USD',
      customer_id: customerId.trim() || undefined,
    })
  }

  const startEdit = (account: Account) => {
    setEditingAccount(account)
    setEditName(account.name)
    setEditCurrency(account.currency || '')
    setEditCustomerId(account.customer_id || '')
    setActionError(null)
  }

  const saveEdit = () => {
    if (!editingAccount) return
    const payload: Record<string, string> = {}
    if (editName.trim()) payload.name = editName.trim()
    if (editCurrency.trim()) payload.currency = editCurrency.trim().toUpperCase()
    if (editCustomerId.trim()) payload.customer_id = editCustomerId.trim()
    updateMutation.mutate({ id: editingAccount.id, data: payload })
  }

  const toggleStatus = (account: Account) => {
    const nextStatus = account.status === 'active' ? 'disabled' : 'active'
    updateMutation.mutate({ id: account.id, data: { status: nextStatus } })
  }

  const handleOAuthConnect = (account: Account) => {
    const params = new URLSearchParams({
      state: `reconnect-${account.id}-${Date.now()}`,
      platform: account.platform,
      name: account.name,
      external_id: account.external_id || '',
      ...(account.currency ? { currency: account.currency } : {}),
      ...(account.customer_id ? { customer_id: account.customer_id } : {}),
    })
    window.location.href = `/api/v1/oauth/connect/${account.platform}?${params.toString()}`
  }

  const platformLabel: Record<string, string> = {
    google: 'Google Ads',
    meta: 'Meta Ads',
    bing: 'Bing Ads',
    tiktok: 'TikTok Ads',
  }

  const authMethodLabel: Record<string, string> = {
    google: 'OAuth 授权',
    meta: 'OAuth 授权',
    bing: 'OAuth 授权',
    tiktok: 'API 密钥 / 手动录入',
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">广告账户管理</h1>
            <p className="text-sm text-text-muted mt-1">管理你的广告平台授权账户</p>
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
          <h1 className="text-2xl font-semibold text-text-primary">广告账户管理</h1>
          <p className="text-sm text-text-muted mt-1">管理你的广告平台授权账户</p>
        </div>
        <div className="flex gap-2">
          <button onClick={() => setShowCreateForm(true)} className="btn btn-primary">授权新账户</button>
          <button onClick={() => refetch()} className="btn btn-secondary">刷新</button>
        </div>
      </div>

      {showCreateForm && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">授权新账户</h2>
              <button onClick={() => { setShowCreateForm(false); setFormError('') }} className="text-text-muted hover:text-text-primary">
                <span className="text-xl leading-none">×</span>
              </button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">账户名称 *</label>
                <input className="input" value={name} onChange={(e) => setName(e.target.value)} placeholder="例如：Google Ads 主账户" />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">平台 *</label>
                <select className="input" value={platform} onChange={(e) => setPlatform(e.target.value)}>
                  <option value="google">Google Ads</option>
                  <option value="meta">Meta Ads</option>
                  <option value="bing">Bing Ads</option>
                  <option value="tiktok">TikTok Ads</option>
                </select>
                <p className="text-xs text-text-muted mt-1">Google / Meta / Bing 使用 OAuth 授权；TikTok 使用 API 密钥录入</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">广告账户 ID *</label>
                <input className="input" value={externalId} onChange={(e) => setExternalId(e.target.value)} placeholder="例如：123-456-7890" />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">币种</label>
                <input className="input" value={currency} onChange={(e) => setCurrency(e.target.value)} placeholder="USD" />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">客户 ID（可选）</label>
                <input className="input" value={customerId} onChange={(e) => setCustomerId(e.target.value)} placeholder="MCC 账户 ID（部分平台需要）" />
              </div>
              {formError && (<div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{formError}</div>)}
              <div className="flex gap-2 justify-end pt-2">
                <button onClick={() => { setShowCreateForm(false); setFormError('') }} className="btn btn-secondary">取消</button>
                <button onClick={handleCreate} disabled={createMutation.isPending} className="btn btn-primary">
                  {createMutation.isPending ? (<><span className="mr-2 inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />授权中...</>) : (OAUTH_PLATFORMS.has(platform) ? '前往授权' : '授权')}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {editingAccount && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">编辑账户</h2>
              <button onClick={() => setEditingAccount(null)} className="text-text-muted hover:text-text-primary"><span className="text-xl leading-none">×</span></button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">账户名称</label>
                <input className="input" value={editName} onChange={(e) => setEditName(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">币种</label>
                <input className="input" value={editCurrency} onChange={(e) => setEditCurrency(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">客户 ID</label>
                <input className="input" value={editCustomerId} onChange={(e) => setEditCustomerId(e.target.value)} />
              </div>
              {actionError && (<div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{actionError}</div>)}
              <div className="flex gap-2 justify-end pt-2">
                <button onClick={() => setEditingAccount(null)} className="btn btn-secondary">取消</button>
                <button onClick={saveEdit} disabled={updateMutation.isPending} className="btn btn-primary">{updateMutation.isPending ? '保存中...' : '保存'}</button>
              </div>
            </div>
          </div>
        </div>
      )}

      {isLoading ? (
        <div className="card p-4 animate-pulse space-y-3">{[1, 2, 3].map((i) => (<div key={i} className="h-10 bg-surface-subtle rounded" />))}</div>
      ) : accounts.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted mb-4">还没有授权任何广告账户</p>
          <button onClick={() => setShowCreateForm(true)} className="btn btn-primary">授权第一个账户</button>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border-default">
                <th className="text-left p-4 font-title-sm text-text-secondary">账户名称</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">平台</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">授权方式</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">状态</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">币种</th>
                <th className="text-right p-4 font-title-sm text-text-secondary">操作</th>
              </tr>
            </thead>
            <tbody>
              {accounts.map((account) => (
                <tr key={account.id} className="border-b border-border-default last:border-0 hover:bg-surface-hover">
                  <td className="p-4 font-medium">{account.name}</td>
                  <td className="p-4 text-text-secondary capitalize">{platformLabel[account.platform] || account.platform}</td>
                  <td className="p-4 text-text-muted text-xs">{authMethodLabel[account.platform] || '未知'}</td>
                  <td className="p-4">
                    <span className={`badge ${account.status === 'active' ? 'bg-success-bg text-success' : account.status === 'disabled' ? 'bg-danger-bg text-danger' : 'bg-warning-bg text-warning'}`}>
                      {account.status}
                    </span>
                  </td>
                  <td className="p-4 text-text-muted">{account.currency || '-'}</td>
                  <td className="p-4 text-right space-x-2">
                    <button onClick={() => startEdit(account)} className="btn btn-ghost text-sm">编辑</button>
                    <button onClick={() => toggleStatus(account)} className="btn btn-ghost text-sm" title={account.status === 'active' ? '禁用' : '启用'}>
                      {account.status === 'active' ? '禁用' : '启用'}
                    </button>
                    {OAUTH_PLATFORMS.has(account.platform) && (
                      <button onClick={() => handleOAuthConnect(account)} className="btn btn-ghost text-sm" title="重新授权">重连</button>
                    )}
                    <button onClick={() => adAccountService.delete(account.id).then(() => queryClient.invalidateQueries({ queryKey: ['ad-accounts'] }))} className="btn btn-ghost text-sm text-danger-500">删除</button>
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