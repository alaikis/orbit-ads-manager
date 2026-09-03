import { apiClient } from '../client'
import type { Workspace, PaginatedResponse } from '../types'

export class WorkspaceService {
  async list(params?: { page?: number; per_page?: number }): Promise<PaginatedResponse<Workspace>> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<PaginatedResponse<Workspace>>(`/workspaces${query}`)
  }

  async getCurrent(): Promise<Workspace | null> {
    return apiClient.get<Workspace | null>('/workspaces/current')
  }

  async get(id: number): Promise<Workspace> {
    return apiClient.get<Workspace>(`/workspaces/${id}`)
  }

  async create(data: { name: string; plan?: string }): Promise<Workspace> {
    return apiClient.post<Workspace>('/workspaces', data)
  }

  async update(id: number, data: { name?: string; plan?: string }): Promise<Workspace> {
    return apiClient.patch<Workspace>(`/workspaces/${id}`, data)
  }

  async delete(id: number): Promise<void> {
    return apiClient.delete<void>(`/workspaces/${id}`)
  }

  async switch(id: number): Promise<{ message: string }> {
    return apiClient.post<{ message: string }>(`/workspaces/${id}/switch`)
  }
}

export const workspaceService = new WorkspaceService()
