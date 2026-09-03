-- Version: 000001
-- Description: Initial schema

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tenants (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    plan VARCHAR(20) NOT NULL DEFAULT 'beta',
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Shanghai',
    auto_approve_high_risk BOOLEAN NOT NULL DEFAULT FALSE,
    settings JSONB
);

CREATE TABLE IF NOT EXISTS tenant_members (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    user_id BIGINT NOT NULL,
    role VARCHAR(30) NOT NULL,
    invited_by BIGINT,
    invited_at TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    user_id BIGINT NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    user_agent VARCHAR(255),
    ip VARCHAR(45)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    actor_user_id BIGINT NOT NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id BIGINT,
    before_snapshot JSONB,
    after_snapshot JSONB,
    ip VARCHAR(45)
);

CREATE TABLE IF NOT EXISTS platform_conns (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    type VARCHAR(20) NOT NULL,
    platform VARCHAR(20) NOT NULL,
    admin_user_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'bound',
    meta JSONB,
    last_synced_at TIMESTAMP,
    store_id BIGINT
);

CREATE TABLE IF NOT EXISTS oauth_tokens (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    conn_id BIGINT NOT NULL UNIQUE,
    token_type VARCHAR(20) NOT NULL,
    encrypted_token1 VARCHAR(512) NOT NULL,
    encrypted_token2 VARCHAR(512),
    expires_at TIMESTAMP,
    scopes VARCHAR(500),
    last_refresh_at TIMESTAMP,
    refresh_error_count INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS store_configs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    store_id BIGINT NOT NULL UNIQUE,
    base_url VARCHAR(255) NOT NULL,
    extra_headers JSONB,
    sync_cursor JSONB,
    webhook_secret VARCHAR(64)
);

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    store_id BIGINT NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    link VARCHAR(500),
    image_url VARCHAR(500),
    brand VARCHAR(255),
    gtin VARCHAR(100),
    google_product_category VARCHAR(255),
    product_type VARCHAR(255),
    price_cents BIGINT NOT NULL,
    compare_at_price_cents BIGINT,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    hash VARCHAR(64),
    variants JSONB
);

CREATE TABLE IF NOT EXISTS product_variants (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    product_id BIGINT NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    sku VARCHAR(255),
    option_string VARCHAR(500),
    price_cents BIGINT NOT NULL,
    inventory_qty INTEGER NOT NULL DEFAULT 0,
    image_url VARCHAR(500)
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    store_id BIGINT NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    order_number VARCHAR(100),
    status VARCHAR(50) NOT NULL,
    total_cents BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    customer_email VARCHAR(255),
    items JSONB,
    placed_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS ad_accounts (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    conn_id BIGINT NOT NULL,
    platform VARCHAR(20) NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    currency VARCHAR(10),
    timezone VARCHAR(50),
    customer_id VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS campaigns (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    ad_account_id BIGINT NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    name VARCHAR(500) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'enabled',
    type VARCHAR(30),
    daily_budget_cents BIGINT,
    bidding_strategy VARCHAR(100),
    raw JSONB
);

CREATE TABLE IF NOT EXISTS ad_groups (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    ad_account_id BIGINT NOT NULL,
    campaign_id BIGINT NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    name VARCHAR(500) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'enabled',
    bid_micros BIGINT,
    targeting JSONB
);

CREATE TABLE IF NOT EXISTS ads (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    ad_account_id BIGINT NOT NULL,
    ad_group_id BIGINT NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    name VARCHAR(500) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'enabled',
    creative_type VARCHAR(50),
    headline_1 VARCHAR(255),
    headline_2 VARCHAR(255),
    description_1 VARCHAR(500),
    description_2 VARCHAR(500),
    raw JSONB
);

CREATE TABLE IF NOT EXISTS daily_stats (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    scope_type VARCHAR(20) NOT NULL,
    scope_id BIGINT NOT NULL,
    date VARCHAR(20) NOT NULL,
    platform VARCHAR(20) NOT NULL,
    metrics JSONB,
    source VARCHAR(20) NOT NULL DEFAULT 'synced'
);

CREATE TABLE IF NOT EXISTS sync_jobs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    type VARCHAR(50) NOT NULL,
    scope VARCHAR(50) NOT NULL,
    object_id BIGINT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    cursor_before JSONB,
    cursor_after JSONB,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    error_msg TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS feeds (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    store_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    format VARCHAR(50) NOT NULL,
    include_condition VARCHAR(100),
    currency_override VARCHAR(10),
    public_token VARCHAR(64) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_generated_at TIMESTAMP,
    generated_rows INTEGER,
    error_count INTEGER
);

CREATE TABLE IF NOT EXISTS feed_runs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    feed_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    rows INTEGER,
    errors JSONB,
    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS feed_product_snapshots (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    feed_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    snapshot_hash VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS conversations (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    user_id BIGINT NOT NULL,
    title VARCHAR(255),
    store_id BIGINT,
    account_id BIGINT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    llm_config JSONB
);

CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    conversation_id BIGINT NOT NULL,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS agent_tool_calls (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    message_id BIGINT NOT NULL,
    tool_name VARCHAR(100) NOT NULL,
    input JSONB,
    output_summary JSONB,
    risk_level VARCHAR(20) NOT NULL DEFAULT 'low',
    status VARCHAR(20) NOT NULL DEFAULT 'success',
    duration_ms BIGINT,
    confirm_action_id BIGINT
);

CREATE TABLE IF NOT EXISTS agent_actions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    action_type VARCHAR(50) NOT NULL,
    target_type VARCHAR(50),
    target_id BIGINT,
    params JSONB,
    risk_level VARCHAR(20) NOT NULL DEFAULT 'low',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    requested_by BIGINT,
    decided_by BIGINT,
    decided_at TIMESTAMP,
    result JSONB,
    expires_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rules (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    scope JSONB,
    condition_spec JSONB,
    action_spec JSONB,
    cooldown_minutes INTEGER NOT NULL DEFAULT 1440,
    created_by BIGINT NOT NULL,
    last_run_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rule_executions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    rule_id BIGINT NOT NULL,
    triggered_by VARCHAR(20) NOT NULL,
    matched_objects JSONB,
    executed_actions JSONB,
    status VARCHAR(20) NOT NULL,
    error_msg TEXT
);

CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    user_id BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT,
    link VARCHAR(500),
    read_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS email_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    to_email VARCHAR(255) NOT NULL,
    subject VARCHAR(500) NOT NULL,
    template VARCHAR(100),
    status VARCHAR(20) NOT NULL,
    error_msg TEXT
);

CREATE TABLE IF NOT EXISTS report_schedules (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL,
    scope_spec JSONB,
    recipients JSONB,
    template JSONB,
    cron_config VARCHAR(100),
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    last_sent_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    store_id BIGINT NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload_hash VARCHAR(64) NOT NULL,
    status VARCHAR(20) NOT NULL,
    processed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tenant_members_user ON tenant_members(user_id);
CREATE INDEX IF NOT EXISTS idx_tenant_members_status ON tenant_members(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_platform_conns_platform ON platform_conns(platform, status);
CREATE INDEX IF NOT EXISTS idx_products_store ON products(store_id, status);
CREATE INDEX IF NOT EXISTS idx_orders_store_placed ON orders(store_id, placed_at);
CREATE INDEX IF NOT EXISTS idx_daily_stats_date ON daily_stats(scope_type, scope_id, date);
CREATE INDEX IF NOT EXISTS idx_campaigns_account ON campaigns(ad_account_id);
CREATE INDEX IF NOT EXISTS idx_ad_groups_campaign ON ad_groups(campaign_id);
CREATE INDEX IF NOT EXISTS idx_ads_group ON ads(ad_group_id);
CREATE INDEX IF NOT EXISTS idx_feeds_store ON feeds(store_id, status);
CREATE INDEX IF NOT EXISTS idx_feed_runs_feed ON feed_runs(feed_id);
CREATE INDEX IF NOT EXISTS idx_conversations_user ON conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_agent_tool_calls_message ON agent_tool_calls(message_id);
CREATE INDEX IF NOT EXISTS idx_agent_actions_status ON agent_actions(status, expires_at);
CREATE INDEX IF NOT EXISTS idx_rules_enabled ON rules(enabled, last_run_at);
CREATE INDEX IF NOT EXISTS idx_rule_executions_rule ON rule_executions(rule_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, read_at);
CREATE INDEX IF NOT EXISTS idx_sync_jobs_status ON sync_jobs(status, type);
