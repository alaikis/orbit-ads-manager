'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import Link from "next/link"
import { productService } from "@/lib/api"
import { useState } from "react"

type ProductDetail = {
  id: number
  title: string
  external_id: string
  price_cents: number
  variants?: { inventory_qty: number }[]
  status: string
  description?: string
  brand?: string
  currency?: string
  product_type?: string
  store_id?: number
}

export default function ProductDetailPage({ params }: { params: { id: string } }) {
  const productId = Number(params.id)
  const queryClient = useQueryClient()
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["product", productId],
    queryFn: () => productService.get(productId),
  })

  const product = data as ProductDetail | undefined

  const [editing, setEditing] = useState(false)
  const [editTitle, setEditTitle] = useState('')
  const [editPrice, setEditPrice] = useState('')
  const [editCurrency, setEditCurrency] = useState('USD')
  const [editStatus, setEditStatus] = useState('active')
  const [editBrand, setEditBrand] = useState('')
  const [editDesc, setEditDesc] = useState('')
  const [actionError, setActionError] = useState('')

  const updateMutation = useMutation({
    mutationFn: (payload: Partial<ProductDetail>) => productService.update(productId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['product', productId] })
      queryClient.invalidateQueries({ queryKey: ['products'] })
      setEditing(false)
      setActionError('')
    },
    onError: (err: unknown) => setActionError(err instanceof Error ? err.message : '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: () => productService.delete(productId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      window.location.href = '/products'
    },
    onError: (err: unknown) => alert('删除失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const startEdit = () => {
    if (!product) return
    setEditTitle(product.title)
    setEditPrice(String(product.price_cents))
    setEditCurrency(product.currency || 'USD')
    setEditStatus(product.status)
    setEditBrand(product.brand || '')
    setEditDesc(product.description || '')
    setEditing(true)
    setActionError('')
  }

  const saveEdit = () => {
    const payload: Partial<ProductDetail> = {}
    if (editTitle.trim()) payload.title = editTitle.trim()
    if (editPrice) payload.price_cents = parseInt(editPrice, 10) || product!.price_cents
    if (editCurrency.trim()) payload.currency = editCurrency.trim().toUpperCase()
    if (editStatus) payload.status = editStatus
    if (editBrand.trim() !== undefined) payload.brand = editBrand
    if (editDesc.trim() !== undefined) payload.description = editDesc
    updateMutation.mutate(payload)
  }

  const confirmDelete = () => {
    if (window.confirm('确定要删除这个商品吗？')) {
      deleteMutation.mutate()
    }
  }

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
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link href="/products" className="btn btn-ghost">返回商品列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">{product.title}</h1>
            <p className="text-sm text-text-muted mt-1">
              SKU：<span className="font-mono">{product.external_id || '-'}</span>
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          <button onClick={startEdit} className="btn btn-secondary">编辑</button>
          <button onClick={confirmDelete} className="btn btn-secondary text-danger-500">删除</button>
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
              <td className="p-4">{product.title}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">SKU</td>
              <td className="p-4 font-mono">{product.external_id || "-"}</td>
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
              <td className="p-4 text-text-muted">品牌</td>
              <td className="p-4">{product.brand || '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">币种</td>
              <td className="p-4">{product.currency || '-'}</td>
            </tr>
            <tr className="border-b border-border-default">
              <td className="p-4 text-text-muted">商品类型</td>
              <td className="p-4">{product.product_type || '-'}</td>
            </tr>
            <tr>
              <td className="p-4 text-text-muted">所属店铺 ID</td>
              <td className="p-4 font-mono">{product.store_id ?? "-"}</td>
            </tr>
          </tbody>
        </table>
      </div>

      {editing && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">编辑商品</h2>
              <button onClick={() => setEditing(false)} className="text-text-muted hover:text-text-primary"><span className="text-xl leading-none">×</span></button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">商品名称</label>
                <input className="input" value={editTitle} onChange={(e) => setEditTitle(e.target.value)} />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-text-secondary mb-1.5">价格（分）</label>
                  <input className="input" type="number" value={editPrice} onChange={(e) => setEditPrice(e.target.value)} />
                </div>
                <div>
                  <label className="block text-sm font-medium text-text-secondary mb-1.5">币种</label>
                  <input className="input" value={editCurrency} onChange={(e) => setEditCurrency(e.target.value)} />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">状态</label>
                <select className="input" value={editStatus} onChange={(e) => setEditStatus(e.target.value)}>
                  <option value="active">active</option>
                  <option value="inactive">inactive</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">品牌</label>
                <input className="input" value={editBrand} onChange={(e) => setEditBrand(e.target.value)} />
              </div>
              <div>
                <label className="block text-sm font-medium text-text-secondary mb-1.5">描述</label>
                <textarea className="input h-24" value={editDesc} onChange={(e) => setEditDesc(e.target.value)} />
              </div>
              {actionError && (<div className="text-sm text-danger-500 bg-danger-bg p-3 rounded-md">{actionError}</div>)}
              <div className="flex gap-2 justify-end pt-2">
                <button onClick={() => setEditing(false)} className="btn btn-secondary">取消</button>
                <button onClick={saveEdit} disabled={updateMutation.isPending} className="btn btn-primary">{updateMutation.isPending ? '保存中...' : '保存'}</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}