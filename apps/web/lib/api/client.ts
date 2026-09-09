const API_BASE = 'https://adsapi.alaikis.com/api/v1'

class ApiError extends Error {
  code: number
  data?: unknown

  constructor(message: string, code: number, data?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.data = data
  }
}

function getAuthHeader(): Record<string, string> {
  if (typeof window === 'undefined') return {}
  try {
    const stored = localStorage.getItem('orbit-auth')
    if (stored) {
      const parsed = JSON.parse(stored)
      const token = parsed?.accessToken
      return token ? { Authorization: `Bearer ${token}` } : {}
    }
  } catch {
    // ignore parse errors
  }
  return {}
}

function clearAuthAndRedirect() {
  if (typeof window === 'undefined') return
  try {
    localStorage.removeItem('orbit-auth')
    window.location.href = '/login?redirect=' + encodeURIComponent(window.location.pathname)
  } catch {
    // ignore
  }
}

async function refreshTokenIfNeeded(): Promise<boolean> {
  if (typeof window === 'undefined') return false
  try {
    const stored = localStorage.getItem('orbit-auth')
    if (!stored) return false
    const parsed = JSON.parse(stored)
    const refreshToken = parsed?.refreshToken
    if (!refreshToken) return false

    const response = await fetch(`${API_BASE}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    })

    if (!response.ok) return false

    const json = await response.json()
    if (json?.code !== 0 || !json?.data) return false

    const data = json.data
    localStorage.setItem(STORAGE_KEY, JSON.stringify({
      user: parsed.user,
      accessToken: data.access_token,
      refreshToken: data.refresh_token || refreshToken,
    }))
    return true
  } catch {
    return false
  }
}

const STORAGE_KEY = 'orbit-auth'

async function request<T>(path: string, options: RequestInit = {}, retry = true): Promise<T> {
  const url = `${API_BASE}${path}`
  const authHeaders = getAuthHeader()
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> | undefined),
    ...authHeaders,
  }

  const response = await fetch(url, {
    ...options,
    headers,
  })

  if (response.status === 401 && retry && path !== '/auth/refresh' && path !== '/auth/login') {
    const refreshed = await refreshTokenIfNeeded()
    if (refreshed) {
      return request<T>(path, options, false)
    }
    clearAuthAndRedirect()
    throw new ApiError('Unauthorized', 401)
  }

  if (response.status === 401) {
    clearAuthAndRedirect()
    throw new ApiError('Unauthorized', 401)
  }

  if (response.status === 403) {
    throw new ApiError('Forbidden', 403)
  }

  if (response.status === 204) {
    return null as T
  }

  const json = await response.json()

  if (!response.ok) {
    const message = json?.message || `HTTP ${response.status}`
    throw new ApiError(message, response.status, json)
  }

  // Return the data field from the API response
  return json?.data as T
}

export const apiClient = {
  get: <T>(path: string, query?: Record<string, unknown>) => {
    const qs = query ? '?' + new URLSearchParams(Object.entries(query).reduce<Record<string, string>>((acc, [key, value]) => {
      if (value !== undefined && value !== null) acc[key] = String(value)
      return acc
    }, {})).toString() : ''
    return request<T>(`${path}${qs}`)
  },
  post: <T>(path: string, data?: unknown) => request<T>(path, {
    method: 'POST',
    body: data ? JSON.stringify(data) : undefined,
  }),
  patch: <T>(path: string, data?: unknown) => request<T>(path, {
    method: 'PATCH',
    body: data ? JSON.stringify(data) : undefined,
  }),
  put: <T>(path: string, data?: unknown) => request<T>(path, {
    method: 'PUT',
    body: data ? JSON.stringify(data) : undefined,
  }),
  delete: <T>(path: string) => request<T>(path, { method: 'DELETE' }),
}

export { ApiError }
