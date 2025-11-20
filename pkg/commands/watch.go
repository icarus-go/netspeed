package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/icarus-go/net-speed/pkg/command"
	"github.com/icarus-go/net-speed/pkg/config"
	"github.com/icarus-go/net-speed/pkg/output"
	"github.com/icarus-go/net-speed/pkg/tester"
)

// WatchCommand 持续监控命令
type WatchCommand struct {
	interval *int
}

// NewWatchCommand 创建监控命令
func NewWatchCommand() *WatchCommand {
	return &WatchCommand{}
}

// Name 返回命令名称
func (c *WatchCommand) Name() string {
	return "watch"
}

// Description 返回命令描述
func (c *WatchCommand) Description() string {
	return "持续监控模式，指定刷新间隔（秒）"
}

// DefineFlags 定义命令的 flag 参数
func (c *WatchCommand) DefineFlags(flags *flag.FlagSet) {
	c.interval = flags.Int("watch", 0, "持续监控模式，指定刷新间隔（秒）")
}

// Execute 执行命令
func (c *WatchCommand) Execute(ctx *command.Context) error {
	if *c.interval <= 0 {
		return nil
	}

	fmt.Printf("👀 持续监控模式 (每 %d 秒刷新)\n", *c.interval)
	fmt.Println("按 Ctrl+C 退出")
	fmt.Println()

	// 加载站点配置
	loader := config.NewLoader()
	sites, err := loader.LoadSites(ctx.ConfigFile)
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 创建测试器
	timeout := time.Duration(ctx.Timeout) * time.Second
	t := tester.NewTester(ctx.HTTPClient, timeout)

	// 首次立即执行
	c.runTest(t, sites)

	// 定期执行
	ticker := time.NewTicker(time.Duration(*c.interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 清屏（跨平台方式）
		fmt.Print("\033[2J\033[H")

		fmt.Printf("⏰ 最后更新: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println()
		c.runTest(t, sites)
	}

	return nil
}

// Priority 返回命令优先级
func (c *WatchCommand) Priority() int {
	return 30
}

// runTest 执行一次测试
func (c *WatchCommand) runTest(t *tester.Tester, sites []tester.Site) {
	fmt.Printf("🚀 开始测试 %d 个网站...\n", len(sites))
	fmt.Println()

	results := t.TestAll(sites)

	output.PrintResultsTable(results)
	output.PrintSummary(results)
}
