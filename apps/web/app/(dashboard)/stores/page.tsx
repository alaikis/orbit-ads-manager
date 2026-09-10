'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { storeService } from '@/lib/api'
import Link from 'next/link'
import { useState } from 'react'

type Store = { id: number; name: string; platform: string; status: string; last_synced_at?: string; base_url?: string; api_key?: string; api_secret?: string }

const PLATFORM_LABEL: Record<string, string> = {
  woocommerce: 'WooCommerce',
  shopify: 'Shopify',
}

export default function StoresPage() {
  const queryClient = useQueryClient()
  const [showForm, setShowForm] = useState(false)
  const [editingStore, setEditingStore] = useState<Store | null>(null)
  const [name, setName] = useState('')
  const [platform, setPlatform] = useState('woocommerce')
  const [baseUrl, setBaseUrl] = useState('')
  const [apiKey, setApiKey] = useState('')
  const [apiSecret, setApiSecret] = useState('')
  const [formError, setFormError] = useState('')
  const [actionError, setActionError] = useState<string | null>(null)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['stores'],
    queryFn: () => storeService.list(),
  })

  const stores: Store[] = (data as any)?.items || []

  const createMutation = useMutation({
    mutationFn: (payload: { name: string; platform: string; base_url: string; api_key: string; api_secret?: string }) =>
      storeService.create(payload as any),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stores'] })
      setShowForm(false)
      setName(''); setBaseUrl(''); setApiKey(''); setApiSecret('')
      setFormError('')
    },
    onError: (err: unknown) => setFormError(err instanceof Error ? err.message : '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: (payload: { id: number; data: Record<string, string> }) =>
      storeService.update(payload.id, payload.data as any),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stores'] })
      setEditingStore(null)
      setActionError(null)
    },
    onError: (err: unknown) => setActionError(err instanceof Error ? err.message : '更新失败'),
  })

  const handleCreate = () => {
    setFormError('')
    if (!name.trim() || !baseUrl.trim() || !apiKey.trim()) {
      setFormError('请填写店铺名称、URL 和 API Key')
      return
    }
    createMutation.mutate({
      name: name.trim(),
      platform,
      base_url: baseUrl.trim(),
      api_key: apiKey.trim(),
      api_secret: apiSecret.trim() || undefined,
    })
  }

  const startEdit = (store: Store) => {
    setEditingStore(store)
    setName(store.name)
    setPlatform(store.platform)
    setBaseUrl(store.base_url || '')
    setApiKey(store.api_key || '')
    setApiSecret(store.api_secret || '')
    setActionError(null)
  }

  const saveEdit = () => {
    if (!editingStore) return
    const payload: Record<string, string> = {}
    if (name.trim()) payload.name = name.trim()
    if (baseUrl.trim()) payload.base_url = baseUrl.trim()
    if (apiKey.trim()) payload.api_key = apiKey.trim()
    if (apiSecret.trim()) payload.api_secret = apiSecret.trim()
    updateMutation.mutate({ id: editingStore.id, data: payload })
  }

  const toggleStatus = (store: Store) => {
    const nextStatus = store.status === 'bound' ? 'disabled' : 'bound'
    updateMutation.mutate({ id: store.id, data: { status: nextStatus } })
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">店铺管理</h1>
            <p className="text-sm text-text-muted mt-1">管理你的电商店铺连接与同步</p>
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
          <h1 className="text-2xl font-semibold text-text-primary">店铺管理</h1>
          <p className="text-sm text-text-muted mt-1">管理你的电商店铺连接与同步</p>
        </div>
        <div className="flex gap-2">
          <button onClick={() => setShowForm(true)} className="btn btn-primary">添加店铺</button>
          <button onClick={() => refetch()} className="btn btn-secondary">刷新</button>
        </div>
      </div>

      {showForm && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">添加店铺</h2>
              <button onClick={() => { setShowForm(false); setFormError('') }} className="text-text-muted hover:text-text-primary">
                <span className="text-xl leading-none">×</span>
              </button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">店铺名称 *</label>
                <input className="input" value={name} onChange={(e) => setName(e.target.value)} placeholder="例如：我的 Shopify 店铺" />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">平台 *</label>
                <select className="input" value={platform} onChange={(e) => setPlatform(e.target.value)}>
                  <option value="woocommerce">WooCommerce</option>
                  <option value="shopify">Shopify</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">店铺 URL *</label>
                <input className="input" value={baseUrl} onChange={(e) => setBaseUrl(e.target.value)} placeholder="https://example.com" />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">API Key *</label>
                <input className="input" value={apiKey} onChange={(e) => setApiKey(e.target.value)} placeholder="ck_xxxxxxxxxxxxx" />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">API Secret（可选）</label>
                <input className="input" value={apiSecret} onChange={(e) => setApiSecret(e.target.value)} placeholder="cs_xxxxxxxxxxxxx" />
              </div>
              {formError && (<div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{formError}</div>)}
              <div className="flex gap-2 justify-end pt-2">
                <button onClick={() => { setShowForm(false); setFormError('') }} className="btn btn-secondary">取消</button>
                <button onClick={handleCreate} disabled={createMutation.isPending} className="btn btn-primary">
                  {createMutation.isPending ? (<><span className="mr-2 inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />创建中...</>) : '创建'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {editingStore && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">编辑店铺</h2>
              <button onClick={() => setEditingStore(null)} className="text-text-muted hover:text-text-primary"><span className="text-xl leading-none">×</span></button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">店铺名称</label>
                <input className="input" value={name} onChange={(e) => setName(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">店铺 URL</label>
                <input className="input" value={baseUrl} onChange={(e) => setBaseUrl(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">API Key</label>
                <input className="input" value={apiKey} onChange={(e) => setApiKey(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">API Secret</label>
                <input className="input" value={apiSecret} onChange={(e) => setApiSecret(e.target.value)} />
              </div>
              {actionError && (<div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{actionError}</div>)}
              <div className="flex gap-2 justify-end pt-2">
                <button onClick={() => setEditingStore(null)} className="btn btn-secondary">取消</button>
                <button onClick={saveEdit} disabled={updateMutation.isPending} className="btn btn-primary">{updateMutation.isPending ? '保存中...' : '保存'}</button>
              </div>
            </div>
          </div>
        </div>
      )}

      {isLoading ? (
        <div className="space-y-4">{[1, 2, 3].map((i) => (<div key={i} className="card p-4 animate-pulse h-12 bg-surface-subtle rounded" />))}</div>
      ) : stores.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted mb-4">还没有绑定任何店铺</p>
          <button onClick={() => setShowForm(true)} className="btn btn-primary">绑定第一个店铺</button>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border-default">
                <th className="text-left p-4 font-title-sm text-text-secondary">店铺名称</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">平台</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">状态</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">最近同步</th>
                <th className="text-right p-4 font-title-sm text-text-secondary">操作</th>
              </tr>
            </thead>
            <tbody>
              {stores.map((store) => (
                <tr key={store.id} className="border-b border-border-default last:border-0 hover:bg-surface-hover">
                  <td className="p-4 font-medium">
                    <Link href={`/stores/${store.id}`} className="text-primary-600 hover:underline">
                      {store.name}
                    </Link>
                  </td>
                  <td className="p-4 text-text-secondary capitalize">{PLATFORM_LABEL[store.platform] || store.platform}</td>
                  <td className="p-4">
                    <span className={`badge ${store.status === 'bound' || store.status === 'active' ? 'bg-success-bg text-success' : 'bg-warning-bg text-warning'}`}>
                      {store.status || 'pending'}
                    </span>
                  </td>
                  <td className="p-4 text-text-muted">{store.last_synced_at || '-'}</td>
                  <td className="p-4 text-right space-x-2">
                    <button onClick={() => startEdit(store)} className="btn btn-ghost text-sm">编辑</button>
                    <button onClick={() => toggleStatus(store)} className="btn btn-ghost text-sm" title={store.status === 'bound' ? '禁用' : '启用'}>
                      {store.status === 'bound' ? '禁用' : '启用'}
                    </button>
                    <button onClick={() => storeService.delete(store.id).then(() => queryClient.invalidateQueries({ queryKey: ['stores'] }))} className="btn btn-ghost text-sm text-danger-500">删除</button>
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