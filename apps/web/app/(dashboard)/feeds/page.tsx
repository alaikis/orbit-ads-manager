'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { feedService } from '@/lib/api'
import { useState } from 'react'
import Link from 'next/link'

type Feed = { id: number; name: string; format: string; status: string; store_id: number; last_generated_at?: string; generated_rows: number; error_count: number }

const FORMAT_OPTIONS = ['xml', 'csv', 'json']

const STATUS_BADGE: Record<string, string> = {
  active: 'bg-success-bg text-success',
  paused: 'bg-warning-bg text-warning',
  error: 'bg-danger-bg text-danger',
}

function formatDate(dateString: string) {
  if (!dateString) return '从未'
  return new Date(dateString).toLocaleString('zh-CN')
}

export default function FeedsPage() {
  const queryClient = useQueryClient()
  const [showCreate, setShowCreate] = useState(false)
  const [editingFeed, setEditingFeed] = useState<Feed | null>(null)
  const [name, setName] = useState('')
  const [format, setFormat] = useState('xml')
  const [storeId, setStoreId] = useState('')
  const [formError, setFormError] = useState('')
  const [actionError, setActionError] = useState<string | null>(null)

  const { data, isLoading, error, refetch } = useQuery({ queryKey: ['feeds'], queryFn: () => feedService.list() })
  const feeds: Feed[] = (data as any)?.items || []

  const createMutation = useMutation({
    mutationFn: (payload: { name: string; store_id: number; format?: string }) =>
      feedService.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feeds'] })
      setShowCreate(false)
      setName(''); setFormat('xml'); setStoreId('')
      setFormError('')
    },
    onError: (err: unknown) => setFormError(err instanceof Error ? err.message : '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: (payload: { id: number; data: { name?: string; status?: string; format?: string } }) =>
      feedService.update(payload.id, payload.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feeds'] })
      setEditingFeed(null)
      setActionError(null)
    },
    onError: (err: unknown) => setActionError(err instanceof Error ? err.message : '更新失败'),
  })

  const syncMutation = useMutation({
    mutationFn: (id: number) => feedService.sync(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feeds'] })
    },
    onError: (err: unknown) => alert('同步失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => feedService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feeds'] })
    },
    onError: (err: unknown) => alert('删除失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const handleCreate = () => {
    setFormError('')
    if (!name.trim() || !storeId.trim()) {
      setFormError('请填写 Feed 名称和店铺 ID')
      return
    }
    createMutation.mutate({
      name: name.trim(),
      store_id: parseInt(storeId, 10),
      format,
    })
  }

  const startEdit = (feed: Feed) => {
    setEditingFeed(feed)
    setName(feed.name)
    setFormat(feed.format)
    setActionError(null)
  }

  const saveEdit = () => {
    if (!editingFeed) return
    const payload: { name?: string; status?: string; format?: string } = {}
    if (name.trim()) payload.name = name.trim()
    if (format.trim()) payload.format = format
    updateMutation.mutate({ id: editingFeed.id, data: payload })
  }

  const toggleStatus = (feed: Feed) => {
    const nextStatus = feed.status === 'active' ? 'paused' : 'active'
    updateMutation.mutate({ id: feed.id, data: { status: nextStatus } })
  }

  const confirmDelete = (id: number) => {
    if (window.confirm('确定要删除这个 Feed 吗？')) {
      deleteMutation.mutate(id)
    }
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">Feed 管理</h1>
          <p className="text-sm text-text-muted mt-1">管理商品 Feed 生成与发布</p>
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
          <h1 className="text-2xl font-semibold text-text-primary">Feed 管理</h1>
          <p className="text-sm text-text-muted mt-1">管理商品 Feed 生成与发布</p>
        </div>
        <div className="flex gap-2">
          <button onClick={() => refetch()} className="btn btn-secondary">刷新</button>
          <button onClick={() => setShowCreate(true)} className="btn btn-primary">新建 Feed</button>
        </div>
      </div>

      {showCreate && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">新建 Feed</h2>
              <button onClick={() => { setShowCreate(false); setFormError('') }} className="text-text-muted hover:text-text-primary"><span className="text-xl leading-none">×</span></button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">Feed 名称 *</label>
                <input className="input" value={name} onChange={(e) => setName(e.target.value)} placeholder="例如：Google Shopping Feed" />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">格式</label>
                <select className="input" value={format} onChange={(e) => setFormat(e.target.value)}>
                  {FORMAT_OPTIONS.map((f) => (<option key={f} value={f}>{f.toUpperCase()}</option>))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">店铺 ID *</label>
                <input className="input" type="number" value={storeId} onChange={(e) => setStoreId(e.target.value)} placeholder="关联的店铺 ID" />
              </div>
              {formError && (<div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{formError}</div>)}
              <div className="flex gap-2 justify-end pt-2">
                <button onClick={() => { setShowCreate(false); setFormError('') }} className="btn btn-secondary">取消</button>
                <button onClick={handleCreate} disabled={createMutation.isPending} className="btn btn-primary">{createMutation.isPending ? '创建中...' : '创建'}</button>
              </div>
            </div>
          </div>
        </div>
      )}

      {editingFeed && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">编辑 Feed</h2>
              <button onClick={() => setEditingFeed(null)} className="text-text-muted hover:text-text-primary"><span className="text-xl leading-none">×</span></button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">Feed 名称</label>
                <input className="input" value={name} onChange={(e) => setName(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">格式</label>
                <select className="input" value={format} onChange={(e) => setFormat(e.target.value)}>
                  {FORMAT_OPTIONS.map((f) => (<option key={f} value={f}>{f.toUpperCase()}</option>))}
                </select>
              </div>
              {actionError && (<div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{actionError}</div>)}
              <div className="flex gap-2 justify-end pt-2">
                <button onClick={() => setEditingFeed(null)} className="btn btn-secondary">取消</button>
                <button onClick={saveEdit} disabled={updateMutation.isPending} className="btn btn-primary">{updateMutation.isPending ? '保存中...' : '保存'}</button>
              </div>
            </div>
          </div>
        </div>
      )}

      {isLoading ? (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead><tr className="border-b border-border-default">{['名称','格式','状态','商品数量','最近同步'].map(h => (<th key={h} className="text-left p-4 font-title-sm text-text-secondary">{h}</th>))}</tr></thead>
            <tbody>{[1,2,3].map(i => (<tr key={i} className="border-b border-border-default last:border-0">{[1,2,3,4,5].map(j => (<td key={j} className="p-4"><div className="h-4 bg-surface-subtle rounded animate-pulse" style={{width: `${50 + Math.random()*40}%`}} /></td>))}</tr>))}</tbody>
          </table>
        </div>
      ) : feeds.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted mb-2">尚无 Feed</p>
          <p className="text-sm text-text-muted">请先绑定店铺后生成</p>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border-default">
                <th className="text-left p-4 font-title-sm text-text-secondary">名称</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">格式</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">状态</th>
                <th className="text-right p-4 font-title-sm text-text-secondary">商品数量</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">最近同步</th>
                <th className="text-right p-4 font-title-sm text-text-secondary">操作</th>
              </tr>
            </thead>
            <tbody>
              {feeds.map((feed) => (
                <tr key={feed.id} className="border-b border-border-default last:border-0 hover:bg-surface-hover">
                  <td className="p-4"><Link href={`/feeds/${feed.id}`} className="text-primary-600 hover:underline font-medium">{feed.name}</Link></td>
                  <td className="p-4 text-text-secondary">{feed.format}</td>
                  <td className="p-4"><span className={`badge ${STATUS_BADGE[feed.status] || 'bg-surface-subtle text-text-muted'}`}>{feed.status}</span></td>
                  <td className="p-4 text-right">{feed.generated_rows ?? '-'}</td>
                  <td className="p-4 text-text-muted">{formatDate(feed.last_generated_at || '')}</td>
                  <td className="p-4 text-right space-x-2">
                    <button onClick={() => startEdit(feed)} className="btn btn-ghost text-sm">编辑</button>
                    <button onClick={() => toggleStatus(feed)} className="btn btn-ghost text-sm" title={feed.status === 'active' ? '暂停' : '启用'}>
                      {feed.status === 'active' ? '暂停' : '启用'}
                    </button>
                    <button onClick={() => syncMutation.mutate(feed.id)} disabled={syncMutation.isPending} className="btn btn-ghost text-sm">同步</button>
                    <button onClick={() => confirmDelete(feed.id)} disabled={deleteMutation.isPending} className="btn btn-ghost text-sm text-danger-500">删除</button>
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
