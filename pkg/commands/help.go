package commands

import (
	"flag"

	"github.com/icarus-go/netspeed/pkg/command"
)

// HelpCommand 帮助命令
type HelpCommand struct {
	enabled *bool
}

// NewHelpCommand 创建帮助命令
func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

// Name 返回命令名称
func (c *HelpCommand) Name() string {
	return "help"
}

// Description 返回命令描述
func (c *HelpCommand) Description() string {
	return "显示帮助信息"
}

// DefineFlags 定义命令的 flag 参数
func (c *HelpCommand) DefineFlags(flags *flag.FlagSet) {
	c.enabled = flags.Bool("help", false, "显示帮助信息")
}

// Execute 执行命令
func (c *HelpCommand) Execute(ctx *command.Context) error {
	if !*c.enabled {
		return nil
	}

	c.showHelp()
	return nil
}

// Priority 返回命令优先级
func (c *HelpCommand) Priority() int {
	return 1 // 最高优先级
}

// showHelp 显示帮助信息
func (c *HelpCommand) showHelp() {
	println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	println("🚀 NetSpeed - 跨平台网络质量检测工具")
	println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	println()
	println("用法:")
	println("  netspeed [选项]")
	println()
	println("选项:")
	println("  -test             测试网站速度并以表格形式输出")
	println("  -ip               获取当前 IP 地理信息")
	println("  -purity           检测 IP 纯净度和风险评分")
	println("  -proxy <url>      设置代理 (支持 http://, socks5://, https://)")
	println("  -watch <秒>       持续监控模式，指定刷新间隔（秒）")
	println("  -config <文件>    自定义测试站点配置文件（JSON 格式）")
	println("  -timeout <秒>     请求超时时间（默认 10 秒）")
	println("  -help             显示此帮助信息")
	println()
	println("示例:")
	println("  netspeed -test")
	println("  netspeed -ip")
	println("  netspeed -purity")
	println("  netspeed -test -watch 30")
	println("  netspeed -test -proxy socks5://127.0.0.1:1080")
	println("  netspeed -test -proxy http://proxy.example.com:8080")
	println("  netspeed -test -config sites.example.json")
	println("  netspeed -test -timeout 5")
	println()
	println("配置文件格式 (JSON):")
	println(`  [
    {"Name": "Google", "URL": "https://www.google.com"},
    {"Name": "GitHub", "URL": "https://github.com"}
  ]`)
	println()
	println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
