#!/bin/bash
# Orbit Build Script for Linux/macOS
# 自动检测环境、安装依赖、构建前后端、输出到 release/

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
RELEASE_DIR="$ROOT_DIR/release"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

log_info() { echo -e "${CYAN}[INFO]${NC} $1"; }
log_ok() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

show_help() {
    cat << EOF
Orbit 自动化构建脚本

用法:
  ./cmd/build.sh [选项]

选项:
  --skip-frontend   跳过前端构建
  --skip-backend    跳过后端构建
  --skip-deps       跳过依赖安装
  --clean           清理构建产物后重新构建
  --help            显示此帮助信息

示例:
  ./cmd/build.sh                # 完整构建（依赖+前端+后端）
  ./cmd/build.sh --clean        # 清理后完整构建
  ./cmd/build.sh --skip-frontend # 仅构建后端
EOF
    exit 0
}

# Parse arguments
SKIP_FRONTEND=false
SKIP_BACKEND=false
SKIP_DEPS=false
CLEAN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-frontend) SKIP_FRONTEND=true; shift ;;
        --skip-backend) SKIP_BACKEND=true; shift ;;
        --skip-deps) SKIP_DEPS=true; shift ;;
        --clean) CLEAN=true; shift ;;
        --help) show_help ;;
        *) log_error "未知选项: $1"; show_help ;;
    esac
done

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) log_error "不支持的架构: $ARCH"; exit 1 ;;
esac

log_info "检测到系统: $OS / $ARCH"

# Check Go version
GO_VERSION=$(go version 2>/dev/null | grep -oP 'go\K[0-9]+\.[0-9]+' || echo "")
if [ -z "$GO_VERSION" ]; then
    log_error "未检测到 Go，请安装 Go 1.22+ (https://go.dev/dl/)"
    exit 1
fi

GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)
if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 22 ]); then
    log_error "Go 版本过低 ($GO_VERSION)，需要 Go 1.22+"
    exit 1
fi
log_ok "Go $GO_VERSION"

# Check Node.js version
NODE_VERSION=$(node --version 2>/dev/null | sed 's/v//' || echo "")
if [ -z "$NODE_VERSION" ]; then
    log_error "未检测到 Node.js，请安装 Node.js 18+ (https://nodejs.org/)"
    exit 1
fi

NODE_MAJOR=$(echo $NODE_VERSION | cut -d. -f1)
if [ "$NODE_MAJOR" -lt 18 ]; then
    log_error "Node.js 版本过低 ($NODE_VERSION)，需要 Node.js 18+"
    exit 1
fi
log_ok "Node.js $NODE_VERSION"

NPM_VERSION=$(npm --version 2>/dev/null || echo "")
log_ok "NPM $NPM_VERSION"

# Check PostgreSQL (optional)
if ! pg_isready -q 2>/dev/null; then
    log_warn "PostgreSQL 未检测到运行中，构建将继续但本地测试需要数据库"
fi

# Check Redis (optional)
if ! redis-cli ping 2>/dev/null | grep -q PONG; then
    log_warn "Redis 未检测到运行中，构建将继续但本地测试需要 Redis"
fi

log_ok "环境检查通过"

# Clean
if [ "$CLEAN" = true ]; then
    log_info "清理构建产物..."
    rm -rf "$ROOT_DIR/apps/api/bin"
    rm -rf "$ROOT_DIR/apps/web/.next"
    rm -rf "$ROOT_DIR/apps/web/out"
    rm -rf "$ROOT_DIR/apps/web/dist"
    rm -rf "$RELEASE_DIR"
    log_ok "清理完成"
fi

# Create release directory
mkdir -p "$RELEASE_DIR/bin"

# Backend build
if [ "$SKIP_BACKEND" = false ]; then
    log_info "安装后端依赖..."
    cd "$ROOT_DIR/apps/api"
    go mod tidy
    log_ok "后端依赖安装完成"

    log_info "构建后端二进制文件..."
    cd "$ROOT_DIR/apps/api"

    # Build unified server binary (includes API + Worker + Scheduler)
    log_info "  构建 unified server..."
    go build -ldflags "-s -w" -o "$RELEASE_DIR/bin/orbit-server-$OS-$ARCH" ./cmd/server
    log_ok "  server -> bin/orbit-server-$OS-$ARCH"

    # Copy .env.example
    cp "$ROOT_DIR/apps/api/.env.example" "$RELEASE_DIR/.env.example"
    log_ok "后端构建完成"
fi

# Frontend build
if [ "$SKIP_FRONTEND" = false ]; then
    log_info "安装前端依赖..."
    cd "$ROOT_DIR/apps/web"
    if [ ! -d "node_modules" ]; then
        npm install
    else
        log_info "  node_modules 已存在，跳过安装"
    fi
    log_ok "前端依赖安装完成"

    log_info "构建前端..."
    npm run lint 2>/dev/null || log_warn "lint 发现问题，继续构建..."
    npm run typecheck 2>/dev/null || log_warn "typecheck 发现问题，继续构建..."
    npm run build
    log_ok "前端构建完成"

    # Copy build output
    log_info "复制前端产物..."
    FRONTEND_OUTPUT="$RELEASE_DIR/web"
    rm -rf "$FRONTEND_OUTPUT"
    mkdir -p "$FRONTEND_OUTPUT"

    if [ -d "$ROOT_DIR/apps/web/.next" ]; then
        cp -r "$ROOT_DIR/apps/web/.next" "$FRONTEND_OUTPUT/"
    fi
    if [ -d "$ROOT_DIR/apps/web/public" ]; then
        cp -r "$ROOT_DIR/apps/web/public" "$FRONTEND_OUTPUT/"
    fi
    cp "$ROOT_DIR/apps/web/package.json" "$FRONTEND_OUTPUT/"
    cp "$ROOT_DIR/apps/web/.env.example" "$RELEASE_DIR/.env.example" 2>/dev/null || true
    log_ok "前端产物已复制到 $FRONTEND_OUTPUT"
fi

# Generate deployment files
log_info "生成部署文件..."

# Systemd service files
mkdir -p "$RELEASE_DIR/systemd"

cat > "$RELEASE_DIR/systemd/orbit.service" << 'EOF'
[Unit]
Description=Orbit Unified Server (API + Worker + Scheduler)
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/orbit
EnvironmentFile=/opt/orbit/.env
ExecStart=/opt/orbit/bin/orbit-server
Restart=always
RestartSec=5
LimitNOFILE=65536
StandardOutput=journal
StandardError=journal
SyslogIdentifier=orbit-server

[Install]
WantedBy=multi-user.target
EOF

# Quick start script (foreground, no systemd)
cat > "$RELEASE_DIR/start.sh" << 'START_SCRIPT'
#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

if [ ! -f "$ROOT_DIR/.env" ]; then
    echo "[ERROR] .env file missing, copy from .env.example"
    exit 1
fi

set -a
source "$ROOT_DIR/.env"
set +a

mkdir -p "$ROOT_DIR/logs"

echo "[INFO] Starting unified orbit server (API + Worker + Scheduler)..."
nohup "$ROOT_DIR/bin/orbit-server" > "$ROOT_DIR/logs/server.log" 2>&1 &
echo $! > "$ROOT_DIR/logs/server.pid"
echo "[OK] Server started (PID: $!)"

echo ""
echo "URLs:"
echo "  - Frontend: http://localhost:3000"
echo "  - Backend: http://localhost:8080/api/v1"
echo "  - Health: http://localhost:8080/health/ready"
echo ""
echo "Logs: tail -f $ROOT_DIR/logs/server.log"
START_SCRIPT

# Stop script
cat > "$RELEASE_DIR/stop.sh" << 'STOP_SCRIPT'
#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo "[INFO] Stopping orbit server..."

if [ -f "$ROOT_DIR/logs/server.pid" ]; then
    pid=$(cat "$ROOT_DIR/logs/server.pid")
    if ps -p "$pid" > /dev/null 2>&1; then
        kill "$pid" 2>/dev/null || true
        sleep 1
        if ps -p "$pid" > /dev/null 2>&1; then
            kill -9 "$pid" 2>/dev/null || true
        fi
        echo "[OK] Server stopped"
    else
        echo "[INFO] Server not running"
    fi
    rm -f "$ROOT_DIR/logs/server.pid"
else
    echo "[INFO] No PID file found"
fi

echo "[OK] All services stopped"
STOP_SCRIPT

chmod +x "$RELEASE_DIR/start.sh" "$RELEASE_DIR/stop.sh"

log_ok "Deployment scripts generated"

echo ""
echo "=========================================="
log_ok "Build Complete!"
echo "=========================================="
echo ""
log_info "Output: $RELEASE_DIR"
log_info "Contents:"
log_info "  - bin/orbit-server-*-*    (Unified Server: API + Worker + Scheduler)"
log_info "  - web/                    (Frontend Static Files)"
log_info "  - .env.example            (Environment Template)"
log_info "  - systemd/orbit.service   (Systemd service file)"
log_info "  - start.sh / stop.sh      (Quick Start/Stop)"
echo ""
log_info "Deploy Options:"
log_info "  1. Binary deployment:  ./cmd/deploy.sh (systemd services)"
log_info "  2. Quick start:        ./cmd/start.sh (foreground, PID files)"
log_info "  3. Docker:             ./cmd/deploy-linux.sh (Docker Compose)"
echo ""
log_info "Binary Deploy Steps:"
log_info "  1. Upload release/ directory to Linux server"
log_info "  2. Copy .env.example to .env and configure JWT_SECRET / ENCRYPTION_MASTER_KEY"
log_info "  3. Run ./cmd/deploy.sh for systemd or ./cmd/start.sh for quick start"
echo ""
