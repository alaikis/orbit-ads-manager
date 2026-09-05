# Platform Integration Spec

> 主规格 - 平台 SDK 适配层（按 platform 注册 client）

## Purpose

PlatformIntegration 定义每个外部平台的 HTTP 客户端契约。所有平台 client 实现统一接口 `PlatformClient`（test/sync/get_token_info/refresh_token）。

## Requirements

### Requirement: Platform registry

系统 **MUST** 维护 9 个平台的 `PlatformSchema` 注册表：

- `id` (platform key)
- `display_name` (i18n key, EN/ZH)
- `auth_flow` (api_key | developer_token | oauth)
- `oauth_provider` (google | microsoft | meta | tiktok | generic)
- `scopes` ([]string)
- `fields` ([]Field, 表单 schema)
- `post_auth_actions` ([]Action)

#### Scenario: 列出可用平台
- **WHEN** 前端 `GET /api/v1/connections/platforms`
- **THEN** 系统返回 9 个平台 schema 数组
- **AND** 前端按 schema 渲染 Stepper

### Requirement: Platform client interface

每个平台 **MUST** 实现：

```go
type PlatformClient interface {
    Test(ctx, conn) error
    Sync(ctx, conn) (resources, error)
    RefreshToken(ctx, conn) (token, error)
}
```

#### Scenario: WooCommerce test
- **WHEN** client.Test(ctx, conn)
- **THEN** Basic Auth 调 `GET /wp-json/wc/v3/system_status`
- **AND** 200 → nil；其他 → wrapped error

### Requirement: OAuth provider abstraction

OAuth 平台 **MUST** 通过 `OAuthProvider` 抽象处理 token endpoint、refresh endpoint 差异。

- Google: `https://oauth2.googleapis.com/token`
- Microsoft: `https://login.microsoftonline.com/common/oauth2/v2.0/token`
- Meta: `https://graph.facebook.com/v25.0/oauth/access_token` (短 token) + `grant_type=fb_exchange_token` (长 token)
- TikTok: `https://open.tiktokapis.com/v2/oauth/token/`
- Bing: `https://login.microsoftonline.com/common/oauth2/v2.0/token` (同 Microsoft)

#### Scenario: Google refresh
- **WHEN** `RefreshToken` 调用
- **THEN** POST `https://oauth2.googleapis.com/token` with `grant_type=refresh_token&refresh_token=...&client_id=...&client_secret=...`
- **AND** 解析响应，更新 expires_at

## File Structure

```
apps/api/internal/connection/
  registry.go   # 9 schema 注册 + PlatformClient interface
  clients.go    # 9 平台 client 实现
  oauth.go      # OAuthProvider abstraction
```

## Platform Catalog

见 `connection/spec.md` 第 5 节 Platform Catalog (9 平台)。
