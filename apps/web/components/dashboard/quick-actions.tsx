'use client'

import Link from 'next/link'
import { Plus, BarChart3, Bot, Store } from 'lucide-react'

const actions = [
  { icon: BarChart3, label: '看报表', href: '/reports', desc: '查看跨平台数据报表' },
  { icon: Bot, label: 'Ask Agent', href: '/agent', desc: '让 AI 帮你分析广告数据' },
  { icon: Store, label: '添加店铺', href: '/stores', desc: '绑定新的电商店铺' },
]

export function QuickActions() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      {actions.map((action) => (
        <Link key={action.href} href={action.href} className="card p-6 hover:shadow-md transition-shadow">
          <div className="flex items-start gap-4">
            <div className="p-2 rounded-md bg-primary-50 text-primary-600">
              <action.icon size={20} />
            </div>
            <div>
              <h3 className="font-title-md text-text-primary">{action.label}</h3>
              <p className="text-sm text-text-muted mt-1">{action.desc}</p>
            </div>
          </div>
        </Link>
      ))}
    </div>
  )
}
