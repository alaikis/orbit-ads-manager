import { apiClient } from '../client'
import type { Product, PaginatedResponse } from '../types'

export class ProductService {
  async list(params?: { 
    page?: number
    per_page?: number
    store_id?: number
    status?: string
    search?: string 
  }): Promise<PaginatedResponse<Product>> {
    const query = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.get<PaginatedResponse<Product>>(`/products${query}`)
  }

  async get(id: number): Promise<Product> {
    return apiClient.get<Product>(`/products/${id}`)
  }

  async sync(storeId: number): Promise<{ message: string }> {
    return apiClient.post<{ message: string }>(`/products/sync`, { store_id: storeId })
  }
}

export const productService = new ProductService()
