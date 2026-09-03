'use client'

export default function ReportDetailPage({ params }: { params: { id: string } }) {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-text-primary">看板详情</h1>
        <p className="text-sm text-text-muted mt-1">ID: {params.id}</p>
      </div>
      <div className="card p-6">
        <p className="text-text-secondary">自定义看板功能即将推出</p>
      </div>
    </div>
  )
}
