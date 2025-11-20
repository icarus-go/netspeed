package command

import (
	"fmt"
	"sort"
)

// Registry 命令注册中心
type Registry struct {
	commands map[string]Command
}

// NewRegistry 创建新的命令注册中心
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

// Register 注册命令
func (r *Registry) Register(cmd Command) error {
	name := cmd.Name()
	if _, exists := r.commands[name]; exists {
		return fmt.Errorf("命令 '%s' 已注册", name)
	}
	r.commands[name] = cmd
	return nil
}

// Get 获取命令
func (r *Registry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

// All 获取所有命令（按优先级排序）
func (r *Registry) All() []Command {
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}

	// 按优先级排序
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].Priority() < cmds[j].Priority()
	})

	return cmds
}

// Has 检查命令是否存在
func (r *Registry) Has(name string) bool {
	_, ok := r.commands[name]
	return ok
}

// Count 返回已注册命令数量
func (r *Registry) Count() int {
	return len(r.commands)
}
