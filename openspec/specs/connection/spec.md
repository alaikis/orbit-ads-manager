# Connection Spec

> 主规格 - 统一凭据管理与按平台授权
> 来源: `openspec/changes/archive/2026-09-03-orbit-platform` (已归档)
> 适用 tenant_id 隔离模式（项目无 workspace context）

## Purpose

为每个租户提供统一的"连接"抽象，封装平台凭证、token 生命周期、OAuth 授权流。连接 (connection) 是租户与外部平台（WooCommerce/Shopify/Google/Meta/Bing/TikTok/SMTP/LLM）的信任边界。

## Requirements

### Requirement: Connection entity model

系统 **MUST** 为每个租户维护独立的 connection 集合，按 `tenant_id` 隔离。每个 connection 持有：

- `id` (PK), `tenant_id`, `platform` (枚举), `name` (用户可读), `status` (active/error/disabled), `last_error`, `created_at`, `updated_at`
- 可选 `conn_id` (null 表示旧 Provider 兼容模式)
- 关联 `connection_credential` (0..n) 存储 AES-256-GCM 加密的 api_key/secret/developer_token
- 关联 `connection_token` (0..1) 存储加密的 access_token/refresh_token + expires_at + scope
- 关联 `oauth_state` (0..n) 短期 state (5 分钟 TTL) 用于 OAuth callback 校验

#### Scenario: 创建 api_key 类型连接
- **WHEN** 用户提交 `POST /api/v1/connections` 携带 `{platform, name, fields: {base_url, consumer_key, consumer_secret}}`
- **THEN** 系统加密字段写 `connection_credential`，创建 `connection` 状态 `active`，返回 201
- **AND** 立即触发 `test_connection` 验证可连通性，更新 `status` 或 `last_error`

#### Scenario: 创建 oauth 类型连接
- **WHEN** 用户点击「Google 授权」按钮，前端调 `POST /api/v1/connections/oauth/start` 携带 `{platform, scopes, name, draft_fields}`
- **THEN** 系统生成 PKCE `state` 写入 `oauth_state` (5min TTL)，返回 `authorize_url`
- **AND** 前端 `window.location = authorize_url`，用户在 Google 同意后回调 `GET /api/v1/oauth/callback?code=...&state=...`
- **AND** 后端用 code 换 token，加密写 `connection_token`，创建 connection，redirect 到 `/settings/connections?connected=<id>`

### Requirement: Token encryption at rest

系统 **MUST** 使用 AES-256-GCM (12 字节 nonce, 32 字节 key) 加密所有 `connection_credential` 与 `connection_token` 字段。

- `ENCRYPTION_MASTER_KEY` 从环境变量读取，长度严格 32 字节
- nonce 每次加密唯一（`crypto/rand`）
- 密文格式: `base64(nonce || ciphertext || tag)`
- 解密失败 **MUST** 返回 500 而非泄漏明文

#### Scenario: 加密写入
- **WHEN** 任意 `credential` 或 `token` 字段写入
- **THEN** 系统使用 AES-256-GCM 加密后存储
- **AND** 永不写明文到 DB

### Requirement: Token refresh worker

系统 **MUST** 定时 (cron `*/10 * * * *`) 扫描即将过期 (expires_at < now + 5min) 的 token，自动 refresh。

- 支持 Google/Microsoft/TikTok OAuth refresh
- Refresh 失败重试 3 次，写 `status=error` + `last_error`
- Bing refresh_token 永久有效（首次 OAuth 颁发），只需检测 access_token 过期

#### Scenario: Google token 自动续期
- **WHEN** `connection_token.expires_at < now + 5min` AND `platform IN (google_ads, google_shopping)`
- **THEN** 后台 worker 用 refresh_token 调 Google token endpoint
- **AND** 成功后更新 `access_token` (加密) + `expires_at`
- **AND** 失败累计 3 次后 `connection.status='error'`

### Requirement: Connection test & sync

每个 connection **MUST** 支持 `test` 和 `sync` 操作。

- `test`: 调用平台轻量 endpoint 验证连通性
  - WooCommerce: `GET /wp-json/wc/v3/system_status`
  - Shopify: `GET /admin/api/2024-04/shop.json`
  - Google Shopping: `GET /content/v2.1/accounts/authinfo`
  - Meta: `GET /v25.0/me/adaccounts`
  - Bing/TikTok: 调 OAuth userinfo endpoint
- `sync`: 拉取账户/店铺列表写入 `store` 或 `ad_account` 表 (带 `conn_id` 关联)

#### Scenario: WooCommerce test
- **WHEN** 用户点击「Test」按钮
- **THEN** 系统 Basic Auth 调用 `/wp-json/wc/v3/system_status`
- **AND** 200 → `status=active`；非 200 → `last_error` 包含响应码和消息

#### Scenario: Shopify sync stores
- **WHEN** 用户点击「Sync stores」
- **THEN** 系统拉取 `/admin/api/2024-04/shop.json` (单店铺) 或 `/admin/api/2024-04/shops.json` (Plus)
- **AND** 写入 `store` 表，绑定 `conn_id` + `tenant_id`

### Requirement: Connection deletion cascade

系统 **MUST** 在删除 connection 时级联清理关联资源。

- 删除 `connection` → 删除 `connection_credential` + `connection_token` + `oauth_state` (where conn_id = ?)
- `store` / `ad_account` 关联 `conn_id` 的记录 **MUST** 检查：若该 store 有关键数据 (campaigns/products)，阻止删除并提示
- 否则级联 `SET NULL` (保留历史数据但断关联)

#### Scenario: 删除无数据的连接
- **WHEN** 用户 `DELETE /api/v1/connections/:id` 且该 conn_id 无 store/ad_account
- **THEN** 系统删除 connection + credential + token，返回 204

#### Scenario: 删除有关联店铺的连接
- **WHEN** 该 conn_id 关联 store 且 store 有 campaigns
- **THEN** 系统返回 409 Conflict，body 包含 `{blocked_by: ["stores:5"]}`

### Requirement: Multi-tenant isolation

所有 connection 查询 **MUST** 强制 `tenant_id` 过滤，由中间件注入到 context。

- 路径参数 `id` **MUST** 校验属于当前 tenant 的 connection
- 越权访问返回 404（不泄漏存在性）
- 共享租户内不区分角色（owner/admin 均可操作）

#### Scenario: 越权访问
- **WHEN** tenant A 用户用 tenant B 的 connection id 请求
- **THEN** 系统返回 404 not found

## File Structure

```
apps/api/internal/connection/
  models.go        # Connection/Credential/Token/OAuthState structs
  registry.go      # 9 平台 PlatformSchema 注册
  service.go       # CRUD + test/sync 业务逻辑
  handlers.go      # HTTP handlers
  oauth.go         # OAuth start/callback/refresh
  clients.go       # 9 平台 HTTP client (test/sync 实现)
  refresh.go       # cron refresh worker

apps/api/pkg/crypto/
  aesgcm.go        # AES-256-GCM 包装

apps/api/migrations/
  000003_connection_platform.sql   # 4 表 + 1 索引 + 1 字段

apps/web/app/(dashboard)/settings/connections/
  page.tsx         # 完整 UI (Stepper + tabs + 卡片)
```

## Platform Catalog (9 平台)

| platform | auth_flow | key fields | oauth scopes | post_auth_actions |
|----------|-----------|------------|--------------|-------------------|
| woocommerce | api_key | base_url, consumer_key, consumer_secret | - | - |
| shopify | api_key | shop_domain, access_token | - | sync_stores |
| google_ads | oauth | developer_token | adwords | list_ad_accounts |
| google_shopping | oauth | (none) | content | sync_accounts, sync_products |
| meta | oauth | (none) | ads_management, business_management | list_ad_accounts |
| bing | oauth | (none) | ads.read | list_ad_accounts |
| tiktok | oauth | (none) | user.info.basic, ads.read, ads.management | list_ad_accounts |
| llm | api_key | provider, api_key, base_url, model | - | - |
| smtp | api_key | host, port, username, password, from | - | - |
