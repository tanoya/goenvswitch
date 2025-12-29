#!/bin/bash

# ============================================================================
# GoEnv-Switch 快速测试脚本
# 用于验证工程是否可以正常编译和运行
# ============================================================================

set -e

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "========================================"
echo "  GoEnv-Switch 工程测试"
echo "========================================"
echo ""

# 1. 检查 Go 环境
echo -e "${YELLOW}[1/6]${NC} 检查 Go 环境..."
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗${NC} 未找到 Go 编译器"
    exit 1
fi
echo -e "${GREEN}✓${NC} Go 版本: $(go version)"
echo ""

# 2. 检查文件结构
echo -e "${YELLOW}[2/6]${NC} 检查文件结构..."
if [ ! -f "cmd/main.go" ]; then
    echo -e "${RED}✗${NC} 未找到 cmd/main.go"
    exit 1
fi
if [ ! -f "config/config.yaml" ]; then
    echo -e "${RED}✗${NC} 未找到 config/config.yaml"
    exit 1
fi
echo -e "${GREEN}✓${NC} 文件结构正确"
echo ""

# 3. 下载依赖
echo -e "${YELLOW}[3/6]${NC} 下载依赖..."
go mod tidy
go mod download
echo -e "${GREEN}✓${NC} 依赖下载完成"
echo ""

# 4. 测试编译（当前平台）
echo -e "${YELLOW}[4/6]${NC} 测试编译..."
mkdir -p build
go build -o build/goenv-switch ./cmd
echo -e "${GREEN}✓${NC} 编译成功: build/goenv-switch"
echo ""

# 5. 测试运行
echo -e "${YELLOW}[5/6]${NC} 测试运行..."
./build/goenv-switch --help > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} 程序可以正常运行"
else
    echo -e "${RED}✗${NC} 程序运行失败"
    exit 1
fi
echo ""

# 6. 测试 Makefile
echo -e "${YELLOW}[6/6]${NC} 测试 Makefile..."
if command -v make &> /dev/null; then
    make clean > /dev/null 2>&1
    make build > /dev/null 2>&1
    if [ -f "build/goenv-switch" ]; then
        echo -e "${GREEN}✓${NC} Makefile 工作正常"
    else
        echo -e "${RED}✗${NC} Makefile 编译失败"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠${NC} 未安装 make，跳过测试"
fi
echo ""

echo "========================================"
echo -e "${GREEN}✓ 所有测试通过！${NC}"
echo "========================================"
echo ""
echo "工程状态: 可以正常编译和运行"
echo ""
echo "快速开始:"
echo "  make build        # 编译当前平台"
echo "  make install      # 安装到 GOPATH/bin"
echo "  ./build.sh build  # 使用 build.sh 编译"
echo ""
