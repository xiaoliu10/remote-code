#!/bin/bash

# Remote Code Unified Startup Script
# Remote Code 统一启动脚本
#
# Usage: ./start.sh [options]
# Options:
#   --no-frp      Disable FRP even if enabled in config
#   --frp         Enable FRP even if disabled in config
#   --help        Show this help message

set -e

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Config directory (user's home)
CONFIG_DIR="$HOME/.remote-code"
CONFIG_FILE="$CONFIG_DIR/config.ini"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Parse command line arguments
FRP_OVERRIDE=""
while [[ $# -gt 0 ]]; do
    case $1 in
        --no-frp)
            FRP_OVERRIDE="false"
            shift
            ;;
        --frp)
            FRP_OVERRIDE="true"
            shift
            ;;
        --help|-h)
            echo "Remote Code Startup Script"
            echo ""
            echo "Usage: ./start.sh [options]"
            echo ""
            echo "Options:"
            echo "  --no-frp      Disable FRP even if enabled in config"
            echo "  --frp         Enable FRP even if disabled in config"
            echo "  --help, -h    Show this help message"
            echo ""
            echo "Config location: $CONFIG_FILE"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}"
echo "🚀 Remote Code"
echo "================================"
echo -e "${NC}"

# Check dependencies
check_dependencies() {
    local missing=()

    if ! command -v tmux &> /dev/null; then
        missing+=("tmux")
    fi

    if ! command -v go &> /dev/null; then
        missing+=("go")
    fi

    if ! command -v node &> /dev/null; then
        missing+=("node")
    fi

    if ! command -v npm &> /dev/null; then
        missing+=("npm")
    fi

    if [ ${#missing[@]} -ne 0 ]; then
        echo -e "${RED}❌ Missing dependencies: ${missing[*]}${NC}"
        echo ""
        echo "Install instructions:"
        echo "  macOS: brew install ${missing[*]}"
        echo "  Ubuntu: sudo apt-get install ${missing[*]}"
        exit 1
    fi
}

check_dependencies

# Initialize config directory
init_config() {
    if [ ! -d "$CONFIG_DIR" ]; then
        echo -e "${YELLOW}📁 Creating config directory: $CONFIG_DIR${NC}"
        mkdir -p "$CONFIG_DIR"
    fi

    if [ ! -f "$CONFIG_FILE" ]; then
        echo -e "${YELLOW}📝 Creating default config: $CONFIG_FILE${NC}"

        # Get default values
        local home_dir="$HOME"
        local allowed_dir="$home_dir/projects"

        # Generate random password and JWT secret
        local admin_password=$(openssl rand -base64 12 | tr -d '/+=' | head -c 16)
        local jwt_secret=$(openssl rand -base64 32 | tr -d '/+=' | head -c 32)

        cat > "$CONFIG_FILE" << EOF
# Remote Code Configuration
# 配置文件路径: $CONFIG_FILE
#
# 修改此文件后重启服务生效

# ==================== Backend Configuration ====================
# 后端服务配置

# Backend port (后端服务端口)
BACKEND_PORT=9090

# JWT secret (JWT 密钥)
JWT_SECRET=$jwt_secret

# Admin password (管理员密码)
ADMIN_PASSWORD=$admin_password

# Allowed working directory (会话允许的工作目录)
ALLOWED_DIR=$allowed_dir

# Rate limiting (速率限制)
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPS=10
RATE_LIMIT_BURST=20

# ==================== Frontend Configuration ====================
# 前端服务配置

# Frontend dev server port (前端开发服务器端口)
FRONTEND_PORT=5173

# ==================== FRP Configuration ====================
# FRP 内网穿透配置

# Enable FRP tunnel (是否启用 FRP 内网穿透)
FRP_ENABLED=false

# FRP server address (FRP 服务器地址)
FRP_SERVER_ADDR=

# FRP server port (FRP 服务器端口)
FRP_SERVER_PORT=7000

# FRP authentication token (FRP 认证令牌)
FRP_TOKEN=

# ==================== Advanced Settings ====================
# 高级设置

# Log directory (日志目录)
LOG_DIR=$CONFIG_DIR/logs

# PID file directory (PID 文件目录)
PID_DIR=$CONFIG_DIR/logs
EOF

        echo -e "${GREEN}✅ Created config file with random password${NC}"
        echo -e "${YELLOW}   Admin Password: $admin_password${NC}"
        echo -e "${YELLOW}   Please save this password!${NC}"
        echo ""
    fi
}

# Load configuration
load_config() {
    init_config

    # Source the config file (export variables)
    set -a
    source "$CONFIG_FILE"
    set +a

    # Apply command line override
    if [ -n "$FRP_OVERRIDE" ]; then
        FRP_ENABLED="$FRP_OVERRIDE"
    fi

    # Set defaults if not configured
    BACKEND_PORT=${BACKEND_PORT:-9090}
    FRONTEND_PORT=${FRONTEND_PORT:-5173}
    LOG_DIR=${LOG_DIR:-$CONFIG_DIR/logs}
    PID_DIR=${PID_DIR:-$CONFIG_DIR/logs}
}

load_config

# Create necessary directories
mkdir -p "$LOG_DIR"
mkdir -p "$PID_DIR"

# Generate backend .env from config
generate_backend_env() {
    local env_file="$SCRIPT_DIR/backend/.env"

    cat > "$env_file" << EOF
# Server Configuration - Auto-generated from config.ini
PORT=$BACKEND_PORT

# JWT Secret - CHANGE THIS IN PRODUCTION!
JWT_SECRET=$JWT_SECRET

# Admin Credentials
ADMIN_PASSWORD=$ADMIN_PASSWORD

# Security
ALLOWED_DIR=$ALLOWED_DIR

# Rate Limiting
RATE_LIMIT_ENABLED=$RATE_LIMIT_ENABLED
RATE_LIMIT_RPS=$RATE_LIMIT_RPS
RATE_LIMIT_BURST=$RATE_LIMIT_BURST
EOF

    echo -e "${GREEN}✅ Generated backend/.env${NC}"
}

# Generate frontend .env from config
generate_frontend_env() {
    local env_file="$SCRIPT_DIR/frontend/.env"

    cat > "$env_file" << EOF
# Frontend Configuration - Auto-generated from config.ini
VITE_BACKEND_PORT=$BACKEND_PORT
EOF

    echo -e "${GREEN}✅ Generated frontend/.env${NC}"
}

# Generate FRP config if enabled
generate_frp_config() {
    local frp_config="$CONFIG_DIR/frpc.ini"

    if [ "$FRP_ENABLED" = "true" ]; then
        # Check if frpc.ini already has valid config
        if [ -f "$frp_config" ] && [ -n "$(grep -E '^serverAddr = "[0-9]+\.' "$frp_config" 2>/dev/null)" ]; then
            # Has valid IP address config, don't overwrite
            echo -e "${GREEN}✅ Using existing $frp_config${NC}"
            return
        fi

        cat > "$frp_config" << EOF
# FRP Client Configuration - Auto-generated from config.ini
serverAddr = "$FRP_SERVER_ADDR"
serverPort = $FRP_SERVER_PORT

# Authentication token
auth.token = "$FRP_TOKEN"

# Frontend service
[[proxies]]
name = "remote-code-frontend"
type = "tcp"
localIP = "127.0.0.1"
localPort = $FRONTEND_PORT
remotePort = $FRONTEND_PORT

# Backend API + WebSocket service
[[proxies]]
name = "remote-code-backend"
type = "tcp"
localIP = "127.0.0.1"
localPort = $BACKEND_PORT
remotePort = $BACKEND_PORT
EOF
        echo -e "${GREEN}✅ Generated $frp_config${NC}"
    fi
}

# Generate environment files
echo -e "${BLUE}📝 Generating configuration files...${NC}"
generate_backend_env
generate_frontend_env
generate_frp_config
echo ""

# Stop any existing services
stop_existing() {
    local pid_file="$PID_DIR/pids.txt"
    if [ -f "$pid_file" ]; then
        echo -e "${YELLOW}⚠️  Stopping existing services...${NC}"
        read -ra PIDS < "$pid_file"
        for pid in "${PIDS[@]}"; do
            if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
                kill "$pid" 2>/dev/null || true
            fi
        done
        rm -f "$pid_file"
        sleep 1
    fi
}

stop_existing

# Start backend
start_backend() {
    echo -e "${BLUE}🔧 Starting backend service (port $BACKEND_PORT)...${NC}"

    cd "$SCRIPT_DIR/backend"

    # Check if node_modules exists for frontend (to ensure we're in right place)
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}❌ Backend directory structure invalid${NC}"
        exit 1
    fi

    # Detect OS and architecture for compiled binary
    local os="linux"
    local arch="amd64"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        os="macos"
        if [[ "$(uname -m)" == "arm64" ]]; then
            arch="apple"
        else
            arch="intel"
        fi
    fi

    # Look for compiled binary
    local binary_name="remote-code-${os}-${arch}"
    local binary_path="$SCRIPT_DIR/build/$binary_name"

    if [ -f "$binary_path" ]; then
        echo -e "${GREEN}   Using compiled binary: $binary_name${NC}"
        "$binary_path" > "$LOG_DIR/backend.log" 2>&1 &
    else
        echo -e "${YELLOW}   No compiled binary found, using 'go run'${NC}"
        echo -e "${YELLOW}   Run './build.sh' to compile for faster startup${NC}"
        go run cmd/server/main.go > "$LOG_DIR/backend.log" 2>&1 &
    fi

    BACKEND_PID=$!

    cd "$SCRIPT_DIR"
    echo "$BACKEND_PID" > "$PID_DIR/backend.pid"
    echo -e "${GREEN}✅ Backend started (PID: $BACKEND_PID)${NC}"
}

# Start frontend
start_frontend() {
    echo -e "${BLUE}🎨 Starting frontend service (port $FRONTEND_PORT)...${NC}"

    cd "$SCRIPT_DIR/frontend"

    # Install dependencies if needed
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}   Installing dependencies...${NC}"
        npm install --silent
    fi

    npm run dev > "$LOG_DIR/frontend.log" 2>&1 &
    FRONTEND_PID=$!

    cd "$SCRIPT_DIR"
    echo "$FRONTEND_PID" > "$PID_DIR/frontend.pid"
    echo -e "${GREEN}✅ Frontend started (PID: $FRONTEND_PID)${NC}"
}

# Start FRP
start_frp() {
    if [ "$FRP_ENABLED" != "true" ]; then
        return
    fi

    echo -e "${BLUE}🌐 Starting FRP client...${NC}"

    # Check for frpc binary
    local frpc_bin=""
    if [ -f "$CONFIG_DIR/frpc" ]; then
        frpc_bin="$CONFIG_DIR/frpc"
    elif [ -f "$SCRIPT_DIR/frp/frpc" ]; then
        frpc_bin="$SCRIPT_DIR/frp/frpc"
    elif command -v frpc &> /dev/null; then
        frpc_bin="frpc"
    else
        echo -e "${YELLOW}⚠️  FRP client (frpc) not found, downloading...${NC}"

        # Detect OS and architecture
        local os="linux"
        local arch="amd64"
        if [[ "$OSTYPE" == "darwin"* ]]; then
            os="darwin"
            if [[ "$(uname -m)" == "arm64" ]]; then
                arch="arm64"
            fi
        fi

        local frp_version="0.61.1"
        local frp_url="https://github.com/fatedier/frp/releases/download/v${frp_version}/frp_${frp_version}_${os}_${arch}.tar.gz"

        echo -e "${YELLOW}   Downloading from: $frp_url${NC}"
        curl -sL "$frp_url" | tar xz -C "$CONFIG_DIR" --strip-components=1 "frp_${frp_version}_${os}_${arch}/frpc"
        chmod +x "$CONFIG_DIR/frpc"
        frpc_bin="$CONFIG_DIR/frpc"
        echo -e "${GREEN}✅ Downloaded FRP client to $CONFIG_DIR${NC}"
    fi

    "$frpc_bin" -c "$CONFIG_DIR/frpc.ini" > "$LOG_DIR/frp.log" 2>&1 &
    FRP_PID=$!
    echo "$FRP_PID" > "$PID_DIR/frp.pid"
    echo -e "${GREEN}✅ FRP client started (PID: $FRP_PID)${NC}"
    echo -e "${GREEN}   Server: $FRP_SERVER_ADDR:$FRP_SERVER_PORT${NC}"
}

# Start all services
start_backend
sleep 2
start_frontend
start_frp

# Save all PIDs
echo "$BACKEND_PID $FRONTEND_PID $FRP_PID" > "$PID_DIR/pids.txt"

# Get local IP
LOCAL_IP=$(hostname -I 2>/dev/null | awk '{print $1}' || echo "localhost")

# Display success message
echo ""
echo -e "${GREEN}"
echo "================================"
echo "✨ Started successfully!"
echo "================================"
echo -e "${NC}"
echo ""
echo -e "${BLUE}Access URLs:${NC}"
echo -e "  Local:   ${GREEN}http://localhost:$FRONTEND_PORT${NC}"
if [ "$LOCAL_IP" != "localhost" ]; then
    echo -e "  Network: ${GREEN}http://$LOCAL_IP:$FRONTEND_PORT${NC}"
fi

if [ "$FRP_ENABLED" = "true" ]; then
    echo ""
    echo -e "${BLUE}Remote Access (via FRP):${NC}"
    echo -e "  ${GREEN}http://$FRP_SERVER_ADDR:$FRONTEND_PORT${NC}"
fi

echo ""
echo -e "${BLUE}Default Login:${NC}"
echo -e "  Username: ${YELLOW}admin${NC}"
echo -e "  Password: ${YELLOW}$ADMIN_PASSWORD${NC}"
echo ""
echo -e "${BLUE}Management:${NC}"
echo -e "  Stop:     ${YELLOW}./stop.sh${NC}"
echo -e "  Logs:     ${YELLOW}tail -f $LOG_DIR/*.log${NC}"
echo -e "  Config:   ${YELLOW}$CONFIG_FILE${NC}"
echo ""
