# Tasks: oauth-platform-config

- [x] T1. 创建迁移 SQL `000004_oauth_platform_providers.sql` (1 张表)
- [x] T2. 修改 `apps/api/internal/connection/oauth.go` 添加 `LoadOAuthProvider(tenantID, platform)` 函数
- [x] T3. 创建 `apps/api/.env.example` 含 6 平台占位配置
- [x] T4. 部署到 portal.gusty.top 并 live-verify (DB + env 双路径均 OK)
