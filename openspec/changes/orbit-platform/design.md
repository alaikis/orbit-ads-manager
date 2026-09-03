# Design: 统一凭据管理与按平台授权体系

> 高层架构决策，详细设计见 `docs/superpowers/specs/2026-09-03-connection-platform-design.md`（design 阶段产出）。

## 架构选型

### 1. 实体模型

```
workspace (1) ──< connection (1) ──< connection_credential (n)
                                  ──< connection_token (n, 0..1)
                                  ──< oauth_state (n, callback 校验)
        connection (1) ──< store (n)        # 一个连接可派生多个店铺（如 Meta Business 多 Page）
        connection (1) ──< ad_account (n)   # 一个连接可派生多个广告账户（如 Google Ads MCC 子账号）
```

### 2. 三种 auth flow

| type | 例子 | 前端动作 | 后端动作 |
|------|------|---------|---------|
| `api_key` | WooCommerce, Shopify | 直接表单提交 `consumer_key/secret` 或 access_token | 加密入 `connection_credential` |
| `developer_token` | Bing Ads dev_token | 表单提交 client_id + client_secret + dev_token | 加密入 `connection_credential` |
| `oauth` | Google Ads/Shopping, Meta, Bing, TikTok | 「授权」按钮 → 跳到 provider authorize URL → callback 路由回后端 | 启动时生成 `state`（含 conn_id 草稿），callback 用 code 换 token，加密入 `connection_token`，创建/更新 `connection` |

OAuth 步骤：
1. 前端调 `POST /api/v1/connections/oauth/start` 携带 `{platform, scopes, name, draft_fields}` → 返回 `authorize_url` + 临时 `state`；
2. 前端 `window.location = authorize_url`；
3. 用户在 Google/Meta/Bing/TikTok 同意，回调到 `GET /api/v1/oauth/callback?code=...&state=...`；
4. 后端用 code 换 `access_token` + `refresh_token`（若有）→ 加密写 `connection_token` → 创建 `connection` → 重定向到前端 `/settings/connections?connected=<id>`。

### 3. 平台 schema 注册

`apps/api/internal/platform/registry.go`：

```go
var Schemas = map[string]PlatformSchema{
  "woocommerce": {
    AuthFlow: "api_key",
    Fields: []Field{{Key:"base_url", Type:"text", Required:true}, {Key:"consumer_key", Type:"text", Required:true}, {Key:"consumer_secret", Type:"password", Required:true}},
  },
  "shopify": {
    AuthFlow: "api_key",
    Fields: []Field{{Key:"shop_domain", Type:"text", Required:true, Hint:"*.myshopify.com"}, {Key:"access_token", Type:"password", Required:true, Hint:"Admin API access token"}},
  },
  "google_ads": {
    AuthFlow: "oauth",
    OAuthProvider: "google",
    Scopes: []string{"https://www.googleapis.com/auth/adwords"},
    Fields: []Field{{Key:"developer_token", Type:"password", Required:true}},
    PostAuthActions: []string{"list_ad_accounts"},
  },
  "google_shopping": {
    AuthFlow: "oauth",
    OAuthProvider: "google",
    Scopes: []string{"https://www.googleapis.com/auth/content"},
    Fields: []Field{},
    PostAuthActions: []string{"list_merchant_accounts"},
  },
  "meta": {
    AuthFlow: "oauth",
    OAuthProvider: "meta",
    Scopes: []string{"ads_management","ads_read","business_management","catalog_management","commerce_account_manage_orders","commerce_account_read_orders"},
    Fields: []Field{},
    PostAuthActions: []string{"list_ad_accounts"},
  },
  "bing": {
    AuthFlow: "oauth",
    OAuthProvider: "microsoft",
    Scopes: []string{"https://ads.microsoft.com/ads.manage"},
    Fields: []Field{{Key:"developer_token", Type:"password", Required:true}},
  },
  "tiktok": {
    AuthFlow: "oauth",
    OAuthProvider: "tiktok",
    Scopes: []string{"user.info.basic","ads.read","ads.management"},
    Fields: []Field{},
  },
  "llm": {AuthFlow: "api_key", Fields: []Field{{Key:"api_key", Type:"password", Required:true}, {Key:"base_url", Type:"text"}, {Key:"model", Type:"text", Required:true}}},
  "smtp": {AuthFlow: "api_key", Fields: []Field{{Key:"host", Type:"text", Required:true}, {Key:"port", Type:"text", Required:true}, {Key:"user", Type:"text", Required:true}, {Key:"pass", Type:"password", Required:true}}},
}
```

### 4. Token 生命周期

- `connection_token.expires_at` 存绝对时间；
- asynq 调度器每 5 分钟扫描即将到期（< 10 min）的 token，按平台调用 refresh：
  - Google：`https://oauth2.googleapis.com/token`（refresh_token grant）
  - Meta：长 token 已无 refresh 机制，仅在 60 天后提醒用户续期；
  - Microsoft：标准 OAuth refresh；
  - TikTok：标准 OAuth refresh；
- 失败 3 次标记 `connection.status=error` + 写通知。

### 5. /settings/connections 页面

- 顶部 Tabs：全部 / 店铺 / 广告平台 / AI / 邮件；
- 主区卡片网格：每卡片 = platform icon + name + status + 最近刷新时间 + 「测试 / 编辑 / 删除」按钮；
- 「+ 新建连接」按钮 → Stepper 模态：
  1. 选 platform（图标的卡片）
  2. 填 schema fields（参数型）或「授权」按钮（OAuth）
  3. 选 scope（仅 OAuth）
  4. 完成 → 跳回 / 触发 PostAuthActions

### 6. Store/AdAccount 引入

- `POST /api/v1/stores` 接受 `conn_id`（推荐）；不传则旧字段保留；
- `POST /api/v1/ad-accounts` 接受 `conn_id`；不传则旧字段保留；
- `GET /api/v1/stores?conn_id=X` 列出某 connection 派生的店铺；
- `GET /api/v1/ad-accounts?conn_id=X` 同上；
- 触发 `POST /api/v1/connections/{id}/sync` 走 PostAuthActions（例：列出 Google Ads 子账号、列出 Google Merchant Center 账号、列出 Woo shops）。

### 7. 兼容性 / 降级

- 旧 `Provider` 表保留（用于 LLM/SMTP 兼容，3 个月内迁移提示）；
- 旧 `Store/AdAccount` 不带 conn_id 仍然可用；
- `/settings/providers` 顶部加横幅「已迁移到 /settings/connections，点击跳转」。

## 文件结构（新增/修改）

```
apps/api/
├── internal/
│   ├── connection/
│   │   ├── models.go         # Connection / Credential / Token / OAuthState
│   │   ├── handlers.go       # CRUD + /oauth/start + /oauth/callback
│   │   ├── crypto.go         # AES-GCM helpers (复用 config.MasterKey)
│   │   ├── service.go        # 业务：refresh、sync、validate
│   │   ├── registry.go       # 平台 schema 注册（先硬编码 Go map，v2 再迁 DB）
│   │   └── oauth_*.go        # 各 provider authorize/exchange 实现
│   ├── platform/
│   │   ├── google_ads/       # listAccessibleCustomers、campaigns CRUD
│   │   ├── google_shopping/  # products、shippingsettings、returnpolicyonline
│   │   ├── meta/             # /me/adaccounts、catalogs、orders
│   │   ├── bing/             # customer accounts
│   │   ├── tiktok/           # /advertiser/list
│   │   ├── woocommerce/      # REST client
│   │   └── shopify/          # Admin API client
│   ├── store/handlers.go     # +conn_id 字段
│   ├── adplatform/handlers/accounts.go # +conn_id
│   └── scheduler/token_refresh.go  # asynq 任务
└── pkg/crypto/aesgcm.go

apps/web/app/(dashboard)/settings/
├── connections/page.tsx      # 新：统一凭据管理
├── providers/page.tsx        # 改：重定向到 /settings/connections?tab=ai
apps/web/app/(dashboard)/accounts/page.tsx  # 改：+ "从已连接平台导入"
apps/web/app/(dashboard)/stores/page.tsx     # 改：+ "从已连接 Woo/Shopify 导入"
apps/web/lib/api/services/connection.service.ts  # 新
apps/web/lib/api/services/store.service.ts        # 改：+conn_id
apps/web/lib/api/services/ad-account.service.ts   # 改：+conn_id
```

## 风险与缓解

| 风险 | 缓解 |
|------|------|
| OAuth callback 域名/Caddy 配置 | design 阶段确认 Caddyfile 允许 `GET /api/v1/oauth/callback` |
| Token 加密 key 不一致 | 启动时校验 `Encryption.MasterKey` 长度=32，否则拒启 |
| 旧 Provider 迁移数据 | 在 Connection 表加迁移脚本：`INSERT INTO connection ... SELECT FROM provider WHERE type IN ('llm','smtp')` |
| 多平台 schema 扩展 | registry 集中，新增 platform 仅 1 文件 |
| 长时间运行的 access_token 过期 | token_refresh 任务 + 失败告警 + 401 自动降级提示重连 |
| Google Shopping 范围（shippingsettings/returnpolicyonline）复杂 | v1 只实现 OAuth + products.list + accounts.list；shipping/return 列为 v1.1 |
