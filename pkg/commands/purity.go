package commands

import (
	"flag"
	"fmt"

	"github.com/icarus-go/net-speed/pkg/command"
	"github.com/icarus-go/net-speed/pkg/ipinfo"
)

// IPScoreCommand IP 纯净度检测命令
type IPScoreCommand struct {
	enabled *bool
}

// NewIPScoreCommand 创建 IP 纯净度检测命令
func NewIPScoreCommand() *IPScoreCommand {
	return &IPScoreCommand{}
}

// Name 返回命令名称
func (c *IPScoreCommand) Name() string {
	return "purity"
}

// Description 返回命令描述
func (c *IPScoreCommand) Description() string {
	return "检测 IP 纯净度和风险评分"
}

// DefineFlags 定义命令的 flag 参数
func (c *IPScoreCommand) DefineFlags(flags *flag.FlagSet) {
	c.enabled = flags.Bool("purity", false, "检测 IP 纯净度和风险评分")
}

// Execute 执行命令
func (c *IPScoreCommand) Execute(ctx *command.Context) error {
	if !*c.enabled {
		return nil
	}

	fmt.Println("🔍 正在检测 IP 纯净度...")
	fmt.Println()

	detector := ipinfo.NewDetector(ctx.HTTPClient)
	score, err := detector.DetectScore()
	if err != nil {
		return fmt.Errorf("检测 IP 纯净度失败: %v", err)
	}

	c.displayIPScore(score)
	return nil
}

// Priority 返回命令优先级
func (c *IPScoreCommand) Priority() int {
	return 15
}

// displayIPScore 显示 IP 纯净度信息
func (c *IPScoreCommand) displayIPScore(score *ipinfo.IPScore) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 IP 纯净度报告")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📍 IP 地址:       %s\n", score.IP)
	fmt.Printf("🏢 ASN:           %s\n", score.ASN)
	fmt.Printf("🏭 组织:          %s\n", score.ASNOrg)
	fmt.Println()

	// 纯净度评分
	fmt.Printf("✨ 纯净度评分:    %.1f/100  (%s)\n", score.Score, score.GetQualityDescription())
	fmt.Printf("⚠️  风险等级:     %s\n", score.RiskLevel)
	fmt.Println()

	// 检测结果
	fmt.Println("🔎 检测结果:")
	c.printCheckResult("VPN", score.IsVPN)
	c.printCheckResult("代理", score.IsProxy)
	c.printCheckResult("Tor", score.IsTor)
	c.printCheckResult("数据中心", score.IsDatacenter)
	c.printCheckResult("黑名单", score.IsBlacklisted)

	fmt.Println()

	// 建议
	c.printRecommendation(score)

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// printCheckResult 打印检测项结果
func (c *IPScoreCommand) printCheckResult(name string, detected bool) {
	if detected {
		fmt.Printf("  ❌ %s: 检测到\n", name)
	} else {
		fmt.Printf("  ✅ %s: 未检测到\n", name)
	}
}

// printRecommendation 打印建议
func (c *IPScoreCommand) printRecommendation(score *ipinfo.IPScore) {
	fmt.Println("💡 建议:")
	switch {
	case score.Score >= 80:
		fmt.Println("  ✓ IP 纯净度很高，适合大多数场景使用")
	case score.Score >= 60:
		fmt.Println("  ⚠ IP 纯净度中等，部分网站可能会有限制")
	case score.Score >= 40:
		fmt.Println("  ⚠ IP 纯净度较低，建议更换 IP 或节点")
	default:
		fmt.Println("  ❌ IP 纯净度很低，强烈建议更换")
	}
}
