# Brainstorm Summary — orbit-platform

> 恢复检查点，非最终 Design Doc。最终 Design Doc 见 `docs/superpowers/specs/2026-09-03-connection-platform-design.md`。
> Change: orbit-platform | Phase: design | handoff_hash: f18b035a1ecd6aa1a5766ba36b227fe41630f8eb0653e59702492aa303f98bd2

## 已确认方案

### 实体模型
```
workspace (1) ──< connection (1) ──< connection_credential (n)
                                  ──< connection_token (n, 0..1)
                                  ──< oauth_state (n)
        connection (1) ──< store (n)
        connection (1) ──< ad_account (n)
```

### 三种 auth flow
| type | 例子 | 前端动作 | 后端动作 |
|------|------|---------|---------|
| api_key | Woo, Shopify | 直接表单提交 consumer_key/secret | 加密入 connection_credential |
| developer_token | Bing dev_token | 表单提交 client_id + secret + dev_token | 加密入 connection_credential |
| oauth | Google Ads/Shopping, Meta, Bing, TikTok | 授权按钮 → 跳 authorize URL → callback | 生成 state → code 换 token → 加密 → 创建 connection |

### OAuth 回调地址
`https://ads.alaikis.com/api/v1/oauth/callback`（Caddy 必须放过 GET）

### Token 加密
复用 `config.AppCfg.Encryption.MasterKey`（AES-256-GCM），新增 `pkg/crypto/aesgcm.go`

### 平台 schema 来源
**候选 C（JSON config 文件 + DB 缓存）** — 可版本控制，启动时 seed 到 `platform_schema` 表

### 优先范围
Google Shopping > Meta > WooCommerce > Google Ads > Bing > TikTok > Shopify

### 降级策略
`/settings/providers` 顶部加横幅「已迁移到 /settings/connections，3 个月后下线」

## 关键取舍

| 决策 | 选择 | 理由 |
|------|------|------|
| Connection 与 Store/AdAccount 关系 | 候选 B（弱耦合，conn_id 可选） | 向后兼容，不破坏现有数据 |
| WooCommerce 鉴权 | Basic Auth header | 比 query string 更安全 |
| Meta 两种 token 支持 | 先做 User access token，System User v2 | 减少 scope，简化首次迭代 |
| Google Shopping vs Google Ads | 视为两个独立 platform | 各自 scope 不同，互不干扰 |
| shippingSettings/returnpolicyonline | v1 只做 OAuth + products.list | 范围太大，分 v1/v1.1 |

## 待 Spec Patch（delta spec 补充）

- `specs/connection/spec.md`：Connection CRUD + Test + Sync 接口语义
- `specs/platform-integration/spec.md`：OAuth 启动/回调/刷新/失败告警
- `specs/store/spec.md`：conn_id 字段 + 从 Connection 同步
- `specs/ad-account/spec.md`：conn_id 字段 + PostAuth 派生候选列表

## 测试策略

- 后端单元：Connection Service（CRUD + Test + Sync）
- 后端集成：OAuth mock provider（Google/Meta/Bing/TikTok 的 token exchange）
- E2E：Woo 走通 → create connection → sync stores → list products
- 前端：/settings/connections Stepper 模态（参数型表单 + OAuth 跳转）
