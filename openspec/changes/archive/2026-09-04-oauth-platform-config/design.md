# Design: OAuth 平台配置表

## 加载优先级

1. DB `platform_providers` 表（按 `tenant_id + platform_key` 查）
2. 环境变量 `OAUTH_<PLATFORM>_CLIENT_ID` / `OAUTH_<PLATFORM>_CLIENT_SECRET`
3. nil（返回 401 "OAuth client not configured"）

## 表 schema

```sql
CREATE TABLE platform_providers (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  tenant_id BIGINT UNSIGNED NULL,         -- NULL 表示全局默认
  platform VARCHAR(32) NOT NULL,          -- google / microsoft / meta / tiktok / bing / shopify
  client_id VARCHAR(255) NOT NULL,
  client_secret_enc VARBINARY(512) NULL,  -- AES-256-GCM 加密（与 connection_token 共享 key）
  redirect_uri VARCHAR(512) NOT NULL,
  scopes JSON NULL,                       -- ["scope1", "scope2"]
  is_active TINYINT(1) DEFAULT 1,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_tenant_platform (tenant_id, platform)
);
```

## 加载逻辑

```go
func (s *Service) GetOAuthProvider(tenantID int64, platform string) (*OAuthProvider, error) {
    // 1. DB
    if p, err := s.repo.GetProvider(tenantID, platform); err == nil {
        return p, nil
    }
    // 2. Env fallback
    clientID := os.Getenv("OAUTH_" + strings.ToUpper(platform) + "_CLIENT_ID")
    clientSecret := os.Getenv("OAUTH_" + strings.ToUpper(platform) + "_CLIENT_SECRET")
    if clientID == "" || clientSecret == "" {
        return nil, errors.New("OAuth client not configured for platform " + platform)
    }
    return &OAuthProvider{ClientID: clientID, ClientSecret: clientSecret, RedirectURI: defaultRedirect}, nil
}
```

## 迁移

- 文件: `apps/api/migrations/000004_oauth_platform_providers.sql`
- 启动时自动执行（项目约定）

## 占位配置

`.env.example` 模板：

```bash
# Google OAuth (https://console.cloud.google.com/apis/credentials)
OAUTH_GOOGLE_CLIENT_ID=
OAUTH_GOOGLE_CLIENT_SECRET=

# Microsoft (https://portal.azure.com/#blade/Microsoft_AAD_RegisteredApps)
OAUTH_MICROSOFT_CLIENT_ID=
OAUTH_MICROSOFT_CLIENT_SECRET=

# Meta (https://developers.facebook.com/apps/)
OAUTH_META_CLIENT_ID=
OAUTH_META_CLIENT_SECRET=

# TikTok (https://ads.tiktok.com/marketing_api/homepage)
OAUTH_TIKTOK_CLIENT_ID=
OAUTH_TIKTOK_CLIENT_SECRET=

# Bing (同 Microsoft, 用同一 app)
OAUTH_BING_CLIENT_ID=
OAUTH_BING_CLIENT_SECRET=

# Shopify (https://shopify.dev/apps)
OAUTH_SHOPIFY_CLIENT_ID=
OAUTH_SHOPIFY_CLIENT_SECRET=
```

## 修改文件

1. `apps/api/migrations/000004_oauth_platform_providers.sql` (新)
2. `apps/api/internal/connection/oauth.go` (改 GetProvider 逻辑)
3. `apps/api/.env.example` (新)
4. `docs/deploy/portal.gusty.top.md` (新，部署说明)

总计 4 个文件，符合 tweak 范围。
