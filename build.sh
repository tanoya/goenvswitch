#!/bin/bash

# ============================================================================
# GoEnv-Switch 编译脚本
# 支持多平台交叉编译
# ============================================================================

set -e

# 项目信息
APP_NAME="goenv-switch"
VERSION="${VERSION:-1.0.0}"
BUILD_TIME=$(date "+%Y-%m-%d %H:%M:%S")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 目录设置
ROOT_DIR=$(cd "$(dirname "$0")" && pwd)
BUILD_DIR="${ROOT_DIR}/build"
DIST_DIR="${ROOT_DIR}/dist"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# 打印横幅
print_banner() {
    echo ""
    echo "=============================================="
    echo "  GoEnv-Switch 编译脚本"
    echo "  版本: ${VERSION}"
    echo "  提交: ${GIT_COMMIT}"
    echo "=============================================="
    echo ""
}

# 检查 Go 环境
check_go() {
    if ! command -v go &> /dev/null; then
        error "未找到 Go 编译器，请先安装 Go"
    fi
    
    GO_VERSION=$(go version | awk '{print $3}')
    info "Go 版本: ${GO_VERSION}"
}

# 清理构建目录
clean() {
    info "清理构建目录..."
    rm -rf "${BUILD_DIR}"
    rm -rf "${DIST_DIR}"
    success "清理完成"
}

# 下载依赖
download_deps() {
    info "下载依赖..."
    cd "${ROOT_DIR}"
    go mod tidy
    go mod download
    success "依赖下载完成"
}

# 运行测试
run_tests() {
    info "运行测试..."
    cd "${ROOT_DIR}"
    go test -v ./...
    success "测试通过"
}

# 代码检查
lint() {
    info "代码检查..."
    
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run ./...
        success "代码检查通过"
    else
        warn "未安装 golangci-lint,跳过代码检查"
    fi
}

# 编译单个平台
build_single() {
    local os=$1
    local arch=$2
    local output_name=$3
    
    info "编译 ${os}/${arch}..."
    
    # 设置输出目录
    local output_dir="${BUILD_DIR}/${os}_${arch}"
    mkdir -p "${output_dir}"
    
    # 设置可执行文件名
    local binary_name="${APP_NAME}"
    if [ "${os}" = "windows" ]; then
        binary_name="${APP_NAME}.exe"
    fi
    
    # 编译
    CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build \
        -ldflags "-s -w \
            -X 'main.Version=${VERSION}' \
            -X 'main.BuildTime=${BUILD_TIME}' \
            -X 'main.GitCommit=${GIT_COMMIT}'" \
        -o "${output_dir}/${binary_name}" \
        "${ROOT_DIR}/main.go"
    
    # 复制配置文件示例
    cp "${ROOT_DIR}/config.yaml" "${output_dir}/" 2>/dev/null || true
    cp "${ROOT_DIR}/README.md" "${output_dir}/" 2>/dev/null || true
    
    success "编译完成: ${output_dir}/${binary_name}"
}

# 编译当前平台
build_current() {
    info "编译当前平台..."
    
    mkdir -p "${BUILD_DIR}"
    
    local binary_name="${APP_NAME}"
    if [ "$(go env GOOS)" = "windows" ]; then
        binary_name="${APP_NAME}.exe"
    fi
    
    go build \
        -ldflags "-s -w \
            -X 'main.Version=${VERSION}' \
            -X 'main.BuildTime=${BUILD_TIME}' \
            -X 'main.GitCommit=${GIT_COMMIT}'" \
        -o "${BUILD_DIR}/${binary_name}" \
        "${ROOT_DIR}/main.go"
    
    success "编译完成: ${BUILD_DIR}/${binary_name}"
}

# 编译所有平台
build_all() {
    info "开始多平台编译..."
    
    # 定义目标平台
    local platforms=(
        "linux/amd64"
        "linux/arm64"
        "linux/386"
        "darwin/amd64"
        "darwin/arm64"
        "windows/amd64"
        "windows/386"
        "windows/arm64"
    )
    
    for platform in "${platforms[@]}"; do
        local os="${platform%/*}"
        local arch="${platform#*/}"
        build_single "${os}" "${arch}"
    done
    
    success "所有平台编译完成"
}

# 打包发布文件
package() {
    info "打包发布文件..."
    
    mkdir -p "${DIST_DIR}"
    
    # 遍历构建目录
    for dir in "${BUILD_DIR}"/*; do
        if [ -d "${dir}" ]; then
            local platform=$(basename "${dir}")
            local archive_name="${APP_NAME}_${VERSION}_${platform}"
            
            info "打包 ${platform}..."
            
            cd "${BUILD_DIR}"
            
            if [[ "${platform}" == windows_* ]]; then
                # Windows 使用 zip
                zip -r "${DIST_DIR}/${archive_name}.zip" "${platform}"
            else
                # 其他平台使用 tar.gz
                tar -czvf "${DIST_DIR}/${archive_name}.tar.gz" "${platform}"
            fi
        fi
    done
    
    success "打包完成，发布文件位于: ${DIST_DIR}"
}

# 生成校验和
generate_checksums() {
    info "生成校验和..."
    
    cd "${DIST_DIR}"
    
    if command -v sha256sum &> /dev/null; then
        sha256sum * > checksums.txt
    elif command -v shasum &> /dev/null; then
        shasum -a 256 * > checksums.txt
    else
        warn "未找到 sha256sum 或 shasum，跳过校验和生成"
        return
    fi
    
    success "校验和已生成: ${DIST_DIR}/checksums.txt"
}

# 安装到 GOPATH/bin
install_local() {
    info "安装到 GOPATH/bin..."
    
    go install \
        -ldflags "-s -w \
            -X 'main.Version=${VERSION}' \
            -X 'main.BuildTime=${BUILD_TIME}' \
            -X 'main.GitCommit=${GIT_COMMIT}'" \
        "${ROOT_DIR}"
    
    success "安装完成: $(go env GOPATH)/bin/${APP_NAME}"
}

# 显示帮助信息
show_help() {
    echo "用法: $0 <命令>"
    echo ""
    echo "命令:"
    echo "  build         编译当前平台"
    echo "  build-all     编译所有平台"
    echo "  build-linux   编译 Linux 平台 (amd64, arm64)"
    echo "  build-darwin  编译 macOS 平台 (amd64, arm64)"
    echo "  build-windows 编译 Windows 平台 (amd64, 386)"
    echo "  clean         清理构建目录"
    echo "  deps          下载依赖"
    echo "  test          运行测试"
    echo "  lint          代码检查"
    echo "  package       打包发布文件"
    echo "  checksum      生成校验和"
    echo "  install       安装到 GOPATH/bin"
    echo "  release       完整发布流程 (clean + deps + test + build-all + package + checksum)"
    echo "  help          显示帮助信息"
    echo ""
    echo "环境变量:"
    echo "  VERSION       设置版本号 (默认: 1.0.0)"
    echo ""
    echo "示例:"
    echo "  $0 build                    # 编译当前平台"
    echo "  $0 build-all                # 编译所有平台"
    echo "  VERSION=2.0.0 $0 release    # 使用指定版本号发布"
}

# 编译 Linux 平台
build_linux() {
    build_single "linux" "amd64"
    build_single "linux" "arm64"
}

# 编译 macOS 平台
build_darwin() {
    build_single "darwin" "amd64"
    build_single "darwin" "arm64"
}

# 编译 Windows 平台
build_windows() {
    build_single "windows" "amd64"
    build_single "windows" "386"
}

# 完整发布流程
release() {
    clean
    download_deps
    run_tests
    build_all
    package
    generate_checksums
    
    echo ""
    echo "=============================================="
    success "发布完成!"
    echo "  版本: ${VERSION}"
    echo "  发布目录: ${DIST_DIR}"
    echo "=============================================="
    echo ""
    
    # 列出发布文件
    ls -lh "${DIST_DIR}"
}

# 主函数
main() {
    print_banner
    check_go
    
    cd "${ROOT_DIR}"
    
    case "${1:-build}" in
        build)
            download_deps
            build_current
            ;;
        build-all)
            download_deps
            build_all
            ;;
        build-linux)
            download_deps
            build_linux
            ;;
        build-darwin)
            download_deps
            build_darwin
            ;;
        build-windows)
            download_deps
            build_windows
            ;;
        clean)
            clean
            ;;
        deps)
            download_deps
            ;;
        test)
            run_tests
            ;;
        lint)
            lint
            ;;
        package)
            package
            ;;
        checksum)
            generate_checksums
            ;;
        install)
            download_deps
            install_local
            ;;
        release)
            release
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            error "未知命令: $1\n运行 '$0 help' 查看帮助"
            ;;
    esac
}

# 执行主函数
main "$@"