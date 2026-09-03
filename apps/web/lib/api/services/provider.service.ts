import { apiClient } from '../client'
import type { Provider, PaginatedResponse } from '../types'

export class ProviderService {
  async list(params?: { page?: number; per_page?: number; type?: string }): Promise<PaginatedResponse<Provider>> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<PaginatedResponse<Provider>>(`/providers${query}`)
  }

  async get(id: number): Promise<Provider> {
    return apiClient.get<Provider>(`/providers/${id}`)
  }

  async create(data: { name: string; type: string; config?: Record<string, unknown> }): Promise<Provider> {
    return apiClient.post<Provider>('/providers', data)
  }

  async update(id: number, data: { name?: string; config?: Record<string, unknown> }): Promise<Provider> {
    return apiClient.patch<Provider>(`/providers/${id}`, data)
  }

  async delete(id: number): Promise<void> {
    return apiClient.delete<void>(`/providers/${id}`)
  }

  async test(id: number): Promise<{ success: boolean; message: string }> {
    return apiClient.post<{ success: boolean; message: string }>(`/providers/${id}/test`)
  }

  async getUsage(id: number): Promise<{ used: number; limit: number }> {
    return apiClient.get<{ used: number; limit: number }>(`/providers/${id}/usage`)
  }
}

export const providerService = new ProviderService()
