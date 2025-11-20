package command

import (
	"flag"
	"net/http"
)

// Command 命令接口
type Command interface {
	// Name 返回命令名称
	Name() string

	// Description 返回命令描述
	Description() string

	// DefineFlags 定义命令的 flag 参数
	DefineFlags(flags *flag.FlagSet)

	// Execute 执行命令
	Execute(ctx *Context) error

	// Priority 返回命令优先级（数字越小优先级越高，用于决定执行顺序）
	Priority() int
}

// Context 命令执行上下文
type Context struct {
	// HTTPClient HTTP 客户端（已配置代理）
	HTTPClient *http.Client

	// Flags flag 集合
	Flags *flag.FlagSet

	// ProxyURL 代理 URL
	ProxyURL string

	// Timeout 超时时间（秒）
	Timeout int

	// ConfigFile 配置文件路径
	ConfigFile string
}
