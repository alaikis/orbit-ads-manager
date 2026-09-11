import { apiClient } from '../client'

export class CampaignService {
  async list(params?: { account_id?: number; status?: string }): Promise<{ items: any[] }> {
    const query = params ? '?' + new URLSearchParams(Object.entries(params).filter(([,v]) => v !== undefined).reduce((a,[k,v]) => ({ ...a, [k]: String(v) }), {} as Record<string, string>)).toString() : ''
    return apiClient.get<{ items: any[] }>(`/advertising/campaigns${query}`)
  }

  async get(id: number): Promise<any> {
    return apiClient.get<any>(`/advertising/campaigns/${id}`)
  }

  async pause(id: number): Promise<{ message: string }> {
    return apiClient.post<{ message: string }>(`/advertising/campaigns/${id}/pause`, {})
  }

  async resume(id: number): Promise<{ message: string }> {
    return apiClient.post<{ message: string }>(`/advertising/campaigns/${id}/resume`, {})
  }

  async updateBudget(id: number, budgetCents: number): Promise<{ message: string; new_budget_cents: number }> {
    return apiClient.post<{ message: string; new_budget_cents: number }>(`/advertising/campaigns/${id}/budget`, { new_budget_cents: budgetCents })
  }

  async listAdGroups(id: number): Promise<{ items: any[] }> {
    return apiClient.get<{ items: any[] }>(`/advertising/campaigns/${id}/ad-groups`)
  }
}

export const campaignService = new CampaignService()
