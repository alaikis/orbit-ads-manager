import { apiClient } from '../client'
import type { Report, PaginatedResponse } from '../types'

export class ReportService {
  async list(params?: { page?: number; per_page?: number; type?: string }): Promise<PaginatedResponse<Report>> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<PaginatedResponse<Report>>(`/reports${query}`)
  }

  async get(id: number): Promise<Report> {
    return apiClient.get<Report>(`/reports/${id}`)
  }

  async create(data: { name: string; type: string; filters?: Record<string, unknown> }): Promise<Report> {
    return apiClient.post<Report>('/reports', data)
  }

  async delete(id: number): Promise<void> {
    return apiClient.delete<void>(`/reports/${id}`)
  }

  async download(id: number): Promise<Blob> {
    const response = await fetch(`/api/v1/reports/${id}/download`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('orbit-auth') ? JSON.parse(localStorage.getItem('orbit-auth')!).state.accessToken : ''}`,
      },
    })
    return response.blob()
  }
}

export const reportService = new ReportService()
