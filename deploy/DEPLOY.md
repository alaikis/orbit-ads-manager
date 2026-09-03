# Orbit 部署指南

## 前期部署方案

采用 **单二进制 + 多 goroutine** 模式：
- 一个 `orbit-server` 二进制
- 内部同时运行：Gin API + Asynq Worker + Cron Scheduler
- 简化部署，减少进程管理复杂度

## 构建

```bash
# 本地构建
./cmd/build.sh

# 或手动构建
cd apps/api
go build -ldflags="-s -w" -o orbit-server ./cmd/server
```

## 部署到服务器

### 方式一：Binary 部署（推荐）

```bash
# 1. 上传 release/ 目录到服务器
scp -r release/ root@portal.gusty.top:/apphub/

# 2. SSH 到服务器
ssh root@portal.gusty.top

# 3. 配置环境变量
cd /apphub
cp .env.example .env
nano .env  # 编辑 JWT_SECRET, ENCRYPTION_MASTER_KEY, DATABASE_URL 等

# 4. 前台启动（调试）
./cmd/start.sh

# 5. 或 systemd 部署（生产）
./cmd/deploy.sh
```

### 方式二：Docker 部署

```bash
./cmd/deploy-linux.sh
```

## Caddy 配置

服务器上已部署 Caddy，配置文件位置：
- `/etc/caddy/Caddyfile`

推荐配置：
```
falcon.alaikis.com {
    reverse_proxy /api/* localhost:8080
    root * /apphub/web
    file_server
}

ads.alaikis.com {
    reverse_proxy /* localhost:8080
}
```

部署后访问：
- 前端：https://falcon.alaikis.com
- API：https://falcon.alaikis.com/api/v1
- 健康检查：https://falcon.alaikis.com/health/ready

## 目录结构

```
/apphub/
├── bin/
│   └── orbit-server          # 统一二进制
├── web/                      # 前端静态文件
├── .env                      # 环境变量
├── logs/                     # 日志目录
├── cmd/
│   ├── start.sh              # 前台启动
│   └── stop.sh               # 停止服务
└── systemd/
    └── orbit.service         # systemd 服务文件
```

## 管理命令

```bash
# systemd 方式
sudo systemctl status orbit
sudo journalctl -u orbit -f
sudo systemctl restart orbit

# PID 文件方式
tail -f /apphub/logs/server.log
kill $(cat /apphub/logs/server.pid)
```
