# Proposal: OAuth 平台配置表 (oauth-platform-config)

## 动机

当前 OAuth 平台（Google/Microsoft/Meta/TikTok/Bing）的 `client_id` / `client_secret` 在代码中硬编码或在 `oauth.go` 中按 platform 字符串分支加载。当需要为不同环境（dev/staging/prod）使用不同的 OAuth 应用，或需要支持租户自定义 OAuth 凭据时，硬编码不可持续。

## 目标

1. 新增 `platform_providers` 表存储 OAuth 客户端配置（client_id, client_secret, redirect_uri, scopes）
2. `oauth.go` 通过 `platform` + `tenant_id` 查找 provider 配置（多租户可覆盖）
3. 提供 6 个平台（Google/Microsoft/Meta/TikTok/Bing/Shopify）的占位配置（env vars 形式）
4. env var 格式 `OAUTH_<PLATFORM>_CLIENT_ID` / `OAUTH_<PLATFORM>_CLIENT_SECRET`

## 范围

- 新增 1 张表 + 1 个 Go 工具函数 + 1 个 .env.example 模板
- 修改 `oauth.go` 加载逻辑（按 platform 找 provider，否则降级到 env）
- 6 平台占位 `client_id` / `client_secret` 留空（待用户在 .env 填入真实值）

## 不在范围

- 多租户 OAuth 自定义（每个租户独立 OAuth app）— 留作未来 phase
- Token 加密存储（沿用现有 AES-GCM 流程）
- UI 界面（通过 .env 文件管理）
