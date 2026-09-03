'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { ruleService } from '@/lib/api'
import { Plus, RefreshCw } from 'lucide-react'
import { useState } from 'react'

type Rule = {
  id: number
  name: string
  condition: string
  action: string
  status: 'active' | 'inactive'
  enabled?: boolean
  cooldown_minutes?: number
  created_at?: string
}

export default function RulesPage() {
  const [showCreate, setShowCreate] = useState(false)
  const [name, setName] = useState('')
  const [condition, setCondition] = useState('{"metric":"roas","operator":"lt","value":1.5}')
  const [action, setAction] = useState('{"type":"pause","target":"self"}')
  const [createError, setCreateError] = useState('')
  const queryClient = useQueryClient()

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['rules'],
    queryFn: () => ruleService.list(),
  })

  const rules: Rule[] = (data as any)?.items || []

  const createMutation = useMutation({
    mutationFn: () => ruleService.create({ name, condition, action } as any),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rules'] })
      setShowCreate(false)
      setName('')
      setCreateError('')
    },
    onError: (err: unknown) => {
      setCreateError(err instanceof Error ? err.message : '创建失败')
    },
  })

  const handleCreate = () => {
    setCreateError('')
    if (!name.trim()) {
      setCreateError('请输入规则名称')
      return
    }
    try {
      JSON.parse(condition)
      JSON.parse(action)
    } catch {
      setCreateError('条件或动作 JSON 格式错误')
      return
    }
    createMutation.mutate()
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">规则管理</h1>
            <p className="text-sm text-text-muted mt-1">创建和管理自动化规则</p>
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
          <h1 className="text-2xl font-semibold text-text-primary">规则管理</h1>
          <p className="text-sm text-text-muted mt-1">创建和管理自动化规则</p>
        </div>
        <div className="flex gap-2">
          <button onClick={() => setShowCreate(true)} className="btn btn-primary"><Plus size={16} className="mr-2" />创建规则</button>
          <button onClick={() => refetch()} className="btn btn-secondary"><RefreshCw size={16} /></button>
        </div>
      </div>
      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (<div key={i} className="card p-4 animate-pulse h-16 bg-surface-subtle rounded" />))}
        </div>
      ) : rules.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted mb-4">还没有创建任何规则</p>
          <button onClick={() => setShowCreate(true)} className="btn btn-primary">创建第一条规则</button>
        </div>
      ) : (
        <div className="space-y-4">
          {rules.map((rule) => (
            <div key={rule.id} className="card p-4">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium text-text-primary">{rule.name}</h3>
                  <p className="text-sm text-text-muted mt-1">
                    状态: {rule.status === 'active' ? '已启用' : '已禁用'}
                    {rule.cooldown_minutes ? ` | 冷却: ${rule.cooldown_minutes} 分钟` : ''}
                  </p>
                </div>
                <span className={`badge ${rule.status === 'active' ? 'bg-success-bg text-success' : 'bg-surface-subtle text-text-muted'}`}>
                  {rule.status === 'active' ? '启用' : '禁用'}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
      {showCreate && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <h2 className="text-lg font-semibold mb-4">创建规则</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">规则名称</label>
                <input
                  className="input"
                  placeholder="例如：ROAS 低于阈值暂停"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">条件（JSON）</label>
                <textarea
                  className="input h-24 font-mono text-sm"
                  placeholder='{"metric":"roas","operator":"lt","value":1.5}'
                  value={condition}
                  onChange={(e) => setCondition(e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">动作（JSON）</label>
                <textarea
                  className="input h-24 font-mono text-sm"
                  placeholder='{"type":"pause","target":"self"}'
                  value={action}
                  onChange={(e) => setAction(e.target.value)}
                />
              </div>
              {createError && (
                <div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{createError}</div>
              )}
              <div className="flex gap-2 justify-end">
                <button onClick={() => { setShowCreate(false); setCreateError('') }} className="btn btn-secondary">取消</button>
                <button onClick={handleCreate} disabled={createMutation.isPending} className="btn btn-primary">
                  {createMutation.isPending ? '创建中...' : '创建'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
