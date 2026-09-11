'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { feedService } from '@/lib/api'
import Link from 'next/link'
import { useState } from 'react'

const FORMAT_LABEL: Record<string, string> = {
  xml: 'XML',
  csv: 'CSV',
  json: 'JSON',
}

const STATUS_BADGE: Record<string, string> = {
  active: 'bg-success-bg text-success',
  paused: 'bg-warning-bg text-warning',
  error: 'bg-danger-bg text-danger',
}

const FORMAT_OPTIONS = ['xml', 'csv', 'json']

export default function FeedDetailPage({ params }: { params: { id: string } }) {
  const feedId = Number(params.id)
  const queryClient = useQueryClient()

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['feed', feedId],
    queryFn: () => feedService.get(feedId),
  })

  const feed = data as any

  const [editing, setEditing] = useState(false)
  const [editName, setEditName] = useState('')
  const [editFormat, setEditFormat] = useState('xml')
  const [actionError, setActionError] = useState('')

  const updateMutation = useMutation({
    mutationFn: (payload: { name?: string; format?: string; status?: string }) =>
      feedService.update(feedId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feed', feedId] })
      queryClient.invalidateQueries({ queryKey: ['feeds'] })
      setEditing(false)
      setActionError('')
    },
    onError: (err: unknown) => setActionError(err instanceof Error ? err.message : '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: () => feedService.delete(feedId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feeds'] })
      window.location.href = '/feeds'
    },
    onError: (err: unknown) => alert('删除失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const startEdit = () => {
    if (!feed) return
    setEditName(feed.name)
    setEditFormat(feed.format)
    setEditing(true)
    setActionError('')
  }

  const saveEdit = () => {
    const payload: { name?: string; format?: string } = {}
    if (editName.trim()) payload.name = editName.trim()
    if (editFormat.trim()) payload.format = editFormat
    updateMutation.mutate(payload)
  }

  const confirmDelete = () => {
    if (window.confirm('确定要删除这个 Feed 吗？')) {
      deleteMutation.mutate()
    }
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/feeds" className="btn btn-ghost">返回 Feed 列表</Link>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">{(error as Error).message}</p>
          <button onClick={() => refetch()} className="btn btn-primary">重试</button>
        </div>
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/feeds" className="btn btn-ghost">返回 Feed 列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">Feed 详情</h1>
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
              <div className="h-4 bg-surface-subtle rounded" />
              <div className="h-4 bg-surface-subtle rounded w-3/4" />
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link href="/feeds" className="btn btn-ghost">返回 Feed 列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">{feed?.name || 'Feed 详情'}</h1>
            <p className="text-sm text-text-muted mt-1">格式: {FORMAT_LABEL[feed?.format || ''] || feed?.format || '-'} | 状态: {feed?.status || '-'}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <button onClick={startEdit} className="btn btn-secondary">编辑</button>
          <button onClick={confirmDelete} className="btn btn-secondary text-danger-500">删除</button>
        </div>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <tbody>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted w-40">Feed ID</td>
              <td className="p-4 font-mono">{feed?.id ?? '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">Feed 名称</td>
              <td className="p-4">{feed?.name || '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">格式</td>
              <td className="p-4">{FORMAT_LABEL[feed?.format || ''] || feed?.format || '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">状态</td>
              <td className="p-4">
                <span className={`badge ${STATUS_BADGE[feed?.status || ''] || 'bg-surface-subtle text-text-muted'}`}>
                  {feed?.status || '-'}
                </span>
              </td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">商品数量</td>
              <td className="p-4">{feed?.generated_rows ?? '-'}</td>
            </tr>
            <tr>
              <td className="p-4 text-text-muted">最近同步时间</td>
              <td className="p-4">{feed?.last_generated_at || '-'}</td>
            </tr>
          </tbody>
        </table>
      </div>

      {editing && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">编辑 Feed</h2>
              <button onClick={() => setEditing(false)} className="text-text-muted hover:text-text-primary"><span className="text-xl leading-none">×</span></button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">Feed 名称</label>
                <input className="input" value={editName} onChange={(e) => setEditName(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">格式</label>
                <select className="input" value={editFormat} onChange={(e) => setEditFormat(e.target.value)}>
                  {FORMAT_OPTIONS.map((f) => (<option key={f} value={f}>{f.toUpperCase()}</option>))}
                </select>
              </div>
              {actionError && (<div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{actionError}</div>)}
              <div className="flex gap-2 justify-end pt-2">
                <button onClick={() => setEditing(false)} className="btn btn-secondary">取消</button>
                <button onClick={saveEdit} disabled={updateMutation.isPending} className="btn btn-primary">{updateMutation.isPending ? '保存中...' : '保存'}</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}