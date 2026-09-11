'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'
import { useState } from 'react'

type Tab = 'profile' | 'members' | 'notifications' | 'auto_approve' | 'security' | 'connections'

export default function SettingsPage() {
  const queryClient = useQueryClient()
  const [tab, setTab] = useState<Tab>('profile')
  const [autoApprove, setAutoApprove] = useState(false)
  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['tenant'],
    queryFn: () => api.get('/tenants/current'),
  })

  const [saveMessage, setSaveMessage] = useState('')

  const updateMutation = useMutation({
    mutationFn: (payload: Record<string, unknown>) => api.patch('/tenants/current', payload),
    onSuccess: () => {
      setSaveMessage('保存成功')
      queryClient.invalidateQueries({ queryKey: ['tenant'] })
      setTimeout(() => setSaveMessage(''), 3000)
    },
    onError: () => setSaveMessage('保存失败'),
  })

  const handleSaveProfile = () => {
    const form = document.getElementById('settings-profile-form') as HTMLFormElement | null
    if (!form) return
    const formData = new FormData(form)
    updateMutation.mutate({
      name: formData.get('name'),
      timezone: formData.get('timezone'),
    })
  }

  const handleAutoApproveToggle = () => {
    const next = !autoApprove
    setAutoApprove(next)
    updateMutation.mutate({ auto_approve_high_risk: next })
  }

  const handleChangePassword = () => {
    if (!currentPassword || !newPassword) {
      setSaveMessage('请填写当前密码和新密码')
      return
    }
    updateMutation.mutate({
      current_password: currentPassword,
      new_password: newPassword,
    })
    setCurrentPassword('')
    setNewPassword('')
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">系统设置</h1>
          <p className="text-sm text-text-muted mt-1">管理租户、成员与系统偏好</p>
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
          <h1 className="text-2xl font-semibold text-text-primary">系统设置</h1>
          <p className="text-sm text-text-muted mt-1">管理租户、成员与系统偏好</p>
        </div>
        <div className="flex gap-2">
          <button onClick={() => refetch()} className="btn btn-secondary">刷新</button>
        </div>
      </div>
      <div className="flex gap-6">
        <nav className="w-48 space-y-1">
          {[
            { id: 'profile', label: '租户资料' },
            { id: 'members', label: '成员与角色' },
            { id: 'notifications', label: '通知偏好' },
            { id: 'auto_approve', label: '自动执行' },
            { id: 'security', label: '安全' },
            { id: 'connections', label: '连接凭据' },
          ].map((item) => (
            <button
              key={item.id}
              onClick={() => setTab(item.id as Tab)}
              className={`block w-full text-left px-3 py-2 rounded-md text-sm transition-colors ${
                tab === item.id ? 'bg-primary-50 text-primary-700 font-medium' : 'text-text-secondary hover:bg-surface-hover'
              }`}
            >
              {item.label}
            </button>
          ))}
        </nav>
        <div className="flex-1 card p-6">
          {isLoading ? (
            <div className="space-y-4">
              <div className="h-6 bg-surface-subtle rounded w-1/3 animate-pulse" />
              <div className="h-10 bg-surface-subtle rounded animate-pulse" />
              <div className="h-10 bg-surface-subtle rounded animate-pulse" />
            </div>
          ) : (
            <>
              {tab === 'profile' && (
                <div className="space-y-4" id="settings-profile-form">
                  <h2 className="font-title-md text-text-primary">租户资料</h2>
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-1.5">租户名称</label>
                    <input className="input" name="name" placeholder="输入租户名称" defaultValue={(data as any)?.name || ''} />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-1.5">时区</label>
                    <select className="input" name="timezone" defaultValue={(data as any)?.timezone || 'Asia/Shanghai'}>
                      <option value="Asia/Shanghai">Asia/Shanghai</option>
                      <option value="America/New_York">America/New_York</option>
                      <option value="Europe/London">Europe/London</option>
                    </select>
                  </div>
                  <button onClick={handleSaveProfile} disabled={updateMutation.isPending} className="btn btn-primary">{updateMutation.isPending ? '保存中...' : '保存'}</button>
                  {saveMessage && <p className="text-sm text-success-500">{saveMessage}</p>}
                </div>
              )}
              {tab === 'members' && (
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <h2 className="font-title-md text-text-primary">成员与角色</h2>
                    <button className="btn btn-primary text-sm">邀请成员</button>
                  </div>
                  <p className="text-sm text-text-muted">成员管理功能即将推出</p>
                </div>
              )}
              {tab === 'notifications' && (
                <div className="space-y-4">
                  <h2 className="font-title-md text-text-primary">通知偏好</h2>
                  <div className="space-y-3">
                    {['同步失败', 'Token 过期', '预算预警', 'Agent 待确认', '报表完成'].map((event) => (
                      <label key={event} className="flex items-center gap-2">
                        <input type="checkbox" className="rounded" defaultChecked />
                        <span className="text-sm text-text-primary">{event}</span>
                      </label>
                    ))}
                  </div>
                </div>
              )}
              {tab === 'auto_approve' && (
                <div className="space-y-4">
                  <h2 className="font-title-md text-text-primary">自动执行</h2>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium text-text-primary">允许高风险自动执行</p>
                      <p className="text-xs text-text-muted mt-1">开启后 Agent 和规则可自动执行高风险操作</p>
                    </div>
                    <button
                      onClick={handleAutoApproveToggle}
                      className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                        autoApprove ? 'bg-primary-500' : 'bg-surface-subtle'
                      }`}
                    >
                      <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${autoApprove ? 'translate-x-6' : 'translate-x-1'}`} />
                    </button>
                  </div>
                </div>
              )}
              {tab === 'security' && (
                <div className="space-y-4">
                  <h2 className="font-title-md text-text-primary">安全</h2>
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-1.5">当前密码</label>
                    <input type="password" className="input" placeholder="••••••••" value={currentPassword} onChange={(e) => setCurrentPassword(e.target.value)} />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-1.5">新密码</label>
                    <input type="password" className="input" placeholder="••••••••" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} />
                  </div>
                  <button onClick={handleChangePassword} className="btn btn-primary">修改密码</button>
                </div>
              )}
              {tab === 'connections' && (
                <div className="space-y-4">
                  <h2 className="font-title-md text-text-primary">连接凭据</h2>
                  <p className="text-sm text-text-muted">管理店铺、广告平台、AI、邮件等接入凭据。Token 经过 AES-256-GCM 加密存储。</p>
                  <a href="/settings/connections" className="btn btn-primary">前往连接管理</a>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
