import { apiClient } from '../client'
import type { AdAccount, PaginatedResponse } from '../types'

export class AdAccountService {
  async list(params?: { page?: number; per_page?: number; platform?: string }): Promise<PaginatedResponse<AdAccount>> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<PaginatedResponse<AdAccount>>(`/ad-accounts${query}`)
  }

  async get(id: number): Promise<AdAccount> {
    return apiClient.get<AdAccount>(`/ad-accounts/${id}`)
  }

  async create(data: { name: string; platform: string; external_id: string; currency?: string; customer_id?: string }): Promise<AdAccount> {
    return apiClient.post<AdAccount>('/ad-accounts', data)
  }

  async update(id: number, data: { name?: string; currency?: string; customer_id?: string; status?: string }): Promise<AdAccount> {
    return apiClient.patch<AdAccount>(`/ad-accounts/${id}`, data)
  }

  async authorize(platform: string): Promise<{ authorization_url: string; state: string }> {
    return apiClient.post<{ authorization_url: string; state: string }>(`/oauth/connect/${platform}`, {})
  }

  async delete(id: number): Promise<void> {
    return apiClient.delete<void>(`/ad-accounts/${id}`)
  }

  async refresh(id: number): Promise<AdAccount> {
    return apiClient.post<AdAccount>(`/ad-accounts/${id}/refresh`)
  }
}

export const adAccountService = new AdAccountService()