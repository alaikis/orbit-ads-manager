#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
RELEASE_DIR="$ROOT_DIR/release"

echo "=========================================="
echo "  Orbit Linux Binary Deployment"
echo "=========================================="
echo ""

# Check prerequisites
command -v systemctl >/dev/null 2>&1 || { echo "[ERROR] systemctl not found, use ./cmd/start.sh for manual deployment"; exit 1; }

# Check if release directory exists
if [ ! -d "$RELEASE_DIR" ]; then
    echo "[ERROR] Release directory not found: $RELEASE_DIR"
    echo "       Run ./cmd/build.sh first to build binaries"
    exit 1
fi

# Check if .env exists
if [ ! -f "$RELEASE_DIR/.env" ]; then
    echo "[WARN] .env file not found in release/"
    if [ -f "$RELEASE_DIR/.env.example" ]; then
        echo "       Copying .env.example to .env"
        cp "$RELEASE_DIR/.env.example" "$RELEASE_DIR/.env"
        echo "       Please edit $RELEASE_DIR/.env and configure JWT_SECRET and ENCRYPTION_MASTER_KEY"
        read -p "Continue anyway? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        echo "[ERROR] .env.example not found, cannot continue"
        exit 1
    fi
fi

# Detect binary suffix
BIN_SUFFIX=""
if [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "arm64" ]; then
    BIN_SUFFIX="-arm64"
elif [ "$(uname -m)" = "x86_64" ]; then
    BIN_SUFFIX="-amd64"
fi

echo "[INFO] Detected architecture: $(uname -m) (suffix: $BIN_SUFFIX)"

# Check if binaries exist
if [ ! -f "$RELEASE_DIR/bin/orbit-server$BIN_SUFFIX" ]; then
    echo "[ERROR] Binary not found: orbit-server$BIN_SUFFIX"
    exit 1
fi

# Create app directory
APP_DIR="/opt/orbit"
echo "[INFO] Creating application directory: $APP_DIR"
sudo mkdir -p "$APP_DIR"/{bin,logs,data}

# Copy binaries
echo "[INFO] Copying binaries..."
sudo cp "$RELEASE_DIR/bin/orbit-server$BIN_SUFFIX" "$APP_DIR/bin/orbit-server"
sudo chmod +x "$APP_DIR/bin/orbit-server"

# Copy .env
echo "[INFO] Copying environment file..."
sudo cp "$RELEASE_DIR/.env" "$APP_DIR/.env"

# Create systemd service files
echo "[INFO] Creating systemd service files..."

sudo tee /etc/systemd/system/orbit.service > /dev/null << 'EOF'
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

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
echo "[INFO] Reloading systemd..."
sudo systemctl daemon-reload

# Enable and start services
echo "[INFO] Enabling services..."
sudo systemctl enable orbit

echo "[INFO] Starting services..."
sudo systemctl start orbit

# Wait a moment for services to start
sleep 3

# Check status
echo "[INFO] Checking service status..."
sudo systemctl status orbit --no-pager || true

echo ""
echo "=========================================="
echo "  Deployment Complete!"
echo "=========================================="
echo ""
echo "Services:"
echo "  - API Server:   http://localhost:8080/api/v1"
echo "  - Health Check: http://localhost:8080/health/ready"
echo ""
echo "Management:"
echo "  - Status:  sudo systemctl status orbit"
echo "  - Logs:    sudo journalctl -u orbit -f"
echo "  - Restart: sudo systemctl restart orbit"
echo "  - Stop:    sudo systemctl stop orbit"
echo ""
echo "Configuration:"
echo "  - Env file: /opt/orbit/.env"
echo "  - Binary:   /opt/orbit/bin/orbit-server"
echo "  - Logs:     journalctl -u orbit"
echo ""
