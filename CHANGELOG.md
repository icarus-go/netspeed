# NetSpeed 更新日志

## v1.0.4 - 增强测试覆盖率和代理配置 (2025-12-01)

### ✨ 新增功能

**大幅提升测试覆盖率**
- 新增 `pkg/commands/ip_test.go` - 282 行单元测试，覆盖 IP 检测命令
- 新增 `pkg/ipinfo/detector_http_test.go` - 376 行 HTTP 检测测试
- 新增 `pkg/ipinfo/score_test.go` - 481 行纯净度评分测试
- 新增 `pkg/ipinfo/detector_test.go` - 187 行基础检测测试
- 总计新增 1300+ 行测试代码，显著提升代码质量和可维护性

**优化代理配置**
- 改进 `HTTP.Transport.Proxy` 默认值，默认读取系统环境变量
- 简化代理配置逻辑，自动感知系统代理设置
- 确保在没有显式代理配置时也能正确使用环境变量

### 📊 测试统计

- ✅ 新增 4 个测试文件，总计 1300+ 行测试代码
- ✅ 覆盖 IP 检测、纯净度评分、HTTP 请求等核心功能
- ✅ 所有单元测试通过
- ✅ 继续支持跨平台构建（Linux, macOS, Windows）

### 📝 验证方式

测试通过以下命令运行：

```bash
# 运行所有测试
go test ./...

# 查看测试覆盖率
go test -cover ./...

# 测试特定模块
go test ./pkg/ipinfo
go test ./pkg/commands
```

### 🚀 安装方式

```bash
# 重新编译以获得最佳性能
go build -o netspeed ./cmd/netspeed

# 或使用 go install
go install github.com/icarus-go/netspeed/cmd/netspeed@latest
```

### 📦 变更文件

- `pkg/commands/ip_test.go` - 新增
- `pkg/ipinfo/detector_http_test.go` - 新增
- `pkg/ipinfo/score_test.go` - 新增
- `pkg/ipinfo/detector_test.go` - 新增
- `pkg/proxy/proxy.go` - 优化代理配置
- `Makefile` - 更新构建脚本

---

## v1.0.3 - 修复导入路径问题 (2025-11-28)

### 🐛 Bug 修复

**修复模块导入路径不一致问题**
- 统一所有导入路径从 `github.com/icarus-go/net-speed` 修改为 `github.com/icarus-go/netspeed`
- 解决了 `go install` 和 `go build` 失败的问题
- 确保项目能够通过 `go install github.com/icarus-go/netspeed/cmd/netspeed@latest` 正确安装

### 📦 影响的文件

修复了以下 8 个文件的导入路径：
- `cmd/netspeed/main.go`
- `pkg/commands/ip.go`
- `pkg/commands/purity.go`
- `pkg/commands/test.go`
- `pkg/commands/watch.go`
- `pkg/commands/help.go`
- `pkg/config/loader.go`
- `pkg/output/table.go`

### ✅ 测试验证

- ✅ 所有单元测试通过
- ✅ 跨平台编译成功（Linux, macOS, Windows）
- ✅ 模块路径与 go.mod 保持一致

### 📝 安装方式

现在可以通过以下方式安装：

```bash
# 使用 go install（推荐）
go install github.com/icarus-go/netspeed/cmd/netspeed@latest

# 或从源码编译
git clone https://github.com/icarus-go/netspeed.git
cd netspeed
go build -o netspeed ./cmd/netspeed
```

---

## v2.0.0 - 架构重构与功能增强 (2025-11-19)

### 🎯 重大更新

#### 1. 架构重构
- ✅ 采用**命令注册模式**，实现插件化架构
- ✅ **模块化分包**设计，15个独立模块
- ✅ 单文件 460 行 → 模块化 15 个文件
- ✅ 主程序入口仅 90 行，极简设计

#### 2. 新增功能

**IP 纯净度检测 (`-purity`)**
- 智能检测 IP 类型（VPN/代理/Tor/数据中心）
- 0-100 分纯净度评分
- 风险等级评估（低/中/高）
- 基于 ISP/组织名称的智能分析
- 个性化建议

**ping0.cc/geo 支持**
- 新增 ping0.cc 作为首选 IP 检测 API
- 响应速度更快，国内访问友好
- 完整解析文本格式：
  - IP 地址
  - 地区信息（国家/省份/城市）
  - ASN 编号
  - 网络运营商完整名称
- 智能解析中英文格式

#### 3. 命令系统

新的命令注册系统支持：

| 命令 | 优先级 | 功能 |
|------|--------|------|
| `-help` | 1 | 显示帮助信息 |
| `-ip` | 10 | 获取 IP 地理信息 |
| `-purity` | 15 | **新增**: IP 纯净度检测 |
| `-test` | 20 | 网站速度测试 |
| `-watch` | 30 | 持续监控模式 |

### 📦 项目结构

```
net-speed/
├── cmd/netspeed/main.go           # 主入口 (90行)
├── pkg/
│   ├── command/                   # 命令注册系统
│   ├── commands/                  # 5个命令实现
│   ├── tester/                    # 网站测试模块
│   ├── ipinfo/                    # IP检测模块 + 纯净度
│   ├── proxy/                     # 代理配置
│   ├── output/                    # 输出格式化
│   └── config/                    # 配置管理
├── ARCHITECTURE.md                # 架构文档
├── README.md                      # 用户文档
└── CHANGELOG.md                   # 本文件
```

### 🔧 IP 检测提供商

优先级排序（带故障转移）：

1. **ping0.cc** (新增) - 文本格式，响应快
2. ipapi.co - JSON 格式
3. ipinfo.io - JSON 格式
4. ip-api.com - JSON 格式

### 📊 ping0.cc/geo 响应格式

```
14.153.68.158                          # IP 地址
中国 广东省深圳市福田中国电信               # 地区 + ISP
AS4134                                  # ASN 编号
CHINANET Guangdong province network     # 网络运营商完整名称
```

### 🧪 测试

- ✅ 新增单元测试 `pkg/ipinfo/detector_test.go`
- ✅ 测试中英文格式解析
- ✅ 测试 ASN 提取
- ✅ 所有测试通过

### 🎨 使用示例

**IP 检测（使用 ping0.cc）**
```bash
$ netspeed -ip

📍 IP 地址:    14.153.68.158
🌍 国家:       中国
📌 地区:       广东省
🏙️  城市:       深圳市
🔌 ISP:        福田中国电信
🏢 组织:       AS4134 CHINANET Guangdong province network
```

**IP 纯净度检测**
```bash
$ netspeed -purity

📊 IP 纯净度报告
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📍 IP 地址:       14.153.68.158
🏭 组织:          福田中国电信

✨ 纯净度评分:    100.0/100  (优秀)
⚠️  风险等级:     低风险

🔎 检测结果:
  ✅ VPN: 未检测到
  ✅ 代理: 未检测到
  ✅ Tor: 未检测到
  ✅ 数据中心: 未检测到

💡 建议:
  ✓ IP 纯净度很高，适合大多数场景使用
```

### 🚀 性能提升

- ✅ ping0.cc 响应速度更快
- ✅ 模块化设计提升可维护性
- ✅ 命令注册系统提升扩展性

### 📝 文档

- ✅ 新增 `ARCHITECTURE.md` - 完整架构文档
- ✅ 更新 `README.md` - 用户文档
- ✅ 新增 `CHANGELOG.md` - 更新日志

### 🔄 升级指南

从 v1.x 升级到 v2.0：

1. 重新编译：`go build -o netspeed ./cmd/netspeed`
2. 所有命令保持兼容，无需修改脚本
3. 新增 `-purity` 命令可选使用

### 🛠️ 开发者

**添加新命令（3步）：**

1. 在 `pkg/commands/` 创建命令文件
2. 实现 `Command` 接口
3. 在 `main.go` 注册命令

**示例：**
```go
// 1. 创建 pkg/commands/dns.go
type DNSCommand struct {
    enabled *bool
}

func (c *DNSCommand) Name() string { return "dns" }
func (c *DNSCommand) Execute(ctx *command.Context) error {
    // 实现逻辑
}

// 2. 在 main.go 注册
commands.NewDNSCommand(),  // 一行代码完成
```

### 📚 相关链接

- [架构文档](ARCHITECTURE.md)
- [快速开始](QUICKSTART.md)
- [用户文档](README.md)

---

## v1.0.0 - 初始版本

- 基础网站测试功能
- IP 检测功能
- 代理支持
- 持续监控模式
- 自定义配置文件
