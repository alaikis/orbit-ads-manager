# Store Spec

> 主规格 - 店铺实体（多平台店铺/店铺型连接目标）

## Purpose

Store 代表租户在某个平台上拥有的店铺（WooCommerce shop, Shopify store, Google Shopping account）。支持按 `conn_id` 关联到 connection 实体，实现"凭据 → 店铺"的可追溯链路。

## Requirements

### Requirement: Store entity

系统 **MUST** 维护 `store` 表：

- `id` (PK), `tenant_id`, `platform` (枚举), `external_id` (平台侧 ID), `name`, `currency`, `timezone`, `status` (active/disabled), `last_synced_at`
- 可选 `conn_id` (nullable, FK → connection.id)
- 唯一索引 `(tenant_id, platform, external_id)`

#### Scenario: 同步 Shopify 店铺
- **WHEN** 用户触发 `POST /api/v1/stores/sync?conn_id=X`
- **THEN** 系统从 conn_id 关联的 connection 拉取 Shopify shop 信息
- **AND** 写入 store 表，状态 active，更新 `last_synced_at`

### Requirement: Legacy mode (without conn_id)

为兼容旧版直接 `tenant_id + provider` 模式，store 允许 `conn_id IS NULL`。

- 旧 Provider 路径 (llm/smtp 之外) **MUST** 通过 migrate-legacy 工具迁移到 conn_id 关联模式
- 旧路径查询 **MUST** 仍可用（带 `?legacy=true` 或自动检测）

#### Scenario: 旧 store 迁移
- **WHEN** migrate-legacy 工具运行
- **THEN** 旧 Provider 的 store/ad_account 创建对应 connection，迁移 `conn_id`
- **AND** 旧表保留 3 个月观察期后下线

## File Structure

```
apps/api/internal/tenant/handlers/stores.go    # store CRUD + sync
apps/api/internal/model/models.go              # model.Store / model.AdAccount
```

## Schema

```sql
CREATE TABLE stores (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  tenant_id BIGINT UNSIGNED NOT NULL,
  platform VARCHAR(32) NOT NULL,
  external_id VARCHAR(128) NOT NULL,
  name VARCHAR(255),
  currency CHAR(3),
  timezone VARCHAR(64),
  status VARCHAR(16) DEFAULT 'active',
  conn_id BIGINT UNSIGNED NULL,        -- NEW: 关联 connection
  last_synced_at DATETIME NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_tenant_platform_ext (tenant_id, platform, external_id),
  KEY idx_conn (conn_id),
  KEY idx_tenant (tenant_id)
);
```
