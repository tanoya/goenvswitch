# GoEnv-Switch

一个专业的 Go 语言环境配置切换工具，支持通过配置文件管理多个环境配置，可以在不同的 GOPRIVATE、GOPROXY、GOSUMDB 等配置之间快速切换。

## 功能特性

- 🔧 **配置文件驱动** - 通过 YAML 配置文件管理多个环境
- 🌐 **支持多环境** - 可定义任意数量的环境配置（公司内网、公共环境、混合环境等）
- 📦 **完整配置项** - 支持 GOPRIVATE、GOPROXY、GOSUMDB、GONOPROXY、GONOSUMDB
- 💡 **友好提示** - 清晰的命令行输出和错误提示
- 🎯 **灵活配置** - 支持指定配置文件路径
- 🚀 **一键初始化** - 自动创建默认配置文件模板

## 安装

### 从源码编译

```bash
# 克隆项目
git clone https://github.com/yourusername/goenv-switch.git
cd goenv-switch

# 下载依赖
go mod tidy

# 编译
go build -o goenv-switch .

# 安装到 GOPATH/bin（可选，方便全局使用）
go install .
```

### 验证安装

```bash
./goenv-switch --help
```

## 快速开始

### 1. 初始化配置文件

```bash
# 在当前目录创建默认配置文件
./goenv-switch init

# 或指定路径
./goenv-switch init /path/to/config.yaml
```

### 2. 编辑配置文件

根据你的需求修改 `config.yaml` 文件，添加或修改环境配置。

### 3. 切换环境

```bash
# 切换到公司内网环境
./goenv-switch switch company

# 切换到公共环境
./goenv-switch switch public
```

## 命令参考

| 命令 | 说明 | 示例 |
|------|------|------|
| `list` | 列出所有可用环境 | `goenv-switch list` |
| `show <环境名>` | 显示指定环境的详细配置 | `goenv-switch show company` |
| `switch <环境名>` | 切换到指定环境 | `goenv-switch switch public` |
| `current` | 显示当前 Go 环境配置 | `goenv-switch current` |
| `init` | 在当前目录创建默认配置文件 | `goenv-switch init` |
| `help` | 显示帮助信息 | `goenv-switch help` |

### 命令选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `-c, --config` | 指定配置文件路径 | `goenv-switch -c /path/to/config.yaml switch company` |

## 配置文件说明

配置文件使用 YAML 格式，默认查找顺序：

1. 当前目录下的 `config.yaml`
2. `~/.goenv-switch/config.yaml`

### 配置文件示例

```yaml
# Go 环境配置切换工具配置文件

environments:
  # 公司内网环境
  company:
    name: "公司内网环境"
    goprivate: "git.company.com,gitlab.internal.com"
    goproxy: "https://goproxy.company.com,direct"
    gosumdb: "off"
    gonoproxy: "git.company.com"
    gonosumdb: "git.company.com"

  # 公共环境
  public:
    name: "公共环境"
    goprivate: ""
    goproxy: "https://goproxy.cn,https://goproxy.io,direct"
    gosumdb: "sum.golang.org"
    gonoproxy: ""
    gonosumdb: ""

  # 混合环境
  hybrid:
    name: "混合环境"
    goprivate: "github.com/mycompany/*,git.internal.com"
    goproxy: "https://goproxy.cn,direct"
    gosumdb: "sum.golang.org"
    gonoproxy: "github.com/mycompany/*"
    gonosumdb: "github.com/mycompany/*"

```

### 配置项说明

| 配置项 | 说明 |
|--------|------|
| `name` | 环境的友好名称，用于显示 |
| `goprivate` | 私有模块路径模式，多个用逗号分隔 |
| `goproxy` | Go 模块代理地址，多个用逗号分隔 |
| `gosumdb` | 校验和数据库地址，设为 `off` 可禁用 |
| `gonoproxy` | 不使用代理的模块路径模式 |
| `gonosumdb` | 不进行校验和验证的模块路径模式 |

## 使用示例

### 查看所有可用环境

```bash
$ ./goenv-switch list

可用的环境配置:
--------------------------------------------------
  company         - 公司内网环境
● public          - 公共环境
  hybrid          - 混合环境
```

### 查看当前配置

```bash
$ ./goenv-switch current

当前 Go 环境配置:
--------------------------------------------------
  GOPRIVATE    = 
  GOPROXY      = https://goproxy.cn,direct
  GOSUMDB      = sum.golang.org
  GONOPROXY    = 
  GONOSUMDB    = 
```

### 查看环境详情

```bash
$ ./goenv-switch show company

环境 [company] 的详细配置:
--------------------------------------------------
  名称:       公司内网环境
  GOPRIVATE:  git.company.com,gitlab.internal.com
  GOPROXY:    https://goproxy.company.com,direct
  GOSUMDB:    off
  GONOPROXY:  git.company.com
  GONOSUMDB:  git.company.com
```

### 切换环境

```bash
$ ./goenv-switch switch company

正在切换到环境: company (公司内网环境)
--------------------------------------------------
  ✓ GOPRIVATE = git.company.com,gitlab.internal.com
  ✓ GOPROXY = https://goproxy.company.com,direct
  ✓ GOSUMDB = off
  ✓ GONOPROXY = git.company.com
  ✓ GONOSUMDB = git.company.com

切换完成!
```

### 使用指定配置文件

```bash
$ ./goenv-switch -c ~/my-config.yaml switch company
```

## 常见场景

### 场景一：公司内网开发

公司使用私有 Git 仓库和内部 Go 代理：

```yaml
company:
  name: "公司内网环境"
  goprivate: "git.company.com,gitlab.internal.com"
  goproxy: "https://goproxy.company.com,direct"
  gosumdb: "off"
  gonoproxy: "git.company.com"
  gonosumdb: "git.company.com"
```

### 场景二：开源项目开发

使用公共代理和校验和数据库：

```yaml
public:
  name: "公共环境"
  goprivate: ""
  goproxy: "https://goproxy.cn,https://goproxy.io,direct"
  gosumdb: "sum.golang.org"
  gonoproxy: ""
  gonosumdb: ""
```

### 场景三：混合开发

同时访问公司私有仓库和公共仓库：

```yaml
hybrid:
  name: "混合环境"
  goprivate: "github.com/mycompany/*,git.internal.com"
  goproxy: "https://goproxy.cn,direct"
  gosumdb: "sum.golang.org"
  gonoproxy: "github.com/mycompany/*"
  gonosumdb: "github.com/mycompany/*"
```

## 项目结构

```
goenv-switch/
├── 📁 .git/                          # Git 版本控制目录
├── 📄 .gitignore                     # Git 忽略文件配置
├── 📄 README.md                      # 项目详细说明文档
├── 📄 go.mod                         # Go 模块依赖管理文件
├── 📄 go.sum                         # Go 依赖校验文件
├── 📄 makefile                       # 构建脚本配置
├── 📄 build.sh                       # 构建脚本
├── 📁 cmd/                           # 主程序目录
│   └── 📄 main.go                    # 程序入口文件
└── 📁 config/                        # 配置文件目录
    └── 📄 config.yaml                # 默认配置文件示例
```

## 依赖

- Go 1.21 或更高版本
- [gopkg.in/yaml.v3](https://github.com/go-yaml/yaml) - YAML 解析库

## 常见问题

### Q: 配置切换后不生效？

A: 本工具使用 `go env -w` 命令设置环境变量，这会修改 Go 的全局配置。如果你的 shell 中设置了相同的环境变量，shell 环境变量会覆盖全局配置。请检查你的 `.bashrc`、`.zshrc` 等配置文件。

### Q: 如何添加新的环境配置？

A: 编辑 `config.yaml` 文件，在 `environments` 下添加新的配置块即可。

### Q: 配置文件放在哪里最好？

A: 建议放在 `~/.goenv-switch/config.yaml`，这样可以全局使用，不受当前目录影响。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
