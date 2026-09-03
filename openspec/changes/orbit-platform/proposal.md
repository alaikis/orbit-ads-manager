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
