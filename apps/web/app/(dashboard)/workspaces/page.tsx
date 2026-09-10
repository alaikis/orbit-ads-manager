'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { workspaceService } from '@/lib/api'
import { useState } from 'react'
import { Plus, RefreshCw } from 'lucide-react'
import Link from 'next/link'

type Workspace = { id: number; name: string; plan: string; status: string; created_at?: string }

const PLAN_BADGE: Record<string, string> = { beta: 'badge-info', pro: 'badge-primary', enterprise: 'badge-warning' }

export default function WorkspacesPage() {
  const queryClient = useQueryClient()
  const [showForm, setShowForm] = useState(false)
  const [name, setName] = useState('')
  const [plan, setPlan] = useState('beta')

  const { data, isLoading, error, refetch } = useQuery({ queryKey: ['workspaces'], queryFn: () => workspaceService.list() })
  const workspaces: Workspace[] = (data as any)?.items || []

  const createMutation = useMutation({
    mutationFn: (payload: { name: string; plan: string }) => workspaceService.create(payload),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['workspaces'] }); setShowForm(false); setName(''); setPlan('beta') },
    onError: (err: unknown) => alert('创建失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => workspaceService.delete(id),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['workspaces'] }) },
    onError: (err: unknown) => alert('删除失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const handleCreate = () => { if (!name.trim()) return; createMutation.mutate({ name, plan }) }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">工作区管理</h1>
          <p className="text-sm text-text-muted mt-1">管理工作区、计划与团队隔离</p>
        </div>
        <div className="flex gap-2">
          <button onClick={() => refetch()} className="btn btn-secondary"><RefreshCw size={16} /></button>
          <button onClick={() => setShowForm(!showForm)} className="btn btn-primary"><Plus size={16} /> 新建工作区</button>
        </div>
      </div>

      {showForm && (
        <div className="card p-6 space-y-4">
          <h2 className="font-title-md text-text-primary">新建工作区</h2>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-text-secondary mb-1.5">名称</label>
              <input className="input" value={name} onChange={(e) => setName(e.target.value)} placeholder="工作区名称" />
            </div>
            <div>
              <label className="block text-sm font-medium text-text-secondary mb-1.5">计划</label>
              <select className="input" value={plan} onChange={(e) => setPlan(e.target.value)}>
                <option value="beta">Beta</option>
                <option value="pro">Pro</option>
                <option value="enterprise">Enterprise</option>
              </select>
            </div>
          </div>
          <div className="flex gap-2">
            <button onClick={handleCreate} className="btn btn-primary" disabled={createMutation.isPending}>创建</button>
            <button onClick={() => setShowForm(false)} className="btn btn-secondary">取消</button>
          </div>
        </div>
      )}

      {error && (<div className="card p-6 text-danger-500">加载失败：{(error as Error).message}</div>)}

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-surface-subtle">
            <tr>
              <th className="text-left px-4 py-3 font-medium text-text-secondary">ID</th>
              <th className="text-left px-4 py-3 font-medium text-text-secondary">名称</th>
              <th className="text-left px-4 py-3 font-medium text-text-secondary">计划</th>
              <th className="text-left px-4 py-3 font-medium text-text-secondary">状态</th>
              <th className="text-left px-4 py-3 font-medium text-text-secondary">创建时间</th>
              <th className="text-right px-4 py-3 font-medium text-text-secondary">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border-default">
            {isLoading ? (
              Array.from({ length: 3 }).map((_, i) => (
                <tr key={i}>
                  <td className="px-4 py-3"><div className="h-4 w-8 rounded bg-surface-subtle animate-pulse" /></td>
                  <td className="px-4 py-3"><div className="h-4 w-32 rounded bg-surface-subtle animate-pulse" /></td>
                  <td className="px-4 py-3"><div className="h-5 w-16 rounded-full bg-surface-subtle animate-pulse" /></td>
                  <td className="px-4 py-3"><div className="h-5 w-14 rounded-full bg-surface-subtle animate-pulse" /></td>
                  <td className="px-4 py-3"><div className="h-4 w-40 rounded bg-surface-subtle animate-pulse" /></td>
                  <td className="px-4 py-3 text-right"><div className="ml-auto h-7 w-14 rounded bg-surface-subtle animate-pulse" /></td>
                </tr>
              ))
            ) : workspaces.length > 0 ? (
              workspaces.map((ws) => (
                <tr key={ws.id} className="hover:bg-surface-hover">
                  <td className="px-4 py-3 text-text-secondary">{ws.id}</td>
                  <td className="px-4 py-3"><Link href={`/workspaces/${ws.id}`} className="font-medium text-primary-600 hover:underline">{ws.name}</Link></td>
                  <td className="px-4 py-3"><span className={PLAN_BADGE[ws.plan] || 'badge-info'}>{ws.plan}</span></td>
                  <td className="px-4 py-3"><span className={`badge ${ws.status === 'active' ? 'bg-success-bg text-success' : 'bg-surface-subtle text-text-muted'}`}>{ws.status}</span></td>
                  <td className="px-4 py-3 text-text-secondary">{ws.created_at ? new Date(ws.created_at).toLocaleString('zh-CN') : '-'}</td>
                  <td className="px-4 py-3 text-right"><button onClick={() => deleteMutation.mutate(ws.id)} className="px-2 py-1 text-danger-500 hover:bg-danger-bg rounded-md text-xs transition-colors">删除</button></td>
                </tr>
              ))
            ) : (
              <tr><td colSpan={6} className="px-4 py-12 text-center text-text-muted">暂无工作区</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
