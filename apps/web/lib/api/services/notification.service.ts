import { apiClient } from '../client'
import type { Notification, PaginatedResponse } from '../types'

export class NotificationService {
  async list(params?: { page?: number; per_page?: number; unread_only?: boolean }): Promise<PaginatedResponse<Notification>> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<PaginatedResponse<Notification>>(`/notifications${query}`)
  }

  async markAsRead(id: number): Promise<void> {
    return apiClient.patch<void>(`/notifications/${id}/read`)
  }

  async markAllAsRead(): Promise<void> {
    return apiClient.post<void>('/notifications/read-all')
  }

  async getUnreadCount(): Promise<{ count: number }> {
    return apiClient.get<{ count: number }>('/notifications/unread-count')
  }
}

export const notificationService = new NotificationService()
