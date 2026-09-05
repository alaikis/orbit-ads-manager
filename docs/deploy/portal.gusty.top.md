# 部署说明: portal.gusty.top

## OAuth 平台配置

**首次部署或更换 OAuth 应用时，需要在 `/apphub/ads.alakis.com/.env` 填入 6 平台 client_id/secret。**

获取地址（6 平台）：

| 平台 | 获取地址 | 说明 |
|------|----------|------|
| Google | https://console.cloud.google.com/apis/credentials | 创建 OAuth 2.0 Client ID，Authorized redirect URI: `https://ads.alaikis.com/api/v1/connections/oauth/callback` |
| Microsoft / Bing | https://portal.azure.com/#blade/Microsoft_AAD_RegisteredApps | 注册应用，redirect URI 同上。Microsoft 和 Bing 可共用一个 app |
| Meta | https://developers.facebook.com/apps/ | 创建应用，添加 Marketing API 产品，配置 OAuth redirect |
| TikTok | https://ads.tiktok.com/marketing_api/homepage | 申请 Marketing API 权限 |
| Shopify | https://shopify.dev/apps | 创建 public app |

**填入位置**（复制 `.env.example` 模板）：
```bash
ssh root@portal.gusty.top
cp /apphub/ads.alakis.com/apps/api/.env.example /apphub/ads.alakis.com/.env
nano /apphub/ads.alakis.com/.env   # 填入 OAUTH_*_CLIENT_ID/SECRET
systemctl restart orbit
```

**加载顺序**：
1. DB `platform_providers` 表（按 tenant_id 优先，全局 NULL 兜底）— 高级用法：每个租户独立 OAuth app
2. 环境变量 `OAUTH_*` — 简单用法：所有租户共享同一 OAuth app
3. 返回 401 "OAuth client not configured" — 提示运维补配置

## 占位状态

截至 2026-09-04，`.env` 中 6 平台 OAUTH_* 字段为**空**，OAuth 流程返回：
```
401 {"code":1001,"message":"OAuth client not configured"}
```
这是**预期**行为 — 代码已就绪，等待真实凭据。
