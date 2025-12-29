package main

// import (
// 	"flag"
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"

// 	"gopkg.in/yaml.v3"
// )

// // Version 信息（编译时注入）
// var (
// 	Version   = "dev"
// 	BuildTime = "unknown"
// 	GitCommit = "unknown"
// )

// // GoEnvConfig 单个环境的配置
// type GoEnvConfig struct {
// 	Name      string `yaml:"name"`
// 	GoPrivate string `yaml:"goprivate"`
// 	GoProxy   string `yaml:"goproxy"`
// 	GoSumDB   string `yaml:"gosumdb"`
// 	GoNoProxy string `yaml:"gonoproxy"`
// 	GoNoSumDB string `yaml:"gonosumdb"`
// }

// // Config 完整配置文件结构
// type Config struct {
// 	Environments map[string]GoEnvConfig `yaml:"environments"`
// 	DefaultEnv   string                 `yaml:"default_env"`
// }

// // ConfigManager 配置管理器
// type ConfigManager struct {
// 	config     Config
// 	configPath string
// }

// // NewConfigManager 创建配置管理器
// func NewConfigManager(configPath string) (*ConfigManager, error) {
// 	cm := &ConfigManager{configPath: configPath}
// 	if err := cm.loadConfig(); err != nil {
// 		return nil, err
// 	}
// 	return cm, nil
// }

// // loadConfig 加载配置文件
// func (cm *ConfigManager) loadConfig() error {
// 	data, err := os.ReadFile(cm.configPath)
// 	if err != nil {
// 		return fmt.Errorf("读取配置文件失败: %w", err)
// 	}

// 	if err := yaml.Unmarshal(data, &cm.config); err != nil {
// 		return fmt.Errorf("解析配置文件失败: %w", err)
// 	}

// 	// 验证配置
// 	if len(cm.config.Environments) == 0 {
// 		return fmt.Errorf("配置文件中没有定义任何环境")
// 	}

// 	return nil
// }

// // ListEnvironments 列出所有可用环境
// func (cm *ConfigManager) ListEnvironments() {
// 	fmt.Println("\n可用的环境配置:")
// 	fmt.Println(strings.Repeat("-", 60))

// 	for key, env := range cm.config.Environments {
// 		defaultMark := ""
// 		if key == cm.config.DefaultEnv {
// 			defaultMark = " (默认)"
// 		}
// 		fmt.Printf("  %-20s - %s%s\n", key, env.Name, defaultMark)
// 	}
// 	fmt.Println()
// }

// // ShowEnvironmentDetail 显示环境详细配置
// func (cm *ConfigManager) ShowEnvironmentDetail(envName string) error {
// 	env, exists := cm.config.Environments[envName]
// 	if !exists {
// 		return fmt.Errorf("环境 '%s' 不存在", envName)
// 	}

// 	fmt.Printf("\n环境 [%s] 的详细配置:\n", envName)
// 	fmt.Println(strings.Repeat("-", 60))
// 	fmt.Printf("  名称:       %s\n", env.Name)
// 	fmt.Printf("  GOPRIVATE:  %s\n", formatValue(env.GoPrivate))
// 	fmt.Printf("  GOPROXY:    %s\n", formatValue(env.GoProxy))
// 	fmt.Printf("  GOSUMDB:    %s\n", formatValue(env.GoSumDB))
// 	fmt.Printf("  GONOPROXY:  %s\n", formatValue(env.GoNoProxy))
// 	fmt.Printf("  GONOSUMDB:  %s\n", formatValue(env.GoNoSumDB))
// 	fmt.Println()

// 	return nil
// }

// // SwitchEnvironment 切换到指定环境
// func (cm *ConfigManager) SwitchEnvironment(envName string) error {
// 	env, exists := cm.config.Environments[envName]
// 	if !exists {
// 		return fmt.Errorf("环境 '%s' 不存在", envName)
// 	}

// 	fmt.Printf("\n正在切换到环境: %s (%s)\n", envName, env.Name)
// 	fmt.Println(strings.Repeat("-", 60))

// 	// 设置各项配置
// 	settings := map[string]string{
// 		"GOPRIVATE": env.GoPrivate,
// 		"GOPROXY":   env.GoProxy,
// 		"GOSUMDB":   env.GoSumDB,
// 		"GONOPROXY": env.GoNoProxy,
// 		"GONOSUMDB": env.GoNoSumDB,
// 	}

// 	for key, value := range settings {
// 		if err := setGoEnv(key, value); err != nil {
// 			return fmt.Errorf("设置 %s 失败: %w", key, err)
// 		}
// 		fmt.Printf("  ✓ %s = %s\n", key, formatValue(value))
// 	}

// 	fmt.Println("\n切换完成!")
// 	return nil
// }

// // ShowCurrentConfig 显示当前 Go 环境配置
// func (cm *ConfigManager) ShowCurrentConfig() error {
// 	fmt.Println("\n当前 Go 环境配置:")
// 	fmt.Println(strings.Repeat("-", 60))

// 	envVars := []string{"GOPRIVATE", "GOPROXY", "GOSUMDB", "GONOPROXY", "GONOSUMDB"}

// 	for _, env := range envVars {
// 		value, err := getGoEnv(env)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Printf("  %-12s = %s\n", env, formatValue(value))
// 	}
// 	fmt.Println()

// 	return nil
// }

// // formatValue 格式化输出值
// func formatValue(value string) string {
// 	if value == "" {
// 		return "(未设置)"
// 	}
// 	return value
// }

// // setGoEnv 设置 Go 环境变量
// func setGoEnv(key, value string) error {
// 	cmd := exec.Command("go", "env", "-w", fmt.Sprintf("%s=%s", key, value))
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("%s: %s", err, string(output))
// 	}
// 	return nil
// }

// // getGoEnv 获取 Go 环境变量
// func getGoEnv(key string) (string, error) {
// 	cmd := exec.Command("go", "env", key)
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return "", fmt.Errorf("获取 %s 失败: %w", key, err)
// 	}
// 	return strings.TrimSpace(string(output)), nil
// }

// // getDefaultConfigPath 获取默认配置文件路径
// func getDefaultConfigPath() string {
// 	// 优先使用当前目录的配置文件
// 	if _, err := os.Stat("config.yaml"); err == nil {
// 		return "config.yaml"
// 	}

// 	// 其次使用用户主目录下的配置文件
// 	home, err := os.UserHomeDir()
// 	if err == nil {
// 		configPath := filepath.Join(home, ".goenv-switch", "config.yaml")
// 		if _, err := os.Stat(configPath); err == nil {
// 			return configPath
// 		}
// 	}

// 	return "config.yaml"
// }

// // createDefaultConfig 创建默认配置文件
// func createDefaultConfig(path string) error {
// 	defaultConfig := `# Go 环境配置切换工具配置文件

// environments:
//   # 公司内网环境
//   company:
//     name: "公司内网环境"
//     goprivate: "git.company.com"
//     goproxy: "https://goproxy.company.com,direct"
//     gosumdb: "off"
//     gonoproxy: "git.company.com"
//     gonosumdb: "git.company.com"

//   # 公共环境
//   public:
//     name: "公共环境"
//     goprivate: ""
//     goproxy: "https://goproxy.cn,https://goproxy.io,direct"
//     gosumdb: "sum.golang.org"
//     gonoproxy: ""
//     gonosumdb: ""

// # 默认使用的环境
// default_env: public
// `
// 	// 创建目录
// 	dir := filepath.Dir(path)
// 	if err := os.MkdirAll(dir, 0755); err != nil {
// 		return err
// 	}

// 	return os.WriteFile(path, []byte(defaultConfig), 0644)
// }

// // printUsage 打印使用说明
// func printUsage() {
// 	fmt.Printf(`
// GoEnv-Switch - Go 环境配置切换工具 v%s

// 用法:
//   goenv-switch <命令> [参数]

// 命令:
//   list              列出所有可用环境
//   show <环境名>     显示指定环境的详细配置
//   switch <环境名>   切换到指定环境
//   current           显示当前 Go 环境配置
//   init              在当前目录创建默认配置文件
//   version           显示版本信息

// 选项:
//   -c, --config      指定配置文件路径
//   -h, --help        显示帮助信息
//   -v, --version     显示版本信息

// 示例:
//   goenv-switch list
//   goenv-switch switch company
//   goenv-switch current
//   goenv-switch show company
//   goenv-switch init
//   goenv-switch -c /path/to/config.yaml switch company

// `, Version)
// }

// // printVersion 打印版本信息
// func printVersion() {
// 	fmt.Printf("GoEnv-Switch v%s\n", Version)
// 	fmt.Printf("Build Time: %s\n", BuildTime)
// 	fmt.Printf("Git Commit: %s\n", GitCommit)
// }

// func main() {
// 	// 定义命令行标志
// 	configPath := flag.String("c", "", "配置文件路径")
// 	flag.StringVar(configPath, "config", "", "配置文件路径（长形式）")
// 	helpFlag := flag.Bool("h", false, "显示帮助信息")
// 	flag.BoolVar(helpFlag, "help", false, "显示帮助信息（长形式）")
// 	versionFlag := flag.Bool("v", false, "显示版本信息")
// 	flag.BoolVar(versionFlag, "version", false, "显示版本信息（长形式）")

// 	flag.Usage = printUsage
// 	flag.Parse()

// 	// 处理帮助标志
// 	if *helpFlag {
// 		printUsage()
// 		return
// 	}

// 	// 处理版本标志
// 	if *versionFlag {
// 		printVersion()
// 		return
// 	}

// 	args := flag.Args()

// 	// 如果没有命令，显示帮助
// 	if len(args) == 0 {
// 		printUsage()
// 		return
// 	}

// 	command := args[0]
// 	commandArgs := args[1:]

// 	// 处理 init 命令（不需要加载配置）
// 	if command == "init" {
// 		path := *configPath
// 		if path == "" {
// 			path = "config.yaml"
// 		}
// 		if err := createDefaultConfig(path); err != nil {
// 			fmt.Fprintf(os.Stderr, "错误: 创建配置文件失败: %v\n", err)
// 			os.Exit(1)
// 		}
// 		fmt.Printf("✓ 配置文件已创建: %s\n", path)
// 		return
// 	}

// 	// 处理 version 命令
// 	if command == "version" {
// 		printVersion()
// 		return
// 	}

// 	// 获取配置文件路径
// 	if *configPath == "" {
// 		*configPath = getDefaultConfigPath()
// 	}

// 	// 加载配置管理器
// 	cm, err := NewConfigManager(*configPath)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
// 		fmt.Fprintf(os.Stderr, "提示: 运行 'goenv-switch init' 创建默认配置文件\n")
// 		os.Exit(1)
// 	}

// 	// 执行命令
// 	switch command {
// 	case "list":
// 		cm.ListEnvironments()

// 	case "show":
// 		if len(commandArgs) == 0 {
// 			fmt.Fprintf(os.Stderr, "错误: 请指定环境名称\n")
// 			os.Exit(1)
// 		}
// 		if err := cm.ShowEnvironmentDetail(commandArgs[0]); err != nil {
// 			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
// 			os.Exit(1)
// 		}

// 	case "switch":
// 		if len(commandArgs) == 0 {
// 			fmt.Fprintf(os.Stderr, "错误: 请指定要切换的环境名称\n")
// 			cm.ListEnvironments()
// 			os.Exit(1)
// 		}
// 		if err := cm.SwitchEnvironment(commandArgs[0]); err != nil {
// 			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
// 			os.Exit(1)
// 		}

// 	case "current":
// 		if err := cm.ShowCurrentConfig(); err != nil {
// 			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
// 			os.Exit(1)
// 		}

// 	case "help":
// 		printUsage()

// 	default:
// 		fmt.Fprintf(os.Stderr, "错误: 未知命令: %s\n", command)
// 		printUsage()
// 		os.Exit(1)
// 	}
// }
