---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: ccdb5f1afb85a6c9dc58ee10b06ac953_adcbddd5a1c411f1a413525400287e28
    ReservedCode1: Ept8DSkQok47wm1kvJ7nwOcIter0n+Dr1d9wvY/SqA91YDkB5aqJBxsaylnLeZxZAfeQlgJGiMJLPBYbYm7RAx3yn6AYAX1eQqa98EO1zyOjXUmXied1yBHnaHqUVajDUXQUHAwUH6QgEv6JLU5iGxDcWsSSyMm7ee6KUkQDPigjKhx4nXt2fYAgwIw=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: ccdb5f1afb85a6c9dc58ee10b06ac953_adcbddd5a1c411f1a413525400287e28
    ReservedCode2: Ept8DSkQok47wm1kvJ7nwOcIter0n+Dr1d9wvY/SqA91YDkB5aqJBxsaylnLeZxZAfeQlgJGiMJLPBYbYm7RAx3yn6AYAX1eQqa98EO1zyOjXUmXied1yBHnaHqUVajDUXQUHAwUH6QgEv6JLU5iGxDcWsSSyMm7ee6KUkQDPigjKhx4nXt2fYAgwIw=
---

# Orbit 广告智能助手系统
## 开发需求文档（V2.0 · 可执行版）

**版本**：V2.0
**上一版本**：V1.1
**日期**：2026年8月27日
**主域名**：orbit.elstella.com
**API 域名**：api.orbit.elstella.com
**文档定位**：本文档在 V1.1 基础上逐条补齐落地细节，目标读者为开发团队（前端 / 后端 / 运维 / 测试），达到"拿起来即可开工"的可执行标准。

---

## 0. V1.1 评审摘要与本版修订说明

### 0.1 V1.1 主要问题清单

| 编号 | 类别 | 问题描述 | 本版处理 |
|------|------|----------|----------|
| R01 | 不一致 | 标题称"前后端一体"，但 1.4 技术原则为"前后端分离"，自相矛盾 | 全文统一为"前后端分离 + Monorepo" |
| R02 | 不一致 | 1.2 产品定位将"Google Shopping"与 Google Ads / Meta / Bing Ads 并列，但 Shopping 是 Google Ads 的广告类型而非独立平台 | 修正定位描述，Shopping 归入 Google Ads 能力 |
| R03 | 技术选型未定 | 后端框架"Gin/Fiber"未定，团队成员无法开工 | 明确：Golang + Gin（生态成熟、中间件丰富、团队经验匹配） |
| R04 | 缺失 | 无数据库表结构、无字段级定义，无法建模 | 新增 5.3 数据库设计（29 张核心表） |
| R05 | 缺失 | 无 API 接口清单、无统一响应与错误码规范，前后端无法联调 | 新增 5.4（API 清单）、5.5（响应与错误码） |
| R06 | 缺失 | Agent 仅列出职责，无 Tool 定义、无编排流程、无人工确认协议 | 新增 5.6（Agent 与 Tool 定义）、3.5 细化 |
| R07 | 缺失 | OAuth 仅一句"一键授权"，未涉及回调、Token 存储、刷新、失败恢复 | 新增 5.7（OAuth 与 Token 生命周期） |
| R08 | 缺失 | 数据同步无任务设计（触发 / 频率 / 游标 / 幂等 / 失败处理） | 新增 5.8（数据同步任务设计） |
| R09 | 缺失 | Feed 无格式 / 字段映射 / 更新策略规范 | 新增 5.9（Feed 生成规范） |
| R10 | 缺失 | 规则引擎无条件结构 / 动作定义 / 触发与防抖机制 | 新增 5.10（规则引擎设计） |
| R11 | 缺失 | 无多租户隔离实现方案（仅一句"严格隔离"） | 新增 2.4（多租户隔离实现方案） |
| R12 | 缺失 | 安全仅列要点，无 JWT 策略、密钥管理、加密算法、限流、审计字段 | 新增 5.11（安全设计） |
| R13 | 缺失 | 无部署与配置清单（镜像、环境变量、依赖服务） | 新增 5.12（部署与配置清单） |
| R14 | 模糊 | 广告指标（ROAS / CPA / CTR）无统一计算口径 | 新增 0.4（指标口径定义），全文统一 |
| R15 | 模糊 | 一期 P0/P1 平台边界不清（OpenCart / Shopyy / UeeShop 一期是否做？） | 明确：一期仅 WooCommerce + Shopify（P0），P1 平台列入二期 |
| R16 | 不可执行 | 7 章交付标准为自然语言，无法验收 | 细化为可勾选验收清单（见第 7 章） |
| R17 | 不可执行 | 8 章节奏无里程碑、无验收标准、无任务拆分 | 细化为周计划 + 里程碑验收（见第 8 章） |
| R18 | 模糊 | 角色仅有名称，无权限点定义 | 新增 3.1 权限矩阵（RBAC） |
| R19 | 缺失 | 无通知与日志的具体载体 / 字段 / 保留策略 | 3.8 与 5.11.5 补齐 |
| R20 | 缺失 | 无测试策略、无 CI/CD、无可观测性设计 | 新增 5.13（测试与质量）及 6 章补充 |

### 0.2 V2.0 新增/重构章节索引

- 新增：0（评审与口径）、2.4（多租户隔离实现）、5.3（数据库设计）、5.4（API 清单）、5.5（响应与错误码）、5.6（Agent 与 Tool）、5.7（OAuth）、5.8（数据同步任务）、5.9（Feed 规范）、5.10（规则引擎）、5.11（安全）、5.12（部署与配置）、5.13（测试与质量）
- 重构：3.2 / 3.3 / 3.4 / 4.1 / 6 / 7 / 8

### 0.3 范围与边界（明确一期范围）

一期（可上线）范围内：

1. 平台授权：Google Ads、Meta（Facebook/Instagram）、Microsoft Advertising（Bing Ads）
2. 电商对接：WooCommerce、Shopify（P0 必须），OpenCart、Shopyy、UeeShop（P1，列入二期）
3. 广告管理：Campaign / Ad Group / Ad 统一视图、跨平台聚合报表
4. Agent 体系：Sync / Report / Optimize / Rule / Creative / Ops / Notify 七个 Agent 全部进入一期，但能力以 5.6 定义的 Tool 清单为上限
5. 自动化规则：预算预警、ROAS/CPA 阈值、平滑消耗辅助

一期明确不在范围内：

1. 多语言界面：一期仅中文，英文为二期（3.8 原有表述已降级）
2. 移动端 App：仅 Web 响应式
3. 客户端广告投放（自动创建广告素材并投送）：Creative Agent 一期只出"文案与素材建议"，不直接创建广告
4. 跨币种自动结算 / 多币种报表：一期按"账户原始币种 + 租户展示币种汇率换算"两条线展示
5. 全托管投放（任 Agent 无确认自动执行大批量操作）：默认人工确认，自动执行仅在租户级开关开启后生效

### 0.4 指标口径定义（全文统一）

| 指标 | 缩写 | 计算口径 |
|------|------|----------|
| 广告花费 | spend | 期内各平台上报花费之和（Google：cost_micros/1e6；Meta：spend；Bing：Spend） |
| 展示 | impressions | 期内展示次数之和 |
| 点击 | clicks | 期内点击次数之和 |
| 点击率 | CTR | clicks ÷ impressions ×100% |
| 平均点击成本 | CPC | spend ÷ clicks |
| 转化数 | conversions | 期内报告转化数（各平台口径），不跨平台合并（因归因模型不同），跨平台报表用"转化行数"并列展示 | 
| 转化成本 | CPA | spend ÷ conversions |
| 广告支出回报率 | ROAS | 归因收入 ÷ spend。收入来源：平台口径（Google Conversions value / Meta purchase_roas）优先；电商平台回传不重复叠加 |
| 预算消耗率 | budget_rate | 当日 spend ÷ 日预算 ×100% |
| 库存同步滞后 | — | 店铺当前库存时间戳 − 最近成功同步时间，≤ 30 分钟视为健康 |

> 约定事项：凡报表同时展示多平台转化 / ROAS 时，必须标注"各平台归因口径不同，数值不可直接相加"，并允许按来源拆分下钻。

---

## 1. 项目概述

### 1.1 产品名称
Orbit（广告智能中枢）

### 1.2 产品定位
AI Agent 驱动的多平台广告自动化与智能化管理平台，支持 **Google Ads**（含 Search / Shopping / Performance Max）、**Meta**（Facebook / Instagram）与 **Microsoft Advertising（Bing Ads）**，深度对接 **WooCommerce** 与 **Shopify** 电商平台，服务多租户、多店铺、多广告账户。一期服务对象为中小电商卖家（含 pawpel.com），以 Agent 对话为核心交互替代传统广告后台的重复操作。

### 1.3 核心目标
- 一期快速上线，覆盖 Google / Meta / Bing 三大广告平台 + WooCommerce / Shopify 两大电商平台
- 支持多租户、多店铺、多广告账户，客户数据严格隔离
- 以 Agent 为核心实现智能化与自动化，传统后台为辅
- 优先服务真实业务（含 pawpel.com 等网站），于 Q4 旺季前（目标：2026-09 中旬内测、10 月前可上线）产生效果

### 1.4 技术原则
- 前后端分离 + Monorepo 管理（apps/web、apps/api）
- 后端 Golang（Gin），前端 Next.js（App Router）
- 前端部署 Vercel，后端 Docker 自部署（api.orbit.elstella.com）
- Agent 优先，传统后台为辅；所有 Agent 操作可追溯
- 高风险操作（改价 / 改预算 / 暂停 / 删除）默认人工确认，可配置自动执行

### 1.5 术语与命名约定

| 术语 | 含义 |
|------|------|
| 租户（Tenant） | 一个客户组织，隔离单元；拥有若干店铺与广告账户 |
| 店铺（Store） | 一个电商平台接入实例（某 WooCommerce/Shopify 站点） |
| 广告账户（Ad Account） | 某广告平台下的一个投放账户，隶属于某个商店/租户 |
| 范围（Scope） | 指标查询作用于的粒度：tenant / store / account / campaign / ad_group |
| Agent | 系统内执行特定职责的 AI 编排单元，通过 Tool 与外部交互 |
| Tool | Agent 可调用的具备确定性入参/出参的功能（后端函数） |
| P0 / P1 | 需求优先级：P0=一期必须交付；P1=二期 |
| OAuth Token | 各广告/电商平台签发的访问令牌（含 refresh token 的平台一并持久化） |

---

## 2. 系统架构

### 2.1 整体架构

```
用户浏览器
    ↓
前端（Next.js @ Vercel，/app 路由 + SSE 流式 Agent 对话）
    ↕  HTTPS  REST  /api/v1  +  WebSocket(SSE 用于 Agent 流式)
后端 API（Golang Gin，Docker 自部署 api.orbit.elstella.com）
    ├─ API 服务（REST + SSE）
    ├─ Worker（Asynq 任务：数据同步 / Feed 生成 / 规则评估 / 报告发送）
    ├─ Scheduler（cron：周期任务编排）
    └─ 共享层（适配器 / 服务 / 存储 / 外部 API 客户端）
    ↓
PostgreSQL（主数据 + Feed 快照）  Redis（缓存 + Asynq 队列 + 会话 + 限流）
    ↓
广告平台 API（Google Ads / Meta / Bing）+ 电商平台 API（WooCommerce / Shopify）
```

### 2.2 技术栈（确定版本基线，锁定可复现）

| 层级 | 技术选型 | 版本基线 / 说明 |
|------|----------|-----------------|
| 前端 | Next.js（App Router）+ TypeScript + Tailwind CSS + shadcn/ui | Next.js 14+；React 18+；Tailwind 3.x |
| 前端数据层 | TanStack Query（SWR 可选）+ Zustand（轻量状态） | SSR 页面用 Server Components，动态数据走 Client fetch |
| 后端 | Go 1.22+ + Gin | Gin v1.9+；ORM 采用 GORM（PostgreSQL dialect） |
| 任务队列 | Asynq（Redis 驱动）+ gocron（scheduler） | worker 与 scheduler 独立进程，便于横向扩容 |
| 数据库 | PostgreSQL 15+ | 主库；Feed 商品快照；WAL 归档备份 |
| 缓存 | Redis 7+ | 缓存、Asynq 队列、限流计数、刷新 token 黑名单 |
| AI | 外部大模型 API（支持 Function Calling）+ 本系统 Tool 注册表 | 模型供应商与密钥经环境变量配置，Agent 不做模型内决策持久化 |
| 认证 | JWT（access + refresh）+ 各平台 OAuth2 | Token 加密存储（AES-256-GCM） |
| 部署 | 前端 Vercel；后端 docker-compose（API / Worker / Scheduler） | 健康检查 + 日志采集，见 5.12 |
| 观测 | 结构化 JSON 日志 + Prometheus 指标 + 可选 Grafana | 错误率 / 队列深度 / 同步延迟看板 |

### 2.3 Monorepo 结构（细化）

```
orbit/
├── apps/
│   ├── web/                       # Next.js 前端
│   │   ├── app/                   # App Router 页面与路由
│   │   ├── components/            # shadcn/ui + 业务组件
│   │   ├── lib/                   # API client、类型、工具函数
│   │   └── .env.example
│   └── api/                       # Golang 后端
│       ├── cmd/
│       │   ├── server/            # API 服务入口
│       │   ├── worker/            # Asynq worker 入口
│       │   └── scheduler/         # cron scheduler 入口
│       ├── internal/
│       │   ├── auth/              # JWT / RBAC / 租户上下文
│       │   ├── tenant/            # 租户与多租户隔离中间件
│       │   ├── store/             # 电商平台适配器（woocommerce/shopify）
│       │   ├── adplatform/        # 广告平台适配器（google/meta/bing）
│       │   ├── oauth/             # OAuth 授权与 token 管理
│       │   ├── sync/              # 数据同步任务
│       │   ├── feed/              # Feed 生成与发布
│       │   ├── agent/             # Agent 编排、Tool 注册表、人工确认
│       │   ├── rule/              # 规则引擎
│       │   ├── report/            # 报表与看板
│       │   ├── notify/            # 站内通知 + 邮件
│       │   ├── audit/             # 审计日志
│       │   ├── model/             # GORM 数据模型
│       │   └── httputil/          # 统一响应 / 错误码 / 中间件
│       ├── pkg/                   # 可复用公共库
│       └── config/                # 配置加载（env）
├── packages/
│   └── shared/                    # 共享类型、常量（指标口径、错误码常量定义）
├── docker/                        # Dockerfile、docker-compose.yml、init.sql
├── docs/                          # 接口文档、ADR、部署手册
└── Makefile                       # 常用命令（dev / build / test / migrate）
```

### 2.4 多租户隔离实现方案（落地条款）

1. **数据层（PostgreSQL）**：所有业务表强制携带 `tenant_id` 列，并建立复合索引 `(tenant_id, 其他业务键)`。查询统一经 Gin 中间件注入的 `x-tenant-id`（取自 JWT claims）拼接 `WHERE tenant_id = ?`。**一期采用"应用层过滤 + 数据库级 RLS 双保险"**：
   - 应用层：`TenantScope()` 中间件强制注入租户过滤，Map 任一业务 repo 查询，杜绝漏 WHERE。
   - 数据库层：对业务表启用 Row Level Security，策略为 `tenant_id = current_setting('app.current_tenant_id')`；前端连接每次事务前 `SET LOCAL app.current_tenant_id = $1`。
   - 所有写入操作在服务层断言"资源.tenant_id == 会话.tenant_id"，不一致返回 `403 FORBIDDEN`。
2. **Redis 层**：键统一加 `tenant:{id}:` 前缀（如 `tenant:12:token:refresh:blacklist`、`tenant:12:cache:feed:*`）；任务队列消息中携带 tenant_id 字段并校验处理器侧租户上下文。
3. **任务与异步层**：入队任务带 `tenant_id`，worker 消费时从任务上下文恢复租户过滤；公共资源（如平台 API 调用配额）按租户维度限流。
4. **产品与账号层**：用户可属于多租户，但每次会话明确一个"当前租户"；切换租户 = 重新取租户级 JWT claims（`tid` 字段）；跨租户数据复制 / 导入必须由超级管理员显式操作并记审计日志。
5. **隔离验证（测试强制项）**：单测必须包含"租户 A 无法通过改 tenant_id 参数读取租户 B 数据"的越权用例；安全测试按 5.13 执行。
6. **备份与恢复**：按租户维度支持导出 / 删除（GDPR 类删除诉求预留），删除时级联清理 PG 与 Redis 及对象存储产物。

---

## 3. 功能需求（一期完整覆盖）

### 3.1 用户与权限（多租户）

#### 3.1.1 账号体系
- 注册：邮箱 + 密码（密码强度：≥8 位，含大小写字母与数字）；注册后自动创建默认租户（租户名 = 邮箱前缀或用户自定义）。
- 登录：邮箱 + 密码；支持后续微信 / Google 第三方 OAuth 登录（二期）。
- 会话：JWT，access token 15 分钟 + refresh token 7 天（详见 5.11.1）。
- 账号安全：登录失败 5 次 / 10 分钟 锁定 15 分钟；支持修改密码、找回密码（邮件验证码，有效期 30 分钟）。

#### 3.1.2 租户与角色
角色定义与权限矩阵（RBAC，权限点以资源:动作命名）：

| 资源 | 动作 | 超级管理员 | 客户管理员 | 操作员 | 只读观察者 |
|------|------|:---:|:---:|:---:|:---:|
| tenant | read / update | ✅ | ✅ | — | — |
| member | read / invite / remove / change_role | ✅ | ✅ | read | read(仅自己) |
| store | read / create / update / delete / connect | ✅ | ✅ | read / create / update | read |
| ad_account | read / connect / update / delete | ✅ | ✅ | read / update | read |
| ad_campaign | read / update(pause/bid/budget) | ✅ | ✅ | ✅(受租户自动执行开关约束) | read |
| feed | read / create / update / regenerate / delete | ✅ | ✅ | ✅ | read |
| rule | read / create / update / enable / disable / delete | ✅ | ✅ | ✅ | read |
| rule_execution | read | ✅ | ✅ | ✅ | read |
| report | read / export / schedule | ✅ | ✅ | ✅ | read |
| agent | chat / approve_action / reject_action | ✅ | ✅ | ✅ | chat(只读,不可确认动作) |
| notification | read / mark_read | ✅ | ✅ | ✅ | ✅ |
| audit_log | read | ✅ | ✅ | — | — |
| admin(租户级) | read/update tenant config | ✅ | — | — | — |

- 超级管理员：平台级，可查看 / 管理所有租户；客户管理员：租户内最高权限；操作员：日常运营（有确认权限）；只读观察者：仅可查看报表和 Agent 对话，不能确认动作。
- 成员管理：客户管理员可邀请（邮箱邀请链接，48 小时有效）、移除成员、变更角色；成员离开租户前需移交或删除其名下待确认动作。

### 3.2 平台授权与账户管理

#### 3.2.1 支持平台与授权方式

| 平台 | OAuth 类型 | 关键附加要求 | 授权粒度 |
|------|-----------|--------------|----------|
| Google Ads | Google OAuth2（Authorization Code + offline） | 需开发者令牌 Developer Token + 客户 ID（MCC 或普通客户） | manager 下多账户批量授权 / 单账户授权 |
| Meta | Facebook Login（ads 相关 permission） | App 需上线审核通过 ads_management / ads_read / read_insights；Long-lived token（60 天）可后台刷新 | 用户 + Business 名下广告账户 |
| Bing Ads | Microsoft Identity Platform OAuth2 v2 | 需 developer token（Microsoft Advertising）；offline_access scope 获取 refresh token | 登录账号名下广告账户列表 |

#### 3.2.2 授权流程（统一协议，细节见 5.7）
1. 租户内用户点击"授权" → 后端生成带 `state`（随机 + 绑定 redirect 目标）的授权 URL。
2. 用户完成平台授权并跳转回调 `/api/v1/oauth/callback?platform=google&state=...&code=...`。
3. 后端校验 `state`（防止 CSRF）→ 用 `code` 换 token → 拉取账户列表 → 展示待绑定账户。
4. 用户勾选账户与本租户下的店铺 → 绑定 → Token 加密落库（AES-256-GCM，见 5.11.3）。
5. 绑定成功后触发一次全量元数据同步（账户、Campaign 层级）。

#### 3.2.3 Token 生命周期与状态
- Token 状态机：`pending`（未绑定）→ `bound`（有效）→ `expiring`（<7 天过期）→ `expired` / `revoked` / `error`（刷新失败）。
- 自动刷新：后台定时任务每 10 分钟扫描 `expires_at < now()+7d` 的 Token；刷新时 Redis 分布式锁防止并发双刷；刷新失败标记 `error`、连续 3 次失败告警并邮件通知客户管理员。
- 账户绑定状态：`active` / `sync_error`（最近同步失败） / `auth_expired`；仪表盘状态栏实时展示。
- Token 访问：进内存与 DB 均按 tenant 隔离；日志脱敏（不打印 refresh_token / secret）。

### 3.3 电商平台对接（一期 P0 为 WooCommerce + Shopify）

| 平台 | 一期优先级 | 连接方式 | 同步内容 | 备注 |
|------|-----------|----------|----------|------|
| WooCommerce | P0 | REST API（Consumer Key/Secret 或 OAuth1） | 商品（含变体）、库存、价格、订单 | 需管理员支持自定义请求频率；分页 offset 或 cursor |
| Shopify | P0 | Admin REST/GraphQL + 商店专属 Access Token | 商品（含变体）、库存、订单 | GraphQL cursor 分页；支持 incremental sync（since_id / updated_at） |
| OpenCart / Shopyy / UeeShop | P1（二期） | REST / API 对接 | 同左（标准化） | 一期仅预留 Store Adapter 接口抽象 |

**Store Adapter 统一接口（代码契约）**：

```
type StoreAdapter interface {
    Connect(ctx, config) (storeInfo, error)
    TestConnection(ctx, store) error
    SyncProducts(ctx, store, cursor) ([]Product, nextCursor, error)   // 全量 + 增量
    SyncOrders(ctx, store, since, orderCursor) ([]Order, nextCursor, error)
    FetchInventory(ctx, store, productIDs) (map[string]int, error)     // 库存与价格
    GetStoreInfo(ctx, store) (StoreInfo, error)
}
```

**通用能力**：
- 多店铺绑定：一个租户可绑定多个 WooCommerce / Shopify 店铺，店铺可绑定多个广告账户。
- 商品同步：全量（首次）＋增量（每 30 分钟轮询或 webhook 可选）；商品纳入标准 Product 模型（见 5.3）；含变体（sku、库存、价格、图、描述、分类、GTIN/brand 可选）。
- 库存 / 价格同步：增量拉取；价格变动写入 Feed 快照并触发 Feed 更新。
- 订单回流：订单数据用于转化归因展示（不参与 ROAS 叠加，见 0.4），支持"订单→转化"匹配报表（按 SKU / 来源标记）。
- 连接健康：TestConnection 每 15 分钟执行一次；失败 ≥3 次置 `sync_error` 并通知。

### 3.4 广告管理与数据同步

#### 3.4.1 统一视图
- 层级：`Campaign → Ad Group → Ad`，跨平台抽象字段统一：名称、状态（enabled/paused/removed 归一化）、预算、出价策略、指标（0.4 口径）。
- 平台差异落在"原始字段"栏（如 Google bidding strategy、Meta optimization_goal），不做强行一致。
- 列表过滤：按账户 / 状态 / 时间范围 / 搜索词；支持排序（spend、roas、cpa）。

#### 3.4.2 数据同步（详见 5.8）
- 同步维度：实体（campaign/ad_group/ad）+ 指标（日聚合）。
- 频次：实体元数据每 3 小时；日指标每天 08:00（Asia/Shanghai）+ 手动刷新随时。
- 粒度：日粒度（date+account+id 主键）；小时粒度仅 PMax/Behavioral 预留（二期）。
- 手动刷新：前端按钮 + API（POST /adsync/{scope}），并发性由队列保证。

#### 3.4.3 跨平台聚合
- 维度对齐：账户 → 店铺 → 租户；报表可按店铺聚合多个平台账户。
- 聚合规则：花费 / 展示 / 点击可跨平台直接求和；转化 / ROAS / CPA 不支持跨平台求和，按平台并列展示（0.4 约定）。
- 时区统一：平台数据存放用 UTC，前端展示按租户配置时区（默认 Asia/Shanghai）。

### 3.5 Agent 智能化体系（核心）

#### 3.5.1 Agent 矩阵（职责 + 核心能力）

| Agent | 职责 | 输入 | 输出 | 关键 Tool（见 5.6） |
|-------|------|------|------|----------------------|
| Sync Agent | 广告 + 电商数据同步 | 范围、动作指令 | 同步结果 / 差异摘要 | sync_platform_data |
| Report Agent | 报表 / 看板 / 日报周报 | 范围、时间、指标 | 指标卡片、趋势图数据、文章式总结、导出文件 | get_metrics、generate_report、export_report |
| Optimize Agent | 低效诊断 / 优化建议 / 确认执行 | 范围、目标 | 诊断列表 + 建议（含预期影响） | get_metrics、get_campaigns、pause_*、update_bid、update_budget、get_ad_groups |
| Rule Agent | 执行用户规则 | 规则触发上下文 | 评估结果、动作执行记录 | evaluate_rule（内部）、trigger 动作 |
| Creative Agent | 商品文案与素材建议 | 商品 / 平台 | 文案（标题/描述/卖点）、素材构图建议（输出 Markdown 建议，不调用媒体服务） | get_product、generate_creative_suggestion |
| Ops Agent | 自然语言指令处理与路由 | 用户文本 | 路由结果 / 工具调用计划 / 对话回复 | 编排层（见 3.5.3），无独立业务 Tool |
| Notify Agent | 异常 / 报告 / 操作通知 | 事件 | 站内通知 + 邮件 | send_notification（内部） |

#### 3.5.2 高风险操作与人工确认协议
- 风险分级：
  - 高风险（必须确认）：暂停 / 启用 campaign 或 ad group、修改出价、修改预算、删除资源。
  - 中风险（默认确认，可配置自动）：批量同步、Feed 重新生成。
  - 低风险（直接执行）：查询 / 报表 / 导出 / 生成文案建议。
- 确认流：Agent 工具调用阶段若命中高风险 → 创建 `agent_actions`（状态 pending）→ 经 SSE 推送前端 → 用户在对话卡片点击"确认 / 拒绝"（附理由）→ 后端锁定执行 → 结果回写到对话。待确认动作 24 小时未处理按"放弃"处理并通知。
- 自动执行开关：租户级配置 `auto_approve_high_risk`（默认 false）。开启后仅对操作员与客户管理员生效，且每次自动执行前仍校验资源一致性（见 5.10.4），并全量记审计日志。

#### 3.5.3 Agent 编排流程（Ops Agent 路由）
```
用户消息 → Ops Agent（LLM，无工具）：意图识别 + 提取参数（范围/时间/目标）
  → 路由到目标 Agent（Sync/Report/Optimize/Creative/Rule）
  → Agent 按需调用 Tool（Function Calling，LLM 决策调用顺序）
  → 高风险工具 → 人工确认（3.5.2）
  → 结果结构化返回 → Ops Agent 汇总生成自然语言回复（SSE 流式）
```
- 每个对话会话（conversation）绑定当前租户 + 当前店铺/账户上下文（会话切换时作为 system prompt 注入）。
- 所有工具调用记录 `agent_tool_calls` 表（入参、出参摘要、耗时、风险分级、确认状态），用于审计与排障。

### 3.6 报表与看板

- 跨平台统一报表：按 0.4 口径；维度支持 时间 / 平台 / 账户 / 店铺 / campaign。
- 基础自定义看板：指标卡片（本期 6 个核心卡片：spend、impressions、clicks、CTR、CPC、CPA、ROAS 中选配展示）、趋势图（折线）、表格（明细）；支持保存为"看板方案"，可设为默认。
- 数据导出：CSV / Excel；导出走异步任务（Asynq），完成后站内通知 + 邮件下载链接（有效期 24 小时）。
- 定时报告：支持日报 / 周报，投递到指定邮箱；模板由 Report Agent 生成（含自动化文字解读）。
- 指标冷启动：新绑定账户无历史数据时显示 N/A 并提示"数据积累中"，不填 0 避免误判。

### 3.7 自动化规则

- 创建：表单式（条件 + 动作），支持预览"若现在满足则哪些对象将触发"。
- 条件与动作的完整定义见 5.10；一期支持条件维度：账户 / 店铺 / campaign / ad_group，指标：spend、roas、cpa、ctr、cvr、impressions、budget_rate（0.4 口径），时间窗口：昨日 / 近 7 天 / 本周累计 / 连续 N 天。
- 动作：pause、update_bid(+/-百分比 或 绝对值)、update_budget、notify。
- 触发：每次相关数据同步完成后评估 + 每日 09:00 例行评估（防抖见 5.10.4）。
- 预算预警与平滑消耗：内置规则模板（预算消耗率 > 80% 提醒、> 100% 告警；平滑消耗建议将预算按日分摊提示），用户可一键启用/禁用模板。

### 3.8 通知与日志

- 通知渠道：站内通知（必达）+ 邮件（重要事件）。邮件支持 SMTP，Html 模板；站内通知经 SSE 实时推送。
- 事件类型：连接状态变化、同步失败、Token 即将过期、规则触发执行、待确认动作创建 / 超时、报告就绪、系统告警。
- 通知设置：租户级开关（哪些事件发邮件）、用户级免打扰（可选，一期支持全局开关）。
- 审计日志：所有写操作与 Agent 工具调用记录，字段见 5.11.5；保留 180 天，支持按用户 / 动作 / 时间检索。
- 平台连接状态监控：账户 / 店铺连接状态看板 + 红黄绿状态色（active / sync_error / auth_expired / disabled）。

---

## 4. 前端需求（Next.js）

### 4.1 核心页面与路由

| 路由 | 页面 | 主要能力 | 权限 |
|------|------|----------|------|
| /login、/register、/forgot-password | 认证 | 登录 / 注册 / 找回密码 | 公开 |
| /auth/callback?platform=... | OAuth 回调中转页 | 显示"授权处理中"，调 API 完成绑定 | 已登录 |
| /（工作台） | Dashboard | 指标卡（spend/CTR/CPA/ROAS）、快速入口、待确认动作角标 | 全部 |
| /stores | 店铺管理 | 店铺列表、绑定向导、连接状态、同步触发 | 见 3.1.2 |
| /stores/:id | 店铺详情 | 商品 / 订单 / 同步历史、Feed 状态 | 同上 |
| /accounts | 广告账户管理 | 账户列表、授权入口、状态、Token 健康 | 同上 |
| /advertising/campaigns | 广告管理 | 统一视图列表（过滤/排序/搜索）、批量选择 | 同上 |
| /advertising/campaigns/:id | 广告详情 | 层级下钻（ad group/ad）、时间序列、Agent 优化入口 | 同上 |
| /reports | 数据报表 | 跨平台报表、维度选择、导出 | read |
| /reports/:id | 看板详情 | 自定义看板（卡片/趋势/表格） | read |
| /agent | Agent 对话 | 流式对话、工具调用卡片、确认/拒绝按钮 | 见 3.1.2 |
| /agent/cards | 待确认动作中心 | 全部 pending/历史 action 列表 | approve 权限 |
| /rules | 规则管理 | 规则列表、创建向导（条件+动作表单）、启停 | 见 3.1.2 |
| /products | 商品管理 | 商品列表、变体、同步状态 | read / update |
| /feeds | Feed 管理 | Feed 列表、URL、手动触发、日志 | 见 3.1.2 |
| /settings | 系统设置与权限 | 租户资料、成员管理、角色、自动执行开关、通知偏好 | 见 3.1.2 |
| /notifications | 通知中心 | 通知列表、已读/未读、筛选 | 全部 |

### 4.2 交互要求
- 店铺 / 账户快速切换：全局顶栏下拉，切换后所有数据视图刷新（上下文写入会话并注入 Agent 对话）。
- Agent 对话流式输出：SSE 逐字渲染；工具调用以结构化卡片展示（调用中 → 结果 → 风险提示 → 确认/拒绝）。断线重连需支持恢复（携带 conversation_id 与 last_message_id）。
- 重要操作二次确认：前端弹窗"该操作将暂停 X 广告组，影响预估日支出 Y 元，确认继续？"；高风险操作除 Agent 卡片外，传统后台也需二次确认。
- 响应式：优先桌面（≥1280px 完整布局），平板/手机降级（关键操作可用，深度管理建议桌面）。
- 中英双语：一期中文为主，i18n 架构预留（next-intl），英文文案二期启用。

### 4.3 UI 规范
- shadcn/ui + Tailwind CSS；专业数据后台风格：左侧导航 + 顶部全局栏 + 内容区。
- 指标卡片统一组件（标题、数值、环比变化、口径 tooltip）；状态色规范：绿=正常、黄=预警（预算消耗>80%、Token 即将过期）、红=异常。
- 数据加载状态：骨架屏优先；错误态提供"重试"；空态给出引导文案与操作入口。
- 表单校验与接口错误统一映射到错误码文案（见 5.5），前端不硬编码散落文案。

### 4.4 数据请求与状态管理
- API Client：统一封装 fetch，自动附加 access token（403 时静默 refresh 重试一次）。
- 数据获取：React Query（Server State）+ Zustand（Client 全局态：当前租户/店铺/账户上下文）。
- SSR 策略：登录态与静态页走 Server Components；指标、列表等动态数据 Client fetch + 缓存（staleTime 默认 30s，手动刷新按钮强制失效）。
- 通知与待确认：全局 SSE 长连接统一推送（站内通知、待确认动作、同步完成），连接失败指数退避重连。

---

## 5. 后端需求（Golang）

### 5.1 核心模块（与 Monorepo 的 internal/ 对应）
1. 用户与权限服务（auth：注册/登录/JWT/RBAC/租户上下文）
2. 平台授权服务（oauth：授权 URL、回调、Token 加密存取、自动刷新 worker）
3. 广告平台适配器（adplatform/google、adplatform/meta、adplatform/bing）
4. 电商平台适配器（store/woocommerce、store/shopify）
5. 数据同步服务（sync：实体 + 指标 + 电商数据任务）
6. Feed 生成服务（feed：快照、格式渲染、发布、日志）
7. Agent 编排与 Tool 服务（agent：编排、Tool 注册表、确认流）
8. 规则引擎（rule：评估、动作执行、防抖）
9. 报表服务（report：聚合、导出、定时报告）
10. 通知与日志服务（notify + audit）

### 5.2 接口设计原则
- RESTful + SSE（Agent 对话流式，SSE 优于 WebSocket：单向上行简单、Vercel/跨域友好；WebSocket 一期不启用）。
- 统一前缀 `/api/v1`；响应统一包裹（5.5 协议）。
- JWT 鉴权 + 租户上下文中间件（从 claims 取 `uid`、`tid` 注入 gin.Context）。
- 完整错误码与结构化日志（request_id 贯穿，5.5 / 5.11.5）。
- 全部写接口幂等约束：携带客户端生成的 idempotency key（header `Idempotency-Key`）可选，平台写接口（改价/暂停）必须带。

### 5.3 数据存储（PostgreSQL 核心表结构）

> 约定：所有表都含 `id BIGSERIAL PRIMARY KEY`、`tenant_id BIGINT NOT NULL`、`created_at TIMESTAMPTZ NOT NULL DEFAULT now()`、`updated_at TIMESTAMPTZ`；`tenant_id` 与业务键建立组合索引；软删除用 `deleted_at TIMESTAMPTZ`。以下列出差异化字段。

#### 5.3.1 用户 / 租户 / 权限域

| 表 | 说明 | 关键字段 | 索引 / 约束 |
|----|------|----------|-------------|
| users | 用户账号 | email(唯一,小写), password_hash, status(active/locked), last_login_at | uk_email; idx(status) |
| tenants | 租户 | name, plan(beta/pro), timezone(默认Asia/Shanghai), auto_approve_high_risk(bool,默认false), settings(jsonb) | — |
| tenant_members | 租户成员 | user_id, tenant_id, role(super_admin/customer_admin/operator/viewer), invited_by, invited_at, status(pending/active/removed) | uk(tenant_id,user_id); idx(user_id) |
| refresh_tokens | 刷新令牌（也可用 JWT 黑名单替代） | user_id, tenant_id, token_hash, expires_at, revoked_at, user_agent, ip | idx(user_id); idx(revoked_at) |
| audit_logs | 审计日志 | actor_user_id, tenant_id, action, resource_type, resource_id, before_snapshot(jsonb), after_snapshot(jsonb), ip, created_at | idx(tenant_id,created_at); idx(actor_user_id) |

#### 5.3.2 平台连接域

| 表 | 说明 | 关键字段 | 索引 / 约束 |
|----|------|----------|-------------|
| platform_conns | 电商/广告连接统一表 | type(store/ad_account), platform(google/meta/bing/woocommerce/shopify), admin_user_id(授权人), status(bound/expiring/expired/error/revoked), meta(jsonb:平台返回的账号信息), last_synced_at, store_id(可空,广告账户绑定的店铺) | idx(platform,status); idx(tenant_id) |
| oauth_tokens | 加密 Token | conn_id, token_type(refresh/longlived/consumer), encrypted_token1, encrypted_token2(可空,如secret), expires_at, scopes, last_refresh_at, refresh_error_count | uk(conn_id); idx(expires_at) |
| store_configs | 店铺详细配置 | store_id, base_url, extra_headers(jsonb), sync_cursor(jsonb:各资源游标), webhook_secret(可空) | uk(store_id) |

#### 5.3.3 商品 / 订单域

| 表 | 说明 | 关键字段 | 索引 / 约束 |
|----|------|----------|-------------|
| products | 商品（店铺维度） | store_id, external_id(平台商品id), title, description, link, image_url, brand, gtin(可空), google_product_category, product_type, price_cents, compare_at_price_cents, currency, status(active/paused/draft/archived 归一化), hash(同步去重) | uk(store_id,external_id); idx(tenant_id, store_id); GIN(title 分词) |
| product_variants | 商品变体 | product_id, external_id, sku, option_string, price_cents, inventory_qty, image_url | uk(product_id,external_id) |
| orders | 订单 | store_id, external_id, order_number, status, total_cents, currency, customer_email, items(jsonb), placed_at | uk(store_id,external_id); idx(tenant_id, placed_at) |

#### 5.3.4 广告实体与指标域

| 表 | 说明 | 关键字段 | 索引 / 约束 |
|----|------|----------|-------------|
| ad_accounts | 广告账户元数据 | conn_id, platform, external_id(平台账户id), name, currency, timezone, customer_id(google)/business_id(meta)/developer_token_registered, status | uk(tenant_id,platform,external_id) |
| campaigns | Campaign（跨平台归一） | ad_account_id, external_id, name, status(enabled/paused/removed), type(search/pmax/shopping/feed/manual), daily_budget_cents, bidding_strategy, raw(jsonb 平台原始) | uk(ad_account_id,external_id); idx(tenant_id,ad_account_id) |
| ad_groups | Ad Group（含 Meta ad set 归入本表，type 区分） | ad_account_id, campaign_id, external_id, name, status, bid_micros(可选), targeting(jsonb) | uk(ad_account_id,external_id); idx(campaign_id) |
| ads | Ad（Meta ad / Google RSA ad） | ad_account_id, ad_group_id, external_id, name, status, creative_type, headline_1/2(可空), description_1/2(可空), raw(jsonb) | uk(ad_account_id,external_id); idx(ad_group_id) |
| daily_stats | 日粒度指标主表（聚合广告+电商可选列） | scope_type(account/campaign/ad_group/ad), scope_id, date(dt), platform, metrics(jsonb: spend_micros, impressions, clicks, conversions, conversion_value_micros, ctr, cpc, cpa, roas...), source(synced/derived) | pk(scope_type,scope_id,date); uk 防重；idx(tenant_id,date) |
| sync_jobs | 同步任务实例 | type(ad_meta/ad_metrics/products/orders/feed_meta), scope(tenant/store/account), object_id, status(pending/running/success/failed), cursor_before, cursor_after, started_at, finished_at, error_msg, retry_count | idx(tenant_id,status,created_at) |

#### 5.3.5 Feed 域

| 表 | 说明 | 关键字段 | 索引 / 约束 |
|----|------|----------|-------------|
| feeds | Feed 配置 | store_id, name, format(google_shopping_xml/google_shopping_tsv/meta_catalog_csv), include_condition(active 等), currency_override, public_token(URL 访问令牌), status(active/paused), last_generated_at, generated_rows, error_count | uk(tenant_id,store_id,name); uk(public_token) |
| feed_runs | 单次生成记录 | feed_id, status(success/failed/partial), rows, errors(jsonb: 每条错误{sku,field,reason}), started_at, finished_at | idx(feed_id,created_at) |
| feed_product_snapshot | 生成快照（增量差量比对用） | feed_id, product_id, snapshot_hash | uk(feed_id,product_id) |

#### 5.3.6 Agent / 规则 / 通知域

| 表 | 说明 | 关键字段 | 索引 / 约束 |
|----|------|----------|-------------|
| conversations | Agent 会话 | user_id, tenant_id, title, store_id(可空,绑定的店铺上下文), account_id(可空), status(active/closed), llm_config(jsonb) | idx(tenant_id,user_id) |
| messages | Agent 消息 | conversation_id, role(user/assistant/tool/system), content(text), created_at | idx(conversation_id,created_at) |
| agent_tool_calls | 工具调用记录 | message_id, tool_name, input(jsonb), output_summary(jsonb), risk_level, status(success/failed/need_confirm/cancelled), duration_ms, confirm_action_id(可空) | idx(conversation_id); idx(status) |
| agent_actions | 待确认/已执行动作 | tenant_id, action_type(pause/update_bid/update_budget/regenerate_feed/...), target_type, target_id, params(jsonb), risk_level, status(pending/approved/rejected/expired/executed/failed), requested_by(agent 会话), decided_by(人工), decided_at, result(jsonb) | idx(tenant_id,status); idx(expires_at) |
| rules | 自动化规则 | tenant_id, name, enabled, scope(account/campaign/ad_group), condition_spec(jsonb，见5.10), action_spec(jsonb), cooldown_minutes(默认 1440), created_by, last_run_at | idx(tenant_id,enabled) |
| rule_executions | 规则执行记录 | rule_id, triggered_by(manual/sync/daily), matched_objects(jsonb), executed_actions(jsonb), status(success/partial/failed/none), error_msg, created_at | idx(rule_id,created_at) |
| notifications | 站内通知 | user_id, tenant_id, type, title, body, link(可空), read_at, created_at | idx(user_id,read_at) |
| email_logs | 邮件发送记录 | tenant_id, to_email, subject, template, status(sent/failed), error_msg, created_at | idx(tenant_id,created_at) |
| report_schedules | 定时报告 | tenant_id, name, type(daily/weekly), scope_spec(jsonb), recipients(jsonb 邮箱数组), template(config), cron_config, enabled, last_sent_at | idx(tenant_id,enabled) |
| webhook_deliveries | 电商 webhook 投递记录 | store_id, event_type, payload_hash, status(processed/failed), processed_at | idx(store_id,created_at) |

**迁移策略**：使用 golang-migrate 管理 schema 版本，`docker/init.sql` 提供一处引导；生产禁止自动建表，迁移由发布流程执行。

### 5.4 核心 API 接口清单（一期全部接口，前缀 /api/v1）

> 分页统一参数 `page(default 1)` + `page_size(default 20, max 100)`；列表响应 `{items, page, page_size, total}`。

| 模块 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 认证 | POST | /auth/register | 注册 |
| 认证 | POST | /auth/login | 登录，返回 access+refresh |
| 认证 | POST | /auth/refresh | 刷新 access |
| 认证 | POST | /auth/logout | 登出（吊销 refresh） |
| 认证 | GET | /auth/me | 当前用户信息 |
| 认证 | POST | /auth/forgot-password | 发送重置邮件 |
| 认证 | POST | /auth/reset-password | 重置（带 token） |
| 租户 | POST | /tenants | 创建租户（注册可自动建） |
| 租户 | GET / PATCH | /tenants/:id | 查询 / 更新租户设置 |
| 租户 | POST / GET | /tenants/:id/members | 邀请成员 / 成员列表 |
| 租户 | PATCH / DELETE | /tenants/:id/members/:uid | 变更角色 / 移除 |
| OAuth | GET | /oauth/connect/:platform | 生成授权 URL（带上 state） |
| OAuth | GET | /oauth/callback | 统一回调：校验 state → 换 token → 拉账户列表 |
| OAuth | POST | /oauth/bind | 确认绑定所选账户/店铺 |
| 店铺 | GET / POST | /stores | 列表（当前租户）/ 创建（WooCommerce 填 key/secret；Shopify 填 token+store） |
| 店铺 | GET / PATCH / DELETE | /stores/:id | 详情 / 更新 / 删除 |
| 店铺 | POST | /stores/:id/test | 连通性测试 |
| 店铺 | POST | /stores/:id/sync | 手动触发同步（products/orders） |
| 店铺 | GET | /stores/:id/orders | 订单列表 |
| 广告账户 | GET / POST | /ad-accounts | 已绑定账户列表 / 手工登记 |
| 广告账户 | GET / PATCH / DELETE | /ad-accounts/:id | 详情 / 更新 / 解绑 |
| 广告账户 | POST | /ad-accounts/:id/sync | 手动刷新元数据+指标 |
| 广告 | GET | /advertising/campaigns | Campaign 列表（account_id 过滤、搜索、分页） |
| 广告 | GET | /advertising/campaigns/:id | 详情 + 指标 |
| 广告 | GET | /advertising/ad-groups | 按 campaign_id 取 ad group |
| 广告 | GET | /advertising/ads | 按 ad_group_id 取 ads |
| 广告 | POST | /advertising/:type/:id/pause · resume | 暂停 / 启用（写接口，高风险） |
| 广告 | POST | /advertising/ad-groups/:id/bid | 修改出价（写接口，高风险） |
| 广告 | POST | /advertising/campaigns/:id/budget | 修改预算（写接口，高风险） |
| 指标 | GET | /metrics/summary | 汇总指标（scope 数组 + date_range + metrics） |
| 指标 | GET | /metrics/timeseries | 时间序列（step day/week） |
| Feed | GET / POST | /feeds | Feed 列表 / 创建 |
| Feed | GET / PATCH / DELETE | /feeds/:id | 详情 / 更新 / 删除 |
| Feed | POST | /feeds/:id/regenerate | 手动重新生成 |
| Feed | GET | /feeds/:id/logs | 生成历史日志 |
| Feed | GET | /public/feeds/:token | Feed 公开可访问 URL（不鉴权，token 鉴权） |
| Agent | POST | /agent/chat | 发起对话（SSE 流式返回） |
| Agent | GET | /agent/conversations | 会话列表 |
| Agent | GET | /agent/conversations/:id/messages | 消息历史 |
| Agent | GET | /agent/actions | 待确认 / 历史动作列表 |
| Agent | POST | /agent/actions/:id/approve | 确认并执行动作 |
| Agent | POST | /agent/actions/:id/reject | 拒绝动作 |
| 规则 | GET / POST | /rules | 列表 / 创建 |
| 规则 | GET / PATCH / DELETE | /rules/:id | 详情 / 更新 / 删除 |
| 规则 | POST | /rules/:id/enable · disable | 启停 |
| 规则 | POST | /rules/:id/run-now | 立即试跑 |
| 规则 | GET | /rules/:id/executions | 执行记录 |
| 报表 | GET | /reports/summary | 跨平台聚合（维度+指标） |
| 报表 | POST | /reports/export | 发起导出（异步，返回 job_id） |
| 报表 | GET | /reports/export/:job_id | 查询导出任务状态 / 下载链接 |
| 报表 | GET / POST | /reports/schedules | 定时报告列表 / 创建 |
| 报表 | PATCH / DELETE | /reports/schedules/:id | 更新 / 删除 |
| 通知 | GET | /notifications | 通知列表 |
| 通知 | PATCH | /notifications/:id/read | 标记已读 |
| 通知 | GET | /notifications/unread-count | 未读数（SSE 头也推送） |
| 审计 | GET | /audit-logs | 审计日志（管理员） |
| 系统 | GET | /health/live · /health/ready | 存活 / 就绪探针 |

### 5.5 统一响应格式与错误码

#### 5.5.1 响应结构
成功：
```json
{ "code": 0, "message": "ok", "request_id": "req_8f3a...", "data": { } }
```
失败（HTTP 状态用 4xx/5xx，body 同构）：
```json
{ "code": 1003, "message": "resource not found", "request_id": "req_8f3a...", "data": null, "details": { "field": "errors" } }
```

#### 5.5.2 HTTP 状态语义
- 200 成功（含列表）；201 创建成功；202 异步任务已接受；400 参数错误；401 未认证；403 无权限/跨租户越权；404 不存在；409 冲突（重复绑定、重复名）；429 限流；500 内部错误；503 依赖不可用。

#### 5.5.3 业务错误码分段

| 段 | 范围 | 说明 | 示例 |
|----|------|------|------|
| 通用 | 1000-1099 | 参数/未认证/权限/不存在/冲突/限流/内部 | 1000 参数错误；1001 未认证；1002 无权限；1003 资源不存在；1004 冲突；1005 频率限制；1006 内部错误 |
| 认证与租户 | 2000-2099 | 会话与租户 | 2001 凭据错误；2002 token 过期；2003 refresh 失败；2004 租户不存在；2005 租户已禁用；2006 邀请无效/过期 |
| 平台集成 | 3000-3099 | 平台连接与 Token | 3001 授权失败/state 不匹配；3002 token 刷新失败；3003 平台 API 限流；3004 平台校验失败(如 developer token 无效)；3005 授权已过期；3006 平台账户不存在 |
| 数据与同步 | 4000-4099 | 同步与 Feed | 4001 同步任务失败；4002 Feed 生成失败；4003 数据尚未就绪；4004 商品字段映射错误 |
| Agent 与规则 | 5000-5099 | 对话与规则执行 | 5001 规则条件不支持；5002 动作被拒绝/资源状态已变化；5003 待确认动作已过期；5004 Agent 调用超时 |
| 通知与报表 | 6000-6099 | 通知/导出/报告 | 6001 邮件发送失败；6002 导出任务失败；6003 报告模板不存在 |

### 5.6 Agent 与 Tool 定义

#### 5.6.1 Tool 注册表（一期全部，LLM 可见的 function 定义按 JSON Schema 描述）

| Tool 名称 | 描述（LLM 侧） | 关键入参 | 出参 | 风险 |
|-----------|----------------|----------|------|------|
| sync_platform_data | 触发广告/店铺数据同步 | platform, scope_type(account/store), scope_id, mode(full/incremental) | {job_id,status} | 中 |
| get_metrics | 查询指标汇总 | scope_type, scope_ids[], date_range{start,end}, metrics[], group_by[] | {rows:[{dim,metrics}]} | 低 |
| get_timeseries | 查询指标时序 | scope, metrics[], step(day/week) | {points:[{date,metrics}]} | 低 |
| get_campaigns | 列 Campaign | account_id, status_filter | {campaigns:[...]} | 低 |
| get_ad_groups | 列 Ad Group | campaign_id | {ad_groups:[...]} | 低 |
| get_campaign_detail | Campaign 详情+近期指标 | campaign_id | {campaign, metrics} | 低 |
| update_bid | 修改 Ad Group 出价 | ad_group_id, new_bid_cents | 预执行返回 {pending_action_id} | 高 |
| update_budget | 修改 Campaign 日预算 | campaign_id, new_budget_cents | {pending_action_id} | 高 |
| pause_campaign / resume_campaign | 暂停 / 启用 Campaign | campaign_id | {pending_action_id} | 高 |
| pause_ad_group | 暂停 Ad Group | ad_group_id | {pending_action_id} | 高 |
| generate_report | 生成结构化报表（含文字解读） | scope, date_range, dimensions, metrics, format(text/markdown) | {report_doc} | 低 |
| export_report | 异步导出文件 | 同 generate_report + format(csv/xlsx) | {job_id} | 中 |
| get_product | 查商品 | product_external_id 或 query | {product, variants} | 低 |
| generate_creative_suggestion | 基于商品生成文案与素材建议 | product_ids[], platform, tone | {suggestions:[{product, headlines, descriptions, image_hint}]} | 低 |
| create_rule | 创建自动化规则 | name, scope, condition_spec, action_spec | {rule_id} | 中 |
| run_rule_now | 立即试跑规则 | rule_id | {execution} | 中 |
| get_feed_status | 查 Feed 状态与最近日志 | feed_id | {feed,last_run} | 低 |
| trigger_feed_regenerate | 重新生成 Feed | feed_id | {job_id} | 中 |
| get_connection_status | 查平台连接健康 | platform*, id | {status} | 低 |
| search_notifications | 查询通知 | filters | {notifications} | 低 |

#### 5.6.2 Tool 调用协议（Agent → 后端）
- 后端通过 `/agent/chat` SSE 返回事件序列：`start → [tool_call(tool,args) → tool_result → ...] → text_chunk*doc → end`。
- 高风险工具（risk=高）在任何 Agent 语境下均先创建 `agent_actions(pending)` 挂起，待人工确认后才真正调用平台写接口。
- Tool 执行超时：外部平台调用 30s 上限；LLM Function Calling 整体 60s 上限；超时返回 error 并在消息中说明"已重试/建议手动操作"。

### 5.7 OAuth 授权与 Token 刷新流程（详细设计）

#### 5.7.1 授权发起（GET /oauth/connect/:platform）
1. 生成 `state` = rand(32字节 hex)，Redis 存 `oauth:state:{state}` = {platform, tenant_id, user_id, redirect(path), expire 10min}。
2. 按平台拼授权 URL：
   - Google：scope=`openid email profile https://www.googleapis.com/auth/adwords`，access_type=offline，prompt=consent（确保返回 refresh_token）。
   - Meta：`https://www.facebook.com/v19.0/dialog/oauth`，scope=`ads_management ads_read read_insights pages_show_list`。
   - Bing：Microsoft 授权端点 scope=`https://ads.microsoft.com/msads.manage offline_access`。
3. 返回 {authorization_url, state}，前端跳转。

#### 5.7.2 回调处理（GET /oauth/callback）
1. 校验 `state`（防 CSRF）与 `code`。
2. 用 code 换 token：
   - Google：POST token 端点 → access_token + refresh_token，有效期 access 1h。
   - Meta：短 token → 换取 long-lived token（有效期 60 天），Meta 无 refresh_token，续期 = 重新换 long-lived。
   - Bing：access + refresh（offline_access）。
3. 校验 token：调用平台账户列表 API（Google：AccountManagementService / customers；Meta：/me/adaccounts；Bing：Customers/List）。
4. 落库：`platform_conns` + `oauth_tokens`（加密，见 5.11.3）。将潜在账户列表返回前端。

#### 5.7.3 绑定（POST /oauth/bind）
- 前端回传所选账户/店铺映射（如 Google 账户 → 店铺）。
- 后端为每个账户建 `ad_accounts` 记录并触发元数据全量同步入队（async）。

#### 5.7.4 Token 自动刷新（后台 worker）
- 调度：scheduler 每 10 分钟扫 `expires_at < now()+7d AND status in (bound,expiring)`。
- 刷新：Redis 分布式锁 `token:refresh:{conn_id}`（TTL 2min）防止并发；成功后更新 `expires_at`、`encrypted_token`、`last_refresh_at`，重置 error_count。
- 失败：`refresh_error_count++`；≥3 次置 `status=error`，Notify Agent 发站内 + 邮件告警（群发客户管理员）。
- Google/Meta/Bing 均遵守平台限流（退避 + jitter）。
- Token 失效也可由"平台调用返回 401/403/AUTH_EXPIRED"实时触发主动刷新一次；仍失败则该次调用返回 3002。

#### 5.7.5 解绑与数据保留
- 解绑：删除或吊销平台侧 access（Google/Meta 可 revoke；Bing 清除），`oauth_tokens` 置 revoked；`ad_accounts.status=removed`；历史指标数据保留供报表回溯。

### 5.8 数据同步任务设计

#### 5.8.1 任务类型与调度

| 任务类型 | 触发方式 | 频率 | 幂等键 | 说明 |
|----------|----------|------|--------|------|
| ad_meta（campaign/ad group/ad） | scheduler + 手动 | 每 3 小时；绑定后立即全量 | (tenant, account, job_type, date) | 增量游标 updated_at |
| ad_metrics（日指标） | scheduler | 每天 08:00(租户时区) + 手动 | (tenant, account, date) | 同步昨日完整日，缺天自动补拉最近 3 天 |
| products（商品+变体） | scheduler + 可选 webhook | 每 30 分钟增量；首次全量 | (tenant, store, resource, cursor_hash) | WooCommerce offset/Shopify cursor |
| orders（订单） | scheduler + 可选 webhook | 每 30 分钟增量 | (tenant, store, order_cursor) | 增量 since_id/placed_at |
| inventory（库存价格） | scheduler | 每 30 分钟（随 products 任务附带） | 同 products | 库存差异计算 |
| feed_regenerate | scheduler + 手动 | 每 4 小时；价格/库存变更时按需 | (tenant, store, feed_id, snapshot_hash) | 见 5.9 |

#### 5.8.2 任务执行约定
- 队列复用 Asynq：`critical`（写入平台、授权） / `default`（同步、生成） / `low`（邮件、报表）三队列，worker 并发默认 5/10/10（可配）。
- 每条任务写入 `sync_jobs`，状态机 `pending→running→success/failed`；失败重试：指数退避（1min → 2 → 4 → ... 最多 5 次）后进死信（dead-letter），死信任务标记 status=failed 并告警。
- 幂等：以任务级幂等键 upsert，重复入队不产生脏数据（指标按 (scope,date) 覆盖写，商品按 (store,external_id) 覆盖写）。
- 分页/游标：游标持久化到 `store_configs.sync_cursor`；广告指标 Google 用 SearchStream query 分页 / Meta insights 游标 / Bing 报表下载任务（异步轮询报表状态 30s/次，超 5min 失败）。
- 一致性校验：同步完成后对 `products.count` 与平台返回总数比对，偏差 >5% 记 warning 并通知。
- 手动刷新并发保护：同一 (scope,type) 已有 running 任务时返回"已在同步中"（409/202 语义）。
- 指标口径源：Google 字段 `cost_micros,impressions,clicks,conversions,conversions_value`（query 维度 date+campaign）；Meta `spend,impressions,clicks,actions,action_values`（breakdown by ad/date）；Bing `Spend,Impressions,Clicks,Conversions`。

### 5.9 Feed 生成规范

#### 5.9.1 支持的格式与输出
| 格式 | 用途 | 文件类型 | 字段集 |
|------|------|----------|--------|
| google_shopping_xml | Google Merchant Center | XML（RSS 2.0 商品源） | 见字段映射表 |
| google_shopping_tsv | Google Merchant Center | TSV（tab 分隔） | 同 XML 字段 |
| meta_catalog_csv | Meta Catalog | CSV | 见字段映射表 |

- 输出位置：`/public/feeds/:token`（GET 无需登录，token 随机 32 位，每 Feed 唯一）；也可在 `feeds` 配置中开启"私有化"（仅授权头访问）。
- 更新频率：默认每 4 小时全量重建 + 先生效后替换（先写临时文件→原子重命名→更新 last_generated_at）；价格/库存变更如需即时生效可通过手动 regenerate。

#### 5.9.2 字段映射（Google Shopping，来源 = Product 模型）
| Feed 字段 | 来源字段 / 规则 | 必填 |
|-----------|-----------------|------|
| id | `{平台前缀}:{store_id}:{product_external_id}`（如 wp:12:5678），保证跨平台不冲突 | 是 |
| title | products.title（≥20 字符规范，超长截断按平台规则） | 是 |
| description | products.description（截断 ≥10 字符） | 是 |
| link | products.link | 是 |
| image_link | products.image_url（首选主图） | 是 |
| price | `price_cents / 100` + currency（格式化 `12.99 USD`） | 是 |
| availability | status=active → in stock；else out of stock | 是 |
| condition | 默认 new | 是 |
| brand | products.brand（为空则用店铺名） | 否 |
| gtin | products.gtin（无则省略） | 否 |
| google_product_category | products.google_product_category（可留空由 GMC 自动） | 否 |
| product_type | products.product_type（分类路径） | 否 |
| age_group / gender | 可选默认值 | 否 |
| identifier_exists | 无 gtin 且无 brand+mpn 时置 FALSE | 否 |

Meta Catalog 字段：`id,title,description,availability,condition,price,link,image_link,brand,gtin,google_product_category,product_type,variant_group_id(可选分组变体)`（CSV，UTF-8，转义遵循 RFC4180）。

#### 5.9.3 无效商品与错误处理
- 生成前校验：缺 title/link/image_link/price 或缺库存 → 标记该行 error 不输出，写入 `errors`（含 sku+字段+原因）。
- 部分失败：`feed_runs.status=partial`，仍发布可输出行；`error_count == rows` 时视为失败不发布，保留上一次有效 Feed。
- 发布日志：每次 run 记录 rows / errors / 耗时；前端 Feed 详情展示"最近生成结果"。

#### 5.9.4 校验辅助
- 提供"在 GMC 注册用校验链接"：用户可用生成 URL 直接提交到 Google Merchant Center / Meta Catalog；文档提供接入检查清单（保证 icon 尺寸、价格格式、availability 值域）。

### 5.10 规则引擎设计

#### 5.10.1 规则结构（condition_spec / action_spec，均为 JSON）

```jsonc
{
  "id": 1, "name": "ROAS 低于阈值暂停",
  "scope": { "type": "campaign", "account_ids": [10], "object_ids": [] }, // 空=全选
  "window": { "days": 7, "type": "rolling" },                            // rolling 近7天
  "when": { "metric": "roas", "operator": "lt", "value": 1.5, "min_spend": 1000 },
  "cooldown_minutes": 1440,
  "action": { "type": "pause", "target": "self" }
}
```

#### 5.10.2 条件字段清单
| 字段 | 允许值 |
|------|--------|
| metric | spend / impressions / clicks / ctr / cpc / conversions / cpa / roas / cvr / budget_rate |
| operator | gt / gte / lt / lte / eq / neq |
| value | 数值（cpa/cpc/tr 用小数，金额用分） |
| min_spend | 可选；spend 低于此值不评估（避免小样本误触发） |
| window | {days:1|7|30, type: rolling|continuity_n(连续N天满足)} |
| budget_rate | 时间窗=当日（默认 09:00 后评估有效） |

#### 5.10.3 动作字段清单
| action.type | 参数 | 风险 | 说明 |
|-------------|------|------|------|
| pause | target(self / campaign_ids / ad_group_ids) | 高 | 暂停广告组或 Campaign |
| resume | 同上 | 高 | 恢复 |
| update_bid | new_value（绝对值或 +/-%）| 高 | 修改 ad group 出价 |
| update_budget | new_value（绝对值或 +/-%）| 高 | 修改 campaign 日预算 |
| notify | level(info/warning/critical), message 模板 | 低 | 站内 + 邮件（critical 必发邮件） |

#### 5.10.4 执行语义与防抖
- 评估时机：相关数据同步完成（success）后触发受 scope 影响的规则；每日 09:00 例行全量评估；手动 run-now。
- 在执行动作**前**一致性校验：重新拉取目标最新状态（若已是 paused，则 pause 动作跳过并记为 skipped）。
- 防抖：同一 (rule, target) 的 high 动作 24h 内（cooldown_minutes 可配）不重复执行；防抖依据存在 `rule_executions` 时间窗比对。
- 失败处理：目标平台调用失败 → 按 5.7.4 token 问题处理；动作失败不影响规则继续，记 execution.partial。
- 审计：每次评估与动作进入 `rule_executions` 与 `audit_logs`；高风险动作与 3.5.2 一致——规则触发的高风险动作**默认直接执行但必检一致性**（规则本身是用户已确认的自动化授权，不再二次确认；若租户关闭"规则自动执行"开关，则改为生成 pending 进确认中心）。

#### 5.10.5 内置规则模板
| 模板 | 条件 | 动作 |
|------|------|------|
| ROAS 守护 | 近 7 天 cpa>目标 且 spend≥阈值 | notify +（可选）pause |
| 预算超支预警 | 当日 budget_rate > 80%（warning）/ >100%（critical） | notify |
| 低效曝光止损 | 近 7 天 ctr<0.3% 且 clicks>100 | notify +（选配）pause |
| 平滑消耗提醒 | 当日 spend > 当日应消耗×(1+30%) | notify |

### 5.11 安全设计

#### 5.11.1 认证与会话
- 密码哈希：bcrypt cost=12；禁止明文存储。
- JWT：access token 15 分钟（`HS256`，密钥 = `JWT_SECRET`）；refresh token 7 天（RAM/DB 存 hash，支持吊销）；refresh 通过 `rotation`（每次刷新发新 refresh，旧 refresh 立即失效）。
- Cookie 策略：refresh 存放 httpOnly + SameSite=Lax + Secure（生产）Cookie；access 存内存/localStorage（前端按 4.4 统一封装）。
- 登录 / 注册 / 找回密码接口启用图形验证码（阶段一期用基础算术验证码，防刷注册）。

#### 5.11.2 敏感信息加密
- 主密钥：`ENCRYPTION_MASTER_KEY`（32 字节，Base64，经环境变量/密钥服务注入；生产建议用云 KMS 落盘加密）。
- 算法：AES-256-GCM，每条敏感记录（refresh_token、consumer key/secret、shopify token、developer token）独立随机 nonce；密文格式 `v1:{nonce}:{ciphertext}`。
- 日志与异常脱敏：服务层统一 `redact()` 过滤 secret/token 字段，禁止打印明文。

#### 5.11.3 Web 安全基线
- 传输：全站 HTTPS（Vercel 自动证书 + 后端 Caddy/Nginx 终止 TLS）。
- CORS 白名单：`https://orbit.elstella.com`（前端域名）；`Authorization` 与 `Idempotency-Key` 放行；`Content-Type: application/json` 允许。
- CSRF：同源校验 + SameSite Cookie + 写接口要求自定义头（`X-Requested-With`）。
- 输入校验：所有接口入参经 Gin binding + 自定义 validator；长度/枚举/类型白名单；防注入（GORM 参数化，禁用原生拼接；Feed 输出做 XML/CSV 转义防注入）。
- 限流：登录/注册 5 次/分钟/IP；OAuth 回调 10 次/分钟/IP；写接口（bid/budget/pause）20 次/分钟/用户；全局限流 100 次/分钟/用户；Redis 计数 + 429 + `Retry-After`。
- 越权防护：RBAC 中间件校验资源归属（3.1.2 + 2.4）；IDOR 用例纳入测试。
- 平台侧写操作安全：改价/改预算/暂停前二次校验目标状态与当前值（防陈旧覆盖），高风险动作永久保留在审计日志。

#### 5.11.4 漏洞与依赖
- 依赖漏洞扫描纳入 CI（Go vulncheck、npm audit）；上线前安全清单自查（OWASP Top10 抽查）。
- Prompt 注入缓解：外部内容（商品描述、用户输入）在注入 LLM system prompt 前做隔离标记；Agent 执行写操作必须经 Tool 白名单，禁止自由文本路由到写接口。

#### 5.11.5 可观测性
- 结构化 JSON 日志（`level, ts, service, request_id, tenant_id, user_id, action, duration_ms, code`）；审计日志独立表（见 5.3.1）。
- 指标（Prometheus）：HTTP 错误率、P95/P99 延迟、Asynq 队列深度与死信数、同步 job 成功率、外部平台 API 断路（熔断：连续 5 次失败开启半开 30s）。
- 告警：错误率 > 1%（5 分钟窗口）、同步失败率 > 5%、队列积压 > 1000、Token 刷新失败数 > 0（连续 3 次）→ 邮件/通知。

### 5.12 部署与配置清单

#### 5.12.1 后端 Docker 化
- 镜像：`orbit-api`（server）、`orbit-worker`（worker+scheduler 可合并镜像不同 entrypoint）。
- 编排：docker-compose 含 `postgres:15`、`redis:7`、`api`、`worker`、`scheduler`、`migrate`（一次性迁移 job）；健康检查 GET `/health/ready`。
- 发布：CI（GitHub Actions）→ 构建镜像推私有 registry → 服务器 compose pull && up -d；迁移在启动前 job 执行。

#### 5.12.2 环境变量清单（.env 示例，生产从密钥服务注入）

```
# 基础
DATABASE_URL=postgres://orbit:...@localhost:5432/orbit?sslmode=disable
REDIS_URL=redis://localhost:6379/0
APP_BASE_URL=https://orbit.elstella.com
API_BASE_URL=https://api.orbit.elstella.com
PORT=8080

# 认证与加密
JWT_SECRET=<32+ 随机字节 base64>
REFRESH_TOKEN_TTL=168h
ACCESS_TOKEN_TTL=15m
ENCRYPTION_MASTER_KEY=<32 字节 base64>
BCRYPT_COST=12

# 平台 OAuth
GOOGLE_OAUTH_CLIENT_ID=...
GOOGLE_OAUTH_CLIENT_SECRET=...
GOOGLE_DEVELOPER_TOKEN=...
META_APP_ID=...
META_APP_SECRET=...
BING_CLIENT_ID=...
BING_CLIENT_SECRET=...

# 电商平台（全局默认，注册时可按店铺单独配置）
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=...
SMTP_PASS=...
FEED_PUBLIC_BASE_URL=https://api.orbit.elstella.com

# AI
LLM_API_KEY=...
LLM_MODEL=...
LLM_TIMEOUT_MS=60000

# 运行
LOG_LEVEL=info
TZ=Asia/Shanghai
```

#### 5.12.3 前端（Vercel）
- 环境变量：`NEXT_PUBLIC_API_BASE_URL=https://api.orbit.elstella.com`、`NODE_ENV=production`。
- 注意：SSE 流式在 Vercel Serverless 需用 Edge Runtime 或 `streaming` 设置，确保长连接不超时；一期 Agent 对话接口直接走后端 SSE（不经 Next API 层代理）以避免超时。

#### 5.12.4 上线检查清单（部署方执行）
- [ ] 域名与证书（orbit.elstella.com / api.orbit.elstella.com）
- [ ] PG 备份策略（每日 pg_dump + WAL 归档 7 天）
- [ ] Redis AOF 持久化开启
- [ ] 密钥轮换预案与权限最小化
- [ ] 熔断/限流参数按容量预估配置
- [ ] 灰度开关：pawpel.com 先接入验证一周
- [ ] 监控与告警配置（错误率/队列/同步延迟）

### 5.13 测试与质量
- 单元测试：核心（指标口径计算、规则执行语义、token 加密解密、校验器）。
- 集成测试：Store/广告适配器 mock 平台 API；OAuth 流程（含回调 state 校验、刷新）。
- 多租户越权测试：租户 A 查询租户 B 资源断言 403（强制用例）。
- 端到端（Playwright）：注册→绑店铺→授权广告→看板数据→Agent 对话→确认动作全链路。
- CI 门槛：go vet / golangci-lint / go test 全绿；前端 typecheck + lint + build 通过；覆盖率目标 ≥ 70%（核心包）。
- 性能冒烟：核心读接口（/metrics/summary）压测 P99 < 500ms（1k 并发用户参考）；写接口幂等验证。

---

## 6. 非功能需求（细化）

| 项目 | 量化要求 | 度量方式 |
|------|----------|----------|
| 性能 | 核心读接口（指标/报表）响应 P50 < 200ms、P99 < 500ms；页面首屏 < 2s（Lighthouse 移动/桌面） | APM + Landing 压测 |
| 可用性 | 月可用性 ≥ 99.5%（允许约 3.6h/月 非计划停机）；关键路径无单点（Redis/PG 依赖监控） | 探针 /health/ready |
| 数据正确性 | 同步成功率 ≥ 99%（周维度）；指标口径 100% 符合 0.4 定义；Feed 可用率（GMC 校验通过率社区基线参考） | sync_jobs 统计 |
| 安全 | 见 5.11：OWASP Top10 自查 0 高风险；无明文密钥/令牌落盘；越权攻击测试 0 通过 | 安全扫描 + 渗透测试抽查 |
| 扩展性 | Store/广告适配器接口化，新增平台 ≤ 3 人日（含联调，参照现有 P0 适配器） | 架构评审 |
| 多租户 | 隔离一致性 100%（越权用例全绿）；租户级备份/删除可用 | 测试 + 运维演练 |
| 日志 | 审计日志 100% 覆盖写操作与 Agent 工具调用；保留 180 天 | 抽样核查 |
| 数据保留 | 指标明细保留 24 个月，之后可归档；审计日志 180 天；通知 90 天 | 运维策略 |
| 容灾 | 每日备份可恢复演练 1 次/月；WAL 归档 7 天 | 恢复演练记录 |

---

## 7. 一期交付标准（可上线 · 可勾选验收清单）

> 满足全部 □ 项视为达到可上线标准；每项对应验收人（前端 / 后端 / 测试 / 运维）。

- [ ] U1 用户可注册、登录、找回密码，注册即建租户（后端验收：注册→登录→JWT 下发链路）
- [ ] U2 支持邀请成员并分配四类角色，权限矩阵按 3.1.2 生效（验收：操作员/观察者访问越权资源返回 403）
- [ ] S1 可绑定多个 WooCommerce 与 Shopify 店铺，TestConnection 通过（验收：≥2 店铺，同步商品数 > 0）
- [ ] S2 商品（含变体）/订单/库存增量同步稳定，游标幂等（验收：连续 3 天同步成功率 100%，手动触发可刷新）
- [ ] A1 可授权 Google Ads + Meta 账户，回显示账户列表并完成绑定（验收：用测试账户走通 授权→回调→绑定→元数据同步）
- [ ] A2 Token 加密存储 + 自动刷新（验收：过期前自动刷新成功；强制 401 触发一次刷新链路）
- [ ] A3 广告数据定时同步 + 手动刷新，跨平台聚合报表按 0.4 口径正确展示（验收：数值与平台后台抽样核对一致）
- [ ] F1 商品 Feed 可生成（Google Shopping）+（Meta Catalog）并通过公开 URL 访问（验收：URL 在 GMC 提交校验通过、字段映射无错误）
- [ ] F2 Feed 定时 4h 更新 + 手动触发 + 生成日志可查（验收：手动触发 1 次，日志含 rows/errors）
- [ ] AG1 Agent 对话可用：支持 Sync/Report/Optimize 核心场景（验收：对话中"查询昨日花费+ROAS"返回正确数值）
- [ ] AG2 Optimize Agent 给出建议并支持"确认执行"真实调价/暂停（验收：pending → 确认 → 平台 API 调用成功 → 审计日志留痕）
- [ ] R1 基础规则可创建、启停、执行（验收：内置"预算超支预警"模板触发，通知送达）
- [ ] N1 站内通知 + 邮件送达关键事件（验收：同步失败/Token 过期告警均触达）
- [ ] P1 部署完成：前端 Vercel + 后端 Docker 自部署 + DB/Redis 就绪，探针通过（验收：/health/ready 200，域名证书有效）
- [ ] P2 可接入 pawpel.com 真实店铺与广告账户完成一轮全流程（验收：线上数据联动一轮，无阻断性 Bug）
- [ ] Q1 CI 全绿（单元/集成/E2E/安全越权用例通过），覆盖达标，无高危依赖漏洞

---

## 8. 开发与上线节奏建议（4 周，含里程碑验收）

> 建议团队：前端 2、后端 2、测试 1、运维/DevOps 0.5。关键依赖：W1 需确认 Google/Meta OAuth 应用审核；W2 需 GMC/Meta Catalog 校验账号。

### 里程碑一览

| 里程碑 | 时间 | 验收标准（对应第 7 章清单） |
|--------|------|------------------------------|
| M1 框架与授权 | 第 1 周末 | U1-U2、A1-A2 通过；CI 骨架运行 |
| M2 数据与 Feed | 第 2 周末 | S1-S2、A3、F1-F2 通过；报表看板展示真实数据 |
| M3 Agent 与规则 | 第 3 周末 | AG1-AG2、R1、N1 通过 |
| M4 上线 | 第 4 周末 | P1-P2、Q1 全绿；pawpel.com 灰度 1 周稳定 |

### 每周计划

| 周 | 任务 | 交付物 | 负责人（建议） |
|----|------|--------|----------------|
| W1 | Monorepo 初始化、Makefile、docker-compose、CI 骨架 | 可构建主分支 | DevOps |
| W1 | 用户/租户/RBAC（users/tenants/members）、JWT + 登录注册 | 认证模块 | 后端-甲 |
| W1 | Google OAuth（授权→回调→绑账户）、Token 加密存储 + 刷新 worker | OAuth 模块 | 后端-乙 |
| W1 | Meta 授权接入（long-lived token） | 同上 | 后端-乙 |
| W1 | 前端脚手架、登录/注册/工作台骨架、路由与顶栏店铺切换 | 前端框架 | 前端-甲 |
| W2 | WooCommerce 适配器（商品/订单/库存） | 电商适配器 | 后端-甲 |
| W2 | 风格：WooCommerce 入库 + 增量同步任务（Asynq） | 同步服务 | 后端-甲 |
| W2 | Shopify 适配器（GraphQL cursor） | 电商适配器 | 后端-乙 |
| W2 | 指标同步（Google/Meta/Bing 日指标）+ 统一指标表 | 指标服务 | 后端-乙 |
| W2 | Feed 生成（Google Shopping XML/TSV + Meta CSV）+ 公开 URL | Feed 服务 | 后端-甲 |
| W2 | 报表/看板前端（指标卡、趋势图、跨平台表格）+ 导出 | 前端-甲 |
| W3 | Agent 编排 + Tool 注册表 + SSE 流式对话 | Agent 服务 | 后端-乙 |
| W3 | 高风险动作确认流（agent_actions → 前端卡片 → 执行） | 确认流 | 后端-甲+前端-乙 |
| W3 | Optimize/Report/Creative Tool 实现 | Tool 实现 | 后端-乙 |
| W3 | 规则引擎（条件/动作/防抖/执行记录） | 规则服务 | 后端-甲 |
| W3 | 前端 Agent 对话页 + 动作卡片 + 规则管理页 | 前端-乙 |
| W4 | 多店铺完善、连接状态监控、通知（站内+邮件） | 通知服务 | 后端-甲 |
| W4 | 稳定性：重试/死信/熔断/限流调参、审计日志完善 | 稳定性 | 全部后端 |
| W4 | E2E 测试补齐、安全越权用例、压测 | 测试 | 测试 |
| W4 | 上线部署、灰度（pawpel.com）、监控告警配置 | 上线记录 | DevOps |

### 主要风险与依赖

| 风险 | 触发条件 | 影响 | 预防 / 应急 |
|------|----------|------|-------------|
| 平台应用审核周期 | Google/Meta OAuth App 审核未按时通过 | 授权不可用（阻断 M1） | W0 提前提交审核；应急：用测试账号 + 临时白名单，延期上线延后验收 |
| 平台 API 兼容变化 | Google v18/Meta v19/Bing 接口变更 | 同步失败 | 版本锁定（5.8.2）+ 变更监控；失败按降级与告警 |
| 数据量增长 | 多店铺/多账户同步过慢 | 队列积压 | worker 横向扩容、索引优化、分片任务 |
| 真实店铺权限受限 | pawpel.com 未开放 API 权限 | Feed/同步不完整 | 提前与客户确认 API 凭据与作用域 |
| LLM 成本与幻觉 | 长对话/复杂路径超预算、工具参数错误 | 误操作风险 | 高风险动作强制确认（硬约束）；请求预算上限告警 |

---

## 9. 文档说明

本文档为前后端一体的完整可执行需求（V2.0），在 V1.1 基础上补齐数据库设计、API 清单、Agent Tool、OAuth 流程、数据同步、Feed 规范、规则引擎、错误码、多租户隔离、安全、部署配置与里程碑验收。V1.1 评审结论见第 0 章。

后续仍建议跟进（非一期阻塞）：
- 详细接口文档（OpenAPI 导出，随代码同步）
- UI 高保真原型与交互说明（可在设计阶段产出）
- P1 平台（OpenCart / Shopyy / UeeShop）适配器规格
- 多语言文案库（英文）与移动端适配专题

---

**文档状态**：已达到开发团队可直接开工的可执行标准。建议以本文档为准组织开发排期，并在代码中同步维护 OpenAPI 与迁移脚本。


*（内容由AI生成，仅供参考）*
