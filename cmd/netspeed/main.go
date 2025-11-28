package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/icarus-go/netspeed/pkg/command"
	"github.com/icarus-go/netspeed/pkg/commands"
	"github.com/icarus-go/netspeed/pkg/proxy"
)

func main() {
	// 创建命令注册中心
	registry := command.NewRegistry()

	// 注册所有命令
	registerCommands(registry)

	// 定义全局 flags
	var (
		proxyURL   = flag.String("proxy", "", "设置代理 (支持 http://, socks5://, https://)")
		configFile = flag.String("config", "", "自定义测试站点配置文件（JSON 格式）")
		timeout    = flag.Int("timeout", 10, "请求超时时间（秒）")
	)

	// 让每个命令定义自己的 flags
	for _, cmd := range registry.All() {
		cmd.DefineFlags(flag.CommandLine)
	}

	// 解析 flags
	flag.Parse()

	// 初始化 HTTP 客户端
	httpClient, err := proxy.InitHTTPClient(*proxyURL, time.Duration(*timeout)*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ 代理配置错误: %v\n", err)
		os.Exit(1)
	}

	// 创建命令执行上下文
	ctx := &command.Context{
		HTTPClient: httpClient,
		Flags:      flag.CommandLine,
		ProxyURL:   *proxyURL,
		Timeout:    *timeout,
		ConfigFile: *configFile,
	}

	// 如果没有任何 flag 被设置，显示帮助
	if flag.NFlag() == 0 {
		helpCmd := commands.NewHelpCommand()
		helpFlags := flag.NewFlagSet("help", flag.ExitOnError)
		helpCmd.DefineFlags(helpFlags)
		helpFlags.Set("help", "true")
		helpCmd.Execute(ctx)
		return
	}

	// 执行所有激活的命令
	for _, cmd := range registry.All() {
		if err := cmd.Execute(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "❌ 执行命令失败: %v\n", err)
			os.Exit(1)
		}
	}
}

// registerCommands 注册所有命令
func registerCommands(registry *command.Registry) {
	// 按优先级注册命令
	cmds := []command.Command{
		commands.NewHelpCommand(),     // 优先级 1
		commands.NewIPCommand(),       // 优先级 10
		commands.NewIPScoreCommand(),  // 优先级 15
		commands.NewTestCommand(),     // 优先级 20
		commands.NewWatchCommand(),    // 优先级 30
	}

	for _, cmd := range cmds {
		if err := registry.Register(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "❌ 注册命令失败: %v\n", err)
			os.Exit(1)
		}
	}
}
