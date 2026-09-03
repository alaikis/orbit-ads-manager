'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { providerService, workspaceService } from '@/lib/api'
import { useState, useEffect } from 'react'
import { Plus, RefreshCw, Play } from 'lucide-react'

type Provider = {
  id: number
  type: string
  name: string
  status: string
  workspace_id?: number
  tenant_id?: number
  config?: Record<string, unknown>
  created_at: string
  updated_at: string
}

type Workspace = {
  id: number
  name: string
}

type FieldDef = {
  key: string
  label: string
  type: 'text' | 'password' | 'number' | 'select'
  placeholder?: string
  options?: { value: string; label: string }[]
  required?: boolean
  hint?: string
}

const PROVIDER_SCHEMAS: Record<string, { label: string; description: string; fields: FieldDef[] }> = {
  llm: {
    label: 'LLM 模型',
    description: '用于 AI Agent 推理的大型语言模型',
    fields: [
      { key: 'api_key', label: 'API Key', type: 'password', placeholder: 'sk-...', required: true, hint: 'OpenAI、Anthropic、DeepSeek 等 LLM 服务商的 API Key' },
      { key: 'base_url', label: 'Base URL', type: 'text', placeholder: 'https://api.openai.com/v1', hint: 'API 服务的根地址，留空使用默认' },
      { key: 'model', label: '模型名称', type: 'text', placeholder: 'gpt-4o-mini', required: true, hint: '如：gpt-4o, gpt-4o-mini, claude-3-5-sonnet, deepseek-chat' },
    ],
  },
  smtp: {
    label: 'SMTP 邮件',
    description: '用于发送系统通知邮件',
    fields: [
      { key: 'host', label: 'SMTP 服务器', type: 'text', placeholder: 'smtp.gmail.com', required: true },
      { key: 'port', label: '端口', type: 'text', placeholder: '587', required: true, hint: '常用：587（TLS）、465（SSL）、25（明文）' },
      { key: 'user', label: '用户名', type: 'text', placeholder: 'your@email.com', required: true },
      { key: 'pass', label: '密码', type: 'password', placeholder: '应用专用密码', required: true, hint: 'Gmail 等服务需要使用应用专用密码' },
    ],
  },
  google: {
    label: 'Google Ads',
    description: 'Google Ads 平台接入凭证',
    fields: [
      { key: 'client_id', label: 'Client ID', type: 'text', placeholder: '*.apps.googleusercontent.com', required: true },
      { key: 'client_secret', label: 'Client Secret', type: 'password', required: true },
      { key: 'developer_token', label: 'Developer Token', type: 'password', required: true, hint: 'Google Ads API 开发者令牌' },
      { key: 'redirect_uri', label: 'Redirect URI', type: 'text', placeholder: 'https://your-domain.com/oauth/callback', hint: 'OAuth 回调地址' },
    ],
  },
  meta: {
    label: 'Meta Ads',
    description: 'Meta（Facebook / Instagram）Ads 接入凭证',
    fields: [
      { key: 'app_id', label: 'App ID', type: 'text', required: true },
      { key: 'app_secret', label: 'App Secret', type: 'password', required: true },
      { key: 'redirect_uri', label: 'Redirect URI', type: 'text', placeholder: 'https://your-domain.com/oauth/callback' },
    ],
  },
  bing: {
    label: 'Bing Ads',
    description: 'Microsoft Advertising 接入凭证',
    fields: [
      { key: 'client_id', label: 'Client ID', type: 'text', required: true },
      { key: 'client_secret', label: 'Client Secret', type: 'password', required: true },
      { key: 'developer_token', label: 'Developer Token', type: 'password', required: true },
    ],
  },
}

export default function ProvidersPage() {
  const queryClient = useQueryClient()
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Provider | null>(null)
  const [type, setType] = useState<keyof typeof PROVIDER_SCHEMAS>('llm')
  const [name, setName] = useState('')
  const [workspaceId, setWorkspaceId] = useState<number | ''>('')
  const [formValues, setFormValues] = useState<Record<string, string>>({})
  const [visibleFields, setVisibleFields] = useState<Record<string, boolean>>({})
  const [formError, setFormError] = useState('')

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['providers'],
    queryFn: () => providerService.list(),
  })

  const { data: workspacesData } = useQuery({
    queryKey: ['workspaces-list'],
    queryFn: () => workspaceService.list(),
  })

  const providers: Provider[] = data?.data || []
  const workspaces: Workspace[] = workspacesData?.data || []

  const createMutation = useMutation({
    mutationFn: (payload: { type: string; name: string; workspace_id?: number; config: Record<string, unknown> }) =>
      providerService.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['providers'] })
      closeForm()
    },
    onError: (err: unknown) => setFormError(err instanceof Error ? err.message : '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: { name?: string; status?: string; config?: Record<string, unknown> } }) =>
      providerService.update(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['providers'] })
      closeForm()
    },
    onError: (err: unknown) => setFormError(err instanceof Error ? err.message : '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => providerService.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['providers'] }),
    onError: (err: unknown) => alert('删除失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const testMutation = useMutation({
    mutationFn: (id: number) => providerService.test(id),
    onSuccess: (res: any) => alert(res?.message || '连接测试成功'),
    onError: (err: unknown) => alert('连接测试失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const openCreate = () => {
    setEditing(null)
    setType('llm')
    setName('')
    setWorkspaceId('')
    setFormValues({})
    setVisibleFields({})
    setFormError('')
    setShowForm(true)
  }

  const openEdit = (p: Provider) => {
    setEditing(p)
    setType(p.type as keyof typeof PROVIDER_SCHEMAS)
    setName(p.name)
    setWorkspaceId(p.workspace_id || '')
    setFormValues({})
    setVisibleFields({})
    setFormError('')
    setShowForm(true)
  }

  const closeForm = () => {
    setShowForm(false)
    setEditing(null)
    setFormError('')
  }

  const handleSubmit = () => {
    setFormError('')
    if (!name.trim()) {
      setFormError('请填写名称')
      return
    }

    const schema = PROVIDER_SCHEMAS[type]
    const config: Record<string, string> = {}
    for (const field of schema.fields) {
      const value = formValues[field.key]?.trim()
      if (field.required && !value && !editing) {
        setFormError(`请填写：${field.label}`)
        return
      }
      if (value) {
        config[field.key] = value
      }
    }

    if (editing) {
      const payload: { name?: string; status?: string; config?: Record<string, unknown> } = { name: name.trim() }
      if (Object.keys(config).length > 0) {
        payload.config = config
      }
      updateMutation.mutate({ id: editing.id, payload })
    } else {
      createMutation.mutate({
        type,
        name: name.trim(),
        workspace_id: workspaceId === '' ? undefined : Number(workspaceId),
        config,
      })
    }
  }

  const getProviderLabel = (t: string) => PROVIDER_SCHEMAS[t]?.label || t
  const getProviderDescription = (t: string) => PROVIDER_SCHEMAS[t]?.description || ''

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">Provider 配置</h1>
          <p className="text-sm text-text-muted mt-1">管理 LLM、SMTP、广告平台等外部服务连接</p>
        </div>
        <div className="flex gap-2">
          <button onClick={() => refetch()} className="btn btn-secondary"><RefreshCw size={16} /></button>
          <button onClick={openCreate} className="btn btn-primary"><Plus size={16} className="mr-2" />新建 Provider</button>
        </div>
      </div>

      {showForm && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4 overflow-y-auto">
          <div className="bg-white rounded-lg p-6 w-full max-w-2xl my-8">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">
                {editing ? '编辑 Provider' : '新建 Provider'}
              </h2>
              <button onClick={closeForm} className="text-text-muted hover:text-text-primary">
                <span className="text-xl leading-none">×</span>
              </button>
            </div>

            <div className="space-y-4">
              {!editing && (
                <div>
                  <label className="block text-sm font-medium text-text-secondary mb-1.5">类型 *</label>
                  <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
                    {Object.entries(PROVIDER_SCHEMAS).map(([key, schema]) => (
                      <button
                        key={key}
                        type="button"
                        onClick={() => { setType(key as any); setFormValues({}); setFormError('') }}
                        className={`p-3 rounded-md border text-left transition-colors ${
                          type === key
                            ? 'border-primary-500 bg-primary-50'
                            : 'border-border-default hover:border-primary-300 hover:bg-surface-hover'
                        }`}
                      >
                        <div className="text-sm font-medium text-text-primary">{schema.label}</div>
                        <div className="text-xs text-text-muted mt-0.5">{schema.description}</div>
                      </button>
                    ))}
                  </div>
                </div>
              )}

              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">名称 *</label>
                <input
                  className="input"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder={`例如：我的 ${getProviderLabel(type)}`}
                />
              </div>

              {!editing && workspaces.length > 0 && (
                <div>
                  <label className="block text-sm font-medium text-text-secondary mb-1.5">关联工作空间（可选）</label>
                  <select
                    className="input"
                    value={workspaceId}
                    onChange={(e) => setWorkspaceId(e.target.value ? Number(e.target.value) : '')}
                  >
                    <option value="">不关联（全局）</option>
                    {workspaces.map((ws) => (
                      <option key={ws.id} value={ws.id}>{ws.name}</option>
                    ))}
                  </select>
                </div>
              )}

              <div className="border-t border-border-default pt-4">
                <div className="text-sm font-medium text-text-secondary mb-3">
                  {getProviderLabel(type)} 配置
                </div>
                <div className="space-y-4">
                  {PROVIDER_SCHEMAS[type].fields.map((field) => (
                    <div key={field.key}>
                      <label className="block text-sm font-medium text-text-secondary mb-1.5">
                        {field.label} {field.required && !editing && '*'}
                      </label>
                      <div className="relative">
                        <input
                          className="input pr-10"
                          type={field.type === 'password' && !visibleFields[field.key] ? 'password' : 'text'}
                          value={formValues[field.key] || ''}
                          onChange={(e) => setFormValues({ ...formValues, [field.key]: e.target.value })}
                          placeholder={editing && field.type === 'password' ? '留空表示不修改' : field.placeholder}
                        />
                        {field.type === 'password' && (
                          <button
                            type="button"
                            onClick={() => setVisibleFields({ ...visibleFields, [field.key]: !visibleFields[field.key] })}
                            className="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-text-muted hover:text-text-primary"
                            tabIndex={-1}
                          >
                            {visibleFields[field.key] ? <span className="text-xs">隐藏</span> : <span className="text-xs">显示</span>}
                          </button>
                        )}
                      </div>
                      {field.hint && (
                        <p className="text-xs text-text-muted mt-1">{field.hint}</p>
                      )}
                    </div>
                  ))}
                </div>
              </div>

              {formError && (
                <div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{formError}</div>
              )}

              <div className="flex gap-2 justify-end pt-2 border-t border-border-default">
                <button onClick={closeForm} className="btn btn-secondary">取消</button>
                <button
                  onClick={handleSubmit}
                  disabled={createMutation.isPending || updateMutation.isPending}
                  className="btn btn-primary"
                >
                  {(createMutation.isPending || updateMutation.isPending) ? (
                    <>
                      <span className="mr-2 inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                      保存中...
                    </>
                  ) : (
                    '保存'
                  )}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {error && (
        <div className="card p-6 text-danger-500">加载失败：{(error as Error).message}</div>
      )}

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-surface-subtle">
            <tr>
              <th className="text-left px-4 py-3 font-medium text-text-secondary">ID</th>
              <th className="text-left px-4 py-3 font-medium text-text-secondary">类型</th>
              <th className="text-left px-4 py-3 font-medium text-text-secondary">名称</th>
              <th className="text-left px-4 py-3 font-medium text-text-secondary">状态</th>
              <th className="text-left px-4 py-3 font-medium text-text-secondary">创建时间</th>
              <th className="text-right px-4 py-3 font-medium text-text-secondary">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border-default">
            {isLoading ? (
              <tr><td colSpan={6} className="px-4 py-8 text-center text-text-muted">加载中...</td></tr>
            ) : providers.length > 0 ? (
              providers.map((p) => (
                <tr key={p.id} className="hover:bg-surface-hover">
                  <td className="px-4 py-3">{p.id}</td>
                  <td className="px-4 py-3">
                    <div>
                      <div className="font-medium text-text-primary">{getProviderLabel(p.type)}</div>
                      <div className="text-xs text-text-muted">{getProviderDescription(p.type)}</div>
                    </div>
                  </td>
                  <td className="px-4 py-3 font-medium text-text-primary">{p.name}</td>
                  <td className="px-4 py-3">
                    <span className={`badge ${p.status === 'active' ? 'bg-success-bg text-success' : 'bg-danger-bg text-danger'}`}>
                      {p.status === 'active' ? '已启用' : '已停用'}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-text-secondary">{new Date(p.created_at).toLocaleString('zh-CN')}</td>
                  <td className="px-4 py-3 text-right">
                    <div className="flex justify-end gap-2">
                      <button onClick={() => testMutation.mutate(p.id)} disabled={testMutation.isPending} className="p-2 text-primary-600 hover:bg-primary-50 rounded-md" title="测试连接">
                        <Play size={16} />
                      </button>
                      <button onClick={() => openEdit(p)} className="px-2 py-1 text-text-secondary hover:bg-surface-hover rounded-md text-xs">
                        编辑
                      </button>
                      <button
                        onClick={() => { if (confirm('确定删除该 Provider？')) deleteMutation.mutate(p.id) }}
                        className="p-2 text-danger-500 hover:bg-danger-50 rounded-md"
                        title="删除"
                      >
                        <span className="text-base leading-none">🗑</span>
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            ) : (
              <tr><td colSpan={6} className="px-4 py-8 text-center text-text-muted">暂无 Provider</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
