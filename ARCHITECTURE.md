# NetSpeed 架构文档

## 架构概览

NetSpeed 采用**命令注册模式**和**模块化分包**设计，实现了高度解耦、易扩展的架构。

## 项目结构

```
net-speed/
├── cmd/
│   └── netspeed/
│       └── main.go              # 主入口 (~90行)
├── pkg/
│   ├── command/                 # 命令注册系统
│   │   ├── command.go           # 命令接口定义
│   │   └── registry.go          # 命令注册中心
│   ├── commands/                # 具体命令实现
│   │   ├── help.go              # 帮助命令
│   │   ├── ip.go                # IP检测命令
│   │   ├── purity.go            # IP纯净度检测命令
│   │   ├── test.go              # 网站测试命令
│   │   └── watch.go             # 持续监控命令
│   ├── tester/                  # 网站测试模块
│   │   ├── model.go             # 数据模型
│   │   └── tester.go            # 测试逻辑
│   ├── ipinfo/                  # IP检测模块
│   │   ├── model.go             # 数据模型
│   │   ├── detector.go          # IP检测器
│   │   └── score.go             # 纯净度评分
│   ├── proxy/                   # 代理配置模块
│   │   └── proxy.go             # 代理初始化
│   ├── output/                  # 输出格式化模块
│   │   └── table.go             # 表格输出
│   └── config/                  # 配置管理模块
│       └── loader.go            # 配置加载
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── sites.example.json
```

## 核心设计模式

### 1. 命令注册模式

#### 命令接口
```go
type Command interface {
    Name() string                    // 命令名称
    Description() string             // 命令描述
    DefineFlags(*flag.FlagSet)      // 定义 flag 参数
    Execute(*Context) error         // 执行命令
    Priority() int                  // 优先级
}
```

#### 命令注册流程
```
1. 创建命令注册中心
   ↓
2. 注册所有命令
   ↓
3. 定义全局和命令特定的 flags
   ↓
4. 解析命令行参数
   ↓
5. 按优先级执行激活的命令
```

#### 命令优先级

| 优先级 | 命令 | 说明 |
|--------|------|------|
| 1 | help | 最高优先级，优先显示帮助 |
| 10 | ip | IP 检测 |
| 15 | purity | IP 纯净度检测 |
| 20 | test | 网站测试 |
| 30 | watch | 持续监控 |

### 2. 模块化分包

#### pkg/command - 命令注册系统
- **command.go**: 定义命令接口和执行上下文
- **registry.go**: 实现命令注册中心，支持动态注册和查找

#### pkg/commands - 具体命令实现
每个命令都是独立的，实现 `Command` 接口：
- 自包含所有逻辑
- 可独立测试
- 易于新增命令

#### pkg/tester - 网站测试模块
- **model.go**: 定义 `Site` 和 `TestResult` 数据结构
- **tester.go**: 实现并发测试逻辑，支持降级策略

#### pkg/ipinfo - IP 检测模块
- **model.go**: 定义 `IPInfo` 和 `IPScore` 数据结构
- **detector.go**: 实现 IP 检测和纯净度评分
- **score.go**: IP 纯净度评分算法

#### pkg/proxy - 代理配置模块
- 支持 HTTP/HTTPS/SOCKS5 代理
- 自动设置环境变量
- 统一的 HTTP 客户端初始化

#### pkg/output - 输出格式化模块
- 表格输出
- 统计摘要
- 未来可扩展: JSON、XML、CSV 输出

#### pkg/config - 配置管理模块
- JSON 配置文件加载
- 默认配置支持
- 未来可扩展: YAML、TOML 支持

## 数据流图

### 网站测试流程
```
用户输入: netspeed -test
    ↓
main.go 解析参数
    ↓
TestCommand.Execute()
    ↓
config.Loader.LoadSites()  → 加载测试站点
    ↓
tester.Tester.TestAll()    → 并发测试所有站点
    ↓         ↓
    ↓    goroutine × N (并发测试)
    ↓         ↓
    ↓    TestResult[] (收集结果)
    ↓
output.PrintResultsTable() → 表格输出
    ↓
output.PrintSummary()      → 统计摘要
```

### IP 检测流程
```
用户输入: netspeed -ip
    ↓
main.go 解析参数
    ↓
IPCommand.Execute()
    ↓
ipinfo.Detector.Detect()
    ↓
尝试 Provider 1  ✓ 成功 → 返回 IPInfo
    ↓ 失败
尝试 Provider 2  ✓ 成功 → 返回 IPInfo
    ↓ 失败
尝试 Provider 3  ✓ 成功 → 返回 IPInfo
    ↓
显示 IP 信息
```

### IP 纯净度检测流程
```
用户输入: netspeed -purity
    ↓
main.go 解析参数
    ↓
IPScoreCommand.Execute()
    ↓
ipinfo.Detector.DetectScore()
    ↓
获取基础 IP 信息 (复用 Detect())
    ↓
calculatePuritySimple()  → 基于 ISP/Org 计算分数
    ↓
analyzeIPCharacteristics() → 检测 VPN/代理特征
    ↓
返回 IPScore (分数 + 特征)
    ↓
显示纯净度报告
```

## 扩展性设计

### 添加新命令

1. 在 `pkg/commands/` 创建新命令文件
2. 实现 `Command` 接口
3. 在 `main.go` 的 `registerCommands()` 中注册

```go
// 示例: 添加 DNS 测试命令
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

// 在 main.go 注册
commands.NewDNSCommand(),  // 优先级 25
```

### 添加新输出格式

在 `pkg/output/` 添加新的格式化函数：

```go
// json.go
func PrintResultsJSON(results []tester.TestResult) {
    // JSON 输出实现
}

// csv.go
func PrintResultsCSV(results []tester.TestResult) {
    // CSV 输出实现
}
```

### 添加新测试站点源

在 `pkg/config/` 扩展配置加载：

```go
func (l *Loader) LoadSitesFromURL(url string) ([]tester.Site, error) {
    // 从 URL 加载站点列表
}

func (l *Loader) LoadSitesFromYAML(file string) ([]tester.Site, error) {
    // 从 YAML 加载站点列表
}
```

## 依赖关系

```
main.go
  ├── command (命令系统)
  ├── commands (所有命令)
  │     ├── tester (网站测试)
  │     ├── ipinfo (IP检测)
  │     ├── output (输出格式)
  │     └── config (配置加载)
  └── proxy (代理配置)
```

**特点:**
- 单向依赖，避免循环依赖
- 业务模块 (`tester`, `ipinfo`) 不依赖命令系统
- 命令层 (`commands`) 组合业务模块
- 主程序只依赖命令注册系统

## 并发安全

### 网站测试并发
```go
// tester/tester.go
func (t *Tester) TestAll(sites []Site) []TestResult {
    var wg sync.WaitGroup
    results := make([]TestResult, len(sites))

    for i, site := range sites {
        wg.Add(1)
        go func(idx int, s Site) {
            defer wg.Done()
            results[idx] = t.TestSite(s)  // 每个 goroutine 写入不同的索引
        }(i, site)
    }

    wg.Wait()
    return results
}
```

**安全保证:**
- 预分配结果切片
- 每个 goroutine 写入不同索引
- WaitGroup 同步等待

### HTTP 客户端共享
```go
// proxy/proxy.go
func InitHTTPClient(proxyURL string, timeout time.Duration) (*http.Client, error) {
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     30 * time.Second,
    }
    // http.Client 本身是并发安全的
    return &http.Client{Transport: transport, Timeout: timeout}, nil
}
```

## 配置管理

### 全局配置
- `proxy`: 代理 URL
- `timeout`: 请求超时时间
- `config`: 配置文件路径

### 命令特定配置
- `watch`: 刷新间隔（秒）

### 配置优先级
```
命令行参数 > 配置文件 > 默认值
```

## 错误处理策略

### 1. 命令级错误
```go
func (c *TestCommand) Execute(ctx *command.Context) error {
    if err := something(); err != nil {
        return fmt.Errorf("操作失败: %v", err)
    }
    return nil
}
```

### 2. 降级策略
```go
// tester/tester.go
func (t *Tester) TestSite(site Site) TestResult {
    // 尝试 HEAD 请求
    resp, err := t.client.Do(req)
    if err != nil {
        // 降级: 尝试 GET 请求
        return t.fallbackTest(site, ctx)
    }
    // ...
}
```

### 3. 故障转移
```go
// ipinfo/detector.go
func (d *Detector) Detect() (*IPInfo, error) {
    for _, provider := range d.providers {
        if ipInfo, err := d.fetchFromProvider(provider); err == nil {
            return ipInfo, nil  // 成功则返回
        }
        // 失败继续尝试下一个
    }
    return nil, errors.New("all providers failed")
}
```

## 性能优化

### 1. 并发测试
- 使用 goroutine 并行测试所有站点
- 时间复杂度: O(1) (最慢的站点决定总时间)

### 2. HTTP 连接池
```go
transport := &http.Transport{
    MaxIdleConns:        100,   // 最大空闲连接数
    MaxIdleConnsPerHost: 10,    // 每个 host 最大空闲连接
    IdleConnTimeout:     30 * time.Second,
}
```

### 3. 超时控制
- 每个请求都有精确的超时控制
- 避免长时间等待

## 测试策略

### 单元测试
```bash
# 测试单个模块
go test ./pkg/tester
go test ./pkg/ipinfo
go test ./pkg/command
```

### 集成测试
```bash
# 测试完整流程
go test ./cmd/netspeed
```

### 跨平台测试
```bash
# 使用 Makefile
make build-all    # 编译所有平台
```

## 未来扩展方向

### 1. 新命令
- [x] `-test`: 网站测试
- [x] `-ip`: IP 检测
- [x] `-purity`: IP 纯净度
- [x] `-watch`: 持续监控
- [ ] `-dns`: DNS 解析测试
- [ ] `-traceroute`: 路由追踪
- [ ] `-benchmark`: 性能基准测试

### 2. 新功能
- [ ] JSON/CSV 输出格式
- [ ] 历史记录保存
- [ ] Web UI 界面
- [ ] 测试报告导出
- [ ] 邮件/webhook 告警

### 3. 性能优化
- [ ] 结果缓存
- [ ] 智能超时调整
- [ ] 自适应并发数

## 贡献指南

### 代码规范
- 使用 `gofmt` 格式化代码
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 每个包都有清晰的职责
- 避免循环依赖

### 提交新命令
1. Fork 项目
2. 创建feature分支
3. 在 `pkg/commands/` 实现新命令
4. 在 `main.go` 注册命令
5. 更新 README 和架构文档
6. 提交 Pull Request

---

**设计原则:**
- 单一职责
- 开闭原则 (对扩展开放，对修改封闭)
- 依赖倒置 (依赖抽象而非具体实现)
- 接口隔离
