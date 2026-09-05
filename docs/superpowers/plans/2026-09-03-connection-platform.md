---
change: orbit-platform
design-doc: docs/superpowers/specs/2026-09-03-connection-platform-design.md
base-ref: bd97d0b7fee259a197fa5e55c15cde6401bc1eae
issue-id: null
issue-type: story
workflow: full
archived-with: openspec/changes/archive/2026-09-03-orbit-platform
archived-at: 2026-09-04
---

# Plan: 统一凭据管理与按平台授权体系

> 实施计划，所有 task 完成后实现 Review 窗口 + 按 `impetus/reference/commit-convention.md` 一次最终 build 提交。
> Plan 内容与 Design Doc (`docs/superpowers/specs/2026-09-03-connection-platform-design.md`) 对齐。

## 范围

- v1：WooCommerce、Shopify、Meta Marketing、Google Shopping、LLM、SMTP 6 个平台
- v1.1：Google Ads、Bing、TikTok、shipping settings、return policy online
- 旧 Provider 降级为 LLM/SMTP 子集

## 任务清单

### 阶段 1：基础设施（Model + Crypto）

- [ ] **T1.1** 创建 `migrations/2026090301_connection.up.sql`：`connection`、`connection_credential`、`connection_token`、`oauth_state` 四张表 + 索引（workspace_id、platform、status）
- [ ] **T1.2** 在 `apps/api/pkg/database/migrate.go` 注册新 migration
- [ ] **T1.3** 创建 `apps/api/pkg/crypto/aesgcm.go`：`Encrypt/Decrypt` 包装
- [ ] **T1.4** 在 `apps/api/config/config.go` 启动时校验 `Encryption.MasterKey` 长度=32
- [ ] **T1.5** 创建 `apps/api/internal/connection/models.go`：`Connection`、`ConnectionCredential`、`ConnectionToken`、`OAuthState`

### 阶段 2：Registry（平台 schema）

- [ ] **T2.1** `apps/api/internal/connection/registry.go`：定义 `PlatformSchema{AuthFlow, OAuthProvider, Scopes, Fields[], PostAuthActions[]}` 与 `Schemas` map
- [ ] **T2.2** 注册 `woocommerce`、`shopify`、`llm`、`smtp`（api_key 流程）
- [ ] **T2.3** 注册 `google_ads`、`google_shopping`、`meta`、`bing`、`tiktok`（oauth 流程 + PostAuth 列表）

### 阶段 3：Service + Handler

- [ ] **T3.1** `apps/api/internal/connection/service.go`：`List/Get/Create/Update/Delete/Test/Sync`，全部 workspace 隔离
- [ ] **T3.2** `apps/api/internal/connection/handlers.go`：REST 路由 + JSON 序列化
- [ ] **T3.3** `apps/api/internal/connection/routes.go`：`RegisterConnectionRoutes(api)` 内部挂 auth + tenant 中间件
- [ ] **T3.4** `apps/api/internal/connection/oauth.go`：`Start(platform, scopes, name, draft)` 生成 state + 写 oauth_state + 返回 authorize_url
- [ ] **T3.5** `apps/api/internal/connection/oauth_callback.go`：`Callback(code, state)` 校验 → token exchange → 加密入库 → 创建 connection
- [ ] **T3.6** `apps/api/cmd/server/main.go` 注册 `RegisterConnectionRoutes`

### 阶段 4：平台客户端（PostAuth 实际调用）

- [ ] **T4.1** `apps/api/internal/platform/woocommerce/client.go`：`GET /wp-json/wc/v3/system_status` 验证 + 拿 shop name
- [ ] **T4.2** `apps/api/internal/platform/shopify/client.go`：`GET /admin/api/2024-10/shop.json` 验证
- [ ] **T4.3** `apps/api/internal/platform/google_ads/client.go`：`listAccessibleCustomers` 派生候选 ad_account
- [ ] **T4.4** `apps/api/internal/platform/google_shopping/client.go`：`accounts.list`、`products.list`、`accounts.shippingSettings.get/insert`（etag 流程）
- [ ] **T4.5** `apps/api/internal/platform/google_shopping/return_policy.go`：`returnpolicyonline` list/insert
- [ ] **T4.6** `apps/api/internal/platform/meta/client.go`：`/me/adaccounts`、`/me/businesses` 派生 ad_account 与 catalog
- [ ] **T4.7** `apps/api/internal/platform/bing/client.go`：`GetAccountsInfo`（留 v1.1 接口骨架）
- [ ] **T4.8** `apps/api/internal/platform/tiktok/client.go`：`/advertiser/list/`（留 v1.1 接口骨架）

### 阶段 5：Store / AdAccount 接受 conn_id

- [ ] **T5.1** `model.Store` 加 `ConnID *uint64` 字段 + 迁移
- [ ] **T5.2** `model.AdAccount` 加 `ConnID *uint64` 字段（已存在，验证非空约束）
- [ ] **T5.3** `POST /api/v1/stores` handler 接受 `conn_id`，不传走旧路径
- [ ] **T5.4** `POST /api/v1/ad-accounts` handler 接受 `conn_id`
- [ ] **T5.5** `GET /api/v1/stores?conn_id=X` / `GET /api/v1/ad-accounts?conn_id=X` 过滤
- [ ] **T5.6** `POST /api/v1/connections/{id}/sync` 路由：按 platform 分发 PostAuthAction

### 阶段 6：Token 刷新 worker

- [ ] **T6.1** `apps/api/internal/scheduler/token_refresh.go`：asynq 5 min 周期任务
- [ ] **T6.2** Google refresh：标准 `POST /token` (grant_type=refresh_token)
- [ ] **T6.3** Microsoft refresh：标准 OAuth
- [ ] **T6.4** TikTok refresh：标准 OAuth
- [ ] **T6.5** 失败 3 次 → `connection.status='error'` + 写 `notification`

### 阶段 7：前端 Connection 页面

- [ ] **T7.1** `apps/web/lib/api/services/connection.service.ts`：`list/get/create/update/delete/test/sync + platforms + oauthStart`
- [ ] **T7.2** `apps/web/app/(dashboard)/settings/connections/page.tsx`：tabs + 卡片网格 + 「+ 新建连接」Stepper
- [ ] **T7.3** Stepper 步骤 1：选 platform（icon 卡片）
- [ ] **T7.4** Stepper 步骤 2：动态表单按 `platform.fields` 渲染
- [ ] **T7.5** Stepper 步骤 3：OAuth 选 scope 后调 `oauthStart` 跳 window.location
- [ ] **T7.6** Callback 回来后 `/settings/connections?connected=<id>` 自动 invalidate
- [ ] **T7.7** 卡片：测试 / 编辑 / 删除按钮 + 状态徽章

### 阶段 8：兼容与迁移

- [ ] **T8.1** 旧 `Provider` 表（type IN llm/smtp）一次性脚本迁到 `Connection`
- [ ] **T8.2** `/settings/providers` 顶部加横幅「已迁移到 /settings/connections，3 个月后下线」
- [ ] **T8.3** 旧 Store/AdAccount 数据保留（conn_id=null）继续可用

### 阶段 9：OAuth callback Caddy

- [ ] **T9.1** 确认 `/etc/caddy/sites/ads.alaikis.com.caddyfile` 不重写 `/api/v1/oauth/callback` 路径
- [ ] **T9.2** 添加 `/oauth/*` 子路由不重写（如有 strip_prefix 之类）

## 关键设计决策（与 Design Doc 一致）

- `Connection.conn_id` 与 Store/AdAccount 是弱耦合（候选 B）：conn_id 可选
- 平台 schema 候选 C：JSON config 文件 + 启动 seed 到 DB（先 Go map，v2 再迁）
- Token 加密：复用 `config.AppCfg.Encryption.MasterKey`
- WooCommerce Basic Auth header（不是 query string）
- Meta 先只做 User access token
- Google Shopping 与 Google Ads 视为两个独立 platform

## 验证策略

- **每个 task 完成**：本任务范围内的 curl 烟雾测试 + 必要时单元测试
- **所有 task 完成**：
  - 前端 build: `cd apps/web && npm run build` 必须成功
  - 后端 build: `cd apps/api && GOOS=linux GOARCH=amd64 go build -o orbit-server-new ./cmd/server`
  - 部署：上传到 `portal.gusty.top:/apphub/ads.alakis.com/bin/orbit-server-linux-amd64` + `systemctl restart orbit`
  - E2E：Woo 完整链路：表单 → 创建 connection → /connections/{id}/sync → store 自动出现
- **生产验证**（`migration_test.py` 已存在）：`bash portal.gusty.top:migration_test.sh` 跑全部端点

## 提交策略

按 Impetus 规则：
- `executing-plans` 模式默认一次 change 级 build 提交（在实现 Review 窗口通过后）
- 中途不自动 commit；保留 unstaged 状态直到实现 Review 窗口
- 因为 `issue_id=null`，不自动 commit，需用户在实现 Review 窗口明确同意才做最终 commit

## 不在范围内

- Bing/TikTok v1 阶段只占位（OAuth 流程已设计，但 PostAuth 实现留 v1.1）
- Google Shopping shippingsettings/returnpolicyonline v1.1
- 跨工作空间共享 connection
- Connection RBAC
- 自定义平台插件机制
