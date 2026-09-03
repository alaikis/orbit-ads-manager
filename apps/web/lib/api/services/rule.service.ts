import { apiClient } from '../client'
import type { Rule, CreateRuleRequest, PaginatedResponse } from '../types'

export class RuleService {
  async list(params?: { page?: number; per_page?: number; status?: string }): Promise<PaginatedResponse<Rule>> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<PaginatedResponse<Rule>>(`/rules${query}`)
  }

  async get(id: number): Promise<Rule> {
    return apiClient.get<Rule>(`/rules/${id}`)
  }

  async create(data: CreateRuleRequest): Promise<Rule> {
    return apiClient.post<Rule>('/rules', data)
  }

  async update(id: number, data: Partial<CreateRuleRequest>): Promise<Rule> {
    return apiClient.patch<Rule>(`/rules/${id}`, data)
  }

  async delete(id: number): Promise<void> {
    return apiClient.delete<void>(`/rules/${id}`)
  }

  async toggle(id: number): Promise<Rule> {
    return apiClient.post<Rule>(`/rules/${id}/toggle`)
  }
}

export const ruleService = new RuleService()
