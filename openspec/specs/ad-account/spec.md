# Ad Account Spec

> 主规格 - 广告账户实体（Google Ads MCC, Meta Ad Account, Bing, TikTok 等）

## Purpose

AdAccount 代表租户在广告平台上拥有的广告投放账户。Google Ads MCC 包含多个子账户，每个 AdAccount 行对应一个。

## Requirements

### Requirement: Ad Account entity

系统 **MUST** 维护 `ad_accounts` 表：

- `id` (PK), `tenant_id`, `platform` (google_ads, meta, bing, tiktok), `external_id` (平台账户 ID), `name`, `currency`, `timezone`, `status`
- 可选 `conn_id` (nullable, FK → connection.id)
- 父账户 `parent_id` (MCC 子账户)
- 唯一索引 `(tenant_id, platform, external_id)`

#### Scenario: 同步 Google Ads 子账户
- **WHEN** Google Ads OAuth 完成 + `list_ad_accounts` PostAuthAction
- **THEN** 系统拉取 MCC 下所有子账户（含自身）
- **AND** 写入 ad_accounts，绑定 `conn_id`

#### Scenario: Meta 多账户
- **WHEN** Meta Business 用户授权
- **THEN** `list_ad_accounts` 拉取所有 accessible adaccounts
- **AND** 每个账户写一行 ad_accounts

### Requirement: MCC hierarchy

Google Ads / Meta 等支持账户层级（MCC → 子账户）。

- 父账户 `parent_id` 指向自己 connection 的 root 账户
- 子账户可独立启用/禁用

## File Structure

```
apps/api/internal/adplatform/handlers/accounts.go  # ad_accounts CRUD + sync
```

## Schema

```sql
CREATE TABLE ad_accounts (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  tenant_id BIGINT UNSIGNED NOT NULL,
  platform VARCHAR(32) NOT NULL,
  external_id VARCHAR(128) NOT NULL,
  parent_id BIGINT UNSIGNED NULL,
  name VARCHAR(255),
  currency CHAR(3),
  timezone VARCHAR(64),
  status VARCHAR(16) DEFAULT 'active',
  conn_id BIGINT UNSIGNED NULL,        -- NEW
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_tenant_platform_ext (tenant_id, platform, external_id),
  KEY idx_parent (parent_id),
  KEY idx_conn (conn_id),
  KEY idx_tenant (tenant_id)
);
```
