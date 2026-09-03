'use client'

import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useRouter, useSearchParams } from 'next/navigation'
import { connectionService, type PlatformSchema, type Connection } from '@/lib/api/services/connection.service'

const PLATFORM_ICONS: Record<string, string> = {
  woocommerce: '🛒',
  shopify: '🛍️',
  google_ads: '🔍',
  google_shopping: '🏬',
  meta: '📘',
  bing: '🅱️',
  tiktok: '🎵',
  llm: '🤖',
  smtp: '✉️',
}

const PLATFORM_GROUP: Record<string, string> = {
  woocommerce: 'store',
  shopify: 'store',
  google_ads: 'ad',
  google_shopping: 'ad',
  meta: 'ad',
  bing: 'ad',
  tiktok: 'ad',
  llm: 'ai',
  smtp: 'ai',
}

const GROUP_LABELS: Record<string, string> = {
  all: '全部',
  store: '店铺',
  ad: '广告平台',
  ai: 'AI / 邮件',
}

type Step = 'select' | 'fill' | 'review'

export default function ConnectionsPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const queryClient = useQueryClient()
  const [activeGroup, setActiveGroup] = useState('all')
  const [showModal, setShowModal] = useState(false)
  const [step, setStep] = useState<Step>('select')
  const [selectedPlatform, setSelectedPlatform] = useState<PlatformSchema | null>(null)
  const [formName, setFormName] = useState('')
  const [formFields, setFormFields] = useState<Record<string, string>>({})
  const [testResult, setTestResult] = useState<{ ok: boolean; message?: string } | null>(null)

  const { data: connsData, isLoading: loadingConns, refetch } = useQuery({
    queryKey: ['connections'],
    queryFn: () => connectionService.list(),
  })
  const connections = (connsData as { items?: Connection[] })?.items || []

  const { data: platformsData } = useQuery({
    queryKey: ['connection-platforms'],
    queryFn: () => connectionService.listPlatforms(),
  })
  const platforms = (platformsData as { items?: PlatformSchema[] })?.items || []

  const createMut = useMutation({
    mutationFn: (data: { platform: string; name: string; fields: Record<string, string> }) =>
      connectionService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['connections'] })
      closeModal()
    },
  })

  const deleteMut = useMutation({
    mutationFn: (id: number) => connectionService.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['connections'] }),
  })

  const oauthStartMut = useMutation({
    mutationFn: (data: { platform: string; name: string }) => connectionService.oauthStart(data),
    onSuccess: (resp: unknown) => {
      if (resp && typeof resp === 'object' && 'authorize_url' in resp) {
        window.location.href = (resp as { authorize_url: string }).authorize_url
      }
    },
  })

  useEffect(() => {
    if (searchParams.get('connected')) {
      refetch()
      router.replace('/settings/connections')
    }
    const error = searchParams.get('error')
    if (error) {
      const detail = searchParams.get('detail') || error
      alert(`OAuth 失败: ${detail}`)
      router.replace('/settings/connections')
    }
  }, [searchParams, refetch, router])

  const closeModal = () => {
    setShowModal(false)
    setStep('select')
    setSelectedPlatform(null)
    setFormName('')
    setFormFields({})
    setTestResult(null)
  }

  const startCreate = (platform: PlatformSchema) => {
    setSelectedPlatform(platform)
    setFormName(`${platform.label} 接入`)
    setFormFields({})
    if (platform.auth_flow === 'oauth') {
      setStep('review')
    } else {
      setStep('fill')
    }
  }

  const submit = async () => {
    if (!selectedPlatform) return
    if (selectedPlatform.auth_flow === 'oauth') {
      oauthStartMut.mutate({ platform: selectedPlatform.platform, name: formName })
    } else {
      createMut.mutate({ platform: selectedPlatform.platform, name: formName, fields: formFields })
    }
  }

  const filtered = activeGroup === 'all' ? connections : connections.filter((c) => PLATFORM_GROUP[c.platform] === activeGroup)
  const filteredPlatforms = activeGroup === 'all' ? platforms : platforms.filter((p) => PLATFORM_GROUP[p.platform] === activeGroup)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">统一凭据管理</h1>
          <p className="text-sm text-text-muted mt-1">按平台管理店铺、广告账户、AI、邮件等接入凭据。Token 经过 AES-256-GCM 加密。</p>
        </div>
        <button
          onClick={() => setShowModal(true)}
          className="btn btn-primary"
        >
          + 新建连接
        </button>
      </div>

      <div className="flex gap-2">
        {Object.entries(GROUP_LABELS).map(([id, label]) => (
          <button
            key={id}
            onClick={() => setActiveGroup(id)}
            className={`px-3 py-1.5 text-sm rounded-md transition-colors ${
              activeGroup === id
                ? 'bg-primary-100 text-primary-700 font-medium'
                : 'text-text-secondary hover:bg-surface-hover'
            }`}
          >
            {label}
          </button>
        ))}
      </div>

      {loadingConns ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="card p-5 animate-pulse h-32" />
          ))}
        </div>
      ) : filtered.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted mb-4">还没有{activeGroup === 'all' ? '' : GROUP_LABELS[activeGroup]}连接</p>
          <button onClick={() => setShowModal(true)} className="btn btn-primary">+ 新建连接</button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filtered.map((c) => {
            const schema = platforms.find((p) => p.platform === c.platform)
            return (
              <div key={c.id} className="card p-5 space-y-3">
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <span className="text-2xl">{PLATFORM_ICONS[c.platform] || '🔌'}</span>
                    <div>
                      <p className="font-medium text-text-primary">{c.name}</p>
                      <p className="text-xs text-text-muted">{schema?.label || c.platform}</p>
                    </div>
                  </div>
                  <span
                    className={`px-2 py-0.5 text-xs rounded-full ${
                      c.status === 'active'
                        ? 'bg-success-50 text-success-700'
                        : c.status === 'error'
                          ? 'bg-danger-50 text-danger-700'
                          : 'bg-surface-subtle text-text-muted'
                    }`}
                  >
                    {c.status}
                  </span>
                </div>
                {c.last_error && (
                  <p className="text-xs text-danger-600 line-clamp-2">{c.last_error}</p>
                )}
                <div className="flex items-center gap-2 pt-2 border-t border-border-subtle">
                  <button
                    onClick={async () => {
                      const res = await connectionService.test(c.id)
                      if (res) setTestResult({ ok: !!(res as { ok?: boolean }).ok, message: (res as { message?: string }).message })
                    }}
                    className="text-xs text-primary-600 hover:underline"
                  >
                    测试
                  </button>
                  <span className="text-border-default">·</span>
                  <button
                    onClick={async () => {
                      const res = await connectionService.sync(c.id)
                      alert(res ? '同步已触发' : '同步失败')
                      refetch()
                    }}
                    className="text-xs text-primary-600 hover:underline"
                  >
                    同步
                  </button>
                  <span className="text-border-default">·</span>
                  <button
                    onClick={() => {
                      if (confirm(`确定删除连接「${c.name}」？`)) {
                        deleteMut.mutate(c.id)
                      }
                    }}
                    className="text-xs text-danger-600 hover:underline"
                  >
                    删除
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      )}

      {showModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50" onClick={closeModal}>
          <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <div className="p-6 space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-lg font-semibold">
                  {step === 'select' && '选择平台'}
                  {step === 'fill' && `配置 ${selectedPlatform?.label}`}
                  {step === 'review' && `授权 ${selectedPlatform?.label}`}
                </h2>
                <button onClick={closeModal} className="text-text-muted hover:text-text-primary">
                  <span className="text-xl leading-none">×</span>
                </button>
              </div>

              {step === 'select' && (
                <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                  {filteredPlatforms.map((p) => (
                    <button
                      key={p.platform}
                      onClick={() => startCreate(p)}
                      className="p-4 border border-border-default rounded-lg hover:border-primary-400 hover:bg-primary-50 transition-colors text-left"
                    >
                      <span className="text-2xl block mb-2">{PLATFORM_ICONS[p.platform] || '🔌'}</span>
                      <p className="font-medium text-sm">{p.label}</p>
                      <p className="text-xs text-text-muted mt-1 line-clamp-2">{p.description}</p>
                      <p className="text-xs text-primary-600 mt-2">
                        {p.auth_flow === 'oauth' ? 'OAuth 授权' : p.auth_flow === 'developer_token' ? '开发者令牌' : 'API Key'}
                      </p>
                    </button>
                  ))}
                </div>
              )}

              {step === 'fill' && selectedPlatform && (
                <div className="space-y-3">
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-1.5">连接名称</label>
                    <input
                      className="input w-full"
                      value={formName}
                      onChange={(e) => setFormName(e.target.value)}
                      placeholder="如：主店铺"
                    />
                  </div>
                  {selectedPlatform.fields.map((f) => (
                    <div key={f.key}>
                      <label className="block text-sm font-medium text-text-secondary mb-1.5">
                        {f.label}
                        {f.required && <span className="text-danger-500"> *</span>}
                      </label>
                      <input
                        type={f.type === 'password' ? 'password' : 'text'}
                        className="input w-full"
                        value={formFields[f.key] || ''}
                        onChange={(e) => setFormFields({ ...formFields, [f.key]: e.target.value })}
                        placeholder={f.placeholder}
                      />
                      {f.hint && <p className="text-xs text-text-muted mt-1">{f.hint}</p>}
                    </div>
                  ))}
                  <div className="flex justify-between pt-3">
                    <button onClick={() => setStep('select')} className="btn btn-secondary">上一步</button>
                    <button
                      onClick={submit}
                      disabled={createMut.isPending || !formName}
                      className="btn btn-primary disabled:opacity-50"
                    >
                      {createMut.isPending ? '创建中...' : '创建'}
                    </button>
                  </div>
                  {createMut.error && (
                    <p className="text-sm text-danger-600">{(createMut.error as Error).message}</p>
                  )}
                </div>
              )}

              {step === 'review' && selectedPlatform && (
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-1.5">连接名称</label>
                    <input
                      className="input w-full"
                      value={formName}
                      onChange={(e) => setFormName(e.target.value)}
                    />
                  </div>
                  <div className="p-4 bg-primary-50 rounded-lg text-sm text-text-secondary space-y-2">
                    <p className="font-medium text-text-primary">即将跳转到 {selectedPlatform.label} 授权</p>
                    {selectedPlatform.scopes && selectedPlatform.scopes.length > 0 && (
                      <p>授权范围：{selectedPlatform.scopes.join(', ')}</p>
                    )}
                    <p className="text-xs">授权完成后浏览器会自动跳回本页。</p>
                  </div>
                  <div className="flex justify-between pt-3">
                    <button onClick={() => setStep('select')} className="btn btn-secondary">上一步</button>
                    <button
                      onClick={submit}
                      disabled={oauthStartMut.isPending || !formName}
                      className="btn btn-primary disabled:opacity-50"
                    >
                      {oauthStartMut.isPending ? '启动中...' : `前往 ${selectedPlatform.label} 授权 →`}
                    </button>
                  </div>
                  {oauthStartMut.error && (
                    <p className="text-sm text-danger-600">{(oauthStartMut.error as Error).message}</p>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {testResult && (
        <div
          className="fixed bottom-6 right-6 card p-4 max-w-sm shadow-lg"
          onClick={() => setTestResult(null)}
        >
          <p className={`text-sm font-medium ${testResult.ok ? 'text-success-700' : 'text-danger-700'}`}>
            {testResult.ok ? '连接成功' : '连接失败'}
          </p>
          {testResult.message && <p className="text-xs text-text-muted mt-1 line-clamp-3">{testResult.message}</p>}
        </div>
      )}
    </div>
  )
}
