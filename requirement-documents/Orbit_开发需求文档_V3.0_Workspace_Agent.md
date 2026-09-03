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



# Orbit 广告智能中枢
## 技术开发需求文档（V3.0 · Workspace + Agent 实现版）

**版本**：V3.0
**上一版本**：V2.1
**日期**：2026年8月29日
**文档定位**：本文档在 V2.1 基础上新增 Workspace 多工作空间隔离体系、统一 Provider 配置体系、可插拔 Agent 架构的完整技术实现规范。目标读者为开发团队（前端 / 后端 / 运维 / 测试），达到"拿起来即可开工"的可执行标准。

---

## 0. 版本修订说明

### 0.1 V2.1 主要问题与本版处理

| 编号 | 类别 | 问题描述 | 本版处理 |
|------|------|----------|----------|
| R01 | 架构缺失 | 无 Workspace 概念，无法实现多项目隔离 | 新增 2.5 Workspace 架构 |
| R02 | 配置缺失 | OAuth/LLM/SMTP 凭证硬编码在 .env，无法在线配置 | 新增 5.14 Provider 统一配置体系 |
| R03 | Agent 缺失 | Agent 无真实 LLM 集成，无工具执行能力 | 新增 5.15 Agent 可插拔架构 |
| R04 | 缺失 | 无调用统计与成本分摊机制 | 新增 5.16 使用计量与限额 |
| R05 | 缺失 | 无前端 Provider 配置界面 | 新增 4.11 Provider 配置页面 |
| R06 | 缺失 | 无 Workspace 管理界面 | 新增 4.12 Workspace 管理页面 |

### 0.2 范围与边界

**V3.0 新增范围内：**

1. Workspace 多工作空间隔离体系
2. Provider 统一配置体系（LLM / SMTP / 广告平台 / 电商平台）
3. 可插拔 Agent 架构（Builtin + Remote 双实现）
4. Agent 工具注册表与执行引擎
5. 使用计量与限额管理
6. Provider 配置前端界面
7. Workspace 管理前端界面

**V3.0 明确不在范围内：**

1. 视频生成接口（已预留，待第三方服务接入）
2. Agent 独立部署为微服务（架构已预留，V3.1 实现）
3. 多语言支持（仍为二期）
4. 移动端 App（仍为二期）

---

## 1. 项目概述

### 1.1 产品名称
Orbit（广告智能中枢）

### 1.2 产品定位
AI Agent 驱动的多平台广告自动化与智能化管理平台，支持 **Google Ads**、**Meta**、**Microsoft Advertising**，深度对接 **WooCommerce** 与 **Shopify**，服务多租户、多工作空间、多店铺、多广告账户。以 Agent 对话为核心交互，支持多团队、多项目独立配置与计量。

### 1.3 核心目标
- 支持 Workspace 多工作空间隔离，不同项目组独立管理
- 统一 Provider 配置体系，LLM/SMTP/广告平台凭证在线管理
- 可插拔 Agent 架构，支持 Builtin 与 Remote 双模式
- 按 Workspace/Tenant 维度的使用统计与成本分摊
- 保留 `.env` 作为 System Default fallback，保证向后兼容

---

## 2. 系统架构

### 2.1 整体架构（V3.0 更新）

```
用户浏览器
    ↓
前端（Next.js @ Vercel）
    ↕  HTTPS  REST  /api/v1  +  SSE
后端 API（Golang Gin）
    ├─ API 服务（REST + SSE）
    ├─ Agent 服务（可插拔：Builtin / Remote）
    │   ├─ LLM Provider（Agnes / OpenAI / Anthropic / 本地）
    │   ├─ Tool 注册表与执行引擎
    │   ├─ 对话记忆管理
    │   └─ 人工确认流
    ├─ Worker（Asynq 任务）
    ├─ Scheduler（cron）
    └─ 共享层
        ├─ Provider 配置服务（统一凭证管理）
        ├─ Workspace 服务（多租户隔离）
        ├─ 适配器（Store / AdPlatform / Affiliate）
        └─ Intelligence 模块（受众分析 / 策略 / 出价 / 素材）
    ↓
PostgreSQL（主数据 + 会话 + 审计）
Redis（缓存 + 队列 + 限流）
    ↓
外部服务
    ├─ 广告平台 API（Google / Meta / Bing）
    ├─ 电商平台 API（WooCommerce / Shopify）
    ├─ LLM API（Agnes / OpenAI / Anthropic）
    └─ SMTP 服务（SendGrid / AWS SES / 企业邮箱）
```

### 2.2 Workspace 架构（核心新增）

```
Workspace（工作空间）
├── 属性：name, owner_user_id, plan, settings, billing
├── 关联：Tenants（一个 Workspace 可包含多个 Tenant）
├── 配置：Providers（workspace_id = NULL 为全局共享）
│   ├── LLM Provider（模型 / API Key / Base URL）
│   ├── SMTP Provider（发件服务 / 域名 / 凭证）
│   ├── Google Ads Provider（OAuth / Developer Token）
│   ├── Meta Provider（App ID / Secret / Access Token）
│   └── Bing Provider（Client ID / Secret / Access Token）
└── 计量：usage_tokens, usage_requests, limits

Tenant（租户）
├── 关联：Workspace（workspace_id）
├── 成员：Users + Roles（RBAC）
├── 资源：Stores + AdAccounts + Campaigns + Feeds + Rules
└── 配置：继承 Workspace Provider，可独立覆盖

User（用户）
├── 可属于多个 Tenant（每次会话选一个"当前租户"）
└── 每个 Tenant 内有独立 Role
```

### 2.3 Provider 配置优先级

```
请求 LLM/SMTP/平台配置时，按以下顺序解析：

1. Tenant 专属 Provider（tenant_id = X, workspace_id = Y）
   ↓ 未找到
2. Workspace 级别 Provider（workspace_id = Y, tenant_id = NULL）
   ↓ 未找到
3. 全局共享 Provider（workspace_id = NULL, tenant_id = NULL）
   ↓ 未找到
4. System Default（.env 环境变量，向后兼容）
```

### 2.4 Agent 可插拔架构

```
AgentBackend 接口（port.go）
├── Chat(ctx, req) → (*ChatResponse, error)
├── Stream(ctx, req) → (<-chan StreamEvent, error)
└── Tools(ctx) → ([]ToolDefinition, error)

实现 A：BuiltinAgent（当前进程内）
├── LLM：AgnesLLM / OpenAILLM / LocalLLM
├── Memory：PostgresMemory / RedisMemory
├── ToolRunner：DefaultToolRunner
└── 适合：单进程部署，低延迟

实现 B：RemoteAgent（V3.1 独立服务）
├── 通信：HTTP / gRPC
├── 协议：JSON-RPC / Protobuf
└── 适合：多 Agent 协作，独立扩容

业务代码只依赖 AgentBackend 接口，切换实现无需改动上层逻辑。
```

---

## 3. 功能需求

### 3.1 Workspace 管理

#### 3.1.1 Workspace 模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| name | varchar(255) | 工作空间名称 |
| owner_user_id | bigint | 创建者/所有者 |
| plan | varchar(20) | 套餐：beta / pro / enterprise |
| settings | jsonb | 自定义配置（时区、语言、品牌等） |
| billing | jsonb | 计费信息（额度、用量、账单） |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

#### 3.1.2 Workspace 权限

- 仅 Workspace Owner 可删除/转让 Workspace
- Workspace Admin 可管理成员和 Provider
- Workspace Member 可继承 Provider 配置使用资源

#### 3.1.3 Workspace 隔离

- 所有数据查询必须携带 `workspace_id` 或通过 `tenant_id` 关联到 workspace
- Provider 配置按 workspace 隔离，不同 workspace 可使用不同 LLM/SMTP
- 使用统计按 workspace 维度汇总

### 3.2 Provider 统一配置体系

#### 3.2.1 Provider 类型

| Type | 说明 | 配置字段 |
|------|------|----------|
| llm | 大语言模型 | api_key, base_url, model, timeout_ms, max_tokens |
| smtp | 邮件服务 | host, port, user, pass, from_email, from_name, encryption |
| google | Google Ads | client_id, client_secret, developer_token, access_token, refresh_token, scopes |
| meta | Meta Ads | app_id, app_secret, access_token, expires_at |
| bing | Bing Ads | client_id, client_secret, access_token, refresh_token |
| woocommerce | WooCommerce | base_url, api_key, api_secret |
| shopify | Shopify | base_url, access_token, api_version |

#### 3.2.2 Provider 状态机

```
pending → active → error → active（自动恢复）
               ↓
            expired → refreshing → active（刷新成功）
                              ↓
                           error（刷新失败，告警）
```

#### 3.2.3 Provider 操作

- **创建**：管理员在设置页面新增 Provider，填写配置后在线测试连通性
- **测试**：调用平台 `/me` 或 `/test` 接口验证凭证有效性
- **编辑**：修改配置后即时生效，无需重启
- **删除**：软删除，关联的 AdAccount/Token 标记为 orphan
- **共享**：全局 Provider（workspace_id = NULL）所有 Workspace 可用

### 3.3 Agent 可插拔架构

#### 3.3.1 AgentBackend 接口

```go
// internal/agent/port.go
type AgentBackend interface {
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
    Stream(ctx context.Context, req ChatRequest) (<-chan StreamEvent, error)
    Tools(ctx context.Context) ([]ToolDefinition, error)
}

type ChatRequest struct {
    ConversationID uint64
    Message        string
    Context        map[string]interface{}
    WorkspaceID    uint64
    TenantID       uint64
}

type ChatResponse struct {
    Messages []AgentMessage
    Actions  []PendingAction
}

type PendingAction struct {
    Type     string
    TargetID uint64
    Params   map[string]interface{}
    Risk     string
    Preview  interface{}
}
```

#### 3.3.2 BuiltinAgent 实现

```go
// internal/agent/builtin/builtin.go
type BuiltinAgent struct {
    llm        LLMProvider
    memory     ConversationMemory
    toolRunner ToolRunner
    registry   *ToolRegistry
}

func NewBuiltinAgent(cfg *config.Config) *BuiltinAgent {
    llm := NewAgnesLLM(cfg)  // 默认使用 Agnes，可配置
    return &BuiltinAgent{
        llm:        llm,
        memory:     NewPostgresMemory(),
        toolRunner: NewDefaultToolRunner(),
        registry:   NewToolRegistry(),
    }
}
```

#### 3.3.3 Agent 工具注册表

```go
// internal/agent/tools/registry.go
type ToolDefinition struct {
    Name        string
    Description string
    Parameters  map[string]interface{}
    Risk        string
    Handler     ToolHandler
}

var BuiltinTools = []ToolDefinition{
    {Name: "get_metrics",         Handler: GetMetricsTool, Risk: "low"},
    {Name: "get_campaigns",       Handler: GetCampaignsTool, Risk: "low"},
    {Name: "pause_campaign",      Handler: PauseCampaignTool, Risk: "medium"},
    {Name: "resume_campaign",     Handler: ResumeCampaignTool, Risk: "medium"},
    {Name: "update_budget",       Handler: UpdateBudgetTool, Risk: "medium"},
    {Name: "generate_creative",   Handler: GenerateCreativeTool, Risk: "low"},
    {Name: "analyze_audience",    Handler: AnalyzeAudienceTool, Risk: "low"},
    {Name: "optimize_bidding",    Handler: OptimizeBiddingTool, Risk: "medium"},
    {Name: "create_campaign",     Handler: CreateCampaignTool, Risk: "high"},
    {Name: "send_notification",   Handler: SendNotificationTool, Risk: "low"},
}
```

#### 3.3.4 Agent 编排流程

```
用户消息
    ↓
Ops Agent（LLM，无工具）：意图识别 + 参数提取
    ↓
路由到目标 Agent / Tool
    ├─ 低风险（low）→ 直接执行 → 返回结果
    ├─ 中风险（medium）→ 执行 + 审计日志 → 返回结果
    └─ 高风险（high）→ 创建 PendingAction → 等待人工确认
            ↓
        用户确认 → 执行 → 审计日志
        用户拒绝 → 记录拒绝原因
```

### 3.4 使用计量与限额

#### 3.4.1 计量维度

| 维度 | 说明 | 存储位置 |
|------|------|----------|
| 按 Workspace | 总调用量、Token 消耗、API 调用次数 | providers.meta.usage |
| 按 Tenant | 继承 Workspace 额度，可单独设置上限 | 预留字段 |
| 按 LLM 模型 | 不同模型的 Token 消耗分别统计 | providers.meta.usage_by_model |
| 按时间周期 | 月度/日度/小时级统计 | providers.meta.usage_periods |

#### 3.4.2 限额策略

- **软限额**：达到 80% 额度时，Agent 回复中提示"本月 LLM 用量已用 80%"
- **硬限额**：达到 100% 额度时，LLM 调用降级为规则引擎模式
- **超额**：支持紧急提升限额（需 Workspace Owner 确认）

---

## 4. 数据库设计

### 4.1 新增 Workspace 相关表

```sql
-- 工作空间表
CREATE TABLE workspaces (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_user_id BIGINT NOT NULL,
    plan VARCHAR(20) NOT NULL DEFAULT 'beta',
    settings JSONB,
    billing JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workspaces_owner ON workspaces(owner_user_id);

-- 租户关联工作空间
ALTER TABLE tenants ADD COLUMN workspace_id BIGINT;
CREATE INDEX idx_tenants_workspace ON tenants(workspace_id);

-- 统一 Provider 表
CREATE TABLE providers (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL,           -- llm, smtp, google, meta, bing, woocommerce, shopify
    name VARCHAR(255) NOT NULL,           -- 显示名称
    workspace_id BIGINT,                  -- NULL=全局共享
    tenant_id BIGINT,                     -- 冗余加速查询
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    config JSONB NOT NULL,                -- 加密存储凭证
    meta JSONB,                           -- 使用统计、偏好设置
    last_test_at TIMESTAMP,
    error_msg TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_providers_type_workspace ON providers(type, workspace_id);
CREATE INDEX idx_providers_tenant ON providers(tenant_id);

-- 使用计量表
CREATE TABLE provider_usage (
    id BIGSERIAL PRIMARY KEY,
    provider_id BIGINT NOT NULL,
    workspace_id BIGINT NOT NULL,
    tenant_id BIGINT,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    tokens_used BIGINT NOT NULL DEFAULT 0,
    requests_count INT NOT NULL DEFAULT 0,
    errors_count INT NOT NULL DEFAULT 0,
    cost_cents BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_provider_usage_provider ON provider_usage(provider_id, period_start, period_end);
CREATE INDEX idx_provider_usage_workspace ON provider_usage(workspace_id, period_start, period_end);
```

### 4.2 现有表结构调整

```sql
-- tenants 表新增 workspace_id
ALTER TABLE tenants ADD COLUMN workspace_id BIGINT;
ALTER TABLE tenants ADD CONSTRAINT fk_tenants_workspace 
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id);

-- platform_conns 表增加 provider_id（可选，关联到统一 provider）
ALTER TABLE platform_conns ADD COLUMN provider_id BIGINT;
ALTER TABLE platform_conns ADD CONSTRAINT fk_platform_conns_provider 
    FOREIGN KEY (provider_id) REFERENCES providers(id);
```

### 4.3 数据迁移策略

1. **Phase 1**：创建新表，保留 `.env` 配置不变
2. **Phase 2**：将 `.env` 中的全局凭证迁移到 `providers` 表（workspace_id = NULL）
3. **Phase 3**：前端配置界面支持在线编辑 Provider
4. **Phase 4**：`.env` 中的 LLM/SMTP 配置变为可选 fallback

---

## 5. API 设计

### 5.1 Workspace API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/workspaces | 列出用户所属工作空间 |
| POST | /api/v1/workspaces | 创建工作空间 |
| GET | /api/v1/workspaces/:id | 工作空间详情 |
| PATCH | /api/v1/workspaces/:id | 更新工作空间 |
| DELETE | /api/v1/workspaces/:id | 删除工作空间 |
| POST | /api/v1/workspaces/:id/switch | 切换当前工作空间 |

### 5.2 Provider API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/providers | 列出当前工作空间的 Providers |
| POST | /api/v1/providers | 创建 Provider |
| GET | /api/v1/providers/:id | Provider 详情（脱敏） |
| PATCH | /api/v1/providers/:id | 更新 Provider |
| DELETE | /api/v1/providers/:id | 删除 Provider |
| POST | /api/v1/providers/:id/test | 测试连通性 |
| GET | /api/v1/providers/:id/usage | 使用统计 |

### 5.3 Agent API（现有，保留）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/agent/chat | Agent 对话（SSE 流式） |
| GET | /api/v1/agent/conversations | 会话列表 |
| GET | /api/v1/agent/conversations/:id/messages | 消息历史 |
| GET | /api/v1/agent/actions | 待确认动作 |
| POST | /api/v1/agent/actions/:id/approve | 确认动作 |
| POST | /api/v1/agent/actions/:id/reject | 拒绝动作 |

### 5.4 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "trace_id": "uuid"
}
```

---

## 6. 前端设计

### 6.1 Workspace 管理页面

**路径**：`/settings/workspaces`

**功能**：
- 工作空间列表（卡片式，显示名称、成员数、套餐、用量）
- 创建工作空间（模态框：名称、套餐、初始成员）
- 切换工作空间（顶部导航下拉）
- 工作空间详情（成员管理、计费信息、使用限额）

### 6.2 Provider 配置页面

**路径**：`/settings/providers`

**功能**：
- Provider 类型分组展示（LLM / SMTP / 广告平台 / 电商平台）
- 每个 Provider 卡片显示：名称、类型、状态、最后测试时间、用量统计
- 创建/编辑 Provider（表单：类型特定字段 + 在线测试按钮）
- 删除 Provider（二次确认）
- 全局 Provider 与 Workspace Provider 区分标识

### 6.3 Agent 对话页面（现有，增强）

**路径**：`/agent`

**增强**：
- 工作空间上下文感知（自动使用当前 workspace 的 Provider）
- 工具调用时显示使用的 Provider（如"使用 Agnes AI 生成文案"）
- 用量提示（如"本月 LLM 用量：450K / 1M tokens"）

---

## 7. 非功能需求

### 7.1 安全性

- Provider 凭证加密存储（AES-256-GCM，key 来自 `.env`）
- 日志脱敏：不打印 api_key / access_token / refresh_token
- 接口权限：Provider 配置仅 Workspace Admin 可访问
- 在线测试连通性时临时使用凭证，不持久化返回的 token

### 7.2 可用性

- Provider 配置变更即时生效，无需重启服务
- LLM 调用失败自动降级到规则引擎模式
- SMTP 失败重试 3 次后告警

### 7.3 可观测性

- Provider 操作审计日志（谁、何时、修改了什么）
- LLM 调用延迟 P99 < 2s
- SMTP 投递成功率 > 95%
- 平台 API 调用错误率 < 5%

### 7.4 兼容性

- 保留 `.env` 作为 System Default fallback
- 无 Provider 配置时自动使用 `.env`
- 数据库迁移向后兼容，支持回滚

---

## 8. 交付标准

### 8.1 后端交付

- [ ] W1 Workspace 表 + Provider 表 + Usage 表迁移脚本
- [ ] W2 Provider CRUD API + 在线测试 + 脱敏
- [ ] W3 AgentBackend 接口 + BuiltinAgent 实现 + Tool 注册表
- [ ] W4 LLM Provider 从 `.env` 迁移到数据库配置
- [ ] W5 SMTP Provider 从 `.env` 迁移到数据库配置
- [ ] W6 使用计量与限额管理
- [ ] W7 审计日志与安全加固

### 8.2 前端交付

- [ ] W1 Workspace 管理页面（列表/创建/切换/详情）
- [ ] W2 Provider 配置页面（分组展示/创建/编辑/测试/删除）
- [ ] W3 Agent 对话页面增强（工作空间上下文/用量提示）
- [ ] W4 设置页面整合（Workspace + Provider 统一入口）

### 8.3 测试交付

- [ ] W1 Workspace 隔离测试（A 无法访问 B 的 Provider）
- [ ] W2 Provider 优先级解析测试（Tenant > Workspace > Global > .env）
- [ ] W3 Agent 可插拔测试（Builtin / Remote 切换无需改业务代码）
- [ ] W4 端到端测试（创建 Workspace → 配置 Provider → Agent 使用 → 查看用量）

### 8.4 文档交付

- [ ] W1 API 文档（OpenAPI 3.0）
- [ ] W2 部署文档（环境变量 + 数据库迁移）
- [ ] W3 运维手册（Provider 故障排查、限额调整、审计日志查询）

---

## 9. 里程碑计划

| 里程碑 | 时间 | 交付物 |
|--------|------|--------|
| M1 Workspace + Provider 基础 | 第 1 周 | 数据库设计、CRUD API、前端页面 |
| M2 Agent 可插拔架构 | 第 2 周 | AgentBackend 接口、BuiltinAgent、Tool 注册表 |
| M3 LLM/SMTP 迁移 | 第 3 周 | 从 .env 迁移到 Provider、在线配置生效 |
| M4 计量与限额 | 第 4 周 | 使用统计、限额管理、告警 |
| M5 集成测试与文档 | 第 5 周 | 端到端测试、API 文档、部署文档 |

---

## 附录

### A. 术语表

| 术语 | 含义 |
|------|------|
| Workspace | 工作空间，顶层隔离单元，包含多个 Tenant |
| Provider | 统一外部服务配置（LLM/SMTP/广告平台/电商平台） |
| AgentBackend | Agent 后端接口，支持多实现 |
| BuiltinAgent | 进程内 Agent 实现 |
| RemoteAgent | 远程 Agent 服务实现（V3.1） |
| Tool | Agent 可调用的功能单元 |
| PendingAction | 待人工确认的高风险操作 |

### B. 参考文档

- V2.0 开发需求文档
- V2.1 开发需求文档
- Orbit API 接口规范
- Orbit 数据库设计文档

---

**文档结束**
