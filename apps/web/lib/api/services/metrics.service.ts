import { apiClient } from '../client'
import type { MetricsSummary } from '../types'

export class MetricsService {
  async getSummary(params?: { scope_type?: string; scope_id?: number }): Promise<MetricsSummary> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<MetricsSummary>(`/metrics/summary${query}`)
  }

  async getTrend(params: { 
    scope_type?: string
    scope_id?: number
    start_date: string
    end_date: string 
  }): Promise<{ date: string; spend: number; clicks: number }[]> {
    const query = '?' + new URLSearchParams(
      Object.entries(params).reduce<Record<string, string>>((acc, [key, value]) => {
        if (value !== undefined) acc[key] = String(value)
        return acc
      }, {})
    ).toString()
    return apiClient.get<{ date: string; spend: number; clicks: number }[]>(`/metrics/trend${query}`)
  }
}

export const metricsService = new MetricsService()
