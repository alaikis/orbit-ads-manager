'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import {
  LayoutDashboard,
  BarChart3,
  Settings,
  Bell,
  Bot,
  FileText,
  Package,
  Rss,
  ListChecks,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react'
import { useState } from 'react'
import { Topbar } from './topbar'

const menuGroups = [
  { title: '总览', items: [
    { href: '/', label: '工作台', icon: LayoutDashboard },
    { href: '/reports', label: '数据报表', icon: BarChart3 },
  ]},
  { title: '资产', items: [
    { href: '/products', label: '商品', icon: Package },
    { href: '/feeds', label: 'Feed', icon: Rss },
  ]},
  { title: '规则与自动化', items: [
    { href: '/rules', label: '规则管理', icon: ListChecks },
  ]},
  { title: 'AI 能力', items: [
    { href: '/agent', label: 'Agent 对话', icon: Bot },
    { href: '/agent/cards', label: '待确认动作', icon: FileText },
  ]},
  { title: '系统', items: [
    { href: '/settings', label: '系统设置', icon: Settings },
    { href: '/notifications', label: '通知中心', icon: Bell },
    { href: '/workspaces', label: '工作区管理', icon: Package },
    { href: '/settings/connections', label: '连接凭据', icon: Settings },
  ]},
]

export function AppShell({ children }: { children: React.ReactNode }) {
  const [collapsed, setCollapsed] = useState(false)

  return (
    <div className="flex min-h-screen bg-surface-page">
      <aside
        className={`fixed inset-y-0 left-0 z-10 bg-white border-r border-border-default transition-all duration-200 ${
          collapsed ? 'w-16' : 'w-60'
        }`}
      >
        <div className="flex h-14 items-center justify-between px-4 border-b border-border-default">
          {!collapsed && <span className="font-title-lg text-text-primary">Orbit</span>}
          <button
            onClick={() => setCollapsed(!collapsed)}
            className="p-1.5 rounded-md hover:bg-surface-hover text-text-muted"
          >
            {collapsed ? <ChevronRight size={18} /> : <ChevronLeft size={18} />}
          </button>
        </div>
        <nav className="p-2 space-y-6 overflow-y-auto">
          {menuGroups.map((group) => (
            <div key={group.title}>
              {!collapsed && (
                <h3 className="px-2 mb-2 text-xs font-semibold text-text-muted uppercase tracking-wider">
                  {group.title}
                </h3>
              )}
              <div className="space-y-1">
                {group.items.map((item) => (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={`flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors ${
                      collapsed ? 'justify-center' : ''
                    } hover:bg-surface-hover text-text-secondary`}
                    title={collapsed ? item.label : undefined}
                  >
                    <item.icon size={18} />
                    {!collapsed && <span>{item.label}</span>}
                  </Link>
                ))}
              </div>
            </div>
          ))}
        </nav>
      </aside>
      <div className={`flex-1 transition-all duration-200 ${collapsed ? 'ml-16' : 'ml-60'}`}>
        <Topbar />
        <main className="p-6">{children}</main>
      </div>
    </div>
  )
}
