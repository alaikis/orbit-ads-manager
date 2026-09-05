# Tenant Spec

> 主规格 - 租户隔离（项目无 workspace context，按 tenant 划分资源）

## Purpose

Tenant 是项目最高层级数据隔离单位。所有业务表都包含 `tenant_id` 字段，由 JWT 中间件注入到请求 context。

## Requirements

### Requirement: Tenant isolation

系统 **MUST** 强制所有业务查询包含 `tenant_id` 过滤。

- JWT middleware 解析 token，注入 `tenant_id` 到 `gin.Context`
- 所有 handler **MUST** 从 context 取 `tenant_id`，不信任 client input
- DB layer 自动 WHERE `tenant_id = ?`（repository pattern）
- 越权返回 404（不泄漏资源存在性）

#### Scenario: 跨租户访问
- **WHEN** tenant A 用户用 tenant B 的 resource id
- **THEN** 系统返回 404

### Requirement: Tenant roles

Tenant 内部不区分角色（owner/admin 同一权限），简化模型。

- 旧 plan 中提到的 owner/admin 角色 **DEPRECATED**
- 鉴权由 JWT 验证 + tenant_id 注入完成

#### Scenario: 旧角色字段
- **WHEN** 旧数据存在 `role` 字段
- **THEN** migrate-legacy 工具 **MAY** 移除该字段（不强制）

## File Structure

```
apps/api/internal/middleware/
  auth.go    # JWT 解析 + tenant_id 注入

apps/api/internal/tenant/
  handlers/  # tenant 范围资源 handler
```

## Schema

```sql
CREATE TABLE users (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  tenant_id BIGINT UNSIGNED NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  KEY idx_tenant (tenant_id)
);
```
