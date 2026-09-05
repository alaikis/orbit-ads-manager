---
impetus_change: orbit-platform
role: technical-design
canonical_spec: openspec
archived_with: openspec/changes/archive/2026-09-03-orbit-platform
status: archived
archived_at: 2026-09-04
---

# Design Doc: 统一凭据管理与按平台授权体系

> 版本：v1 | 状态：design | Change: orbit-platform

## 1. 架构总览

### 1.1 问题与目标

Orbit 目前对店铺、广告账户、AI/邮件接入的方式零散。本方案引入 `Connection` 一等公民实体，统一封装"凭据怎么拿"和"凭据怎么用"两层。

### 1.2 实体关系

```
workspace (1) ──< connection (1) ──< connection_credential (n)
                                  ──< connection_token (n, 0..1)
                                  ──< oauth_state (n, callback 校验)
        connection (1) ──< store (n)
        connection (1) ──< ad_account (n)
```

### 1.3 三种 auth flow

| flow type | 触发 | 前端 | 后端 |
|-----------|------|------|------|
| `api_key` | Woo/Shopify 表单提交 | 直接 POST config | 加密 → connection_credential |
| `developer_token` | Bing 表单提交 | 直接 POST config | 加密 → connection_credential |
| `oauth` | Google/Meta/Bing/TikTok 授权 | POST /oauth/start → window.location 跳转 → callback | state 生成 → code 换 token → connection_token → connection |

OAuth 步骤：
1. `POST /api/v1/connections/oauth/start {platform, scopes, name, draft_fields}` → `authorize_url + state`
2. 前端 `window.location = authorize_url`
3. 用户在 provider 同意，回调 `GET /api/v1/oauth/callback?code=&state=`
4. 后端 code 换 access_token + refresh_token → 加密写入 → 创建 connection → 触发 PostAuthActions → 重定向到 `/settings/connections?connected=<id>`

## 2. 数据模型

### connection
```
id, workspace_id, tenant_id, type[oauth|api_key|developer_token|smtp],
platform[google_ads|google_shopping|meta|bing|tiktok|woocommerce|shopify|llm|smtp],
name, status[active|error|paused], scopes, last_refreshed_at, last_error, created_at, updated_at
```

### connection_credential
```
conn_id, key(字段名如 consumer_key), encrypted_value(AES-GCM), last_4(脱敏尾号)
```

### connection_token（OAuth 专用）
```
conn_id, access_token_enc, refresh_token_enc, expires_at, scopes, system_user_id
```

### oauth_state（callback 防 CSRF）
```
state, platform, workspace_id, name, draft_fields(JSON), expires_at, used[bool]
```

## 3. Token 生命周期

- `expires_at` 存绝对时间
- asynq 调度器每 5 分钟扫描即将到期（< 10 min）token，按平台分发 refresh
  - Google：`POST https://oauth2.googleapis.com/token` (grant_type=refresh_token)
  - Microsoft：标准 OAuth refresh
  - TikTok：标准 OAuth refresh
  - Meta：长 token 60 天自动过期，刷新任务提醒续期（无 refresh_token 机制）
- 失败 3 次 → `connection.status='error'` + 写 `notification`

## 4. 平台注册（registry）

`apps/api/internal/platform/registry.go`：

```go
var Schemas = map[string]PlatformSchema{...}
```

| platform | auth_flow | OAuthProvider | Scopes | PostAuth |
|----------|-----------|---------------|--------|---------|
| woocommerce | api_key | - | - | validate + shop name |
| shopify | api_key | - | - | validate + shop name |
| google_ads | oauth | google | adwords | listAccessibleCustomers |
| google_shopping | oauth | google | content | accounts.list |
| meta | oauth | meta | ads_management,ads_read,... | /me/adaccounts |
| bing | oauth+dev_token | microsoft | ads.manage | GetAccountsInfo |
| tiktok | oauth | tiktok | ads.read,ads.management | /advertiser/list |
| llm | api_key | - | - | - |
| smtp | api_key | - | - | validate connect |

## 5. /settings/connections 页面

- Tabs：全部 / 店铺 / 广告平台 / AI / 邮件
- 卡片网格：platform icon + name + status + 最近刷新 + 操作按钮
- 「+ 新建连接」Stepper：选 platform → 填 fields（参数型）/ 跳授权（OAuth）→ 选 scope → 完成

## 6. Store / AdAccount 改造

- `POST /api/v1/stores` + `POST /api/v1/ad-accounts` 接受可选 `conn_id`
- `GET /api/v1/stores?conn_id=X` 过滤
- `POST /api/v1/connections/{id}/sync/stores` 和 `/sync/ad-accounts` 走 PostAuthActions

## 7. Google Shopping 子资源（v1 范围）

- `GET /api/v1/connections/{id}/google-shopping/products` — 列表商家商品
- `GET/POST /api/v1/connections/{id}/google-shopping/shipping-settings` — etag 流程
- `GET/POST /api/v1/connections/{id}/google-shopping/return-policy` — 退货政策
- `localInventory` 留 v1.1

## 8. 加密实现

```go
// pkg/crypto/aesgcm.go
func Encrypt(plaintext []byte) (ciphertext, nonce []byte, err error)
func Decrypt(ciphertext, nonce []byte) ([]byte, error)
```

复用 `config.AppCfg.Encryption.MasterKey`（32 bytes），启动时校验长度。

## 9. 文件结构（新增/修改）

```
apps/api/
  internal/connection/  (models.go, crypto.go, registry.go, service.go, handlers.go, routes.go)
  internal/platform/    (google_ads/, google_shopping/, meta/, bing/, tiktok/, woocommerce/, shopify/)
  pkg/crypto/aesgcm.go
  internal/scheduler/token_refresh.go
apps/web/
  app/(dashboard)/settings/connections/page.tsx  (新)
  app/(dashboard)/settings/providers/page.tsx    (改：横幅降级)
  app/(dashboard)/accounts/page.tsx              (改：+ 从平台导入)
  app/(dashboard)/stores/page.tsx                (改：+ 从 Woo/Shopify 导入)
  lib/api/services/connection.service.ts  (新)
```

## 10. 风险与缓解

| 风险 | 缓解 |
|------|------|
| OAuth callback 域名/Caddy 配置 | design 阶段确认 Caddyfile 允许 `GET /api/v1/oauth/callback` |
| Token 加密 key 不一致 | 启动时校验 length=32 |
| 旧 Provider 迁移数据 | 一次性脚本 `INSERT INTO connection SELECT FROM provider WHERE type IN ('llm','smtp')` |
| Google Shopping 范围过大 | v1 只做 products.list + accounts.list；shipping/return v1.1 |
| Meta 长 token 60 天过期 | 刷新任务提醒，失败 3 次告警 |
