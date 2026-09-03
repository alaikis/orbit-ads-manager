'use client'

import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'
import { useState } from 'react'
import { RefreshCw } from 'lucide-react'

type Tab = 'profile' | 'members' | 'notifications' | 'auto_approve' | 'security'

export default function SettingsPage() {
  const [tab, setTab] = useState<Tab>('profile')
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['tenant'],
    queryFn: () => api.get('/tenants/1'),
  })

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
        <button onClick={() => refetch()} className="btn btn-secondary"><RefreshCw size={16} /></button>
      </div>
      <div className="flex gap-6">
        <nav className="w-48 space-y-1">
          {[
            { id: 'profile', label: '租户资料' },
            { id: 'members', label: '成员与角色' },
            { id: 'notifications', label: '通知偏好' },
            { id: 'auto_approve', label: '自动执行' },
            { id: 'security', label: '安全' },
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
                <div className="space-y-4">
                  <h2 className="font-title-md text-text-primary">租户资料</h2>
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-1.5">租户名称</label>
                    <input className="input" placeholder="输入租户名称" defaultValue={(data as any)?.name || ''} />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-1.5">时区</label>
                    <select className="input">
                      <option>Asia/Shanghai</option>
                      <option>America/New_York</option>
                      <option>Europe/London</option>
                    </select>
                  </div>
                  <button className="btn btn-primary">保存</button>
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
                    <button className="relative inline-flex h-6 w-11 items-center rounded-full bg-surface-subtle transition-colors">
                      <span className="inline-block h-4 w-4 transform rounded-full bg-white transition-transform translate-x-1" />
                    </button>
                  </div>
                </div>
              )}
              {tab === 'security' && (
                <div className="space-y-4">
                  <h2 className="font-title-md text-text-primary">安全</h2>
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-1.5">当前密码</label>
                    <input type="password" className="input" placeholder="••••••••" />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-1.5">新密码</label>
                    <input type="password" className="input" placeholder="••••••••" />
                  </div>
                  <button className="btn btn-primary">修改密码</button>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
