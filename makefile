# ============================================================================
# GoEnv-Switch Makefile 
# ============================================================================

# 项目信息
APP_NAME := goenv-switch
VERSION ?= 1.0.0
BUILD_TIME := $(shell date "+%Y-%m-%d %H:%M:%S")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go 相关
GO := go
GOFMT := gofmt
GOLINT := golangci-lint

# 目录设置
ROOT_DIR := $(shell pwd)
BUILD_DIR := $(ROOT_DIR)/build
DIST_DIR := $(ROOT_DIR)/dist

# 编译参数
LDFLAGS := -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GitCommit=$(GIT_COMMIT)'

# 目标平台
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	linux/386 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/386 \
	windows/arm64

# 默认目标
.DEFAULT_GOAL := build

# 声明伪目标
.PHONY: all build build-all build-linux build-darwin build-windows \
	clean deps test lint fmt vet \
	package checksum install release \
	help info

# ============================================================================
# 帮助信息
# ============================================================================

help: ## 显示帮助信息
	@echo ""
	@echo "GoEnv-Switch Makefile"
	@echo ""
	@echo "用法: make <目标>"
	@echo ""
	@echo "目标:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "示例:"
	@echo "  make build                    # 编译当前平台"
	@echo "  make build-all                # 编译所有平台"
	@echo "  make release VERSION=2.0.0    # 指定版本号发布"
	@echo ""

info: ## 显示项目信息
	@echo "========================================"
	@echo "项目名称:   $(APP_NAME)"
	@echo "版本:       $(VERSION)"
	@echo "构建时间:   $(BUILD_TIME)"
	@echo "Git 提交:   $(GIT_COMMIT)"
	@echo "Go 版本:    $(shell $(GO) version)"
	@echo "构建目录:   $(BUILD_DIR)"
	@echo "发布目录:   $(DIST_DIR)"
	@echo "========================================"

# ============================================================================
# 依赖管理
# ============================================================================

deps: ## 下载依赖
	@echo ">>> 下载依赖..."
	$(GO) mod tidy
	$(GO) mod download
	@echo ">>> 依赖下载完成"

deps-update: ## 更新依赖
	@echo ">>> 更新依赖..."
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo ">>> 依赖更新完成"

deps-verify: ## 验证依赖
	@echo ">>> 验证依赖..."
	$(GO) mod verify
	@echo ">>> 依赖验证完成"

# ============================================================================
# 代码质量
# ============================================================================

fmt: ## 格式化代码
	@echo ">>> 格式化代码..."
	$(GOFMT) -s -w .
	@echo ">>> 格式化完成"

fmt-check: ## 检查代码格式
	@echo ">>> 检查代码格式..."
	@if [ -n "$$($(GOFMT) -l .)" ]; then \
		echo "以下文件需要格式化:"; \
		$(GOFMT) -l .; \
		exit 1; \
	fi
	@echo ">>> 代码格式检查通过"

vet: ## 静态分析
	@echo ">>> 静态分析..."
	$(GO) vet ./...
	@echo ">>> 静态分析完成"

lint: ## 代码检查 (需要安装 golangci-lint)
	@echo ">>> 代码检查..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run ./...; \
	else \
		echo "警告: golangci-lint 未安装，跳过检查"; \
		echo "安装: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi
	@echo ">>> 代码检查完成"

# ============================================================================
# 测试
# ============================================================================

test: ## 运行测试
	@echo ">>> 运行测试..."
	$(GO) test -v ./...
	@echo ">>> 测试完成"

test-coverage: ## 运行测试并生成覆盖率报告
	@echo ">>> 运行测试并生成覆盖率..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo ">>> 覆盖率报告: coverage.html"

test-race: ## 运行竞态检测测试
	@echo ">>> 运行竞态检测..."
	$(GO) test -race -v ./...
	@echo ">>> 竞态检测完成"

benchmark: ## 运行基准测试
	@echo ">>> 运行基准测试..."
	$(GO) test -bench=. -benchmem ./...
	@echo ">>> 基准测试完成"

# ============================================================================
# 编译
# ============================================================================

build: deps ## 编译当前平台
	@echo ">>> 编译当前平台..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) ./cmd
	@echo ">>> 编译完成: $(BUILD_DIR)/$(APP_NAME)"

build-all: deps $(addprefix build-platform-,$(subst /,-,$(PLATFORMS))) ## 编译所有平台
	@echo ">>> 所有平台编译完成"

# 通用平台编译规则
build-platform-%:
	$(eval OS := $(word 1,$(subst -, ,$*)))
	$(eval ARCH := $(word 2,$(subst -, ,$*)))
	$(eval EXT := $(if $(filter windows,$(OS)),.exe,))
	$(eval OUTPUT_DIR := $(BUILD_DIR)/$(OS)_$(ARCH))
	@echo ">>> 编译 $(OS)/$(ARCH)..."
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build \
		-ldflags "$(LDFLAGS)" \
		-o $(OUTPUT_DIR)/$(APP_NAME)$(EXT) ./cmd
	@cp config/config.yaml $(OUTPUT_DIR)/config.yaml 2>/dev/null || true
	@cp README.md $(OUTPUT_DIR)/ 2>/dev/null || true
	@echo ">>> 完成: $(OUTPUT_DIR)/$(APP_NAME)$(EXT)"

build-linux: deps build-platform-linux-amd64 build-platform-linux-arm64 ## 编译 Linux 平台
	@echo ">>> Linux 平台编译完成"

build-darwin: deps build-platform-darwin-amd64 build-platform-darwin-arm64 ## 编译 macOS 平台
	@echo ">>> macOS 平台编译完成"

build-windows: deps build-platform-windows-amd64 build-platform-windows-386 ## 编译 Windows 平台
	@echo ">>> Windows 平台编译完成"

# ============================================================================
# 打包与发布
# ============================================================================

package: ## 打包发布文件
	@echo ">>> 打包发布文件..."
	@mkdir -p $(DIST_DIR)
	@for dir in $(BUILD_DIR)/*; do \
		if [ -d "$$dir" ]; then \
			platform=$$(basename $$dir); \
			archive_name=$(APP_NAME)_$(VERSION)_$$platform; \
			echo ">>> 打包 $$platform..."; \
			cd $(BUILD_DIR); \
			if echo "$$platform" | grep -q "^windows"; then \
				zip -r $(DIST_DIR)/$$archive_name.zip $$platform; \
			else \
				tar -czvf $(DIST_DIR)/$$archive_name.tar.gz $$platform; \
			fi; \
		fi; \
	done
	@echo ">>> 打包完成: $(DIST_DIR)"

checksum: ## 生成校验和
	@echo ">>> 生成校验和..."
	@cd $(DIST_DIR) && \
	if command -v sha256sum >/dev/null 2>&1; then \
		sha256sum *.tar.gz *.zip 2>/dev/null > checksums.txt || true; \
	elif command -v shasum >/dev/null 2>&1; then \
		shasum -a 256 *.tar.gz *.zip 2>/dev/null > checksums.txt || true; \
	else \
		echo "警告: sha256sum 或 shasum 未找到"; \
	fi
	@echo ">>> 校验和已生成: $(DIST_DIR)/checksums.txt"

# ============================================================================
# 安装
# ============================================================================

install: deps ## 安装到 GOPATH/bin
	@echo ">>> 安装到 GOPATH/bin..."
	$(GO) install -ldflags "$(LDFLAGS)" ./cmd
	@echo ">>> 安装完成: $$($(GO) env GOPATH)/bin/$(APP_NAME)"

uninstall: ## 从 GOPATH/bin 卸载
	@echo ">>> 卸载..."
	@rm -f $$($(GO) env GOPATH)/bin/$(APP_NAME)
	@echo ">>> 卸载完成"

# ============================================================================
# 清理
# ============================================================================

clean: ## 清理构建目录
	@echo ">>> 清理构建目录..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -f coverage.out coverage.html
	@echo ">>> 清理完成"

clean-all: clean ## 清理所有（包括依赖缓存）
	@echo ">>> 清理依赖缓存..."
	$(GO) clean -cache -modcache
	@echo ">>> 全部清理完成"

# ============================================================================
# 完整流程
# ============================================================================

release: clean deps fmt-check vet test build-all package checksum ## 完整发布流程
	@echo ""
	@echo "========================================"
	@echo ">>> 发布完成!"
	@echo "    版本: $(VERSION)"
	@echo "    目录: $(DIST_DIR)"
	@echo "========================================"
	@echo ""
	@ls -lh $(DIST_DIR)

quick-release: clean deps build-all package checksum ## 快速发布（跳过测试）
	@echo ""
	@echo "========================================"
	@echo ">>> 快速发布完成!"
	@echo "    版本: $(VERSION)"
	@echo "    目录: $(DIST_DIR)"
	@echo "========================================"
	@echo ""
	@ls -lh $(DIST_DIR)

# ============================================================================
# Docker 支持（可选）
# ============================================================================

docker-build: ## 构建 Docker 镜像
	@echo ">>> 构建 Docker 镜像..."
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo ">>> Docker 镜像构建完成"

docker-push: ## 推送 Docker 镜像
	@echo ">>> 推送 Docker 镜像..."
	docker push $(APP_NAME):$(VERSION)
	docker push $(APP_NAME):latest
	@echo ">>> Docker 镜像推送完成"

# ============================================================================
# 开发辅助
# ============================================================================

run: ## 运行程序
	$(GO) run ./cmd $(ARGS)

dev: deps fmt vet ## 开发模式（格式化 + 静态分析）
	@echo ">>> 开发检查完成"

ci: deps fmt-check vet lint test ## CI 流程
	@echo ">>> CI 检查完成"

# 监听文件变化自动重新编译（需要安装 air）
watch: ## 监听文件变化
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "请先安装 air: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

# 生成初始配置文件
init-config: ## 生成默认配置文件
	@if [ -f config.yaml ]; then \
		echo "config.yaml 已存在"; \
	else \
		echo ">>> 生成默认配置文件..."; \
		$(GO) run ./cmd init; \
		echo ">>> 配置文件已生成"; \
	fi