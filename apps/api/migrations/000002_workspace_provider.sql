-- Version: 000002
-- Description: Add workspace and provider tables

CREATE TABLE IF NOT EXISTS workspaces (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_user_id BIGINT NOT NULL,
    plan VARCHAR(20) NOT NULL DEFAULT 'beta',
    settings JSONB,
    billing JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_workspaces_owner ON workspaces(owner_user_id);

ALTER TABLE tenants ADD COLUMN IF NOT EXISTS workspace_id BIGINT;
CREATE INDEX IF NOT EXISTS idx_tenants_workspace ON tenants(workspace_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE table_name = 'tenants' AND constraint_name = 'fk_tenants_workspace'
    ) THEN
        ALTER TABLE tenants ADD CONSTRAINT fk_tenants_workspace 
            FOREIGN KEY (workspace_id) REFERENCES workspaces(id);
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS providers (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    workspace_id BIGINT,
    tenant_id BIGINT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    config JSONB NOT NULL,
    meta JSONB,
    last_test_at TIMESTAMP,
    error_msg TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_providers_type_workspace ON providers(type, workspace_id);
CREATE INDEX IF NOT EXISTS idx_providers_tenant ON providers(tenant_id);

CREATE TABLE IF NOT EXISTS provider_usage (
    id BIGSERIAL PRIMARY KEY,
    provider_id BIGINT NOT NULL,
    workspace_id BIGINT NOT NULL,
    tenant_id BIGINT,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    tokens_used BIGINT NOT NULL DEFAULT 0,
    requests_count INT NOT NULL DEFAULT 0,
    errors_count INT NOT NULL DEFAULT 0,
    cost_cents BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_provider_usage_provider ON provider_usage(provider_id, period_start, period_end);
CREATE INDEX IF NOT EXISTS idx_provider_usage_workspace ON provider_usage(workspace_id, period_start, period_end);

ALTER TABLE platform_conns ADD COLUMN IF NOT EXISTS provider_id BIGINT;
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE table_name = 'platform_conns' AND constraint_name = 'fk_platform_conns_provider'
    ) THEN
        ALTER TABLE platform_conns ADD CONSTRAINT fk_platform_conns_provider 
            FOREIGN KEY (provider_id) REFERENCES providers(id);
    END IF;
END $$;
