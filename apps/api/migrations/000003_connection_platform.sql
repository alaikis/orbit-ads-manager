-- Version: 000003
-- Description: Connection platform unified credentials (connection, connection_credential, connection_token, oauth_state)

CREATE TABLE IF NOT EXISTS connections (
    id BIGSERIAL PRIMARY KEY,
    workspace_id BIGINT,
    tenant_id BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL,
    platform VARCHAR(30) NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    scopes JSONB,
    last_refreshed_at TIMESTAMP,
    last_error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_connections_workspace ON connections(workspace_id);
CREATE INDEX IF NOT EXISTS idx_connections_tenant ON connections(tenant_id);
CREATE INDEX IF NOT EXISTS idx_connections_platform ON connections(platform, status);
CREATE INDEX IF NOT EXISTS idx_connections_status ON connections(status);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'connections' AND constraint_name = 'fk_connections_workspace'
    ) THEN
        ALTER TABLE connections ADD CONSTRAINT fk_connections_workspace
            FOREIGN KEY (workspace_id) REFERENCES workspaces(id);
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS connection_credentials (
    id BIGSERIAL PRIMARY KEY,
    conn_id BIGINT NOT NULL,
    key VARCHAR(64) NOT NULL,
    encrypted_value BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    last_4 VARCHAR(8),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (conn_id, key)
);

CREATE INDEX IF NOT EXISTS idx_connection_credentials_conn ON connection_credentials(conn_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'connection_credentials' AND constraint_name = 'fk_connection_credentials_conn'
    ) THEN
        ALTER TABLE connection_credentials ADD CONSTRAINT fk_connection_credentials_conn
            FOREIGN KEY (conn_id) REFERENCES connections(id) ON DELETE CASCADE;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS connection_tokens (
    id BIGSERIAL PRIMARY KEY,
    conn_id BIGINT NOT NULL UNIQUE,
    access_token_enc BYTEA NOT NULL,
    access_token_nonce BYTEA NOT NULL,
    refresh_token_enc BYTEA,
    refresh_token_nonce BYTEA,
    token_type VARCHAR(20) NOT NULL DEFAULT 'Bearer',
    expires_at TIMESTAMP,
    scopes JSONB,
    system_user_id VARCHAR(64),
    last_refresh_at TIMESTAMP,
    refresh_error_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_connection_tokens_expires ON connection_tokens(expires_at);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'connection_tokens' AND constraint_name = 'fk_connection_tokens_conn'
    ) THEN
        ALTER TABLE connection_tokens ADD CONSTRAINT fk_connection_tokens_conn
            FOREIGN KEY (conn_id) REFERENCES connections(id) ON DELETE CASCADE;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS oauth_states (
    id BIGSERIAL PRIMARY KEY,
    state VARCHAR(64) NOT NULL UNIQUE,
    platform VARCHAR(30) NOT NULL,
    workspace_id BIGINT,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    draft_fields JSONB,
    scopes JSONB,
    redirect_to VARCHAR(500),
    used BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_oauth_states_expires ON oauth_states(expires_at);

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'platform_conns') THEN
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_name = 'platform_conns' AND column_name = 'conn_id'
        ) THEN
            ALTER TABLE platform_conns ADD COLUMN conn_id BIGINT;
        END IF;
        CREATE INDEX IF NOT EXISTS idx_platform_conns_conn ON platform_conns(conn_id);
    END IF;
END $$;
