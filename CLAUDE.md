# NetSpeed - AI 开发指南

本文档为 Claude Code 在 NetSpeed 项目上的工作提供指导。NetSpeed 是一个基于 Go 语言开发的跨平台网络质量检测工具。

**强调**
- **所有沟通必须使用`中文`**
- **AI与用户沟通时，需称呼用户为 `Developer`**
- **必须有始有终的开发流程, 不允许开发后没有通过单元测试就提交代码** 
- 
## 项目概述

NetSpeed 是一个用于测试网络质量的 CLI 工具，具有以下特性：
- **网站延迟测试** - 并行测试多个网站的响应速度
- **IP 地理信息检测** - 自动检测公网 IP 和地理位置（支持代理）
- **IP 纯净度检测** - 智能检测 IP 类型和风险评分
- **多 API 支持** - ping0.cc、ipapi.co、ipinfo.io 等多个 API
- **代理支持** - 支持 HTTP、HTTPS、SOCKS5 代理
- **持续监控模式** - 定时刷新网络质量状态
- **自定义配置** - 支持 JSON 格式的自定义测试站点
- **跨平台** - 单个二进制文件，无需依赖
- **模块化架构** - 命令注册模式，易于扩展

## 快速导航

- [项目结构](#项目结构)
- [关键文件](#关键文件)
- [核心命令](#核心命令)
- [开发指南](#开发指南)
- [测试](#测试)
- [快速参考](#快速参考)
- [常见问题](#常见问题)

## 项目结构

```
netspeed/
├── cmd/netspeed/
│   └── main.go                      # 应用程序入口点
├── pkg/
│   ├── command/                     # 命令注册系统
│   │   ├── command.go              # 命令接口定义
│   │   └── registry.go             # 命令注册中心
│   ├── commands/                    # 具体命令实现
│   │   ├── help.go                 # 帮助命令
│   │   ├── ip.go                   # IP 检测命令
│   │   ├── purity.go               # IP 纯净度检测命令
│   │   ├── test.go                 # 网站测试命令
│   │   └── watch.go                # 持续监控命令
│   ├── tester/                      # 网站测试模块
│   │   ├── model.go                # 数据模型（Site, TestResult）
│   │   └── tester.go               # 并发测试逻辑
│   ├── ipinfo/                      # IP 检测模块
│   │   ├── model.go                # 数据模型（IPInfo, IPScore）
│   │   ├── detector.go             # IP 检测与纯净度评分
│   │   └── score.go                # 纯净度评分算法
│   ├── proxy/                       # 代理配置模块
│   │   └── proxy.go                # HTTP/HTTPS/SOCKS5 代理支持
│   ├── output/                      # 输出格式化模块
│   │   └── table.go                # 表格输出格式化
│   └── config/                      # 配置管理
│       └── loader.go                # JSON 配置加载
├── sites.example.json               # 示例配置文件
├── Makefile                         # 构建脚本
├── go.mod                           # Go 模块定义
├── README.md                        # 用户文档
└── ARCHITECTURE.md                  # 技术架构文档（详细）
```

## 关键文件

| 文件路径                      | 作用              | 重要性   |
|---------------------------|-----------------|-------|
| `cmd/netspeed/main.go`    | 应用入口点，参数解析和命令注册 | ⭐⭐⭐⭐⭐ |
| `pkg/command/command.go`  | 命令接口定义          | ⭐⭐⭐⭐⭐ |
| `pkg/command/registry.go` | 命令注册中心          | ⭐⭐⭐⭐  |
| `pkg/commands/test.go`    | 网站测试命令实现        | ⭐⭐⭐⭐⭐ |
| `pkg/commands/ip.go`      | IP 检测命令实现       | ⭐⭐⭐⭐  |
| `pkg/tester/tester.go`    | 并发测试逻辑核心        | ⭐⭐⭐⭐⭐ |
| `pkg/ipinfo/detector.go`  | IP 检测与纯净度评分     | ⭐⭐⭐⭐⭐ |
| `pkg/proxy/proxy.go`      | 代理配置            | ⭐⭐⭐   |
| `pkg/output/table.go`     | 输出格式化           | ⭐⭐⭐   |
| `sites.example.json`      | 配置模板            | ⭐⭐⭐   |
| `ARCHITECTURE.md`         | 详细架构说明          | ⭐⭐⭐   |

## 核心命令

### 基础命令

```bash
# 查看帮助
netspeed -help

# 测试网站速度
netspeed -test

# 获取 IP 信息
netspeed -ip

# 获取原始 IP（绕过代理，检测真实网络出口）
netspeed -ip -origin

# 检测 IP 纯净度
netspeed -purity
```

### 高级用法

```bash
# 使用代理（SOCKS5）
netspeed -test -proxy socks5://127.0.0.1:1080

# 使用代理（HTTP）
netspeed -test -proxy http://proxy.example.com:8080

# 持续监控模式（每 30 秒刷新）
netspeed -test -watch 30

# 使用自定义配置文件
netspeed -test -config sites.json

# 设置超时时间（默认 10 秒）
netspeed -test -timeout 5
```

### 代理配置

**方式 1：参数设置（高优先级）**
```bash
netspeed -test -proxy http://proxy.example.com:8080
```

**方式 2：环境变量（推荐）**
```bash
export http_proxy=http://proxy.example.com:8080
export https_proxy=http://proxy.example.com:8080
netspeed -test
netspeed -ip
```

支持的代理协议：HTTP、HTTPS、SOCKS5

## 开发指南

### 核心设计原则

- **单一职责** - 每个模块都有清晰的职责
- **开闭原则** - 对扩展开放，对修改封闭
- **依赖倒置** - 依赖抽象而非具体实现
- **接口隔离** - 小而专注的接口
- **单向依赖** - 避免循环依赖

### 添加新命令

1. 在 `pkg/commands/` 创建新命令文件
2. 实现 `Command` 接口
3. 在 `main.go` 的 `registerCommands()` 中注册

```go
// 示例：添加 DNS 测试命令
type DNSCommand struct {
    enabled *bool
}

func (c *DNSCommand) Name() string { return "dns" }
func (c *DNSCommand) Description() string { return "DNS 解析测试" }
func (c *DNSCommand) Priority() int { return 25 }

func (c *DNSCommand) DefineFlags(flags *flag.FlagSet) {
    c.enabled = flags.Bool("dns", false, "测试 DNS 解析速度")
}

func (c *DNSCommand) Execute(ctx *command.Context) error {
    if !*c.enabled {
        return nil
    }
    // 实现 DNS 测试逻辑
    return nil
}

// 在 main.go 中注册
commands.NewDNSCommand(),  // 优先级 25
```

### 错误处理策略

1. **命令级错误**
```go
func (c *TestCommand) Execute(ctx *command.Context) error {
    if err := something(); err != nil {
        return fmt.Errorf("操作失败: %v", err)
    }
    return nil
}
```

2. **降级策略** - HEAD 失败自动降级到 GET
3. **故障转移** - IP 检测多 API 故障转移

### 配置优先级

```
命令行参数 > 配置文件 > 默认值
```

## 测试

### 运行测试

```bash
# 测试单个模块
go test ./pkg/tester
go test ./pkg/ipinfo
go test ./pkg/command

# 带覆盖率测试
go test -cover ./...

# 测试完整流程
go test ./cmd/netspeed

# 构建所有平台
make build-all
```

### 测试范围

- **命令层** - 测试命令注册和执行
- **业务逻辑** - 测试 tester、ipinfo 模块
- **集成测试** - 测试端到端工作流
- **跨平台** - 在 Linux、macOS、Windows 上测试

## 快速参考

### 构建项目

```bash
# 构建
go build -o netspeed ./cmd/netspeed

# 跨平台构建
make build-all
# 或手动：
GOOS=linux GOARCH=amd64 go build -o netspeed-linux-amd64 ./cmd/netspeed
GOOS=darwin GOARCH=arm64 go build -o netspeed-darwin-arm64 ./cmd/netspeed
GOOS=windows GOARCH=amd64 go build -o netspeed-windows-amd64.exe ./cmd/netspeed
```

### 配置测试站点

编辑 `sites.json` 文件：

```json
{
  "sites": [
    {"name": "Google", "url": "https://www.google.com"},
    {"name": "Baidu", "url": "https://www.baidu.com"}
  ]
}
```

使用自定义配置：
```bash
netspeed -test -config sites.json
```

### 性能优化要点

1. **并发测试** - 使用 goroutine 并行测试，时间复杂度 O(1)
2. **HTTP 连接池** - 复用连接，减少开销
3. **超时控制** - 对每个请求进行精确超时控制
4. **预分配** - 预分配切片/映射，减少 GC 压力

## 常见问题

### Q: 如何修改测试超时时间？
A: 使用 `-timeout` 参数，默认 10 秒：`netspeed -test -timeout 5`

### Q: 代理不生效怎么办？
A: 检查优先级：1) 命令行参数 2) 环境变量。确保代理地址正确且可达。

### Q: 如何自定义测试站点？
A: 创建 `sites.json` 文件，使用 `netspeed -test -config sites.json` 指定。

### Q: IP 检测失败的原因？
A: 检查网络连接、代理设置，或使用 `-origin` 参数绕过代理获取原始 IP。

### Q: 如何添加新的 IP 检测 API？
A: 在 `pkg/ipinfo/detector.go` 中添加新的 provider，按优先级顺序排列。

### Q: 跨平台构建失败？
A: 确保 Go 版本 >= 1.19，使用 `make build-all` 或手动设置 GOOS/GOARCH。

### Q: 如何禁用某个命令？
A: 在 `main.go` 的 `registerCommands()` 中注释或删除对应命令注册。

### Q: 测试结果不准确？
A: 检查网络环境、代理设置，确保测试站点可达。可尝试多个站点对比。

## 安全注意事项

### 网络安全
- 避免在不受信任的网络中使用真实代理凭据
- 使用 `-origin` 参数时谨慎，避免泄露真实 IP
- 定期更新依赖，避免使用过时 API

### 代码安全
- 所有外部 API 调用都需要错误处理和超时控制
- 配置文件中的站点 URL 需要验证格式
- 代理 URL 必须严格验证协议和格式

### 依赖管理
- 仅使用标准库，无第三方依赖
- 外部 API 仅在运行时调用，非构建时依赖
- 定期检查 API 可用性（ping0.cc, ipapi.co 等）

### 数据隐私
- IP 信息仅用于本地显示，不会上传
- 测试站点列表仅在本地加载
- 建议不要在敏感环境中保存详细测试日志

---

## 贡献指南

### 工作流程

1. Fork 仓库
2. 创建功能分支
3. 在 `pkg/commands/` 实现新命令
4. 在 `main.go` 注册命令
5. 更新 README.md 和 ARCHITECTURE.md
6. 提交 Pull Request

### 测试检查清单

- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 跨平台构建成功
- [ ] 未检测到竞争条件
- [ ] 文档已更新

---

**核心设计理念：**
- Unix 哲学：一件事做好
- 模块化优于单体
- 性能与简洁
- 跨平台兼容性
- 易于扩展，难以破坏

**在这个项目上工作时：**
1. 理解命令注册模式
2. 尊重模块边界
3. 先写测试再实现
4. 遵循现有的错误处理策略
5. 保持新功能简单专注
6. 随变更更新文档