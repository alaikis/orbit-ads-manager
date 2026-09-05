# Impetus Spec Projection

- Change: orbit-platform
- Phase: design
- Mode: beta
- Context hash: f18b035a1ecd6aa1a5766ba36b227fe41630f8eb0653e59702492aa303f98bd2

Generated-by: impetus-handoff.sh

OpenSpec delta specs remain canonical. This file is a deterministic full-text projection for lower-token Build handoff; regenerate it when source artifacts change.

## openspec/changes/orbit-platform/proposal.md

- Role: proposal
- Source: openspec/changes/orbit-platform/proposal.md
- Lines: 1-94
- SHA256: a51042e7c63851e05bd3b2e55ebded5696b4d1468325ee20e80a095d8af0fb3c

```md
# Proposal: 统一凭据管理与按平台授权体系

## 背景

Orbit 平台目前对**店铺（Store）**和**广告账户（AdAccount）**的接入方式零散且不一致：
- 现有 `/settings/providers` 表单直接收集明文 LLM/SMTP/Google/Meta/Bing 的 `client_id/secret` 等参数；
- 现有 `/accounts` 表单只让用户手填 platform + external_id，不与 Provider 关联；
- 现有 `/stores` 表单只让用户手填 base_url + api_key + api_secret，不区分 Woo/Shopify 授权流程；
- 凭据加密复用 Provider 的 `Config` 字段，但后端 schema 缺失（`Provider` 缺乏 `kind=connection/credential` 区分），且 `AdAccount.ConnID` 默认写 0，与 StoreConfig 脱节；
- 缺 OAuth 2.0 authorization_code 流程的回调端点、token 存储、refresh 任务调度；
- Google Shopping（Merchant API / Content API）需要的是商家账号（Merchant Center account_id）+ OAuth scope `content`（新 Merchant API 为 `https://www.googleapis.com/auth/content`），而不是 AdWords/Ads 的 `adwords` scope；
- Meta Marketing API 同时支持 User access token 与 System user access token，需要根据"代表个人"还是"代表业务"分流；
- WooCommerce 与 Shopify 是 API token / Basic Auth，不走 OAuth；
- Bing Ads 是 Microsoft 账户 OAuth（不同于 Google OAuth）；
- TikTok Ads 也是 OAuth 2.0。

业务上，要让用户在一个工作空间内：
- 创建一次商家连接（Woo/Shopify）→ 自动出商品/订单/Flow/Feed；
- 创建一次广告平台连接（Google Ads / Google Shopping / Meta / Bing / TikTok）→ 拉出可投放的广告账户/商品目录；
- 创建一次 LLM/SMTP 连接 → 给 Agent 与报表用。

## 目标

构建一个**统一的 Connection / Credential 抽象层**，把"凭据怎么拿"和"凭据怎么用"解耦：

1. **统一实体**：所有外部平台接入（店铺、广告平台、AI、邮件）落地为 `Connection`（type=oauth|api_key|developer_token|smtp），每个 Connection 持有**加密**的 credentials；
2. **平台 schema 化**：每个 platform 拥有自己的 `auth_schema`（字段定义 + auth flow），前端按 schema 渲染授权表单；OAuth 平台由后端发起 authorization_code 流，本地回填 token；
3. **资源挂接**：Store 与 AdAccount 通过 `conn_id` 引用 Connection，**一个 Connection 可派生多个资源**（例：Google Ads 一次 OAuth 可列出多个 MCC 子账号，Shopify 一个 store 派生一个 Channel）；
4. **优先范围**（按用户指定优先级）：
   - **Google Shopping**（Merchant API：products/accounts/shippingsettings/returnpolicyonline/regionalinventory）— OAuth scope = `content`
   - **Meta Marketing**（User OAuth + System User Token 两种流）— scope = `ads_management,ads_read,catalog_management,commerce_account_*`
   - **WooCommerce**（REST API + Basic Auth，consumer_key/consumer_secret）
   - **Google Ads**（AdWords API：OAuth scope = `adwords`）
   - **Microsoft Bing Ads**（OAuth）
   - **TikTok Ads**（OAuth）
   - **Shopify**（Admin API access token）
5. **Token 生命周期**：access_token 加密存储、定时刷新、刷新失败告警；
6. **统一入口**：`/settings/connections` 取代原 `/settings/providers` 承担"凭据"维度；`/accounts` 接管"广告账户"维度；`/stores` 接管"店铺"维度。

## 范围

**In scope（本次变更）：**

- 后端：
  - 新表 `connection`（id, workspace_id, tenant_id, type[oauth|api_key|developer_token|smtp], platform[google_ads|google_shopping|meta|bing|tiktok|woocommerce|shopify|smtp|llm], name, status, scopes, last_refreshed_at, last_error, created_at, updated_at）
  - 新表 `connection_credential`（conn_id, key, encrypted_value, last_4）— 加密存储 secret
  - 新表 `connection_token`（conn_id, access_token_enc, refresh_token_enc, expires_at, scopes, system_user_id）— OAuth 专用
  - 新表 `platform_schema`（platform_id, field_key, label, type, required, secret, hint, auth_flow 字段，JSON 存储）
  - `Connection` Service：list/get/create/update/delete + 平台 schema 查询 + OAuth 启动 / OAuth 回调 / token 刷新
  - `Store` Service 改造：create 时接受 `conn_id`（不传则按 base_url+api_key 隐式创建 api_key Connection）
  - `AdAccount` Service 改造：list 增加"按 conn_id 拉真实账号"（Google Ads: listAccessibleCustomers, Meta: /me/adaccounts, TikTok: /advertiser/list）
  - Google Shopping 子资源封装：products、accounts.shippingSettings、returnpolicyonline、localinventory
  - WooCommerce REST client（消费者密钥 + 签名 + pagination）
  - Meta Graph client（按 token 调 `/act_<id>/campaigns` 等）
  - Token 刷新 worker（asynq 定时任务，token 即将到期前 10 分钟刷新）
  - encryption 用现有 `config.AppCfg.Encryption.MasterKey`（AES-256-GCM）
- 前端：
  - `/settings/connections` 页面：左侧按平台分组的连接列表 + 右侧「+ 新建连接」按钮 → 弹出按 platform schema 渲染的表单（参数型直接表单，OAuth 跳转）
  - `/accounts` 页面：保持"广告账户"维度，**新增"从已连接平台导入"按钮** → 调真实 API 列出可投放账户
  - `/stores` 页面：保持"店铺"维度，**新增"从已连接 Woo/Shopify 导入"**（拉 Woo shop info）
  - `/settings/providers` 保留为 LLM/SMTP 子集（降级为 LLM/SMTP 配置）；或重定向到 `/settings/connections?tab=ai`。
- 文档：
  - delta spec 4 个 capability：connection、store、ad-account、platform-integration
  - Design Doc：`docs/superpowers/specs/2026-09-03-connection-platform-design.md`

**Out of scope（后续）：**

- 自定义平台插件机制；
- 多租户 marketplace / 公开 OAuth app 发布；
- 跨工作空间共享 connection；
- 内部 Connection 权限 RBAC（先以"workspace 隔离 + admin 可操作"兜底）；
- 离线/定时同步任务（已存在 SyncJob，扩展字段但不重写）；
- 报表/规则/Agent 对 Connection 的新使用（沿用 AdAccount / Store 抽象，间接受益）。

## 关键决策点（需在 design 阶段 brainstorming 确认）

1. **Connection 与 Store/AdAccount 是否解耦？**  
   候选 A：Connection 是一等公民，Store/AdAccount 必须有 conn_id（强一致）。  
   候选 B：Connection 可选，Store/AdAccount 仍可自带 api_key（弱耦合，向后兼容）。  
   候选 C：Connection 即 Store/AdAccount 的"证书视图"，1:1（最简单）。
2. **OAuth 回调地址**：  
   `https://ads.alaikis.com/api/v1/oauth/callback`（Caddy 反代必须放过 GET 路径，否则 405）。
3. **Token 加密复用现有 `Encryption.MasterKey`**：  
   AES-256-GCM（已有 `getEncryptionKey` 实现），新增 `pkg/crypto` 工具方法。
4. **平台 schema 来源**：  
   候选 A：硬编码 Go map（最简，扩展性差）。  
   候选 B：DB `platform_schema` 表 + 启动时 seed（可热更新但需迁移）。  
   候选 C：JSON config 文件 + DB 缓存（推荐：可版本控制）。
5. **WooCommerce 鉴权**（REST API）：  
   `?consumer_key=ck_xxx&consumer_secret=cs_xxx` 作为 query 简化；或 Basic Auth header（推荐 Basic）。
6. **Meta 两种 token 同时支持**？还是先只做 User access token？
7. **Google Shopping 优先级 > Google Ads 单独账号**？  
   Google Shopping 走 Content/Merchant API（不是 AdWords），可与 Google Ads 共存；为简化可视为两个 platform。
8. **降级策略**：`/settings/providers` 是否立即重定向到 `/settings/connections`？
```

## openspec/changes/orbit-platform/design.md

- Role: open-design-notes
- Source: openspec/changes/orbit-platform/design.md
- Lines: 1-161
- SHA256: 69eb4f0b09db05bdd16b2f0f78fede210c3e2b65e5d3ad5ce73009d1d0de0406

```md
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
```

## openspec/changes/orbit-platform/tasks.md

- Role: tasks
- Source: openspec/changes/orbit-platform/tasks.md
- Lines: 1-119
- SHA256: ad84131647eb198586a896717d645de7f8c8abc47375b7391a0870eb8c030140

```md
# Tasks: 统一凭据管理与按平台授权体系

> 状态机：`[ ]` 未开始 · `[x]` 完成 · `[-]` 跳过 · 阻塞时同步 `.impetus.yaml.build_pause`

## 阶段 0：基础设施

- [ ] **T0.1** 启动 schema 迁移文件 `migrations/2026090301_connection.up.sql`：建 `connection`、`connection_credential`、`connection_token`、`oauth_state` 表 + 索引
- [ ] **T0.2** 注册 Go migration runner（`apps/api/pkg/database/migrate.go`）
- [ ] **T0.3** 新建 `apps/api/pkg/crypto/aesgcm.go`（AES-256-GCM 包装）；在 `config.AppCfg.Encryption.MasterKey` 加载时校验长度=32
- [ ] **T0.4** 注册 `internal/connection` 包（`models.go`、`crypto.go`、`registry.go`、`service.go`、`handlers.go`、`routes.go`）

## 阶段 1：Connection 核心 CRUD

- [ ] **T1.1** `Connection` Service：`Create(list params)` / `Get(id)` / `List(workspace_id)` / `Update(id, fields)` / `Delete(id)` / `Test(id)`（解密后调 provider 验证）
- [ ] **T1.2** `Connection` HTTP handlers：`POST/GET/PATCH/DELETE /api/v1/connections` + `POST /api/v1/connections/{id}/test` + `POST /api/v1/connections/{id}/sync`
- [ ] **T1.3** `GET /api/v1/platforms` 返回所有 platform schema（label、auth_flow、fields、scopes）— 前端用于 Stepper

## 阶段 2：平台 schema 注册（registry）

- [ ] **T2.1** `woocommerce`（api_key）：base_url、consumer_key、consumer_secret
- [ ] **T2.2** `shopify`（api_key）：shop_domain、access_token
- [ ] **T2.3** `google_ads`（oauth + developer_token）：scope=adwords
- [ ] **T2.4** `google_shopping`（oauth）：scope=content，PostAuth=ListMerchantAccounts
- [ ] **T2.5** `meta`（oauth）：scope=ads_management,ads_read,business_management,catalog_management
- [ ] **T2.6** `bing`（oauth + developer_token）：scope=https://ads.microsoft.com/ads.manage
- [ ] **T2.7** `tiktok`（oauth）：scope=user.info.basic,ads.read,ads.management
- [ ] **T2.8** `llm`（api_key）：api_key、base_url、model（向下兼容 Provider）
- [ ] **T2.9** `smtp`（api_key）：host、port、user、pass（向下兼容 Provider）

## 阶段 3：OAuth 授权码流程

- [ ] **T3.1** `OAuth.Start(workspace_id, platform, scopes)`：生成 `state`（含 workspace_id、platform、name、nonce）→ 写 `oauth_state`（TTL 10 min）→ 返回 `authorize_url`
- [ ] **T3.2** `OAuth.Callback(code, state)`：校验 state → 调各 provider `token endpoint` → 写 `connection_token` → 创建 `connection` → 触发 PostAuthActions → 重定向到前端 `/settings/connections?connected=<id>`
- [ ] **T3.3** Google token endpoint：`https://oauth2.googleapis.com/token`（grant_type=authorization_code、refresh_token）
- [ ] **T3.4** Meta token endpoint：`https://graph.facebook.com/v25.0/oauth/access_token` + `/oauth/access_token?grant_type=fb_exchange_token`（换长 token）
- [ ] **T3.5** Microsoft token endpoint：`https://login.microsoftonline.com/common/oauth2/v2.0/token`
- [ ] **T3.6** TikTok token endpoint：`https://business-api.tiktok.com/open_api/v1.3/oauth2/token/`
- [ ] **T3.7** 路由：`GET /api/v1/oauth/callback` 走 `auth`+`tenant` 中间件（cors 允许 GET，Caddy 反代需确认无 rewrite）

## 阶段 4：资源挂接（PostAuthActions）

- [ ] **T4.1** Google Ads `ListAccessibleCustomers` → 派生 `ad_account` 候选列表（不直接 create，让用户选）
- [ ] **T4.2** Google Shopping `accounts.list` → 同上（标 `is_merchant_center`）
- [ ] **T4.3** Meta `/me/adaccounts` → 派生 `ad_account` 候选
- [ ] **T4.4** Meta `/me/businesses` → 派生 catalog、page
- [ ] **T4.5** Bing `GetAccountsInfo` → 派生 `ad_account` 候选
- [ ] **T4.6** TikTok `/advertiser/list/` → 派生 `ad_account` 候选
- [ ] **T4.7** WooCommerce `GET /wp-json/wc/v3/system_status` 校验凭据；`GET /wp-json/wc/v3/system_status` 返回 shop name
- [ ] **T4.8** Shopify `GET /admin/api/2024-10/shop.json` 校验凭据

## 阶段 5：Store / AdAccount 接受 conn_id

- [ ] **T5.1** `POST /api/v1/stores` 增加 `conn_id` 字段；不传时旧行为（base_url+api_key）
- [ ] **T5.2** `POST /api/v1/ad-accounts` 增加 `conn_id` 字段
- [ ] **T5.3** `GET /api/v1/stores?conn_id=X` / `GET /api/v1/ad-accounts?conn_id=X` 过滤
- [ ] **T5.4** `POST /api/v1/connections/{id}/sync/ad-accounts` 与 `.../sync/stores` 同步

## 阶段 6：Google Shopping 子资源（v1 范围：products/accounts/shippingsettings/returnpolicyonline 列表 + 写入）

- [ ] **T6.1** `platform/google_shopping/products.list`（GET `/products` 列出商家商品）
- [ ] **T6.2** `platform/google_shopping/shippingsettings.get/insert`（GET/INSERT `/accounts/{id}/shippingSettings`，按文档 etag 流程）
- [ ] **T6.3** `platform/google_shopping/returnpolicyonline.list/insert`（GET/POST `/accounts/{id}/returnpolicyonline`）
- [ ] **T6.4** `platform/google_shopping/localinventory`（可选，v1.1）
- [ ] **T6.5** 前端 `/products?platform=google_shopping` 增加"同步 Google 商家商品"按钮

## 阶段 7：Token 刷新 worker

- [ ] **T7.1** asynq 调度器注册 `token:refresh` 任务（5 min 周期）
- [ ] **T7.2** Worker：扫描 `expires_at < now+10min` 的 token → 按平台分发 refresh
- [ ] **T7.3** Google refresh：`https://oauth2.googleapis.com/token` (refresh_token grant)
- [ ] **T7.4** Microsoft refresh：标准 OAuth
- [ ] **T7.5** TikTok refresh：标准 OAuth
- [ ] **T7.6** 失败 3 次 → `connection.status='error'` + 写 `notification` 提醒重连

## 阶段 8：前端 `/settings/connections`

- [ ] **T8.1** 新建页面骨架：tabs（全部/店铺/广告平台/AI/邮件）、卡片网格、「+ 新建连接」按钮
- [ ] **T8.2** Stepper 模态：选 platform → 填 schema fields（参数型）/ 跳授权页（OAuth）→ 选 scope（OAuth）→ 完成
- [ ] **T8.3** OAuth 跳转完成后回 `/settings/connections?connected=<id>` 自动 toast
- [ ] **T8.4** 卡片操作：测试（POST /test）、编辑（重走 schema 步骤）、删除（带确认）
- [ ] **T8.5** `apps/web/lib/api/services/connection.service.ts`（list/get/create/update/delete/test/sync + platforms schema + oauth start）

## 阶段 9：前端 /accounts 和 /stores 引入

- [ ] **T9.1** `/accounts/page.tsx` 顶部「从已连接平台导入广告账户」按钮 → 弹选择 connection → 调 sync 列出候选 → 多选创建
- [ ] **T9.2** `/stores/page.tsx` 同上（从 Woo/Shopify 导入店铺）
- [ ] **T9.3** 表单 `conn_id` 可选 select（与手填 api_key 互斥）
- [ ] **T9.4** `workspace-guard` 与 stores/accounts 联调（确保 conn_id 鉴权正确）

## 阶段 10：兼容与迁移

- [ ] **T10.1** 旧 `Provider` 表（type IN llm/smtp）写一次性脚本迁到 `Connection`
- [ ] **T10.2** `/settings/providers` 加横幅「已迁移到 /settings/connections，3 个月后下线」
- [ ] **T10.3** 旧 Store/AdAccount 数据保留（conn_id=null）
- [ ] **T10.4** 数据库迁移脚本与回滚脚本

## 阶段 11：delta spec 与 Design Doc

- [ ] **T11.1** `openspec/changes/orbit-platform/specs/connection/spec.md`（Connection + 平台 schema 行为）
- [ ] **T11.2** `openspec/changes/orbit-platform/specs/platform-integration/spec.md`（OAuth 启动/回调/刷新行为）
- [ ] **T11.3** `openspec/changes/orbit-platform/specs/store/spec.md`（Store + conn_id 行为）
- [ ] **T11.4** `openspec/changes/orbit-platform/specs/ad-account/spec.md`（AdAccount + conn_id 行为）
- [ ] **T11.5** `docs/superpowers/specs/2026-09-03-connection-platform-design.md`（完整 RFC：数据模型 + OAuth 流程 + 平台注册 + 错误码）

## 阶段 12：验证与归档

- [ ] **T12.1** 端到端测试（Woo 走通：创建 connection → 同步 store → 列表商品）
- [ ] **T12.2** OAuth 回测（mock provider，验证 state 校验、code 换 token、刷新）
- [ ] **T12.3** 部署到 `portal.gusty.top`，验证 `/settings/connections` 可用
- [ ] **T12.4** 运行 `impetus-verify` 通过后 `/impetus-archive`

## 关键依赖与里程碑

- M1（T0-T1）：后端 Connection 模型 + CRUD + GET /platforms
- M2（T2-T3）：平台 schema + OAuth 启动/回调（不接真实 provider，先 mock）
- M3（T4-T5）：PostAuthActions + Store/AdAccount 接受 conn_id
- M4（T6-T7）：Google Shopping 资源 + token 刷新
- M5（T8-T9）：前端 /settings/connections + /accounts /stores 引入
- M6（T10-T12）：迁移、文档、验证、归档
```

