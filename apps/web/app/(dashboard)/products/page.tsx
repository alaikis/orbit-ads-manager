'use client'

import { useQuery } from '@tanstack/react-query'
import { productService } from '@/lib/api'
import { useState } from 'react'

type Product = { id: number; title: string; external_id: string; price_cents: number; variants: { inventory_qty: number }[] }

export default function ProductsPage() {
  const [search, setSearch] = useState('')
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['products', search],
    queryFn: () => productService.list({ search }),
  })

  const products = (data as any)?.items || []

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">商品管理</h1>
          <p className="text-sm text-text-muted mt-1">查看与管理同步的商品</p>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">加载失败：{(error as Error).message}</p>
          <button onClick={() => refetch()} className="btn btn-primary">重试</button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-text-primary">商品管理</h1>
        <p className="text-sm text-text-muted mt-1">查看与管理同步的商品</p>
      </div>
      <div className="card p-4">
        <input
          className="input max-w-md"
          placeholder="搜索商品名称或 SKU..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>
      {isLoading ? (
        <div className="card p-4 animate-pulse space-y-3">
          {[1, 2, 3].map((i) => (<div key={i} className="h-12 bg-surface-subtle rounded" />))}
        </div>
      ) : products.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted mb-4">暂无商品</p>
          <p className="text-sm text-text-muted">前往店铺触发同步以导入商品</p>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border-default">
                <th className="text-left p-4 font-title-sm text-text-secondary">商品名称</th>
                <th className="text-left p-4 font-title-sm text-text-secondary">SKU</th>
                <th className="text-right p-4 font-title-sm text-text-secondary">价格</th>
                <th className="text-right p-4 font-title-sm text-text-secondary">库存</th>
              </tr>
            </thead>
            <tbody>
              {products.map((product: Product) => (
                <tr key={product.id} className="border-b border-border-default last:border-0 hover:bg-surface-hover">
                  <td className="p-4 font-medium">{product.title}</td>
                  <td className="p-4 text-text-secondary font-mono">{product.external_id}</td>
                  <td className="p-4 text-right font-mono">${(product.price_cents / 100).toFixed(2)}</td>
                  <td className="p-4 text-right">{product.variants?.reduce((sum, v) => sum + v.inventory_qty, 0) || '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
