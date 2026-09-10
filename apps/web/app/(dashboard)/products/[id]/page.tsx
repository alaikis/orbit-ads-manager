'use client'

import { useQuery } from '@tanstack/react-query'
import Link from "next/link"
import { productService } from "@/lib/api"

type ProductDetail = {
  id: number
  name: string
  sku?: string
  price_cents?: number
  variants?: { inventory_qty: number }[]
  status: string
  store_id?: number
}

export default function ProductDetailPage({ params }: { params: { id: string } }) {
  const productId = Number(params.id)
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["product", productId],
    queryFn: () => productService.get(productId),
  })

  const product = data as ProductDetail | undefined

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/products" className="btn btn-ghost">返回商品列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">商品详情</h1>
            <p className="text-sm text-text-muted mt-1">加载中...</p>
          </div>
        </div>
        <div className="card overflow-hidden animate-pulse">
          {[1, 2, 3, 4, 5, 6, 7].map((i) => (
            <div key={i} className="flex items-center border-b border-border-default last:border-0">
              <div className="p-4 w-40"><div className="h-4 bg-surface-subtle rounded w-3/4" /></div>
              <div className="p-4 flex-1"><div className="h-4 bg-surface-subtle rounded w-1/2" /></div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  if (error || !product) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/products" className="btn btn-ghost">返回商品列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">商品详情</h1>
            <p className="text-sm text-text-muted mt-1">加载失败</p>
          </div>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">{(error as Error)?.message || "未找到该商品"}</p>
          <button onClick={() => refetch()} className="btn btn-primary">重试</button>
        </div>
      </div>
    )
  }

  const totalInventory = product.variants?.reduce((sum, v) => sum + v.inventory_qty, 0) ?? 0

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link href="/products" className="btn btn-ghost">返回商品列表</Link>
        <div>
          <h1 className="text-2xl font-semibold text-text-primary">{product.name}</h1>
          <p className="text-sm text-text-muted mt-1">
            状态：<span className={"badge " + (product.status === "active" ? "bg-success-bg text-success" : "bg-warning-bg text-warning")}>{product.status}</span>
          </p>
        </div>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <tbody>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted w-40">商品 ID</td>
              <td className="p-4 font-mono">{product.id}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">商品名称</td>
              <td className="p-4">{product.name}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">SKU</td>
              <td className="p-4 font-mono">{product.sku || "-"}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">价格</td>
              <td className="p-4 font-mono">{"$" + ((product.price_cents ?? 0) / 100).toFixed(2)}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">库存总量</td>
              <td className="p-4">{totalInventory}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">状态</td>
              <td className="p-4">
                <span className={"badge " + (product.status === "active" ? "bg-success-bg text-success" : "bg-warning-bg text-warning")}>
                  {product.status}
                </span>
              </td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">所属店铺 ID</td>
              <td className="p-4 font-mono">{product.store_id ?? "-"}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  )
}