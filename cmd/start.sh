#!/bin/bash
set -e

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
echo "  Orbit Quick Start (Binary Mode)"
echo "=========================================="
echo ""

# Check prerequisites
command -v docker >/dev/null 2>&1 && log_warn "Docker detected but not required for binary mode" || true

# Check if release directory exists
if [ ! -d "$RELEASE_DIR" ]; then
    log_error "Release directory not found: $RELEASE_DIR"
    log_info "Run ./cmd/build.sh first to build binaries"
    exit 1
fi

# Check if .env exists
if [ ! -f "$RELEASE_DIR/.env" ]; then
    log_error ".env file not found in release/"
    log_info "Copy .env.example to .env and configure JWT_SECRET and ENCRYPTION_MASTER_KEY"
    exit 1
fi

# Detect binary suffix
BIN_SUFFIX=""
if [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "arm64" ]; then
    BIN_SUFFIX="-arm64"
elif [ "$(uname -m)" = "x86_64" ]; then
    BIN_SUFFIX="-amd64"
fi

# Check if binary exists
if [ ! -f "$RELEASE_DIR/bin/orbit-server$BIN_SUFFIX" ]; then
    log_error "Binary not found: orbit-server$BIN_SUFFIX"
    exit 1
fi

# Export environment variables
set -a
source "$RELEASE_DIR/.env"
set +a

# Create logs directory
mkdir -p "$RELEASE_DIR/logs"

log_info "Starting unified orbit server (API + Worker + Scheduler)..."
nohup "$RELEASE_DIR/bin/orbit-server$BIN_SUFFIX" > "$RELEASE_DIR/logs/server.log" 2>&1 &
echo $! > "$RELEASE_DIR/logs/server.pid"
log_ok "Server started (PID: $!)"

echo ""
echo "=========================================="
log_ok "Services Started!"
echo "=========================================="
echo ""
echo "URLs:"
echo "  - Frontend: http://localhost:3000"
echo "  - Backend:  http://localhost:8080/api/v1"
echo "  - Health:   http://localhost:8080/health/ready"
echo ""
echo "Logs:"
echo "  - tail -f $RELEASE_DIR/logs/server.log"
echo ""
echo "Stop:"
echo "  - kill \$(cat $RELEASE_DIR/logs/server.pid)"
echo ""
