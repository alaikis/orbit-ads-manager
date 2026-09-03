import { apiClient } from '../client'
import type { Campaign, AdGroup, MetricsSummary, PaginatedResponse } from '../types'

export class CampaignService {
  async list(params?: { 
    page?: number
    per_page?: number
    platform?: string
    status?: string
    search?: string 
  }): Promise<PaginatedResponse<Campaign>> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<PaginatedResponse<Campaign>>(`/advertising/campaigns${query}`)
  }

  async get(id: number): Promise<Campaign> {
    return apiClient.get<Campaign>(`/advertising/campaigns/${id}`)
  }

  async getMetrics(id: number): Promise<MetricsSummary> {
    return apiClient.get<MetricsSummary>(`/advertising/campaigns/${id}/metrics`)
  }

  async getAdGroups(id: number): Promise<AdGroup[]> {
    const res = await apiClient.get<{ items: AdGroup[] } | AdGroup[]>(`/advertising/campaigns/${id}/ad-groups`)
    if (Array.isArray(res)) return res
    return (res as any)?.items || []
  }

  async pause(id: number): Promise<Campaign> {
    return apiClient.post<Campaign>(`/advertising/campaigns/${id}/pause`)
  }

  async resume(id: number): Promise<Campaign> {
    return apiClient.post<Campaign>(`/advertising/campaigns/${id}/resume`)
  }

  async updateBudget(id: number, newBudgetCents: number): Promise<Campaign> {
    return apiClient.post<Campaign>(`/advertising/campaigns/${id}/budget`, { new_budget_cents: newBudgetCents })
  }
}

export const campaignService = new CampaignService()
