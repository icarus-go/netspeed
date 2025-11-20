package commands

import (
	"flag"
	"fmt"

	"github.com/icarus-go/net-speed/pkg/command"
	"github.com/icarus-go/net-speed/pkg/ipinfo"
)

// IPCommand IP 检测命令
type IPCommand struct {
	enabled *bool
}

// NewIPCommand 创建 IP 检测命令
func NewIPCommand() *IPCommand {
	return &IPCommand{}
}

// Name 返回命令名称
func (c *IPCommand) Name() string {
	return "ip"
}

// Description 返回命令描述
func (c *IPCommand) Description() string {
	return "获取当前 IP 地理信息"
}

// DefineFlags 定义命令的 flag 参数
func (c *IPCommand) DefineFlags(flags *flag.FlagSet) {
	c.enabled = flags.Bool("ip", false, "获取当前 IP 地理信息")
}

// Execute 执行命令
func (c *IPCommand) Execute(ctx *command.Context) error {
	if !*c.enabled {
		return nil
	}

	fmt.Println("🌐 正在获取 IP 信息...")
	fmt.Println()

	detector := ipinfo.NewDetector(ctx.HTTPClient)
	info, err := detector.Detect()
	if err != nil {
		return fmt.Errorf("获取 IP 信息失败: %v", err)
	}

	c.displayIPInfo(info)
	return nil
}

// Priority 返回命令优先级
func (c *IPCommand) Priority() int {
	return 10
}

// displayIPInfo 显示 IP 信息
func (c *IPCommand) displayIPInfo(info *ipinfo.IPInfo) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📍 IP 地址:    %s\n", info.IP)
	fmt.Printf("🌍 国家:       %s (%s)\n", info.Country, info.CountryCode)
	if info.Region != "" {
		fmt.Printf("📌 地区:       %s\n", info.Region)
	}
	if info.City != "" {
		fmt.Printf("🏙️  城市:       %s\n", info.City)
	}
	if info.ISP != "" {
		fmt.Printf("🔌 ISP:        %s\n", info.ISP)
	}
	if info.Org != "" {
		fmt.Printf("🏢 组织:       %s\n", info.Org)
	}
	if info.Timezone != "" {
		fmt.Printf("🕐 时区:       %s\n", info.Timezone)
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
