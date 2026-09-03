---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: ccdb5f1afb85a6c9dc58ee10b06ac953_428a4d81a1c911f193c6525400f8a581
    ReservedCode1: 9QVIbgO6QJaCcYbNVfIh86UaZ+/P/eB6xvr/YOg/zbKHq6v/tLwSBsj1+TPnrX9y1qoUYamgKpUXEbqhL67kI0DNiYQU9/Veb5o64X8firIM2r4Doq7y1OazPbdTIT50Udz+0GsAaSdH+VF05cJtotu9XSBoOZ20Ur1XL2R5qLuwBJhamKBr25KPx9M=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: ccdb5f1afb85a6c9dc58ee10b06ac953_428a4d81a1c911f193c6525400f8a581
    ReservedCode2: 9QVIbgO6QJaCcYbNVfIh86UaZ+/P/eB6xvr/YOg/zbKHq6v/tLwSBsj1+TPnrX9y1qoUYamgKpUXEbqhL67kI0DNiYQU9/Veb5o64X8firIM2r4Doq7y1OazPbdTIT50Udz+0GsAaSdH+VF05cJtotu9XSBoOZ20Ur1XL2R5qLuwBJhamKBr25KPx9M=
---



# Orbit 广告智能助手系统
## 开发需求文档（V2.1 · 可执行版）

**版本**：V2.1
**上一版本**：V2.0
**日期**：2026年8月27日
**主域名**：orbit.elstella.com
**API 域名**：api.orbit.elstella.com
**文档定位**：本文档在 V2.0 基础上将第 4 章前端需求由概要规范升级为页面级可执行 UI/UX 规范（设计系统 design tokens、组件目录、逐页面区块/交互/三态/校验、图表选型、新用户引导流程、Agent 对话与工具卡片状态机、响应式与无障碍、前端验收清单），目标读者为开发团队（前端 / 后端 / 运维 / 测试），达到"开发团队拿起来即可开工、无需再自行定夺页面细节"的可执行标准。

---

## 0. 版本评审与本版修订说明

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

### 0.2 V2.1 修订说明（本次）

**V2.1 变更范围**：聚焦第 4 章前端需求扩容，由概要规范升级为**页面级可执行 UI/UX 规范**。第 4 章以外的业务、后端、非功能需求与验收章节未作实质改动，仅同步文档版本号、章节索引与第 7 章前端验收条目。

**新增/重构**：
- 4.1 页面架构与导航（路由表补鉴权/加载/菜单分组元数据）
- 4.2 设计系统（Design Tokens：色彩 / 字体 / 间距 / 圆角 / 阴影 / 边框 / 动效 / z-index 具体值）
- 4.3 组件目录（基础组件 + 业务组件清单、职责、页面复用矩阵）
- 4.4 页面级需求（认证 / OAuth 回调 / 工作台 / 店铺 / 广告账户 / 广告管理 / 广告详情 / 报表 / 看板 / Agent / 待确认动作 / 规则 / 商品 / Feed / 设置 / 通知共 17 页：区块布局、内容清单、交互、加载/错误/空态、必填校验）
- 4.5 图表选型与样式规范（Recharts，趋势 / 柱状 / 环形样式细则）
- 4.6 关键 UX 流程（新用户首次使用路径：注册 → 绑店铺 → 授权账户 → 建广告 → 看报表，分步指引设计）
- 4.7 Agent 对话与工具卡片完整状态机（SSE 流式渲染、工具卡片五状态 + 确认流、断线重连恢复）
- 4.8 响应式断点具体规则与无障碍（a11y）规范
- 4.9 数据请求与状态管理（V2.0 4.4 扩容）
- 4.10 前端验收清单

**V2.0 历史索引（保留）**：
- V2.0 新增：0（评审与口径）、2.4（多租户隔离实现）、5.3（数据库设计）、5.4（API 清单）、5.5（响应与错误码）、5.6（Agent 与 Tool）、5.7（OAuth）、5.8（数据同步任务）、5.9（Feed 规范）、5.10（规则引擎）、5.11（安全）、5.12（部署与配置）、5.13（测试与质量）
- V2.0 重构：3.2 / 3.3 / 3.4 / 4.1 / 6 / 7 / 8

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

## 4. 前端需求（Next.js · 页面级 UI/UX 可执行规范）

> 本章为前端 UI/UX 的**可执行规范**：设计系统、组件目录、逐页面需求、图表选型、关键 UX 流程、Agent 对话状态机、响应式与无障碍、数据请求规范、前端验收清单均已明确到可落地粒度。凡本章未覆盖的细节，以前端实现的合理默认值为准，但不得与本文档其他章节冲突。
>
> 技术基线（沿用 1.4 / 4.2）：Next.js App Router + TypeScript；样式 Tailwind CSS + shadcn/ui；图标 lucide-react；图表 Recharts；状态管理 React Query + Zustand；国际化 next-intl（一期仅中文，结构预留）。

### 4.1 页面架构与导航

#### 4.1.1 整体布局框架

- **标准布局**（除认证页外的所有业务页）：左侧导航栏 Sider（固定 240px，可折叠至 64px 图标栏）+ 顶部全局栏 Topbar（高 56px）+ 内容区（页面背景 surface-page，内边距 24px）。
- **认证页布局**（`/login` `/register` `/forgot-password`）：独立居中卡片布局，不显示 Sider 与 Topbar。
- **OAuth 回调页**（`/auth/callback`）：空白居中"处理中"布局。

**Topbar 内容（从左到右）**：
1. 页面标题面包屑（当前一级菜单 → 当前页面标题）；
2. 全局上下文切换区：店铺选择器（StoreSwitcher）+ 广告账户选择器（AccountSwitcher，级联联动），切换后触发 `contextVersion` 递增、全站数据视图刷新（见 4.9）；
3. 右侧操作区：Agent 待确认动作角标（Badge，接入 /agent/cards 的 pending 数）、通知铃铛（未读数角标）、当前成员/租户菜单（资料、设置、退出）。

**Sider 导航分组**（菜单高亮规则：当前路由前缀最长匹配）：
- 总览：工作台 `/`、数据报表 `/reports`、看板 `/reports/:id`
- 投放管理：广告管理 `/advertising/campaigns`、规则管理 `/rules`
- 资产：店铺 `/stores`、广告账户 `/accounts`、商品 `/products`、Feed `/feeds`
- AI 能力：Agent 对话 `/agent`、待确认动作 `/agent/cards`（带 pending 角标）
- 系统：系统设置 `/settings`、通知中心 `/notifications`

#### 4.1.2 路由表（V2.0 4.1 保留并补充元数据）

| 路由 | 页面 | 菜单分组 | 主要能力 | 权限 | 加载方式 |
|------|------|----------|----------|------|----------|
| /login、/register、/forgot-password | 认证 | —（独立布局） | 登录 / 注册 / 找回密码 | 公开 | SSR 跳转：已登录访问 → 回工作台 |
| /auth/callback?platform=... | OAuth 回调中转页 | —（独立布局） | 显示"授权处理中"，调 API 完成绑定 | 已登录 | Client 一次性任务 |
| /（工作台） | Dashboard | 总览 | 指标卡、快速入口、待确认动作横幅、同步状态 | 全部 | Server 壳 + Client 指标 |
| /stores | 店铺管理 | 资产 | 店铺列表、绑定向导、连接状态、同步触发 | 见 3.1.2 | Client |
| /stores/:id | 店铺详情 | 资产 | 商品 / 订单 / 同步历史、Feed 状态 | 同上 | Client |
| /accounts | 广告账户管理 | 资产 | 账户列表、授权入口、状态、Token 健康 | 同上 | Client |
| /advertising/campaigns | 广告管理 | 投放管理 | 统一视图列表（过滤/排序/搜索）、批量选择 | 同上 | Client |
| /advertising/campaigns/:id | 广告详情 | 投放管理 | 层级下钻（ad group/ad）、时间序列、Agent 优化入口 | 同上 | Client |
| /reports | 数据报表 | 总览 | 跨平台报表、维度选择、导出 | read | Client |
| /reports/:id | 看板详情 | 总览 | 自定义看板（卡片/趋势/表格） | read | Client |
| /agent | Agent 对话 | AI 能力 | 流式对话、工具调用卡片、确认/拒绝 | 见 3.1.2 | Client |
| /agent/cards | 待确认动作中心 | AI 能力 | 全部 pending/历史 action 列表 | approve 权限 | Client |
| /rules | 规则管理 | 投放管理 | 规则列表、创建向导（条件+动作表单）、启停 | 见 3.1.2 | Client |
| /products | 商品管理 | 资产 | 商品列表、变体、同步状态 | read / update | Client |
| /feeds | Feed 管理 | 资产 | Feed 列表、URL、手动触发、日志 | 见 3.1.2 | Client |
| /settings | 系统设置与权限 | 系统 | 租户资料、成员管理、角色、自动执行开关、通知偏好 | 见 3.1.2 | Client |
| /notifications | 通知中心 | 系统 | 通知列表、已读/未读、筛选 | 全部 | Client |

- 未登录访问任意受保护路由 → SSR 重定向到 `/login`，并携带 `redirect` 参数，登录后回跳。
- 无权限访问路由 → 渲染 403 禁止页（模板：无权限图标 + 文案 + "联系管理员"）。
- 所有 Client 页首屏：整页骨架（PageState Loading），数据到达后渐进渲染。

### 4.2 设计系统（Design Tokens）

> 所有 token 以 CSS 变量在 `:root` 定义，Tailwind theme 引用变量；组件与页面**禁止**硬编码色值/圆角/阴影。一期仅实现浅色主题，变量命名同时预留暗色（Phase 2）。

#### 4.2.1 色彩

**品牌色（Indigo 蓝为主色，用于主 CTA、选中态、品牌标识）**

| Token | 值 | 用途 |
|-------|-----|------|
| `--color-primary-50` | `#EEF2FF` | 主色浅底、选中背景 |
| `--color-primary-100` | `#E0E7FF` | 图标底色 |
| `--color-primary-500` | `#6366F1` | 主按钮背景、链接、focus 环 |
| `--color-primary-600` | `#4F46E5` | 主按钮 hover、active |
| `--color-primary-700` | `#4338CA` | 按下态 |
| `--color-primary-900` | `#312E81` | 文字强调（深色文字可用） |

**语义色（状态表达，必须结合图标/文案不单用颜色）**

| Token | 值 | 含义 |
|-------|-----|------|
| `--color-success-500` | `#16A34A` | 正常 / 成功 |
| `--color-warning-500` | `#D97706` | 预警（预算消耗>80%、Token 即将过期） |
| `--color-danger-500` | `#DC2626` | 异常 / 危险 / 失败 |
| `--color-info-500` | `#2563EB` | 信息（同步中、提示） |
| `--color-success-bg` | `#ECFDF5` | 成功浅底（徽标底色） |
| `--color-warning-bg` | `#FFFBEB` | 预警浅底 |
| `--color-danger-bg` | `#FEF2F2` | 异常浅底 |
| `--color-info-bg` | `#EFF6FF` | 信息浅底 |

**中性色（文字 / 边框 / 背景）**

| Token | 值 | 用途 |
|-------|-----|------|
| `--color-text-primary` | `#111827` | 主文字 |
| `--color-text-secondary` | `#4B5563` | 次级文字 |
| `--color-text-muted` | `#9CA3AF` | 弱化文字（说明/占位） |
| `--color-text-disabled` | `#D1D5DB` | 禁用文字 |
| `--color-border-default` | `#E5E7EB` | 默认边框、分隔线 |
| `--color-border-strong` | `#D1D5DB` | 高强度边框（输入框 hover） |
| `--color-surface` | `#FFFFFF` | 卡片 / 内容底 |
| `--color-surface-subtle` | `#F9FAFB` | 表格斑马纹、区块底 |
| `--color-surface-page` | `#F3F4F6` | 页面底色 |
| `--color-surface-hover` | `#F3F4F6` | 行 / 卡片 hover |
| `--color-focus-ring` | `rgba(99,102,241,.35)` | 键盘 focus 环 |

**数据可视化色板（图表专用，8 色循环）**

`#2563EB`（蓝）、`#10B981`（绿）、`#F59E0B`（琥珀）、`#EF4444`（红）、`#8B5CF6`（紫）、`#06B6D4`（青）、`#F97316`（橙）、`#14B8A6`（teal）。
语义映射：正增长 `#16A34A`、负增长 `#DC2626`；目标线 / 对比基线统一用 `#9CA3AF` 虚线 3px。

#### 4.2.2 字体

- **正文字体栈**：`-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", "Noto Sans SC", "Helvetica Neue", Arial, sans-serif`
- **数字 / 指标等宽栈**：`"SF Mono", "Roboto Mono", "JetBrains Mono", ui-monospace, monospace`（用于指标数值、金额、时间戳、SKU、ID）

| Token | 字号/行高/字重 | 用途 |
|-------|----------------|------|
| `--font-display-lg` | 32/40/700 | 页面大标题（仅 Dashboard 首屏） |
| `--font-title-lg` | 20/28/600 | 页面标题 |
| `--font-title-md` | 16/24/600 | 卡片标题、区块标题 |
| `--font-title-sm` | 14/20/600 | 列标题、分组标题 |
| `--font-body-md` | 14/22/400 | 正文默认 |
| `--font-body-sm` | 13/20/400 | 辅助说明、表格副文本 |
| `--font-body-xs` | 12/18/400 | 脚注、时间戳、caption |
| `--font-metric-lg` | 28/36/700 | 指标卡大数值 |
| `--font-metric-md` | 20/28/700 | 卡片次级数值 |
| `--font-button` | 14/20/500 | 按钮文字 |
| `--font-caption` | 12/16/400（色 muted） | 口径提示、单位说明 |

#### 4.2.3 间距

| Token | 值 | 用途 |
|-------|-----|------|
| `--space-1` | 4px | 图标与文字间隔 |
| `--space-2` | 8px | 紧凑元素间隔 |
| `--space-3` | 12px | 表单内左右留白 |
| `--space-4` | 16px | 组件间常规间隔、卡片间距 |
| `--space-6` | 24px | 页面内容区 padding、区块间距 |
| `--space-8` | 32px | 大区块间距 |
| `--space-10` | 40px | 页面顶部留白 |
| `--space-12` | 48px | 极端留白（引导页） |

固定尺寸约定：卡片内边距 20px；表格行高 44px；标准输入框/按钮高 36px（小 32px、大 44px）；表单纵向控件间距 20px、横向 12px；页面内容区 padding 24px、卡片间距 16px。

#### 4.2.4 圆角、阴影与边框

| Token | 值 | 用途 |
|-------|-----|------|
| `--radius-sm` | 4px | 输入框、按钮、标签 |
| `--radius-md` | 6px | 卡片、表格容器、气泡 |
| `--radius-lg` | 8px | 弹窗、下拉面板、Drawer |
| `--radius-full` | 999px | 徽标、头像、状态点 |
| `--shadow-xs` | `0 1px 2px rgba(16,24,40,.05)` | 卡片默认 |
| `--shadow-sm` | `0 1px 3px rgba(16,24,40,.10)` | 行 hover、小型浮层 |
| `--shadow-md` | `0 4px 8px -2px rgba(16,24,40,.10)` | Tooltip / Popover |
| `--shadow-lg` | `0 12px 24px -4px rgba(16,24,40,.15)` | Modal、Drawer、选中态卡片 |
| `--border-w` | `1px` | 常规边框线宽 |

> 约定：常态卡片用**边框**（`--color-border-default`）而非阴影；阴影仅用于浮层、弹出层与选中态，避免页面"浮起感"过重。

#### 4.2.5 动效与层级

- 动效时长：hover/按压 150ms `ease-out`；展开/折叠 200ms `ease`；浮层（Modal/Drawer）300ms `ease`（fade + translateY 8px）。遵循 `prefers-reduced-motion: reduce` 时全部关闭。
- z-index 层级：Sider `10`、Topbar `20`、Dropdown/Popover `30`、Modal `40`、Toast `50`、全局 Drawer `60`。
- 骨架屏使用 shimmer 动画（1.2s 循环），reduced-motion 下改为静态占位块。

### 4.3 组件目录（Component Inventory）

#### 4.3.1 基础组件（shadcn/ui 基底 + 统一规范封装）

| 组件 | 职责 / 约定 | 使用页面 |
|------|-------------|----------|
| Button / IconButton | 主/次/危险/幽灵四类 × 大中小；主按钮仅用于单一主操作 | 全部 |
| Input / Select / Radio / Checkbox / Switch / Textarea | 表单控件：统一高 36px、边框 `--color-border-default`、focus 环 `2px --color-focus-ring` | 认证/设置/规则/表单场景 |
| Tabs | 页面内分区切换，下划线式 | 店铺/账户/设置/通知/卡片中心 |
| Table | 通用数据表格：排序、列宽、行选择、固定表头、横向滚动（冻结首列）、分页或滚动加载 | 列表页通用 |
| Tooltip / Popover | 口径说明、悬浮详情、快捷操作 | 指标卡/列表 |
| Dialog / Drawer | 二次确认、表单抽屉、日志抽屉 | 全部 |
| Badge / Tag / Avatar | 状态徽标、平台标签、成员头像 | 全部 |
| Skeleton / Spinner / Toast | 加载骨架、加载圈、全局轻提示 | 全部 |
| EmptyState / ErrorState / LoadingState | 三态组件（见 4.4 通用规则） | 全部 |

#### 4.3.2 业务组件

| 组件 | 职责 | 使用页面 |
|------|------|----------|
| MetricCard | 指标卡：标题、`metric-lg` 数值、环比变化（涨跌色+箭头+tooltip 口径）、可选 sparkline | 工作台/广告详情/报表 |
| TrendChart | 趋势图封装（折线/面积，见 4.5） | 工作台/广告详情/报表/看板 |
| ComparisonBar | 对比柱状图封装（分组/堆叠，见 4.5） | 报表/广告详情 |
| DonutChart | 环形图封装（占比，见 4.5） | 报表/看板 |
| FilterBar | 统一筛选栏（时间范围/平台/店铺/账户/状态/搜索） | 列表页通用 |
| DateRangePicker | 时间范围：预设（今日/昨日/近7天/近30天/本月）+ 自定义；移动端横向滚动预设 | 列表/报表 |
| ConfirmDialog | 二次确认弹窗：影响预估描述 + 确认/取消；高风险强制 | 全部关键操作 |
| ErrorBanner | 错误横幅：错误码→文案映射（见 5.5）+ 重试 | 全部 |
| StatusBadge | 连接/同步/Token 健康状态（绿/黄/红 + 图标 + 文案） | 店铺/账户/商品 |
| StoreSwitcher / AccountSwitcher | 顶栏全局上下文切换（级联） | 全部（Topbar） |
| AgentBubble | 对话气泡（user/assistant/system） | Agent |
| ToolCard | Agent 工具调用卡片（状态机见 4.7） | Agent/待确认中心 |
| RuleCard | 规则卡片：启停开关、最近命中、生效态 | 规则管理 |
| SyncActivityRow | 同步历史行：状态/耗时/错误摘要 | 店铺详情 |
| NotificationItem | 通知条目：未读高亮、跳转链接 | 通知中心 |
| OnboardingStep | 新用户引导步骤组件（见 4.6） | 首次引导 |
| PageState | 页面级 LoadingPage/ErrorPage/EmptyPage | 全部 |

#### 4.3.3 页面 × 核心组件复用矩阵

| 页面 | Metric | Chart | FilterBar | Table | Confirm | StatusBadge | Switch | ToolCard | Onboard |
|------|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| 工作台 / | ✔ | ✔ | ✔ | — | — | ✔ | — | — | ✔ |
| 店铺 /stores | — | — | — | ✔ | ✔ | ✔ | ✔ | — | ✔ |
| 店铺详情 /stores/:id | ✔ | — | — | ✔ | ✔ | ✔ | — | — | — |
| 广告账户 /accounts | — | — | — | ✔ | ✔ | ✔ | — | — | ✔ |
| 广告管理 /campaigns | ✔ | — | ✔ | ✔ | ✔ | ✔ | — | — | — |
| 广告详情 /campaigns/:id | ✔ | ✔ | — | ✔ | ✔ | — | — | ✔ | — |
| 报表 /reports | ✔ | ✔ | ✔ | ✔ | — | — | — | — | ✔ |
| 看板 /reports/:id | ✔ | ✔ | — | ✔ | — | — | — | — | — |
| Agent /agent | — | — | — | — | ✔ | — | — | ✔ | — |
| 待确认 /agent/cards | — | — | — | ✔ | ✔ | — | — | ✔ | — |
| 规则 /rules | — | — | — | ✔ | ✔ | — | ✔ | — | — |
| 商品 /products | — | — | ✔ | ✔ | — | ✔ | — | — | — |
| Feed /feeds | — | — | — | ✔ | ✔ | ✔ | — | — | — |
| 设置 /settings | — | — | — | ✔ | ✔ | — | ✔ | — | — |
| 通知 /notifications | — | — | — | ✔ | — | — | — | — | — |

### 4.4 页面级需求（逐页区块 / 内容 / 交互 / 三态 / 校验）

#### 4.4.0 通用规则（所有页面适用）

- **三态统一**：Loading＝骨架屏（Skeleton，复用行结构，禁转圈整页）；Error＝ErrorState＋"重试"按钮（重试即 re-fetch）；Empty＝EmptyState＋引导文案＋操作入口（无账户引导授权、无数据引导触发同步等）。
- **接口错误**：统一映射错误码文案（见 5.5），经 ErrorBanner / 字段内联展示，前端不硬编码散落文案。
- **表单校验**：必填、格式（邮箱/URL/数值范围）前置校验；提交中按钮禁用（防重复提交）；失焦即校验；错误以 `aria-describedby` 关联。
- **高频刷新**：列表与报表提供手动刷新按钮（强制 invalidate，见 4.9）。

#### 4.4.1 认证组：/login、/register、/forgot-password

- **布局**：≥1280px 分栏（左品牌区：Logo＋产品名＋标语＋装饰插图；右 400px 表单卡）；<1280px 仅表单（居中，顶部 Logo）。
- **登录**：邮箱 + 密码 + "记住我" + "忘记密码"链接 + 主按钮。交互：提交中按钮 spinner；失败呈错误码文案；成功 → 回跳 `redirect` 或工作台。
- **注册**：邮箱 + 密码 + 确认密码 + 店铺类型（WooCommerce/Shopify/暂不确定）+ 主按钮"创建账户"。校验：邮箱格式；密码 ≥8 位含字母与数字；两次一致。成功 → 进入引导 Step1（见 4.6），并在顶部提示"账户已创建"。
- **找回密码**：邮箱 → 提交后展示"已发送重置链接"成功态（同页替换表单），并提供"返回登录"。
- 三态：登录/注册/提交按钮 loading；忘记密码无需 loading 长态。

#### 4.4.2 OAuth 回调中转页：/auth/callback?platform=...

- **布局**：居中卡片，平台图标 + "正在连接 Google Ads…"文案 + 进度动效。
- **交互**：进入即调后端完成绑定（见 5.7）；成功 → Toast "绑定成功" → 跳回来源页（`state` 带回 origin）或店铺详情；失败（授权取消/Token 失效/权限不足）→ ErrorBanner＋"重试"（重新发起授权）＋"返回"。
- **超时**：15s 无结果 → 提示"连接超时，请返回对应平台重新授权"。
- 三态：处理中 loading；fail 错误态；不适用空态。校验：无表单，但 `platform` 参数非法时展示 error 卡。

#### 4.4.3 工作台 Dashboard：`/`

- **区块**：
  1. 指标卡行（4 张 MetricCard：Spend / CTR / CPA / ROAS，含环比、口径 tooltip、sparkline），时间范围快速切换（今日/昨日/近7天）；
  2. 待确认动作横幅：存在 pending 时置顶显示"有 N 项操作待你确认"＋"去处理"按钮（角标同步更新）；
  3. 同步状态摘要卡：今日同步成功/失败/进行中、最近失败项（点击进店铺详情）；全部健康时显示绿色摘要行；
  4. 快速入口卡组：建广告（/advertising/campaigns）、看报表（/reports）、Ask Agent（/agent）、添加店铺（/stores）；
  5. 预警提醒卡：来自规则的预算超支预警等（黄色 info 卡，含时间与指标值）。
- **交互**：指标卡点击下钻 `/reports?dimension=...`；空账户态：显示引导卡"先授权你的第一个广告账户"→ /accounts。
- **三态**：loading＝4 张指标卡骨架；error＝整页重试；empty（无店铺/账户）＝引导授权。

#### 4.4.4 店铺管理：/stores

- **区块**：①状态筛选 Tabs（全部/正常/预警/异常/断开）；②店铺列表（Columns：店铺名+平台图标、连接状态 StatusBadge、商品/订单同步数量、最近同步时间、操作）；③右上"添加店铺"主按钮。
- **交互**：行点击 → 店铺详情；行内"立即同步"按钮（触发确认：影响范围商品/订单/全部）；断开/异常店铺禁用态 hover 显示原因 tooltip；"添加店铺"打开**绑定向导 Drawer**（Step1：平台选择 → Step2：字段表单 → Step3：TestConnection → Step4：成功回调）。
- **绑定向导表单校验**：店铺 URL（必填、URL 格式、可访问）、API Key / Secret（必填，WooCommerce）或 Access Token（必填，Shopify）；TestConnection 失败将平台原始错误映射为 5.5 文案展示。
- **三态**：列表 loading 骨架行；error 重试；empty＝空态插画＋"绑定第一个店铺"主按钮。

#### 4.4.5 店铺详情：/stores/:id

- **区块**：
  1. 概览卡：店铺名、平台、连接状态、最近同步时间、主要配置摘要；
  2. 同步控制卡：手动触发按钮（商品/订单/全部）、当前进行中 job 状态条；API 凭据"测试连接"；
  3. 商品摘要卡：商品总数 / 近24h新增 / 同步失败数（点击 → /products 带入筛选）；
  4. 订单摘要卡：近7天订单数、同步滞后（≤30min 健康）；
  5. Feed 状态卡：Feed URL、上次生成时间、校验状态（对应 /feeds）；
  6. 同步历史表：SyncActivityRow（时间、类型、状态、耗时、错误摘要，失败行可展开详情）。
- **交互**：手动同步 → ConfirmDialog（同步范围 + 影响提示）；同步中禁用按钮并显示进度；失败行展开错误详情（错误码文案）。
- **三态**：各卡独立骨架；错误卡独立重试；空态（无同步历史）＝"尚未触发同步"＋引导按钮。

#### 4.4.6 广告账户管理：/accounts

- **区块**：①平台 Tabs（全部/Google Ads/Meta/Bing）；②账户列表（Columns：平台图标、账户名、所属店铺、授权状态、Token 健康 StatusBadge、最近同步、操作）；③右上"授权新账户"主按钮。
- **交互**：授权流程＝平台选择 Dialog → 跳转平台 OAuth 授权（新窗口）→ 回跳绑定（同 4.4.2）；Token 即将过期（≤7 天，黄）hover 显示过期时间与"重新授权"；删除账户需 ConfirmDialog（提示将影响报表数据与已有关键操作）；重新授权复用同平台授权流程。
- **三态**：列表 loading 骨架；error 重试；empty＝"没有已授权的广告账户"＋引导授权。

#### 4.4.7 广告管理：/advertising/campaigns

- **区块**：
  1. FilterBar：时间范围、平台、店铺、账户、状态（全部/启用/暂停）、关键词搜索；
  2. 统计摘要条：所选范围的 spend/CTR/CPA/ROAS 合计（可点击查看口径 tooltip）；
  3. Campaign 表格（Columns：campaign 名+平台图标、状态、日预算、spend、clicks、conversions、CPA、ROAS、预算消耗率 badge——绿/黄/红）；
  4. 批量操作工具条：选择行后出现"批量暂停/批次恢复"（ConfirmDialog 列出受影响 N 条与影响预估，来自后端 estimate）；
  5. 导出 CSV。
- **交互**：行点击 → `/advertising/campaigns/:id`；列头排序；列配置（可隐藏列）；筛选变更仅刷新表格区域（保持 FilterBar 状态）；分页或滚动加载（>100 条）。
- **三态**：初次整页骨架、筛选变更仅表格骨架；error 重试；empty＝"没有匹配的广告系列，请清除筛选或先授权广告账户"。

#### 4.4.8 广告详情：/advertising/campaigns/:id

- **区块**：
  1. 概览指标卡行（campaign 级 spend/CTR/CPA/ROAS + 环比）；
  2. 时间序列趋势图（TrendChart，花费/收入 双线或 ROI 面积，时间范围切换）；
  3. 层级下钻：Campaign → Ad Group → Ad 树形导航，切换焦点后下方明细表联动；
  4. 维度明细表（当前层级行：名称、状态、spend、clicks、conversions、CPA、ROAS）；
  5. 操作区：暂停/恢复、改日预算（改预算表单一：数值 >0、上限提示）、"开启 Agent 优化"入口（→ /agent 带入 campaign 上下文）；
  6. Agent 推荐横幅：最近 Agent 建议数（→ /agent 或 /agent/cards）。
- **交互**：暂停/恢复/改预算一律 ConfirmDialog（含后端影响预估："将暂停 X 个广告组，预计影响日支出 Y 元"）；改预算校验数值合法区间并防抖保存。
- **三态**：区块独立骨架；错误独立重试；empty（无数据）＝引导触发同步。

#### 4.4.9 数据报表：/reports

- **区块**：
  1. FilterBar：时间范围、维度（日/周/月）、店铺、账户、平台、指标勾选（spend/clicks/CTR/CPA/conversions/ROAS）；
  2. 汇总指标卡行（按所选维度聚合）；
  3. 趋势图（TrendChart：维度时间序列，多平台分系列）；
  4. 对比柱状图（ComparisonBar：按平台对比 spend）；
  5. 明细表（维度行：平台/店铺/账户 各分层，含转化与 ROAS——标注"各平台归因口径不同，数值不可直接相加"，见 0.4）；
  6. 导出（CSV / PDF 按钮，生成后 Toast + 下载）。
- **交互**：维度切换重新请求；表格行点击下钻；导出格式选择 Dialog。
- **三态**：整页骨架；error 重试；empty＝"所选范围暂无数据"＋调整条件提示。

#### 4.4.10 看板详情：/reports/:id

- **区块**：看板网格（可拖拽卡片：MetricCard / TrendChart / DonutChart / Table / 文本备注）；编辑模式＋查看模式切换；卡片库清单面板（添加卡片）；右上"保存布局"。
- **交互**：编辑模式下拖拽/缩放/删除/新增卡片（react-grid-layout）；保存持久化（PUT /reports/:id/layout）；查看模式只读；支持复制看板链接。
- **三态**：整页骨架；error 重试；empty（无卡片）＝引导从卡片库添加。

#### 4.4.11 Agent 对话：/agent

- **区块**：
  1. 左侧会话列表（窄栏：历史会话、"新对话"按钮）；
  2. 对话主区：历史气泡列表（可滚动，底部吸附）+ ToolCard 串行插入 + 流式输入；
  3. 上下文条：当前店铺/账户（StoreSwitcher/AccountSwitcher 内嵌）+"本次对话上下文"标识；
  4. 输入区：多行文本（Enter 发送 / Shift+Enter 换行）、发送按钮、停止按钮（流式接收中可中止）。
- **交互**：详见图表 4.7（SSE 流式渲染、ToolCard 状态机、断线重连、确认/拒绝）。
- **三态**：首屏 loading（历史会话骨架）；empty＝空态引导文案（"试试问我：昨天花了多少、ROAS 如何？"＋示例问题 chips）；流中断 error＝重连提示，成功后自动补齐。

#### 4.4.12 待确认动作中心：/agent/cards

- **区块**：Tabs（待处理/已确认/已拒绝/全部）；动作列表（Columns：动作描述、涉及 campaign、风险等级 Badge、影响预估、发起 Agent、发起时间、操作）。
- **交互**：确认 → ConfirmDialog（提示"确认执行 {action}"＋影响预估；标注"撤销需手动恢复"）；拒绝 → 可选填原因（供审计）；支持低风险同类型批量确认；处理后的卡片移入对应历史 Tab。
- **三态**：列表骨架；error 重试；empty＝"暂无待确认动作"（待处理 Tab）等。

#### 4.4.13 规则管理：/rules

- **区块**：①规则列表（RuleCard：名称、条件摘要、动作摘要、启停 Switch、命中次数、最近触发时间）；②右上"创建规则"主按钮 → **创建向导 Drawer**（步骤 1 模板 or 自定义 → 步骤 2 条件表单 → 步骤 3 动作表单 → 步骤 4 预览确认）。
- **交互**：启停开关即时生效但弹 ConfirmDialog（"开启后 {规则} 将被自动执行"）；模板选择卡（预算超支预警/ROAS 阈值/CPA 阈值/平滑消耗）；编辑/复制/删除（删除需确认）。
- **条件表单校验**：指标必选、运算符必选、阈值必填（合法区间提示，如预算率 0-999%）、时间窗口必选；**动作表单校验**：动作类型必选、参数必填（如降低预算比例 1-90%）、抑制窗口（防抖）必选。
- **三态**：列表骨架；error 重试；empty＝"还没有规则"＋"创建第一条规则"＋模板推荐。

#### 4.4.14 商品管理：/products

- **区块**：FilterBar（店铺、搜索关键词、同步状态）+ 商品表格（Columns：图缩略、名称、SKU、价格、库存、变体数、同步状态 StatusBadge、最近同步时间）。
- **交互**：行展开查看变体（Variants 子表）；同步状态筛选；价格/库存列数值等宽字体。
- **三态**：骨架行；error 重试；empty＝"暂无商品"＋"前往店铺触发同步"。

#### 4.4.15 Feed 管理：/feeds

- **区块**：Feed 列表（Columns：店铺、平台（Google/Meta）、Feed URL（可复制+可打开）、格式、更新频率、上次生成时间、校验状态 StatusBadge、日志按钮）。
- **交互**："立即生成"按钮 → ConfirmDialog（生成时间提示）；日志按钮打开 **日志 Drawer**（按批次展示生成日志 rows/errors 详单，失败行可展开）。
- **三态**：列表骨架；error 重试；empty＝"尚无 Feed"＋引导生成说明。

#### 4.4.16 系统设置：/settings

- **区块**：Tabs（租户资料 / 成员与角色 / 通知偏好 / 自动执行 / 安全）。
  - 租户资料：租户名、Logo、默认展示币种、时区、联系信息表单（保存 Toast）。
  - 成员与角色：成员表格（姓名、邮箱、角色、状态、操作：改角色/移除）、"邀请成员"按钮（邮箱+角色表单；校验邮箱格式、必填）；移除成员 ConfirmDialog；角色说明 tooltip 引用 3.1.2。
  - 通知偏好：邮件/站内各事件类勾选（同步失败、Token 过期、预算预警、Agent 待确认、报表完成），保存即时。
  - 自动执行：租户级总开关"允许 Agent / 规则自动执行"（高风险开关，开启需 ConfirmDialog 提示风险）；细分：Rule 自动执行、Optimize 自动执行（各自独立开关）。
  - 安全：修改密码、会话管理（查看/注销会话）。
- **校验**：邮箱格式、必填；开关类操作保存后即时反馈 Toast。

#### 4.4.17 通知中心：/notifications

- **区块**：筛选（全部/未读）+ 通知列表（NotificationItem：类型图标、标题、摘要、时间、跳转链接、未读圆点）。
- **交互**：未读高亮；打开页面 2s 后批量标记已读（debounce，POST）；"全部已读"按钮；点击条目跳转关联页（同步详情/Agent 卡/店铺）；顶栏铃铛角标同步。
- **三态**：骨架；error 重试；empty＝"暂无通知"。

### 4.5 图表选型与样式规范

- **指定库**：**Recharts**（React 生态、Tree-shaking、Vercel 兼容）。当单图表数据点 >1000（长周期明细）时评估降采样，仍不足则 Eclipse 备选（Phase 2 再定，不在本期范围）。
- **图型选用**：时间序列→LineChart/AreaChart（趋势）；类别对比→BarChart（分组/堆叠）；占比→PieChart（donut 形态）；漏斗/留存→本期不做。
- **样式统一（Chart 组件级封装，所有页面一致）**：
  - 坐标轴文字 `--color-text-secondary` 12px；网格线 `#E5E7EB` 虚线；无外框。
  - Tooltip：白底 `--shadow-md`、`--radius-md`、padding 12px，数值 `--font-metric-md`。
  - 颜色：数据可视化色板（4.2.1）；正/负变化语义色；目标线/基线 `#9CA3AF` 虚线 3px。
  - 空数据：显示图表空态（占位 + "暂无数据"），不渲染空坐标系。
  - 高度约定：趋势 280px、环形 220px、柱状 260px；容器 width 100%，高度固定防跳动。
- **交互**：tooltip hover/键盘聚焦显示；图例点击切换系列显隐；跨度 >45 天自动按日降采样；数据更新 transition 200ms；入场动画 300ms；`prefers-reduced-motion` 时关闭动画。
- **无障碍**：每个图表提供"查看数据表格"切换（可视可访问数据表兜底，`aria-label`）；状态/系列区分不依赖单一颜色（图例＋线型/图案共同表达）。
- **指标卡 sparkline**：TrendChart 的缩略形态（inline 56×24，无坐标轴、无色值语义，仅趋势形状）。

### 4.6 关键 UX 流程：新用户首次使用路径（Onboarding）

> 流程基线（后端验收链路见 5.13 U）：**注册 → 绑店铺 → 授权账户 → 建广告 → 看报表**。

**总体设计**：注册完成后启动 Onboarding 引导（全屏分层引导或步骤式 Drawer），顶部 5 步进度指示器（完成 ✔ / 当前 ● / 未来 ○）；每步可"跳过"（跳过仍可在 Dashboard 空态重新唤起）；状态持久化于 `localStorage.onboardingState`，完成步骤可复查；任意步骤的主入口同时保留于对应管理页。

- **Step 1 绑定第一个店铺**：欢迎卡（"欢迎使用 Orbit"+ 一分钟搭建引导 → 绑定店铺）。表单＝平台选择（WooCommerce/Shopify）+ URL + 凭据 → TestConnection → 成功即勾选推进 Step 2。跳过 → 工作台空态后续唤起。
- **Step 2 授权广告账户**：平台选择卡片（Google Ads/Meta/Bing 图标 + 优势文案）→ 打开平台授权（新窗口）→ 回跳自动完成（见 4.4.2）→ 成功 Toast + 勾选。跳过同理。
- **Step 3 建广告（首个 campaign）**：检查当前上下文是否有 campaign：无 → 呈现两路引导——'让 Agent 帮你创建'（→ /agent 带入 "帮我为 {店铺} 创建第一个广告系列" 预设提示）或 '手动创建'（→ 打开新建 Dialog：系列名称、平台、目标、日预算、投放目标表单；校验名称必填、预算 > 0）；有 → 直接展示已有 campaign 摘要并勾选。允许跳过。
- **Step 4 看报表**：等待首次数据同步（后台异步），完成时经通知推送引导（"数据已就绪"）；引导卡提供"查看趋势图"（→ /reports）与"查看工作台"（→ /）。点击任一即完成并进入应用。
- **完成态**：显示"完成 🎉"（纯文本，不用 emoji 由文案替代）总结 + "进入工作台"主按钮；此后 30 天内 Dashboard 不再弹引导，可经设置页"重新查看引导"重置。
- **组件**：OnboardingStep（步骤图标、标题、描述、正文插槽、跳过/上一步/下一步、进度条）；移动端全屏化。

### 4.7 Agent 对话与工具卡片状态机

#### 4.7.1 消息角色与渲染

| 角色 | 渲染 |
|------|------|
| user | 右侧 primary 底气泡，按文本/Enter 换行 |
| assistant | 左侧 surface 气泡，**流式逐字追加**（含光标闪烁） |
| system | 居中 muted 文本（连接状态、错误提示） |
| tool | 以 **ToolCard** 渲染，插入对话流对应位置 |

#### 4.7.2 SSE 流式渲染协议（前端侧）

- 连接：`GET /api/v1/agent/conversations/:id/stream`（带 access token；SSE 事件表由 5.6/5.4 定义：`delta` / `message_id` / `tool_call` / `tool_result` / `action` / `done` / `error` / `ping`）。
- 渲染规则：`delta` 追加到当前 assistant 气泡；`message_id` 记录用于断线恢复；收到 `done` 停止光标；`ping` 每 ≤30s 一次判定心跳；无 `delta` 的帧跳过。
- 会话存储：`conversation_id + last_message_id` 维护于 Zustand 会话 store，用于恢复补齐。

#### 4.7.3 工具卡片状态机（核心技术点）

```
sending ──► running ──► succeeded ──► needs_approval ──(确认)──► executing ──► done
  │              │           │                 │◄─(拒绝)──► rejected
  │              │           └─(低风险自动)──► done
  │              ├──► failed ──►(重试)──► running（携带原入参）
  └──(组装失败/取消)──► failed
```

| 状态 | 卡片呈现 | 操作 |
|------|----------|------|
| `sending` | 工具名 + 灰态"组装中" | 无 |
| `running` | 旋转图标 + 工具名 + 已耗时 | 可中止（服务端取消） |
| `succeeded` | 结果摘要（如同步条数、查询数值） | 展开/收起 |
| `needs_approval` | 风险等级 Badge + 影响预估 + 完整入参摘要 + **确认 / 拒绝** 按钮 | 确认→ConfirmDialog→approved 进入 executing；拒绝→填原因→rejected |
| `executing` | 进度条/旋转 + "正在执行" | 无 |
| `done` | 绿色勾 + 最终结果摘要 + 时间戳；自动执行标注"Agent 已自动执行" | 展开 JSON 详情（可复制） |
| `failed` | 红色 + 错误码文案 + **重试** 按钮 | 重试（入参不变） |
| `rejected` | 灰色 + "已拒绝" + 原因 | 无 |

- 卡片字段：工具名、入参摘要（敏感字段脱敏，如凭证显示 `***`）、状态徽标、结果摘要、发起时间、操作按钮。高风险动作同时推送站内通知并出现在 4.4.12 待确认中心。
- 确认后乐观更新：本地置 `executing`，服务端最终结果经 SSE `action→done` 校正，不符则以服务端为准回滚。

#### 4.7.4 断线重连与恢复

- 心跳超时（>30s 无 ping）→ tooltip/横幅提示"连接中断，正在重连"。
- 指数退避重连：1s → 2s → 4s → 8s → 16s → 封顶 30s；页面可见性恢复（hidden→visible）立即重连。
- 恢复请求：`GET .../messages?after=last_message_id` 补齐缺失；`tool_cards` 状态以后端返回为准校正（running 是否已完成等）。

#### 4.7.5 输入区状态

| 状态 | 表现 | 可用操作 |
|------|------|----------|
| 空闲 | 正常输入框 | 输入/发送 |
| 输入中 | 正常 | 发送/换行 |
| 发送中 | 按钮 loading，防重复 | 无 |
| 流式接收中 | 按钮切换为"停止" | 停止（中止当前流） |
| 等待确认 | 输入框上方高亮提示条"Agent 正在等待你的确认" | 仍可输入，可切换会话 |

### 4.8 响应式断点与无障碍（a11y）

#### 4.8.1 响应式断点

| 断点 | 范围 | Sider | Topbar | 内容区 | 表单/交互 |
|------|------|-------|--------|--------|-----------|
| 桌面 XL | ≥1280px | 展开 240px | 完整（切换器/角标/菜单） | 完整多列 | 完整 |
| 平板 MD | 768–1279px | 折叠为 64px 图标栏（hover tooltip） | 压缩（搜索隐藏，切换器仅当前值＋弹层） | 表格横向滚动、FilterBar 换行 | 深度管理可用 |
| 手机 SM | <768px | 隐藏，汉堡按钮开 Drawer | 精简（仅角标与菜单） | 纵向单列、指标卡 2 列网格 | 触控目标 ≥44px 高 |

- **表格**：横向滚动 + 首列 sticky（`position: sticky; left:0`，背景 surface）。
- **指标卡**：SM 2 列网格、MD 4 列、XL 4 列。
- **FilterBar**：SM 纵向堆叠、日期预设横向滚动。
- **危险操作边界**：暂停/改价/删除/确认动作在移动端仍可用（不做深度隐藏），提示以全屏确认页呈现。
- 统一使用 Tailwind 断点（sm 640 / md 768 / lg 1024 / xl 1280 / 2xl 1536）。

#### 4.8.2 无障碍（a11y）

- **键盘可达**：全部功能 Tab 可达，focus 环 `2px --color-focus-ring` 可见；Esc 关闭浮层/弹窗；Enter/Space 触发按钮；Select/菜单支持方向键；拖拽看板提供非拖拽的备用排序方式。
- **语义**：landmark（`header`/`nav`/`main`）齐备；每页唯一 `h1`；所有表单 `label` 关联；icon-only 按钮带 `aria-label`；表格 `th scope`；状态不以颜色为唯一表达（图标＋文案）。
- **对比度**：正文 4.5:1、大字号 3:1、禁用组件 ≥3:1 且配合可见文案。
- **动效**：`prefers-reduced-motion: reduce` 关闭全部动画/自动滚动。
- **动态播报**：Toast/错误用 `aria-live="polite"`；表单错误 `aria-describedby`；流式对话主区 `aria-live="off"`（避免逐字干扰），以 `role="status"` 提示"新消息已到达/已完成"。
- **焦点管理**：Modal 打开聚焦内部、关闭返回触发器；Drawer 同理；路由切换聚焦 `main` 顶部。
- **落地验证**：axe / Lighthouse 无障碍 0 critical；核心流程全键盘可用（E2E 覆盖，见 4.10）。

### 4.9 数据请求与状态管理（V2.0 4.4 扩容）

- **API Client**：统一封装 fetch——baseURL（`NEXT_PUBLIC_API_BASE_URL`）、默认超时 15s、自动附加 access token；`401` → 跳登录（保留 redirect）；`403` → 静默 refresh 重试一次（见 5.7）；错误归一化为 `{ code, message }`（code 见 5.5 映射）；并发重复请求去重（相同 key 合并）。
- **React Query（Server State）**：query key 规范 `['v1','scope(store|account|campaign)','entity','params']`；`staleTime 30s`、`gcTime 5min`、`refetchOnWindowFocus` 关闭（防打断）；手动刷新按钮 `invalidateQueries` 强制失效；写操作（启停、同步）用乐观更新 + 失败回滚。
- **Zustand（Client State）**：全局上下文 store（`tenant / store / account`）；写入会话后 SSR 恢复；`contextVersion` 递增触发依赖上下文的数据刷新；SSE 事件写入（通知数、待确认数、同步状态）。
- **全局 SSE（notifications 与实时推送）**：`/api/v1/events` 长连接；事件类型 `notification` / `action_pending` / `sync_completed` 等；失败指数退避重连（同 4.7.4）；页面 `visibilitychange` 恢复重连；连接状态在 Topbar 以状态点呈现（绿/灰）。
- **错误处理**：根级 ErrorBoundary＋友好错误页；接口错误 → ErrorBanner（错误码→文案映射）；表单提交错误内联到字段。
- **SSR 策略**：登录态与静态页走 Server Components（含首屏 bundle 瘦身）；指标、列表、报表等动态数据 Client fetch + React Query 缓存；SEO 无关页面统一 Client 渲染。
- **权限前端壳**：菜单/按钮级权限指令基于 3.1.2 RBAC 过滤展示；后端接口仍强制鉴权（前端仅 UX 层，不作为安全边界）。

### 4.10 前端验收清单（可勾选）

**A. 设计系统（4.2）**
- [ ] A-1 token 全部以 CSS 变量落地，组件/页面零硬编码色值/圆角/阴影（code review 抽查）
- [ ] A-2 语义色三态在状态徽标统一生效（绿/黄/红 + 图标/文案）
- [ ] A-3 字体/间距/圆角/阴影与 4.2 表格逐项一致（视觉走查）

**B. 组件目录（4.3）**
- [ ] B-1 4.3.2 全部业务组件已实现并被对应页面复用（页面×组件矩阵核对）
- [ ] B-2 通用三态组件（Loading/Error/Empty）被所有页面采用，无手写散落状态

**C. 页面级（4.4）**
- [ ] C-1 17 个页面（含独立布局页）区块齐全、内容清单一致（对照 4.4 逐页）
- [ ] C-2 每页具备 Loading/Error/Empty 三态且可用
- [ ] C-3 表单必填/格式校验收敛、错误码文案统一（无前端散落文案）
- [ ] C-4 高风险操作（暂停/改价/删除/动作确认/自动执行开关）均走 ConfirmDialog

**D. 图表（4.5）**
- [ ] D-1 图表统一 Recharts 封装 + 样式规范一致（坐标轴/网格/tooltip/色板）
- [ ] D-2 空数据态、降采样、图例显隐、数据表兜底均可用

**E. Onboarding（4.6)**
- [ ] E-1 引导全流程可走通：注册→绑店铺→授权账户→建广告→看报表
- [ ] E-2 每步可跳过、可复查、完成态正确、30 天不重复弹出

**F. Agent 与工具卡片（4.7）**
- [ ] F-1 SSE 流式逐字渲染、停止/续传可用
- [ ] F-2 ToolCard 全部状态（sending/running/succeeded/needs_approval/executing/done/failed/rejected）呈现与转换正确
- [ ] F-3 断线重连指数退避 + 按 message_id 补齐 + 状态校正
- [ ] F-4 确认/拒绝操作闭合（含审计原因），与 /agent/cards 联动

**G. 响应式与 a11y（4.8）**
- [ ] G-1 三档断点布局符合 4.8.1（走查 1280/1024/390 三档）
- [ ] G-2 表格横向滚动 + 首列冻结可用
- [ ] G-3 axe/Lighthouse 障碍 0 critical；全站键盘可达（E2E）

**H. 数据请求（4.9）**
- [ ] H-1 全局上下文切换全站刷新（contextVersion 生效）
- [ ] H-2 403 refresh 重试、401 跳登录、错误码映射正确
- [ ] H-3 全局 SSE 通知/待确认数实时更新，重连逻辑可验证

**I. 性能与 E2E**
- [ ] I-1 页面首屏 <2s（Lighthouse 桌面/移动，见第 6 章）
- [ ] I-2 前端 typecheck / lint / build（生产）通过
- [ ] I-3 Playwright E2E：登录、店铺/账户切换、报表、Agent 确认动作关键链路全绿（见 5.13）

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
- [ ] FE1 第 4 章设计系统（design tokens）全部落地，语义色/字体/间距/圆角/阴影按 4.2 规定值实现，组件零硬编码（前端验收，详见 4.10）
- [ ] FE2 顶部全局栏 / 左侧导航 / 内容区布局在各断点正确渲染，店铺/账户切换全局刷新（见 4.1 / 4.8 / 4.9）
- [ ] FE3 每个核心页面具备加载骨架、错误重试、空态引导三个状态（见 4.4）
- [ ] FE4 Agent 对话流式渲染，工具卡片状态机各状态齐全（见 4.7），支持断线重连恢复
- [ ] FE5 新用户引导全流程可走通：注册→绑店铺→授权账户→建广告→看报表（见 4.6，对照 E2E）
- [ ] FE6 a11y 验收通过：全站键盘可操作、对比度达标、表单有 aria 标注（见 4.8，axe 0 critical）
- [ ] FE7 前端 E2E（Playwright）全绿：登录、店铺/账户切换、报表、Agent 确认动作链路（见 5.13 / 4.10）
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

本文档为前后端一体的完整可执行需求（V2.1），在 V2.0 基础上将第 4 章前端需求升级为页面级可执行 UI/UX 规范（设计系统、组件目录、逐页面区块/交互/三态/校验、图表选型、新用户引导、Agent 对话与工具卡片状态机、响应式与无障碍、前端验收清单）。第 4 章之外的业务、后端与非功能需求未作实质改动，各版本评审结论见第 0 章。

后续仍建议跟进（非一期阻塞）：
- 详细接口文档（OpenAPI 导出，随代码同步）
- UI 高保真原型（可按第 4 章逐页面区块直接出稿，不再需要信息架构层面的返工）
- P1 平台（OpenCart / Shopyy / UeeShop）适配器规格
- 多语言文案库（英文）与移动端专项适配验证

---

**文档状态**：已达到开发团队可直接开工的可执行标准，前端页面细节无需再自行定夺。建议以本文档为准组织开发排期，并在代码中同步维护 OpenAPI、设计令牌与迁移脚本。


*（内容由AI生成，仅供参考）*
*（内容由AI生成，仅供参考）*
