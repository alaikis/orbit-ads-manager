-- Version: 000004
-- Description: OAuth platform provider configurations (per-tenant or global default)
-- Lookup priority: 1) platform_providers WHERE tenant_id = ? AND platform = ?
--                 2) platform_providers WHERE tenant_id IS NULL AND platform = ?
--                 3) environment variable OAUTH_<PLATFORM>_CLIENT_ID/SECRET
--                 4) return ErrOAuthNotConfigured

CREATE TABLE IF NOT EXISTS platform_providers (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT,                          -- NULL = global default
    platform VARCHAR(32) NOT NULL,             -- google / microsoft / meta / tiktok / bing / shopify
    client_id VARCHAR(255) NOT NULL,
    client_secret_enc BYTEA,                   -- AES-256-GCM ciphertext (separate column from nonce)
    client_secret_nonce BYTEA,                 -- 12-byte nonce
    redirect_uri VARCHAR(512) NOT NULL,
    scopes JSONB,                              -- ["scope1", "scope2"]
    is_active SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- tenant-scoped row wins over global default; COALESCE makes NULL distinct from 0
CREATE UNIQUE INDEX IF NOT EXISTS uk_platform_providers_tenant_platform
    ON platform_providers(COALESCE(tenant_id, 0), platform);

CREATE INDEX IF NOT EXISTS idx_platform_providers_tenant
    ON platform_providers(tenant_id);

CREATE INDEX IF NOT EXISTS idx_platform_providers_platform
    ON platform_providers(platform, is_active);

COMMENT ON TABLE platform_providers IS
    'OAuth provider credentials. tenant_id NULL = global default. '
    'Lookup order: row with tenant_id=X → row with tenant_id=NULL → env var.';
