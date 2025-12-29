# GoEnv-Switch 快速开始指南

## 当前工程状态

### ✅ 已修复的问题

1. **路径问题** - build.sh 和 Makefile 中的路径已修复
   - `main.go` → `cmd/main.go`
   - `config.yaml` → `config/config.yaml`

2. **Go 代码优化** - 创建了改进版本 `cmd/main_improved.go`
   - 添加了版本变量支持
   - 使用 `flag` 包解析命令行参数
   - 改进了错误处理（使用 stderr）
   - 添加了 `formatValue` 函数
   - 添加了配置验证

### 📁 当前文件结构

```
goenv-switch/
├── cmd/
│   ├── main.go              # 原始版本
│   └── main_improved.go     # 优化版本（推荐使用）
├── config/
│   └── config.yaml          # 配置文件
├── build.sh                 # 编译脚本（已修复）
├── makefile                 # Makefile（已修复）
├── test_build.sh            # 测试脚本（新增）
├── CODE_REVIEW.md           # 代码审查报告（新增）
├── QUICK_START.md           # 本文件
├── README.md                # 项目说明
├── go.mod                   # Go 模块文件
└── go.sum                   # 依赖校验文件
```

## 使用哪个版本？

### 选项 1: 使用优化版本（推荐）

```bash
# 1. 重命名文件
mv cmd/main.go cmd/main_old.go
mv cmd/main_improved.go cmd/main.go

# 2. 编译
make build

# 或
./build.sh build
```

### 选项 2: 使用原始版本

原始版本也可以工作，但缺少一些功能：
- 没有版本信息
- 没有 `--version` 命令
- 命令行参数解析较简单

```bash
# 直接编译
make build
```

## 快速测试

### 运行测试脚本

```bash
chmod +x test_build.sh
./test_build.sh
```

这个脚本会：
1. 检查 Go 环境
2. 检查文件结构
3. 下载依赖
4. 测试编译
5. 测试运行
6. 测试 Makefile

### 手动测试

```bash
# 1. 下载依赖
make deps

# 2. 编译
make build

# 3. 测试运行
./build/goenv-switch --help
./build/goenv-switch list
```

## 编译选项

### 使用 Makefile（推荐）

```bash
# 查看所有命令
make help

# 编译当前平台
make build

# 编译所有平台
make build-all

# 完整发布流程
make release VERSION=1.0.0

# 安装到 GOPATH/bin
make install
```

### 使用 build.sh

```bash
# 赋予执行权限
chmod +x build.sh

# 查看帮助
./build.sh help

# 编译当前平台
./build.sh build

# 编译所有平台
./build.sh build-all

# 完整发布
VERSION=1.0.0 ./build.sh release
```

### 直接使用 go 命令

```bash
# 编译
go build -o goenv-switch ./cmd

# 安装
go install ./cmd

# 运行
go run ./cmd list
```

## 验证编译结果

### 检查可执行文件

```bash
# 查看编译输出
ls -lh build/

# 运行程序
./build/goenv-switch --version
./build/goenv-switch --help
./build/goenv-switch list
```

### 测试功能

```bash
# 1. 初始化配置文件
./build/goenv-switch init

# 2. 查看环境列表
./build/goenv-switch list

# 3. 查看当前配置
./build/goenv-switch current

# 4. 查看环境详情
./build/goenv-switch show public

# 5. 切换环境（需要 Go 环境）
./build/goenv-switch switch public
```

## 常见问题

### Q1: 编译时找不到 main.go

**问题**: `can't load package: package .: no Go files in ...`

**解决**: 确保使用正确的路径
```bash
# 错误
go build .

# 正确
go build ./cmd
```

### Q2: 配置文件路径错误

**问题**: 运行时找不到 config.yaml

**解决**: 
```bash
# 方法 1: 在根目录创建配置文件
./build/goenv-switch init

# 方法 2: 指定配置文件路径
./build/goenv-switch -c config/config.yaml list
```

### Q3: 版本信息显示为 "dev"

**问题**: 版本信息没有正确注入

**解决**: 使用 Makefile 或 build.sh 编译
```bash
# 使用 Makefile
make build VERSION=1.0.0

# 使用 build.sh
VERSION=1.0.0 ./build.sh build
```

### Q4: make 命令不可用

**问题**: 系统没有安装 make

**解决**: 
```bash
# macOS
brew install make

# Ubuntu/Debian
sudo apt-get install build-essential

# 或直接使用 build.sh
./build.sh build
```

## 下一步

### 1. 选择版本并编译

```bash
# 推荐：使用优化版本
mv cmd/main.go cmd/main_old.go
mv cmd/main_improved.go cmd/main.go
make build
```

### 2. 测试功能

```bash
./build/goenv-switch init
./build/goenv-switch list
```

### 3. 安装使用

```bash
make install
goenv-switch --help
```

### 4. 多平台编译（可选）

```bash
make build-all
ls -lh build/
```

### 5. 打包发布（可选）

```bash
make release VERSION=1.0.0
ls -lh dist/
```

## 总结

### 当前状态
- ✅ **可以编译** - 路径问题已修复
- ✅ **可以运行** - 核心功能正常
- ✅ **脚本可用** - build.sh 和 Makefile 都可以使用
- ⚠️ **建议优化** - 使用 main_improved.go 获得更好的体验

### 推荐流程

```bash
# 1. 运行测试
chmod +x test_build.sh
./test_build.sh

# 2. 使用优化版本
mv cmd/main.go cmd/main_old.go
mv cmd/main_improved.go cmd/main.go

# 3. 编译安装
make install

# 4. 开始使用
goenv-switch init
goenv-switch list
goenv-switch switch public
```

## 需要帮助？

查看以下文档：
- `README.md` - 项目说明和使用文档
- `CODE_REVIEW.md` - 详细的代码审查报告
- `make help` - Makefile 命令帮助
- `./build.sh help` - build.sh 命令帮助
