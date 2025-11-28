package commands

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/icarus-go/netspeed/pkg/command"
	"github.com/icarus-go/netspeed/pkg/ipinfo"
)

// IPCommand IP 检测命令
type IPCommand struct {
	enabled *bool
	origin  *bool
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
	c.origin = flags.Bool("origin", false, "获取原始 IP（不使用代理）")
}

// Execute 执行命令
func (c *IPCommand) Execute(ctx *command.Context) error {
	if !*c.enabled {
		return nil
	}

	fmt.Println("🌐 正在获取 IP 信息...")
	fmt.Println()

	// 如果指定了 -origin 参数，创建一个不使用代理的 HTTP 客户端
	var httpClient *http.Client
	if *c.origin {
		transport := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     30 * time.Second,
			// 通过设置一个始终返回 nil 的 Proxy 函数来完全禁用代理
			Proxy: func(req *http.Request) (*url.URL, error) {
				return nil, nil
			},
		}
		httpClient = &http.Client{
			Transport: transport,
			Timeout:   time.Duration(ctx.Timeout) * time.Second,
		}
		fmt.Println("⚠️  已设置为获取原始 IP（不使用代理）")
		fmt.Println()
	} else {
		// 使用默认的 HTTP 客户端（可能已配置代理）
		httpClient = ctx.HTTPClient
	}

	detector := ipinfo.NewDetector(httpClient)
	info, err := detector.Detect()
	if err != nil {
		return fmt.Errorf("获取 IP 信息失败: %v", err)
	}

	c.displayIPInfo(info, *c.origin)
	return nil
}

// Priority 返回命令优先级
func (c *IPCommand) Priority() int {
	return 10
}

// displayIPInfo 显示 IP 信息
func (c *IPCommand) displayIPInfo(info *ipinfo.IPInfo, isOrigin bool) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if isOrigin {
		fmt.Printf("📍 原始 IP:    %s\n", info.IP)
	} else {
		fmt.Printf("📍 IP 地址:    %s\n", info.IP)
	}
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
