export interface User {
  id: number
  email: string
  role: string
  tenant_id: number
}

export interface Store {
  id: number
  name: string
  platform: 'woocommerce' | 'shopify'
  status: string
  last_synced_at?: string
}

export interface Campaign {
  id: number
  name: string
  platform: string
  status: 'enabled' | 'paused' | 'removed'
  daily_budget_cents?: number
  metrics?: {
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
}

export interface MetricCard {
  title: string
  value: string | number
  change?: number
  changeLabel?: string
  tooltip?: string
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
