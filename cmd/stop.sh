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

echo "=========================================="
echo "  Orbit Stop Services"
echo "=========================================="
echo ""

# Stop by PID files
if [ -d "$RELEASE_DIR/logs" ]; then
    for pidfile in "$RELEASE_DIR/logs"/*.pid; do
        if [ -f "$pidfile" ]; then
            service_name=$(basename "$pidfile" .pid)
            pid=$(cat "$pidfile")
            if ps -p "$pid" > /dev/null 2>&1; then
                log_info "Stopping $service_name (PID: $pid)..."
                kill "$pid" 2>/dev/null || true
                sleep 1
                if ps -p "$pid" > /dev/null 2>&1; then
                    log_warn "Force killing $service_name..."
                    kill -9 "$pid" 2>/dev/null || true
                fi
                log_ok "$service_name stopped"
            else
                log_info "$service_name not running"
            fi
            rm -f "$pidfile"
        fi
    done
fi

# Also try systemctl if available
if command -v systemctl >/dev/null 2>&1; then
    log_info "Stopping systemd services..."
    sudo systemctl stop orbit 2>/dev/null || true
    log_ok "Systemd services stopped"
fi

echo ""
log_ok "All services stopped"
echo ""
