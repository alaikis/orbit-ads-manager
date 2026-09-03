import type { Metadata } from 'next'
import Providers from './providers'
import './globals.css'

export const metadata: Metadata = {
  title: 'Orbit - 广告智能中枢',
  description: 'AI Agent 驱动的多平台广告自动化与智能化管理平台',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <body>
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  )
}
