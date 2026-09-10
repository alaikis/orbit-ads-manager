'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { productService } from '@/lib/api'
import { useState } from 'react'
import Link from 'next/link'

type Product = { id: number; title: string; external_id: string; price_cents: number; variants: { inventory_qty: number }[]; status: string; description?: string; link?: string; image_url?: string; brand?: string; currency?: string; product_type?: string }

const formatCurrency = (cents: number) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(cents / 100)

export default function ProductsPage() {
  const queryClient = useQueryClient()
  const [search, setSearch] = useState('')
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [editTitle, setEditTitle] = useState('')
  const [editPrice, setEditPrice] = useState('')
  const [editCurrency, setEditCurrency] = useState('USD')
  const [editStatus, setEditStatus] = useState('active')
  const [editDesc, setEditDesc] = useState('')
  const [editBrand, setEditBrand] = useState('')
  const [actionError, setActionError] = useState<string | null>(null)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['products', search],
    queryFn: () => productService.list({ search }),
  })

  const products: Product[] = (data as any)?.items || []

  const updateMutation = useMutation({
    mutationFn: (payload: { id: number; data: Partial<Product> }) =>
      productService.update(payload.id, payload.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      setEditingProduct(null)
      setActionError(null)
    },
    onError: (err: unknown) => setActionError(err instanceof Error ? err.message : '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => productService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
    },
    onError: (err: unknown) => alert('删除失败：' + (err instanceof Error ? err.message : '未知错误')),
  })

  const startEdit = (product: Product) => {
    setEditingProduct(product)
    setEditTitle(product.title)
    setEditPrice(String(product.price_cents))
    setEditCurrency(product.currency || 'USD')
    setEditStatus(product.status)
    setEditDesc(product.description || '')
    setEditBrand(product.brand || '')
    setActionError(null)
  }

  const saveEdit = () => {
    if (!editingProduct) return
    const payload: Partial<Product> = {}
    if (editTitle.trim()) payload.title = editTitle.trim()
    if (editPrice) payload.price_cents = parseInt(editPrice, 10) || editingProduct.price_cents
    if (editCurrency.trim()) payload.currency = editCurrency.trim().toUpperCase()
    if (editStatus) payload.status = editStatus
    if (editDesc.trim() !== undefined) payload.description = editDesc
    if (editBrand.trim() !== undefined) payload.brand = editBrand
    updateMutation.mutate({ id: editingProduct.id, data: payload })
  }

  const confirmDelete = (id: number) => {
    if (window.confirm('确定要删除这个商品吗？此操作不可撤销。')) {
      deleteMutation.mutate(id)
    }
  }

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
      <div className="card p-4 flex items-center gap-3">
        <input className="input max-w-md" placeholder="搜索商品名称或 SKU..." value={search} onChange={(e) => setSearch(e.target.value)} />
        <button onClick={() => refetch()} className="btn btn-secondary">刷新</button>
      </div>
      {isLoading ? (
        <div className="card overflow-hidden">
          <table className="w-full text-sm">
            <thead><tr className="border-b border-border-default">{['商品名称','SKU','价格','库存','状态'].map(h => (<th key={h} className="text-left p-4 font-title-sm text-text-secondary">{h}</th>))}</tr></thead>
            <tbody>{[1,2,3].map(i => (<tr key={i} className="border-b border-border-default last:border-0">{['商品名称','SKU','价格','库存','状态'].map((_,j) => (<td key={j} className="p-4"><div className="h-4 bg-surface-subtle rounded animate-pulse" style={{width: `${60 + Math.random()*40}%`}} /></td>))}</tr>))}</tbody>
          </table>
        </div>
      ) : products.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-text-muted mb-2">暂无商品</p>
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
                <th className="text-right p-4 font-title-sm text-text-secondary">状态</th>
                <th className="text-right p-4 font-title-sm text-text-secondary">操作</th>
              </tr>
            </thead>
            <tbody>
              {products.map((product) => (
                <tr key={product.id} className="border-b border-border-default last:border-0 hover:bg-surface-hover">
                  <td className="p-4 font-medium">
                    <Link href={`/products/${product.id}`} className="text-primary-600 hover:underline">{product.title}</Link>
                  </td>
                  <td className="p-4 text-text-secondary font-mono">{product.external_id}</td>
                  <td className="p-4 text-right font-mono">{formatCurrency(product.price_cents)}</td>
                  <td className="p-4 text-right">{product.variants?.reduce((sum, v) => sum + v.inventory_qty, 0) ?? '-'}</td>
                  <td className="p-4 text-right"><span className={`badge ${product.status === 'active' ? 'bg-success-bg text-success' : 'bg-warning-bg text-warning'}`}>{product.status}</span></td>
                  <td className="p-4 text-right space-x-2">
                    <button onClick={() => startEdit(product)} className="btn btn-ghost text-sm">编辑</button>
                    <button onClick={() => confirmDelete(product.id)} disabled={deleteMutation.isPending} className="btn btn-ghost text-sm text-danger-500">删除</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {editingProduct && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-text-primary">编辑商品</h2>
              <button onClick={() => setEditingProduct(null)} className="text-text-muted hover:text-text-primary"><span className="text-xl leading-none">×</span></button>
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
                <button onClick={() => setEditingProduct(null)} className="btn btn-secondary">取消</button>
                <button onClick={saveEdit} disabled={updateMutation.isPending} className="btn btn-primary">{updateMutation.isPending ? '保存中...' : '保存'}</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
