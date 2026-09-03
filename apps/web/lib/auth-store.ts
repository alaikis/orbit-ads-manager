'use client'

import { create } from 'zustand'

export interface User {
  id: number
  email: string
  role: string
  tenant_id: number
}

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  setAuth: (user: User, accessToken: string, refreshToken: string) => void
  logout: () => void
  getAccessToken: () => string | null
  getRefreshToken: () => string | null
  isAuthenticated: () => boolean
}

const STORAGE_KEY = 'orbit-auth'

function saveAuth(state: { user: User | null; accessToken: string | null; refreshToken: string | null }) {
  if (typeof window === 'undefined') return
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
  } catch {
    // ignore
  }
}

function loadFromStorage(): { user: User | null; accessToken: string | null; refreshToken: string | null } {
  if (typeof window === 'undefined') return { user: null, accessToken: null, refreshToken: null }
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      const parsed = JSON.parse(stored)
      return {
        user: parsed.user || null,
        accessToken: parsed.accessToken || null,
        refreshToken: parsed.refreshToken || null,
      }
    }
  } catch {
    // ignore
  }
  return { user: null, accessToken: null, refreshToken: null }
}

const initial = loadFromStorage()

// @ts-ignore
export const useAuthStore: any = create((set: any, get: any) => ({
  user: initial.user,
  accessToken: initial.accessToken,
  refreshToken: initial.refreshToken,
  setAuth: (user: User, accessToken: string, refreshToken: string) => {
    saveAuth({ user, accessToken, refreshToken })
    set({ user, accessToken, refreshToken })
  },
  logout: () => {
    saveAuth({ user: null, accessToken: null, refreshToken: null })
    set({ user: null, accessToken: null, refreshToken: null })
    if (typeof window !== 'undefined') {
      window.location.href = '/login'
    }
  },
  getAccessToken: () => get().accessToken,
  getRefreshToken: () => get().refreshToken,
  isAuthenticated: () => !!get().accessToken && !!get().user,
}))
