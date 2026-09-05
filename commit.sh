cd /e/develop-space/ads.alaikis.com
git -c user.name=Kilo -c user.email=kilo@mini.local commit -m "feat(orbit-platform): 统一凭据管理与按平台授权体系

后端 (apps/api)
- 000003_connection_platform.sql: 新增 connection / connection_credential /
  connection_token / oauth_state 四张表，platform_conns 加 conn_id 列
- pkg/crypto/aesgcm.go: AES-256-GCM 包装，复用 ENCRYPTION_MASTER_KEY
- internal/connection: 完整模块（models / registry / service / handlers /
  oauth / clients / refresh）支持 9 平台 schema：
  * api_key: woocommerce, shopify, llm, smtp
  * developer_token: bing
  * oauth: google_ads, google_shopping, meta, tiktok
- 平台客户端实装: Woo test/sync (system_status)、Shopify test
  (shop.json)、Google Shopping sync (merchantapi accounts/products)、
  Meta adaccounts、Meta/Bing/TikTok 占位
- Token 刷新: RefreshExpiringTokens (cron */10 * * * *) 支持 Google /
  Microsoft / TikTok refresh_token grant，失败 3 次写 last_error
- cmd/migrate-legacy: 旧 Provider (llm/smtp) 一次性迁移工具（已部署到
  /apphub/ads.alakis.com/bin/orbit-migrate-legacy）
- model.PlatformConn / model.AdAccount: 接受 conn_id（nullable）

前端 (apps/web)
- app/(dashboard)/settings/connections/page.tsx: 完整页面
  * tabs: 全部 / 店铺 / 广告平台 / AI
  * 卡片网格 + 「+ 新建连接」Stepper 模态
  * 选 platform → 填 schema fields (api_key/dev_token) / 跳授权 (OAuth)
  * 测试/同步/删除操作
  * OAuth 跳转回 ?connected=<id> 自动 toast
- app/(dashboard)/settings/providers/page.tsx: 横幅「已迁移到
  /settings/connections，3 个月后下线」
- lib/api/services/connection.service.ts: TS 客户端

Impetus artifacts
- proposal.md / design.md / tasks.md (12 阶段 47 任务)
- .impetus.yaml: workflow=full, phase=build → verify, build_mode=executing-plans
- docs/superpowers/specs/2026-09-03-connection-platform-design.md
- docs/superpowers/plans/2026-09-03-connection-platform.md
- handoff/spec-context.json + .md (beta mode, hash f18b035a...)" 2>&1
echo "--- log ---"
git log --oneline -5
