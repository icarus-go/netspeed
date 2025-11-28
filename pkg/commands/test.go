package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/icarus-go/netspeed/pkg/command"
	"github.com/icarus-go/netspeed/pkg/config"
	"github.com/icarus-go/netspeed/pkg/output"
	"github.com/icarus-go/netspeed/pkg/tester"
)

// TestCommand 网站测试命令
type TestCommand struct {
	enabled *bool
}

// NewTestCommand 创建测试命令
func NewTestCommand() *TestCommand {
	return &TestCommand{}
}

// Name 返回命令名称
func (c *TestCommand) Name() string {
	return "test"
}

// Description 返回命令描述
func (c *TestCommand) Description() string {
	return "测试网站速度并以表格形式输出"
}

// DefineFlags 定义命令的 flag 参数
func (c *TestCommand) DefineFlags(flags *flag.FlagSet) {
	c.enabled = flags.Bool("test", false, "测试网站速度并以表格形式输出")
}

// Execute 执行命令
func (c *TestCommand) Execute(ctx *command.Context) error {
	if !*c.enabled {
		return nil
	}

	// 加载站点配置
	loader := config.NewLoader()
	sites, err := loader.LoadSites(ctx.ConfigFile)
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	fmt.Printf("🚀 开始测试 %d 个网站...\n", len(sites))
	fmt.Println()

	// 创建测试器并执行测试
	timeout := time.Duration(ctx.Timeout) * time.Second
	t := tester.NewTester(ctx.HTTPClient, timeout)
	results := t.TestAll(sites)

	// 显示结果
	output.PrintResultsTable(results)
	output.PrintSummary(results)

	return nil
}

// Priority 返回命令优先级
func (c *TestCommand) Priority() int {
	return 20
}
