package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Version
const Version = "1.0.1"

// 颜色常量
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

// GoEnvConfig 单个环境的配置
type GoEnvConfig struct {
	Name      string `yaml:"name"`
	GoPrivate string `yaml:"goprivate"`
	GoProxy   string `yaml:"goproxy"`
	GoSumDB   string `yaml:"gosumdb"`
	GoNoProxy string `yaml:"gonoproxy"`
	GoNoSumDB string `yaml:"gonosumdb"`
}

// Config 完整配置文件结构
type Config struct {
	Environments map[string]GoEnvConfig `yaml:"environments"`
	DefaultEnv   string                 `yaml:"default_env"`
}

// ConfigManager 配置管理器
type ConfigManager struct {
	config     Config
	configPath string
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configPath string) (*ConfigManager, error) {
	cm := &ConfigManager{configPath: configPath}
	if err := cm.loadConfig(); err != nil {
		return nil, err
	}
	return cm, nil
}

// loadConfig 加载配置文件
func (cm *ConfigManager) loadConfig() error {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := yaml.Unmarshal(data, &cm.config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// ListEnvironments 列出所有可用环境
func (cm *ConfigManager) ListEnvironments() {
	fmt.Printf("\n%s可用的环境配置:%s\n", ColorBold+ColorCyan, ColorReset)
	fmt.Println(strings.Repeat("-", 50))

	for key, env := range cm.config.Environments {
		defaultMark := ""
		if key == cm.config.DefaultEnv {
			defaultMark = fmt.Sprintf("%s (默认)%s", ColorGreen, ColorReset)
		}
		fmt.Printf("  %s%-15s%s - %s%s\n", ColorYellow, key, ColorReset, env.Name, defaultMark)
	}
	fmt.Println()
}

// ShowEnvironmentDetail 显示环境详细配置
func (cm *ConfigManager) ShowEnvironmentDetail(envName string) error {
	env, exists := cm.config.Environments[envName]
	if !exists {
		return fmt.Errorf("环境 '%s' 不存在", envName)
	}

	fmt.Printf("\n%s环境 [%s] 的详细配置:%s\n", ColorBold+ColorCyan, envName, ColorReset)
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("  %s名称:%s       %s\n", ColorBold, ColorReset, env.Name)
	fmt.Printf("  %sGOPRIVATE:%s  %s\n", ColorBold, ColorReset, env.GoPrivate)
	fmt.Printf("  %sGOPROXY:%s    %s\n", ColorBold, ColorReset, env.GoProxy)
	fmt.Printf("  %sGOSUMDB:%s    %s\n", ColorBold, ColorReset, env.GoSumDB)
	fmt.Printf("  %sGONOPROXY:%s  %s\n", ColorBold, ColorReset, env.GoNoProxy)
	fmt.Printf("  %sGONOSUMDB:%s  %s\n", ColorBold, ColorReset, env.GoNoSumDB)
	fmt.Println()

	return nil
}

// SwitchEnvironment 切换到指定环境
func (cm *ConfigManager) SwitchEnvironment(envName string) error {
	env, exists := cm.config.Environments[envName]
	if !exists {
		return fmt.Errorf("环境 '%s' 不存在", envName)
	}

	fmt.Printf("\n%s正在切换到环境: %s (%s)%s\n", ColorBold+ColorGreen, envName, env.Name, ColorReset)
	fmt.Println(strings.Repeat("-", 50))

	// 设置各项配置
	settings := map[string]string{
		"GOPRIVATE": env.GoPrivate,
		"GOPROXY":   env.GoProxy,
		"GOSUMDB":   env.GoSumDB,
		"GONOPROXY": env.GoNoProxy,
		"GONOSUMDB": env.GoNoSumDB,
	}

	for key, value := range settings {
		if err := setGoEnv(key, value); err != nil {
			return fmt.Errorf("设置 %s 失败: %w", key, err)
		}
		fmt.Printf("  %s✓%s %s = %s\n", ColorGreen, ColorReset, key, value)
	}

	fmt.Printf("\n%s切换完成!%s\n", ColorGreen, ColorReset)
	return nil
}

// ShowCurrentConfig 显示当前 Go 环境配置
func (cm *ConfigManager) ShowCurrentConfig() error {
	fmt.Printf("\n%s当前 Go 环境配置:%s\n", ColorBold+ColorCyan, ColorReset)
	fmt.Println(strings.Repeat("-", 50))

	envVars := []string{"GOPRIVATE", "GOPROXY", "GOSUMDB", "GONOPROXY", "GONOSUMDB"}

	for _, env := range envVars {
		value, err := getGoEnv(env)
		if err != nil {
			return err
		}
		fmt.Printf("  %s%-12s%s = %s\n", ColorBold, env, ColorReset, value)
	}
	fmt.Println()

	return nil
}

// InteractiveSwitch 交互式环境切换
func (cm *ConfigManager) InteractiveSwitch() error {
	if len(cm.config.Environments) == 0 {
		return fmt.Errorf("没有可用的环境配置")
	}

	fmt.Printf("\n%s请选择要切换的环境:%s\n", ColorBold+ColorCyan, ColorReset)
	fmt.Println(strings.Repeat("-", 50))

	// 显示环境列表
	var envNames []string
	index := 1
	for key, env := range cm.config.Environments {
		defaultMark := ""
		if key == cm.config.DefaultEnv {
			defaultMark = fmt.Sprintf("%s (默认)%s", ColorGreen, ColorReset)
		}
		fmt.Printf("  %s%d.%s %-15s - %s%s\n", ColorYellow, index, ColorReset, key, env.Name, defaultMark)
		envNames = append(envNames, key)
		index++
	}

	fmt.Printf("\n%s请输入环境编号 (1-%d) 或环境名称: %s", ColorBold, len(envNames), ColorReset)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var selectedEnv string

	// 检查是否为数字选择
	if num, err := fmt.Sscanf(input, "%d", new(int)); err == nil && num == 1 {
		var choice int
		fmt.Sscanf(input, "%d", &choice)
		if choice < 1 || choice > len(envNames) {
			return fmt.Errorf("无效的选择: %d", choice)
		}
		selectedEnv = envNames[choice-1]
	} else {
		// 直接输入环境名称
		selectedEnv = input
	}

	// 验证环境是否存在
	if _, exists := cm.config.Environments[selectedEnv]; !exists {
		return fmt.Errorf("环境 '%s' 不存在", selectedEnv)
	}

	return cm.SwitchEnvironment(selectedEnv)
}

// setGoEnv 设置 Go 环境变量
func setGoEnv(key, value string) error {
	cmd := exec.Command("go", "env", "-w", fmt.Sprintf("%s=%s", key, value))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

// getGoEnv 获取 Go 环境变量
func getGoEnv(key string) (string, error) {
	cmd := exec.Command("go", "env", key)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getDefaultConfigPath 获取默认配置文件路径
func getDefaultConfigPath() string {
	// 优先使用当前目录的配置文件
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml"
	}

	// 其次使用用户主目录下的配置文件
	return defaultConfigUserPath()
}

// defaultConfigUserPath 获取用户主目录下的配置文件路径
func defaultConfigUserPath() string {
	home, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(home, ".goenv-switch", "config.yaml")

	}
	return "config.yaml"
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(rawContent []byte, path string) error {
	// 创建目录
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, rawContent, 0644)
}

func printUsage() {
	fmt.Printf(`%sGo 环境配置切换工具%s

%s用法:%s
  %sgoenv-switch%s <命令> [参数]

%s命令:%s
  %slist%s              列出所有可用环境
  %sshow%s <环境名>     显示指定环境的详细配置
  %sswitch%s <环境名>   切换到指定环境
  %sinteractive%s       交互式环境切换
  %scurrent%s           显示当前 Go 环境配置
  %sinit%s              在当前目录创建默认配置文件

%s选项:%s
  %s-c, --config%s      指定配置文件路径

%s示例:%s
  %sgoenv-switch list%s
  %sgoenv-switch switch company%s
  %sgoenv-switch interactive%s
  %sgoenv-switch current%s
  %sgoenv-switch show company%s
  %sgoenv-switch init%s
  %sgoenv-switch -c /path/to/config.yaml switch company%s

`,
		ColorBold+ColorCyan, ColorReset,
		ColorBold, ColorReset, ColorGreen, ColorReset,
		ColorBold, ColorReset, ColorYellow, ColorReset, ColorYellow, ColorReset,
		ColorYellow, ColorReset, ColorYellow, ColorReset, ColorYellow, ColorReset,
		ColorYellow, ColorReset,
		ColorBold, ColorReset, ColorYellow, ColorReset,
		ColorBold, ColorReset, ColorGreen, ColorReset, ColorGreen, ColorReset,
		ColorGreen, ColorReset, ColorGreen, ColorReset, ColorGreen, ColorReset,
		ColorGreen, ColorReset, ColorGreen, ColorReset)
}
func printWelcome() {
	fmt.Printf(`%s
 ██████╗  ██████╗ ███████╗███╗   ██╗██╗   ██╗███████╗██╗    ██╗██╗████████╗ ██████╗██╗  ██╗
██╔════╝ ██╔═══██╗██╔════╝████╗  ██║██║   ██║██╔════╝██║    ██║██║╚══██╔══╝██╔════╝██║  ██║
██║  ███╗██║   ██║█████╗  ██╔██╗ ██║██║   ██║█████╗  ██║ █╗ ██║██║   ██║   ██║     ███████║
██║   ██║██║   ██║██╔══╝  ██║╚██╗██║╚██╗ ██╔╝██╔══╝  ██║███╗██║██║   ██║   ██║     ██╔══██║
╚██████╔╝╚██████╔╝███████╗██║ ╚████║ ╚████╔╝ ███████╗╚███╔███╔╝██║   ██║   ╚██████╗██║  ██║
 ╚═════╝  ╚═════╝ ╚══════╝╚═╝  ╚═══╝  ╚═══╝  ╚══════╝ ╚══╝╚══╝ ╚═╝   ╚═╝    ╚═════╝╚═╝  ╚═╝
                                                                                             %s
%sGo 环境配置切换工具 v%s%s
%s`,
		ColorCyan, ColorReset, ColorBold+ColorGreen, Version, ColorReset, ColorReset)
}

func main() {
	printWelcome()
	args := os.Args[1:]

	// 如果没有参数，显示欢迎信息和帮助
	if len(args) == 0 {
		printUsage()
		return
	}

	// 解析参数
	configPath := getDefaultConfigPath()
	var command string
	var commandArgs []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-c", "--config":
			if i+1 < len(args) {
				configPath = args[i+1]
				i++
			} else {
				fmt.Printf("%s错误: 缺少配置文件路径%s\n", ColorRed, ColorReset)
				os.Exit(1)
			}
		case "-h", "--help", "help":
			printUsage()
			return
		case "-v", "--version":
			fmt.Printf("goenv-switch version %s\n", Version)
			return
		default:
			if command == "" {
				command = args[i]
			} else {
				commandArgs = append(commandArgs, args[i])
			}
		}
	}

	// 处理 init 命令（不需要加载配置）
	if command == "init" {
		// 从configPath中读取配置文本内容，当作纯文本读取
		content, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("%s读取配置文件失败: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}
		userPath := defaultConfigUserPath()

		if err := createDefaultConfig(content, userPath); err != nil {
			fmt.Printf("%s创建配置文件失败: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}
		fmt.Printf("%s配置文件已创建: %s%s\n", ColorGreen, userPath, ColorReset)
		return
	}

	// 加载配置管理器
	cm, err := NewConfigManager(configPath)
	if err != nil {
		fmt.Printf("%s错误: %v%s\n", ColorRed, err, ColorReset)
		fmt.Printf("%s提示: 运行 'goenv-switch init' 创建默认配置文件%s\n", ColorYellow, ColorReset)
		os.Exit(1)
	}

	// 执行命令
	switch command {
	case "list":
		cm.ListEnvironments()

	case "show":
		if len(commandArgs) == 0 {
			fmt.Printf("%s错误: 请指定环境名称%s\n", ColorRed, ColorReset)
			os.Exit(1)
		}
		if err := cm.ShowEnvironmentDetail(commandArgs[0]); err != nil {
			fmt.Printf("%s错误: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}

	case "switch":
		if len(commandArgs) == 0 {
			fmt.Printf("%s错误: 请指定要切换的环境名称%s\n", ColorRed, ColorReset)
			cm.ListEnvironments()
			os.Exit(1)
		}
		if err := cm.SwitchEnvironment(commandArgs[0]); err != nil {
			fmt.Printf("%s错误: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}

	case "interactive", "i":
		if err := cm.InteractiveSwitch(); err != nil {
			fmt.Printf("%s错误: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}

	case "current":
		if err := cm.ShowCurrentConfig(); err != nil {
			fmt.Printf("%s错误: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}

	default:
		fmt.Printf("%s未知命令: %s%s\n", ColorRed, command, ColorReset)
		printUsage()
		os.Exit(1)
	}
}
