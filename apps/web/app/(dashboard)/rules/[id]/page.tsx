'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { ruleService } from '@/lib/api'
import Link from 'next/link'
import { useState } from 'react'

const STATUS_BADGE: Record<string, string> = {
  active: 'bg-success-bg text-success',
  inactive: 'bg-surface-subtle text-text-muted',
}

export default function RuleDetailPage({ params }: { params: { id: string } }) {
  const ruleId = Number(params.id)
  const queryClient = useQueryClient()

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['rule', ruleId],
    queryFn: () => ruleService.get(ruleId),
  })

  const rule = data as any

  const [editing, setEditing] = useState(false)
  const [editName, setEditName] = useState('')
  const [editCondition, setEditCondition] = useState('')
  const [editAction, setEditAction] = useState('')
  const [editCooldown, setEditCooldown] = useState('')
  const [actionError, setActionError] = useState('')

  const updateMutation = useMutation({
    mutationFn: (payload: { name?: string; condition?: string; action?: string; cooldown_minutes?: number }) =>
      ruleService.update(ruleId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rule', ruleId] })
      queryClient.invalidateQueries({ queryKey: ['rules'] })
      setEditing(false)
      setActionError('')
    },
    onError: (err: unknown) => setActionError(err instanceof Error ? err.message : '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: () => ruleService.delete(ruleId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rules'] })
      window.location.href = '/rules'
    },
    onError: (err: unknown) => alert('删除失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const startEdit = () => {
    if (!rule) return
    setEditName(rule.name)
    setEditCondition(rule.condition)
    setEditAction(rule.action)
    setEditCooldown(rule.cooldown_minutes != null ? String(rule.cooldown_minutes) : '')
    setEditing(true)
    setActionError('')
  }

  const saveEdit = () => {
    const payload: { name?: string; condition?: string; action?: string; cooldown_minutes?: number } = {}
    if (editName.trim()) payload.name = editName.trim()
    if (editCondition.trim()) payload.condition = editCondition.trim()
    if (editAction.trim()) payload.action = editAction.trim()
    if (editCooldown.trim()) payload.cooldown_minutes = parseInt(editCooldown, 10)
    updateMutation.mutate(payload)
  }

  const confirmDelete = () => {
    if (window.confirm('确定要删除这个规则吗？')) {
      deleteMutation.mutate()
    }
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/rules" className="btn btn-ghost">返回规则列表</Link>
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
          <Link href="/rules" className="btn btn-ghost">返回规则列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">规则详情</h1>
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
          <Link href="/rules" className="btn btn-ghost">返回规则列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">{rule?.name || '规则详情'}</h1>
            <p className="text-sm text-text-muted mt-1">状态: {rule?.status || '-'}</p>
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
              <td className="p-4 text-text-muted w-40">规则 ID</td>
              <td className="p-4 font-mono">{rule?.id ?? '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">规则名称</td>
              <td className="p-4 font-medium">{rule?.name || '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">状态</td>
              <td className="p-4">
                <span className={`badge ${STATUS_BADGE[rule?.status || ''] || 'bg-surface-subtle text-text-muted'}`}>
                  {rule?.status === 'active' ? '启用' : '禁用'}
                </span>
              </td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">条件</td>
              <td className="p-4">
                <pre className="font-mono text-xs bg-surface-subtle p-3 rounded-md overflow-auto">
                  {rule?.condition ? JSON.stringify(JSON.parse(rule.condition), null, 2) : '-'}
                </pre>
              </td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">动作</td>
              <td className="p-4">
                <pre className="font-mono text-xs bg-surface-subtle p-3 rounded-md overflow-auto">
                  {rule?.action ? JSON.stringify(JSON.parse(rule.action), null, 2) : '-'}
                </pre>
              </td>
            </tr>
            <tr>
              <td className="p-4 text-text-muted">创建时间</td>
              <td className="p-4">{rule?.created_at ? new Date(rule.created_at).toLocaleString('zh-CN') : '-'}</td>
            </tr>
          </tbody>
        </table>
      </div>

      {editing && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">编辑规则</h2>
              <button onClick={() => setEditing(false)} className="text-text-muted hover:text-text-primary"><span className="text-xl leading-none">×</span></button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">规则名称</label>
                <input className="input" value={editName} onChange={(e) => setEditName(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">条件（JSON）</label>
                <textarea className="input h-24 font-mono text-sm" value={editCondition} onChange={(e) => setEditCondition(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">动作（JSON）</label>
                <textarea className="input h-24 font-mono text-sm" value={editAction} onChange={(e) => setEditAction(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">冷却时间（分钟）</label>
                <input className="input" type="number" value={editCooldown} onChange={(e) => setEditCooldown(e.target.value)} />
              </div>
              {actionError && (<div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{actionError}</div>)}
              <div className="flex gap-2 justify-end">
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