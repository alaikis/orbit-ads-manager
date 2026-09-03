import { apiClient } from '../client'
import type { Feed, PaginatedResponse } from '../types'

export class FeedService {
  async list(params?: { page?: number; per_page?: number; status?: string }): Promise<PaginatedResponse<Feed>> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<PaginatedResponse<Feed>>(`/feeds${query}`)
  }

  async get(id: number): Promise<Feed> {
    return apiClient.get<Feed>(`/feeds/${id}`)
  }

  async create(data: { name: string; store_id: number }): Promise<Feed> {
    return apiClient.post<Feed>('/feeds', data)
  }

  async update(id: number, data: { name?: string; status?: string }): Promise<Feed> {
    return apiClient.patch<Feed>(`/feeds/${id}`, data)
  }

  async delete(id: number): Promise<void> {
    return apiClient.delete<void>(`/feeds/${id}`)
  }

  async sync(id: number): Promise<{ message: string }> {
    return apiClient.post<{ message: string }>(`/feeds/${id}/sync`)
  }
}

export const feedService = new FeedService()
