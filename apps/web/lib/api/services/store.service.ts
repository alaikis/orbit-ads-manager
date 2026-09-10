import { apiClient } from '../client'
import type { Store, CreateStoreRequest, UpdateStoreRequest, PaginatedResponse } from '../types'

export class StoreService {
  async list(params?: { page?: number; per_page?: number; status?: string }): Promise<PaginatedResponse<Store>> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<PaginatedResponse<Store>>(`/stores${query}`)
  }

  async get(id: number): Promise<Store> {
    return apiClient.get<Store>(`/stores/${id}`)
  }

  async create(data: CreateStoreRequest): Promise<Store> {
    return apiClient.post<Store>('/stores', data)
  }

  async update(id: number, data: UpdateStoreRequest): Promise<Store> {
    return apiClient.patch<Store>(`/stores/${id}`, data)
  }

  async delete(id: number): Promise<void> {
    return apiClient.delete<void>(`/stores/${id}`)
  }

  async sync(id: number): Promise<{ message: string }> {
    return apiClient.post<{ message: string }>(`/stores/${id}/sync`)
  }

  async test(id: number): Promise<{ status: string; message: string }> {
    return apiClient.post<{ status: string; message: string }>(`/stores/${id}/test`)
  }
}

export const storeService = new StoreService()