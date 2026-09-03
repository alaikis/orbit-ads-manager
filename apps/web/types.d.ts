declare module 'next' {
  const Next: any
  export default Next
  export type Metadata = any
  export const Metadata: any
  export function Image(props: any): any
  export function Link(props: any): any
  export function redirect(path: string): never
  export function notFound(): never
}

declare module 'next/link' {
  import { ComponentType } from 'react'
  const Link: ComponentType<any>
  export default Link
}

declare module 'next/navigation' {
  export function useRouter(): any
  export function usePathname(): string
  export function useSearchParams(): any
}

declare module 'next/font/google' {
  export function Inter(options: any): any
}

declare module 'lucide-react' {
  export const TrendingUp: any
  export const TrendingDown: any
  export const Plus: any
  export const Send: any
  export const Square: any
  export const Bot: any
  export const Bell: any
  export const Store: any
  export const Monitor: any
  export const Settings: any
  export const BarChart3: any
  export const Package: any
  export const Rss: any
  export const ListChecks: any
  export const ChevronLeft: any
  export const ChevronRight: any
  export const RefreshCw: any
  export const Search: any
  export const Filter: any
  export const ArrowLeft: any
  export const DollarSign: any
  export const Play: any
  export const Pause: any
  export const Download: any
  export const LayoutDashboard: any
  export const ShoppingBag: any
  export const FileText: any
}

declare module '@tanstack/react-query' {
  export function useQuery(options: any): any
  export function useMutation(options: any): any
  export function useQueryClient(): any
  export const QueryClient: any
  export const QueryClientProvider: any
}

declare module 'zustand' {
  export interface StoreApi<T extends object> {
    getState: () => T
    setState(partial: Partial<T> | ((state: T) => Partial<T>)): void
    subscribe(listener: (state: T) => void): () => void
  }
  export type UseBoundStore<T extends object> = {
    (): T
    (selector: (state: T) => any): any
  } & StoreApi<T>
  export function create<T extends object>(): {
    (fn: (set: { setState: (partial: Partial<T>) => void }) => T): UseBoundStore<T>
  }
}

declare module 'zustand/middleware' {
  export function persist<T extends object>(fn: any, options: any): any
}

declare module 'recharts' {
  export const LineChart: any
  export const Line: any
  export const XAxis: any
  export const YAxis: any
  export const CartesianGrid: any
  export const Tooltip: any
  export const ResponsiveContainer: any
}

declare module 'next-intl' {
  export function useTranslations(): any
  export function useLocale(): string
}
