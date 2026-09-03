// API Request/Response Types

export interface User {
  id: number
  email: string
  role: string
  tenant_id: number
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
  token_type: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
}

export interface RefreshRequest {
  refresh_token: string
}

export interface Store {
  id: number
  name: string
  platform: 'woocommerce' | 'shopify'
  status: string
  last_synced_at?: string
  created_at?: string
  updated_at?: string
}

export interface CreateStoreRequest {
  name: string
  platform: 'woocommerce' | 'shopify'
  base_url?: string
  store_url?: string
  api_key?: string
  api_secret?: string
}

export interface UpdateStoreRequest {
  name?: string
  status?: string
}

export interface Campaign {
  id: number
  name: string
  platform: string
  status: 'enabled' | 'paused' | 'removed'
  daily_budget_cents?: number
  metrics?: CampaignMetrics
  created_at?: string
  updated_at?: string
}

export interface CampaignMetrics {
  spend: number
  impressions: number
  clicks: number
  ctr: number
  cpc: number
  conversions: number
  cpa: number
  roas: number
  budget_rate: number
}

export interface AdGroup {
  id: number
  name: string
  status: string
  bid_micros?: number
  metrics?: CampaignMetrics
}

export interface MetricsSummary {
  spend: number
  ctr: number
  cpa: number
  roas: number
  impressions?: number
  clicks?: number
  conversions?: number
}

export interface Notification {
  id: number
  type: string
  title: string
  body: string
  link?: string
  read_at?: string
  created_at: string
}

export interface AdAccount {
  id: number
  name: string
  platform: string
  status: string
  currency?: string
  created_at?: string
}

export interface Product {
  id: number
  name: string
  sku?: string
  price?: number
  status: string
  store_id?: number
}

export interface Feed {
  id: number
  name: string
  status: string
  last_synced_at?: string
  product_count?: number
}

export interface Report {
  id: number
  name: string
  type: string
  status: string
  created_at: string
  download_url?: string
}

export interface Rule {
  id: number
  name: string
  condition: string
  action: string
  status: 'active' | 'inactive'
  created_at?: string
}

export interface CreateRuleRequest {
  name: string
  condition: string
  action: string
}

export interface Workspace {
  id: number
  name: string
  plan: string
  status: string
  member_count?: number
  created_at?: string
}

export interface Provider {
  id: number
  name: string
  type: string
  status: string
  config?: Record<string, unknown>
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  per_page: number
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}
