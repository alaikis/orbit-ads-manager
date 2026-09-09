#!/bin/bash
set -e

REMOTE_USER="root"
REMOTE_HOST="portal.gusty.top"
APP_DIR="/opt/orbit"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
RELEASE_DIR="$ROOT_DIR/release"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log_info() { echo -e "${CYAN}[INFO]${NC} $1"; }
log_ok() { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

echo "=========================================="
echo "  Orbit Portal Deployment"
echo "=========================================="
echo ""

# Check prerequisites
command -v ssh >/dev/null 2>&1 || { log_error "ssh not found"; exit 1; }
command -v scp >/dev/null 2>&1 || { log_error "scp not found"; exit 1; }

# Check if release directory exists
if [ ! -d "$RELEASE_DIR" ]; then
    log_error "Release directory not found: $RELEASE_DIR"
    log_info "Run ./cmd/build.sh first to build binaries"
    exit 1
fi

# Check if .env exists
if [ ! -f "$RELEASE_DIR/.env" ]; then
    log_warn ".env file not found in release/"
    if [ -f "$RELEASE_DIR/.env.example" ]; then
        log_info "Copying .env.example to .env"
        cp "$RELEASE_DIR/.env.example" "$RELEASE_DIR/.env"
        log_warn "Please edit $RELEASE_DIR/.env and configure JWT_SECRET and ENCRYPTION_MASTER_KEY"
        read -p "Continue anyway? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        log_error ".env.example not found, cannot continue"
        exit 1
    fi
fi

log_info "Deploying to $REMOTE_HOST..."

# Create remote app directory
log_info "Creating remote directory..."
ssh "$REMOTE_USER@$REMOTE_HOST" "mkdir -p $APP_DIR"/{bin,logs,data,caddy}

# Upload .env
log_info "Uploading environment file..."
scp "$RELEASE_DIR/.env" "$REMOTE_USER@$REMOTE_HOST:$APP_DIR/.env"

# Upload Caddy config
log_info "Uploading Caddy config..."
scp "$ROOT_DIR/caddy_config" "$REMOTE_USER@$REMOTE_HOST:$APP_DIR/caddy_config"

# Deploy Docker stack
log_info "Deploying Docker stack..."
ssh "$REMOTE_USER@$REMOTE_HOST" "cd $APP_DIR && docker compose up -d"

# Deploy Caddy config
log_info "Deploying Caddy config..."
ssh "$REMOTE_USER@$REMOTE_HOST" "cp $APP_DIR/caddy_config /etc/caddy/Caddyfile && caddy validate --config /etc/caddy/Caddyfile && systemctl reload caddy"

# Wait for services
log_info "Waiting for services to start..."
sleep 5

# Check status
log_info "Checking service status..."
ssh "$REMOTE_USER@$REMOTE_HOST" "systemctl status orbit --no-pager || true"
ssh "$REMOTE_USER@$REMOTE_HOST" "systemctl status caddy --no-pager || true"

echo ""
echo "=========================================="
log_ok "Deployment Complete!"
echo "=========================================="
echo ""
echo "Services:"
echo "  - Frontend:   https://ads.alaikis.com (Vercel)"
echo "  - Frontend:   https://adms.alaikis.com (Vercel)"
echo "  - API:        https://adsapi.alaikis.com/api/v1"
echo "  - Health:     https://adsapi.alaikis.com/health/ready"
echo ""
echo "Management:"
echo "  - API logs:   ssh $REMOTE_USER@$REMOTE_HOST 'journalctl -u orbit -f'"
echo "  - Caddy logs: ssh $REMOTE_USER@$REMOTE_HOST 'journalctl -u caddy -f'"
echo "  - Restart:    ssh $REMOTE_USER@$REMOTE_HOST 'systemctl restart orbit'"
echo ""
