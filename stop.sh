#!/bin/bash

# Remote Code Stop Script
# Remote Code 停止脚本

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Config directory
CONFIG_DIR="$HOME/.remote-code"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
echo "🛑 Stopping Remote Code"
echo "================================"
echo -e "${NC}"

# Load config for PID directory
PID_DIR="$CONFIG_DIR/logs"
if [ -f "$CONFIG_DIR/config.ini" ]; then
    source "$CONFIG_DIR/config.ini"
    PID_DIR=${PID_DIR:-$CONFIG_DIR/logs}
fi

# Function to stop a service by PID file
stop_service() {
    local name=$1
    local pid_file="$PID_DIR/${name}.pid"

    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            echo -e "${YELLOW}Stopping $name (PID: $pid)...${NC}"
            kill "$pid" 2>/dev/null

            # Wait for process to stop
            local count=0
            while kill -0 "$pid" 2>/dev/null && [ $count -lt 10 ]; do
                sleep 1
                count=$((count + 1))
            done

            # Force kill if still running
            if kill -0 "$pid" 2>/dev/null; then
                echo -e "${YELLOW}Force killing $name...${NC}"
                kill -9 "$pid" 2>/dev/null
            fi

            echo -e "${GREEN}✅ $name stopped${NC}"
        else
            echo -e "${YELLOW}⚠️  $name is not running${NC}"
        fi
        rm -f "$pid_file"
    else
        echo -e "${YELLOW}⚠️  No PID file for $name${NC}"
    fi
}

# Stop all services
stop_service "backend"
stop_service "frontend"
stop_service "frp"

# Also try to stop by reading pids.txt (legacy)
if [ -f "$PID_DIR/pids.txt" ]; then
    read -ra PIDS < "$PID_DIR/pids.txt"
    for pid in "${PIDS[@]}"; do
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            echo -e "${YELLOW}Stopping process (PID: $pid)...${NC}"
            kill "$pid" 2>/dev/null || true
        fi
    done
    rm -f "$PID_DIR/pids.txt"
fi

# Kill orphan processes (from old config location or crashed)
echo ""
echo -e "${YELLOW}🔍 Checking for orphan processes...${NC}"

# Kill vite processes for this project
vite_pids=$(pgrep -f "vite.*$SCRIPT_DIR" 2>/dev/null)
if [ -n "$vite_pids" ]; then
    echo -e "${YELLOW}   Killing orphan vite processes: $vite_pids${NC}"
    echo "$vite_pids" | xargs kill 2>/dev/null || true
fi

# Kill go run processes for this project
go_pids=$(pgrep -f "go run.*$SCRIPT_DIR/backend" 2>/dev/null)
if [ -n "$go_pids" ]; then
    echo -e "${YELLOW}   Killing orphan go processes: $go_pids${NC}"
    echo "$go_pids" | xargs kill 2>/dev/null || true
fi

# Kill frpc processes for this project
frpc_pids=$(pgrep -f "frpc.*$SCRIPT_DIR" 2>/dev/null)
if [ -n "$frpc_pids" ]; then
    echo -e "${YELLOW}   Killing orphan frpc processes: $frpc_pids${NC}"
    echo "$frpc_pids" | xargs kill 2>/dev/null || true
fi

# Also check for frpc using config dir
frpc_pids2=$(pgrep -f "frpc.*$CONFIG_DIR" 2>/dev/null)
if [ -n "$frpc_pids2" ]; then
    echo -e "${YELLOW}   Killing orphan frpc processes: $frpc_pids2${NC}"
    echo "$frpc_pids2" | xargs kill 2>/dev/null || true
fi

echo ""
echo -e "${GREEN}✨ All services stopped${NC}"
