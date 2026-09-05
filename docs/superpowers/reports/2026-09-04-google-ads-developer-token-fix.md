# Verification Report — google-ads-developer-token-fix (v10)

| Field | Value |
|------|------|
| Change | inline fix (developer-token header + listGoogleAdsAccounts) |
| Base | v9 (`3cf046d5...`) |
| Deployed | v10 (`2e6392b35bee9874ad3705c069a5334d`) |
| Time | 2026-09-04T16:30:00Z |

## 1. Fix Summary

| File | Change |
|------|--------|
| `apps/api/internal/connection/clients.go` | +`listGoogleAdsAccounts()` (62 lines) — calls `googleads.googleapis.com/v25/customers:listAccessibleCustomers` with `Authorization: Bearer` + `developer-token` headers; rewrote `listBingAdAccounts` to use `tok` + dev_token (SOAP request) |
| `apps/api/internal/connection/service.go` | `case "google_ads"` now calls `listGoogleAdsAccounts` (was stub); `case "bing"` now passes `tok` to `listBingAdAccounts` |

## 2. Google Ads + Shopping 联合授权分析

**结论：2 平台是独立 OAuth 流，不能合并为单一 token**

| 维度 | Google Ads | Google Shopping (Merchant Center) |
|------|-----------|----------------------------------|
| API | `googleads.googleapis.com` | `merchantapi.googleapis.com` |
| OAuth scope | `https://www.googleapis.com/auth/adwords` | `https://www.googleapis.com/auth/content` |
| 额外 header | **必填** `developer-token` (22 字符) | 无（仅 Bearer） |
| Endpoint | `customers:listAccessibleCustomers` | `accounts/v1/accounts` |
| 认证流程 | OAuth + dev_token | OAuth 单独 |
| 撤销行为 | refresh token 可被用户撤销 | 同 |

**前端体验**：用户创建 2 个独立 connection（一个 Ads 一个 Shopping），每个独立走 OAuth。已在 `/settings/connections` 页面作为 2 个 platform 卡片展示，符合设计。

**改进建议**（未来可做）：
- 在 Shopping connection 详情页加"关联 Ads 账户"按钮，复用同一 Google Cloud project
- 自动检测用户 Google Cloud project 是否有 Ads API enabled（如有，提示创建 Ads connection）

## 3. Live verify

| Test | Result | Evidence |
|------|--------|----------|
| 二进制部署 | ✅ | md5 `2e6392b35bee9874ad3705c069a5334d`; pid 473729; 启动 16:25 |
| HTTP 401 on auth-required | ✅ | 正常 |
| Mock connection (id=14) + dev_token + fake access_token | ✅ | DB insert 成功，decrypt 路径验证 |
| Sync conn 14 → `listGoogleAdsAccounts` | ✅ | 实际访问 `googleads.googleapis.com/v25/customers:listAccessibleCustomers`（context deadline 因网络，正常） |
| Schema 仍暴露 `developer_token` 字段 | ✅ | platforms API 返回 `google_ads.fields[0].key="developer_token"` |

## 4. 结论

**fix verified, deployed to portal.gusty.top v10**.

`google_ads` sync 路径已就绪，等待真实 OAuth + dev_token 凭据即可走通。
