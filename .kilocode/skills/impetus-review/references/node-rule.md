# Node.js / TypeScript / Next.js 规则集 v1.1（增量规则）

> 适用于 `*.ts` `*.tsx` `*.js` `*.jsx` 文件审查。优先级映射见 `common-review-guide.md`。
> 默认针对 Next.js 14 App Router + React Query 项目。

## 1. 架构与逻辑（arch）

### arch-node:1.1.1 客户端/服务端边界（P0）
- `'use client'` 指令**必须** 在需要 hooks / 事件处理的组件顶部
- 服务端组件**禁止** 用 `useState` / `useEffect` / `window`
- 敏感密钥（API key、master key）**禁止** 出现在 client bundle（必须用 `NEXT_PUBLIC_` 前缀且仅用于公开值）

### arch-node:1.1.2 React Hooks 规则（P0）
- `useEffect` deps 数组**必须** 完整（ESLint `react-hooks/exhaustive-deps`）
- 条件分支中**禁止** 调用 hooks
- 自定义 hook 必须以 `use` 开头

### arch-node:1.1.3 Server Component 数据获取（P1）
- 优先 RSC + `fetch` (Next.js 14)
- 客户端获取用 React Query `useQuery`（禁止裸 `useEffect + fetch`）
- Mutation 用 `useMutation` + `onSuccess` 失效 query

### arch-node:1.2.1 类型安全（P1）
- **禁止** `any`（除非注释说明原因：`// eslint-disable-next-line @typescript-eslint/no-explicit-any -- <理由>`）
- 函数参数 / 返回值**必须** 显式类型
- 外部 API 响应用 Zod / io-ts 校验

## 2. 安全（sec）

### sec-node:2.1.1 XSS（P0）
- `dangerouslySetInnerHTML` **禁止** 未经 DOMPurify 净化
- 用户输入**禁止** 直接插入 `href`（`javascript:` 协议）
- URL 用 `encodeURIComponent`

### sec-node:2.1.2 CSRF（P0）
- 状态变更请求（POST/PUT/DELETE）**必须** CSRF token 或 SameSite=Strict cookie
- Mutation endpoint 检查 origin

### sec-node:2.1.3 敏感数据泄露（P0）
- **禁止** `console.log` 包含 token、密码、API key
- 错误对象**禁止** 直接 `JSON.stringify` 给前端（含 stack）
- 错误边界（ErrorBoundary）**必须** 兜底

### sec-node:2.2.1 认证状态（P0）
- Token 存 `httpOnly` cookie（不存 localStorage）
- 登录态读取 SSR cookie + 中间件
- 路由**必须** middleware 鉴权

### sec-node:2.2.2 输入校验（P1）
- 表单**必须** Zod schema
- URL 参数解析用 `useSearchParams` + 校验
- API 客户端用 generated types

## 3. 性能（perf）

### perf-node:3.1.1 不必要的 re-render（P1）
- `React.memo` / `useMemo` / `useCallback` 适度使用（避免过度）
- Context value 用 `useMemo` 包裹
- 列表 key **必须** 稳定 ID（禁止 index）

### perf-node:3.1.2 Bundle 体积（P1）
- 大依赖用 dynamic import：`const Comp = dynamic(() => import('./Comp'))`
- 图表、编辑器、Monaco 等重组件**必须** dynamic
- Lucide-react 图标**禁止** 全量 import（已 tree-shake，确认 v0.300+）

### perf-node:3.2.1 N+1 / 串行请求（P0）
- 独立请求**必须** `Promise.all` 并行
- 列表分页用 server-side pagination（禁止一次性 fetch all）

## 4. 可维护性（maint）

### maint-node:4.1.1 命名（P2）
- 组件 PascalCase，hook camelCase (前缀 `use`)
- 文件名与默认导出名一致
- 常量 UPPER_SNAKE_CASE

### maint-node:4.1.2 组件拆分（P1）
- 单组件 > 300 行需拆分
- 业务逻辑提取为 hook

### maint-node:4.2.1 错误边界（P1）
- 路由级 `error.tsx`（App Router）
- 关键功能用 `<ErrorBoundary>`

## 5. 测试（test）

### test-node:5.1.1 测试覆盖（P0）
- 关键 hooks 单元测试（Vitest / Jest）
- 组件用 React Testing Library
- E2E（Playwright）覆盖登录 + 核心流程

## 6. 风格（style）

### style-node:6.1.1 ESLint / Prettier（P0）
- **必须** `next lint` 通过
- **必须** Prettier 格式化

### style-node:6.2.1 类型 vs interface（P2）
- 优先 `type`（union / intersection）
- `interface` 仅用于可扩展的对象

### style-node:6.2.2 async/await（P2）
- 优先 `async/await` 而非 `.then()`
- `void` 标记 fire-and-forget

---

## 已知误报抑制

| 模式 | 误报原因 | 抑制方法 |
|------|----------|----------|
| `useEffect` deps 不全 | 内部函数引用 | 提取为 useCallback |
| `any` | 第三方库类型缺失 | `// eslint-disable` 注释 |
| `dangerouslySetInnerHTML` | 静态 HTML | 检查 content 是常量 |
