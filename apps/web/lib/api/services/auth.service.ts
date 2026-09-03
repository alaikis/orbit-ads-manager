import { apiClient } from '../client'
import type { AuthResponse, LoginRequest, RegisterRequest, RefreshRequest, User } from '../types'

export class AuthService {
  async register(data: RegisterRequest): Promise<AuthResponse> {
    return apiClient.post<AuthResponse>('/auth/register', data)
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    return apiClient.post<AuthResponse>('/auth/login', data)
  }

  async refresh(data: RefreshRequest): Promise<AuthResponse> {
    return apiClient.post<AuthResponse>('/auth/refresh', data)
  }

  async logout(): Promise<void> {
    return apiClient.post<void>('/auth/logout')
  }

  async me(): Promise<User> {
    return apiClient.get<User>('/auth/me')
  }

  async forgotPassword(email: string): Promise<{ message: string }> {
    return apiClient.post<{ message: string }>('/auth/forgot-password', { email })
  }

  async resetPassword(token: string, password: string): Promise<{ message: string }> {
    return apiClient.post<{ message: string }>('/auth/reset-password', { token, password })
  }
}

export const authService = new AuthService()
