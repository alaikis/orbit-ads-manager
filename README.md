# Orbit - 广告智能中枢

AI Agent 驱动的多平台广告自动化与智能化管理平台。

## 技术栈

- **前端**: Next.js 14 + TypeScript + Tailwind CSS + shadcn/ui + Recharts
- **后端**: Go 1.22 + Gin + PostgreSQL 15 + Redis 7 + Asynq
- **部署**: Linux + Docker Compose（前后端一体化部署）

## 生产环境要求

- **OS**: Linux (Ubuntu 20.04+ / CentOS 8+ / Alpine)
- **Docker**: 20.10+
- **Docker Compose**: 1.29+ / 2.0+
- **数据库**: PostgreSQL 15+（远程 `portal.gusty.top:5432/adshub`）
- **Redis**: 7+（本地容器或远程）

## 快速开始（Linux 生产部署）

### 1. 克隆仓库

```bash
git clone <repo-url>
cd ads.alaikis.com
```

### 2. 构建部署产物

```bash
# 赋予构建脚本执行权限
chmod +x cmd/build.sh

# 执行完整构建（自动检测环境、安装依赖、编译二进制、打包前端）
./cmd/build.sh
```

构建产物将生成在 `release/` 目录：
- `bin/orbit-server-linux-amd64` - API 服务
- `bin/orbit-server-linux-arm64` - API 服务（ARM 架构）
- `bin/orbit-worker-*` - 异步任务 worker
- `bin/orbit-scheduler-*` - 定时任务调度器
- `web/` - 前端静态资源
- `docker-compose.yml` - 生产环境编排
- `start.sh` / `stop.sh` - 快捷启动/停止脚本

### 3. 配置环境变量

```bash
cd release
cp .env.example .env
```

编辑 `.env`，配置以下关键变量：
```env
JWT_SECRET=<生成的32位以上随机字符串>
ENCRYPTION_MASTER_KEY=<生成的32字节Base64密钥>
```

> **安全提示**：生产环境请通过密钥服务或环境变量注入敏感凭证，不要硬编码在配置文件中。

### 4. 启动服务

```bash
# 方式一：使用启动脚本
./start.sh

# 方式二：直接使用 Docker Compose
docker compose up -d
```

### 5. 验证部署

```bash
# 检查服务状态
docker compose ps

# 查看日志
docker compose logs -f

# 健康检查
curl http://localhost:8080/health/ready
```

## 访问地址

- **前端**: http://localhost:3000
- **后端 API**: http://localhost:8080/api/v1
- **健康检查**: http://localhost:8080/health/ready

## 项目结构

```
apps/
├── api/                    # Go 后端服务
│   ├── cmd/
│   │   ├── server/        # API 服务入口
│   │   ├── worker/        # Asynq 任务 worker
│   │   └── scheduler/     # Cron 定时任务
│   ├── internal/
│   │   ├── auth/          # JWT + RBAC + 租户上下文
│   │   ├── tenant/        # 租户与多租户隔离
│   │   ├── store/         # 电商平台适配器
│   │   ├── adplatform/    # 广告平台适配器
│   │   ├── oauth/         # OAuth 授权与 token 管理
│   │   ├── sync/          # 数据同步任务
│   │   ├── feed/          # Feed 生成与发布
│   │   ├── agent/         # Agent 编排与 Tool 服务
│   │   ├── rule/          # 规则引擎
│   │   ├── report/        # 报表与看板
│   │   └── notify/        # 通知与日志
│   └── pkg/               # 可复用公共库
└── web/                    # Next.js 前端应用
    ├── app/               # App Router 页面
    ├── components/        # 业务组件
    └── lib/               # API client、类型、工具

packages/shared/            # 前后端共享类型与常量
docker/                     # Dockerfile 与 docker-compose
cmd/                        # 自动化构建与部署脚本
openspec/                   # OpenSpec 需求与设计文档
```

## 核心功能

- 多平台广告管理（Google Ads / Meta / Bing Ads）
- 电商对接（WooCommerce / Shopify）
- AI Agent 对话式交互
- 自动化规则引擎
- 跨平台报表与看板
- Feed 生成（Google Shopping / Meta Catalog）
- 多租户 RBAC 权限管理

## 构建命令

```bash
# Linux/macOS 完整构建
./cmd/build.sh

# 仅构建后端
./cmd/build.sh --skip-frontend

# 仅构建前端
./cmd/build.sh --skip-backend

# 清理后重新构建
./cmd/build.sh --clean

# Windows 构建
.\cmd\build.ps1
```

## 开发模式（本地）

```bash
# 启动 PostgreSQL 和 Redis
docker compose up -d postgres redis

# 终端 1：启动后端
cd apps/api && go run cmd/server/main.go

# 终端 2：启动前端
cd apps/web && npm run dev
```

## 文档

- [需求文档](./requirement-documents/Orbit_开发需求文档_V2.1.md)
- [OpenSpec 变更](./openspec/changes/)

## 许可证

MIT
