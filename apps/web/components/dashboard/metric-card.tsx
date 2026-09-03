'use client'

import { TrendingUp, TrendingDown } from 'lucide-react'
import { MetricCard as MetricCardType } from '@/lib/types'

interface MetricCardProps {
  title: string
  value: string | number
  format?: 'currency' | 'percent' | 'ratio' | 'number'
  change?: number
  tooltip?: string
}

function formatValue(value: string | number, format?: string): string {
  if (typeof value === 'string') return value
  switch (format) {
    case 'currency':
      return `$${value.toFixed(2)}`
    case 'percent':
      return `${value.toFixed(2)}%`
    case 'ratio':
      return value.toFixed(2)
    default:
      return value.toLocaleString()
  }
}

export function MetricCard({ title, value, format, change, tooltip }: MetricCardProps) {
  return (
    <div className="card p-6">
      <div className="flex items-center justify-between">
        <p className="text-sm font-medium text-text-secondary">{title}</p>
        {tooltip && <span className="text-text-muted cursor-help" title={tooltip}>ⓘ</span>}
      </div>
      <div className="mt-2 flex items-baseline gap-2">
        <p className="text-2xl font-semibold text-text-primary font-metric">{formatValue(value, format)}</p>
        {change !== undefined && (
          <span className={`flex items-center text-sm ${change >= 0 ? 'text-success-500' : 'text-danger-500'}`}>
            {change >= 0 ? <TrendingUp size={14} /> : <TrendingDown size={14} />}
            {Math.abs(change)}%
          </span>
        )}
      </div>
    </div>
  )
}
