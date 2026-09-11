'use client'

import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/lib/api'

type Conversation = { id: number; title: string; status: string; created_at?: string }

export default function AgentPage() {
  const [messages, setMessages] = useState<{ role: string; content: string; tool?: string }[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [pendingActions, setPendingActions] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const [conversationId, setConversationId] = useState<number | null>(null)

  const { data: conversationsData } = useQuery({
    queryKey: ['agent-conversations'],
    queryFn: () => apiClient.get<{ items: Conversation[] }>('/agent/conversations'),
  })
  const conversations = (conversationsData as any)?.items || []

  const handleSend = async () => {
    if (!input.trim()) return
    const userMsg = { role: 'user', content: input }
    setMessages((prev) => [...prev, userMsg])
    const userInput = input
    setInput('')
    setLoading(true)
    setError(null)

    try {
      const payload: any = { message: userInput }
      if (conversationId) payload.conversation_id = conversationId
      const data = await apiClient.post<{ response: string; tool?: string; conversation_id?: number }>('/agent/chat', payload)
      if (data?.response) {
        const toolMsg = data.tool ? { role: 'tool' as const, content: data.response, tool: data.tool } : null
        const finalMsg = { role: 'assistant' as const, content: data.response }
        setMessages((prev) => (toolMsg ? [...prev, toolMsg, finalMsg] : [...prev, finalMsg]))
        if (data.conversation_id && !conversationId) {
          setConversationId(data.conversation_id)
        }
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : '连接中断，请重试'
      setMessages((prev) => [...prev, { role: 'system', content: message }])
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  const handleRetry = () => {
    setError(null)
    setMessages((prev) => prev.filter((m) => m.role !== 'system'))
  }

  const loadConversation = async (id: number) => {
    setConversationId(id)
    setMessages([])
    setError(null)
    try {
      const data = await apiClient.get<{ items: { role: string; content: string; tool?: string }[] }>(`/agent/conversations/${id}/messages`)
      const items = (data as any)?.items || []
      setMessages(items.map((m: any) => ({ role: m.role, content: m.content, tool: m.tool })))
    } catch {
      setError('加载历史消息失败')
    }
  }

  return (
    <div className="flex h-[calc(100vh-56px)]">
      <div className="w-64 border-r border-border-default p-4">
        <button className="btn btn-primary w-full mb-4">新对话</button>
        <div className="text-sm text-text-muted mb-2">历史会话</div>
        <div className="space-y-1">
          {conversations.map((c: Conversation) => (
            <div
              key={c.id}
              onClick={() => loadConversation(c.id)}
              className={`px-3 py-2 rounded-md text-sm cursor-pointer truncate ${
                conversationId === c.id ? 'bg-primary-50 text-primary-700' : 'text-text-secondary hover:bg-surface-hover'
              }`}
            >
              {c.title || `对话 ${c.id}`}
            </div>
          ))}
        </div>
      </div>
      <div className="flex-1 flex flex-col">
        <div className="flex-1 overflow-y-auto p-6 space-y-4">
          {messages.length === 0 && !error && (
            <div className="text-center text-text-muted mt-20">
              <div className="text-4xl mb-4">🤖</div>
              <p className="text-lg mb-2">试试问我：昨天花了多少、ROAS 如何？</p>
            </div>
          )}
          {error && (
            <div className="text-center mt-20">
              <p className="text-danger-500 mb-4">{error}</p>
              <button onClick={handleRetry} className="btn btn-primary">重试</button>
            </div>
          )}
          {messages.map((m, i) => (
            <div key={i} className={`flex ${m.role === 'user' ? 'justify-end' : 'justify-start'}`}>
              <div className={`max-w-[70%] rounded-lg px-4 py-2 ${
                m.role === 'user' ? 'bg-primary-600 text-white' :
                m.role === 'assistant' ? 'bg-white border border-border-default text-text-primary' :
                m.role === 'tool' ? 'bg-info-bg border border-info-500/20 text-text-primary' :
                'bg-surface-subtle text-text-muted text-center w-full'
              }`}>
                {m.tool && <div className="text-xs text-text-muted mb-1">Tool: {m.tool}</div>}
                {m.content}
              </div>
            </div>
          ))}
        </div>
        <div className="border-t border-border-default p-4">
          {pendingActions > 0 && (
            <div className="mb-3 px-3 py-2 bg-warning-bg border border-warning-500/20 rounded-md text-sm text-warning-500">
              Agent 正在等待你的确认 ({pendingActions} 项待处理)
            </div>
          )}
          <div className="flex gap-2">
            <textarea
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleSend() }}}
              placeholder="输入消息，Enter 发送..."
              className="input flex-1 min-h-[44px] max-h-32 resize-none"
              rows={1}
              disabled={loading}
            />
            <button onClick={handleSend} disabled={loading || !input.trim()} className="btn btn-primary">
              {loading ? '停止' : '发送'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
