'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { storeService } from '@/lib/api'
import Link from 'next/link'

type Store = { id: number; name: string; platform: string; status: string; last_synced_at?: string; base_url?: string; api_key?: string; api_secret?: string }

const PLATFORM_LABEL: Record<string, string> = {
  woocommerce: 'WooCommerce',
  shopify: 'Shopify',
}

export default function StoreDetailPage({ params }: { params: { id: string } }) {
  const queryClient = useQueryClient()
  const storeId = Number(params.id)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['store', storeId],
    queryFn: () => storeService.get(storeId),
  })

  const syncMutation = useMutation({
    mutationFn: () => storeService.sync(storeId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['store', storeId] })
      queryClient.invalidateQueries({ queryKey: ['stores'] })
    },
  })

  const testMutation = useMutation({
    mutationFn: () => storeService.test(storeId),
  })

  const store = data as Store | undefined

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Link href="/stores" className="btn btn-ghost">返回店铺列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">店铺详情</h1>
            <p className="text-sm text-text-muted mt-1">加载失败</p>
          </div>
        </div>
        <div className="card p-12 text-center">
          <p className="text-danger-500 mb-4">{(error as Error).message}</p>
          <button onClick={() => refetch()} className="btn btn-primary">重试</button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link href="/stores" className="btn btn-ghost">返回店铺列表</Link>
          <div>
            <h1 className="text-2xl font-semibold text-text-primary">{store?.name || '店铺详情'}</h1>
            <p className="text-sm text-text-muted mt-1">平台: {PLATFORM_LABEL[store?.platform || ''] || store?.platform || '-'} | 状态: {store?.status || '-'}</p>
          </div>
        </div>
        <button onClick={() => refetch()} className="btn btn-secondary">刷新</button>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {[1, 2].map((i) => (
            <div key={i} className="card p-6 animate-pulse">
              <div className="h-6 bg-surface-subtle rounded w-1/3 mb-4" />
              <div className="space-y-2">
                <div className="h-4 bg-surface-subtle rounded" />
                <div className="h-4 bg-surface-subtle rounded" />
                <div className="h-4 bg-surface-subtle rounded" />
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="card p-6">
            <h2 className="font-title-md text-text-primary mb-4">概览</h2>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-text-muted">店铺 ID</span>
                <span className="font-mono">{store?.id || '-'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-text-muted">平台</span>
                <span className="capitalize">{PLATFORM_LABEL[store?.platform || ''] || store?.platform || '-'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-text-muted">连接状态</span>
                <span className={`badge ${store?.status === 'bound' || store?.status === 'active' ? 'bg-success-bg text-success' : 'bg-warning-bg text-warning'}`}>
                  {store?.status || 'pending'}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-text-muted">最近同步</span>
                <span>{store?.last_synced_at || '-'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-text-muted">店铺 URL</span>
                <span className="truncate max-w-[200px]">{store?.base_url || '-'}</span>
              </div>
            </div>
          </div>
          <div className="card p-6">
            <h2 className="font-title-md text-text-primary mb-4">同步控制</h2>
            <div className="space-y-3">
              <button onClick={() => syncMutation.mutate()} disabled={syncMutation.isPending} className="btn btn-secondary w-full">
                {syncMutation.isPending ? (<><span className="mr-2 inline-block w-4 h-4 border-2 border-text-primary border-t-transparent rounded-full animate-spin" />同步中...</>) : '手动触发同步'}
              </button>
              <button onClick={() => testMutation.mutate()} disabled={testMutation.isPending} className="btn btn-secondary w-full">
                {testMutation.isPending ? (<><span className="mr-2 inline-block w-4 h-4 border-2 border-text-primary border-t-transparent rounded-full animate-spin" />测试中...</>) : '测试连接'}
              </button>
              {(syncMutation.data || testMutation.data) && (
                <pre className="text-xs text-text-muted p-2 bg-surface-subtle rounded overflow-auto">
                  {JSON.stringify((syncMutation.data || testMutation.data) as object, null, 2)}
                </pre>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}