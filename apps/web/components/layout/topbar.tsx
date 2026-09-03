'use client'

import { Bell, Store } from 'lucide-react'
import Link from 'next/link'
import { useState, useRef, useEffect } from 'react'
import { useAuthStore } from '@/lib/auth-store'

export function Topbar() {
  const user = useAuthStore((s: any) => s.user)
  const logout = useAuthStore((s: any) => s.logout)
  const [menuOpen, setMenuOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setMenuOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  return (
    <header className="h-14 bg-white border-b border-border-default flex items-center justify-between px-6">
      <div className="flex items-center gap-4">
        <h1 className="font-title-lg text-text-primary">Orbit</h1>
        <nav className="hidden md:flex items-center gap-1 text-sm text-text-secondary">
          <Link href="/" className="hover:text-text-primary">工作台</Link>
          <span className="text-text-muted">/</span>
          <span className="text-text-primary">当前页面</span>
        </nav>
      </div>
      <div className="flex items-center gap-4">
        <div className="flex items-center gap-2 text-sm">
          <Store size={16} className="text-text-muted" />
          <select className="input h-8 text-sm w-40">
            <option>默认店铺</option>
          </select>
        </div>
        <button className="relative p-2 rounded-md hover:bg-surface-hover">
          <Bell size={18} className="text-text-secondary" />
        </button>
        <div className="relative" ref={menuRef}>
          <button
            onClick={() => setMenuOpen(!menuOpen)}
            className="flex items-center gap-2 px-2 py-1 rounded-md hover:bg-surface-hover"
          >
            <div className="w-8 h-8 rounded-full bg-primary-100 text-primary-700 flex items-center justify-center text-sm font-medium">
              {user?.email?.[0]?.toUpperCase() || 'U'}
            </div>
            <span className="text-sm text-text-secondary hidden sm:block">{user?.email || 'User'}</span>
            <span className="text-text-muted text-xs">▼</span>
          </button>
          {menuOpen && (
            <div className="absolute right-0 top-full mt-1 w-48 bg-white border border-border-default rounded-md shadow-lg py-1 z-50">
              <div className="px-3 py-2 border-b border-border-default">
                <div className="text-sm font-medium text-text-primary">{user?.email}</div>
                <div className="text-xs text-text-muted">{user?.role}</div>
              </div>
              <button
                onClick={() => logout()}
                className="w-full flex items-center gap-2 px-3 py-2 text-sm text-text-primary hover:bg-surface-hover"
              >
                登出
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  )
}
