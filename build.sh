#!/bin/bash

# Remote Code Build Script
# 构建不同平台的可执行文件

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
BUILD_DIR="$SCRIPT_DIR/build"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
echo "🔨 Remote Code Build Script"
echo "================================${NC}"

# 创建构建目录
mkdir -p "$BUILD_DIR"

# 进入后端目录
cd "$BACKEND_DIR"

# 版本信息
VERSION=${VERSION:-"1.0.0"}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${YELLOW}Building version: $VERSION${NC}"
echo -e "${YELLOW}Git commit: $GIT_COMMIT${NC}"
echo ""

# 构建函数
build_binary() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT_NAME=$3

    echo -e "${BLUE}Building for $GOOS/$GOARCH...${NC}"

    GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build \
        -ldflags "-s -w -X main.Version=$VERSION -X main.GitCommit=$GIT_COMMIT" \
        -o "$BUILD_DIR/remote-code-$OUTPUT_NAME" \
        ./cmd/server 2>&1

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Built: remote-code-$OUTPUT_NAME${NC}"
    else
        echo -e "${RED}❌ Failed to build for $GOOS/$GOARCH${NC}"
        return 1
    fi
}

# 构建所有平台
echo -e "${BLUE}Building for current platform...${NC}"
build_binary "$(go env GOOS)" "$(go env GOARCH)" "$(go env GOOS)-$(go env GOARCH)"

echo ""
echo -e "${BLUE}Building for macOS Intel...${NC}"
build_binary "darwin" "amd64" "macos-intel"

echo ""
echo -e "${BLUE}Building for macOS Apple Silicon...${NC}"
build_binary "darwin" "arm64" "macos-apple"

echo ""
echo -e "${BLUE}Building for Linux x64...${NC}"
build_binary "linux" "amd64" "linux-x64"

echo ""
echo -e "${BLUE}Building for Linux ARM64...${NC}"
build_binary "linux" "arm64" "linux-arm64"

echo ""
echo -e "${BLUE}Building for Windows x64...${NC}"
build_binary "windows" "amd64" "windows-x64.exe"

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}✨ Build completed!${NC}"
echo -e "${GREEN}================================${NC}"
echo ""
echo -e "Output directory: ${YELLOW}$BUILD_DIR${NC}"
echo ""
ls -lh "$BUILD_DIR"
