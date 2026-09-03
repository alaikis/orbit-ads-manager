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
