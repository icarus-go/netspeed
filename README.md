# NetSpeed - 跨平台网络质量检测工具

🚀 基于 Go 语言开发的命令行网络质量检测工具，支持 Linux、macOS 和 Windows。

## 功能特性

- ✅ **网站延迟测试** - 并行测试多个网站的响应速度
- ✅ **IP 地理信息检测** - 自动检测公网 IP 和地理位置（支持代理）
- ✅ **IP 纯净度检测** - ⭐ **新功能**: 智能检测 IP 类型和风险评分
- ✅ **多 API 支持** - ping0.cc (快速)、ipapi.co、ipinfo.io 等多个 API
- ✅ **代理支持** - 支持 HTTP、HTTPS、SOCKS5 代理
- ✅ **持续监控模式** - 定时刷新网络质量状态
- ✅ **自定义配置** - 支持 JSON 格式的自定义测试站点
- ✅ **表格化输出** - 清晰的表格展示测试结果
- ✅ **跨平台** - 单个二进制文件，无需依赖
- ✅ **模块化架构** - 命令注册模式，易于扩展

## 安装

### 从源码编译

```bash
# 克隆项目
git clone https://github.com/icarus-go/net-speed.git
cd net-speed

# 编译
go build -o netspeed

# Linux/macOS
chmod +x netspeed
sudo mv netspeed /usr/local/bin/

# Windows
# 将 netspeed.exe 添加到系统 PATH
```

### 跨平台编译

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o netspeed-linux-amd64

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o netspeed-darwin-amd64

# macOS ARM64 (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o netspeed-darwin-arm64

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o netspeed-windows-amd64.exe
```

## 使用方法

### 基础命令

```bash
# 查看帮助
netspeed -help

# 测试网站速度
netspeed -test

# 获取 IP 信息
netspeed -ip

# ⭐ 新功能: 检测 IP 纯净度
netspeed -purity
```

### 高级用法

```bash
# 使用代理测试（SOCKS5）
netspeed -test -proxy socks5://127.0.0.1:1080

# 使用代理测试（HTTP）
netspeed -test -proxy http://proxy.example.com:8080

# 持续监控模式（每 30 秒刷新）
netspeed -test -watch 30

# 使用自定义配置文件
netspeed -test -config sites.json

# 设置超时时间（默认 10 秒）
netspeed -test -timeout 5

# 组合使用
netspeed -test -proxy socks5://127.0.0.1:1080 -watch 60
```

### 输出示例

#### 网站测试输出

```
🚀 开始测试 12 个网站...

┌─────────────────┬──────────────┬────────────────────────────┬──────────┐
│ 网站              │ 延迟           │ URL                        │ 状态       │
├─────────────────┼──────────────┼────────────────────────────┼──────────┤
│ ✓ Google        │     45 ms    │ https://www.google.com     │ ✓ 优秀     │
│ ✓ GitHub        │     78 ms    │ https://github.com         │ ✓ 优秀     │
│ ✓ YouTube       │    156 ms    │ https://www.youtube.com    │ ✓ 优秀     │
│ ⚠ Twitter       │    523 ms    │ https://twitter.com        │ ⚠ 一般     │
│ ✗ Facebook      │    Timeout   │ https://www.facebook.com   │ ✗ 超时     │
└─────────────────┴──────────────┴────────────────────────────┴──────────┘

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 统计摘要
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
在线网站: 4/5 (80.0%)
平均延迟: 200 ms
最低延迟: 45 ms (Google)
最高延迟: 523 ms (Twitter)
网络质量: 良好
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### IP 信息输出

```
🌐 正在获取 IP 信息...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📍 IP 地址:    14.153.68.158
🌍 国家:       中国
📌 地区:       广东省
🏙️  城市:       深圳市
🔌 ISP:        福田中国电信
🏢 组织:       AS4134 CHINANET Guangdong province network
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### ⭐ IP 纯净度检测输出

```
🔍 正在检测 IP 纯净度...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 IP 纯净度报告
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📍 IP 地址:       14.153.68.158
🏭 组织:          AS4134 CHINANET Guangdong province network

✨ 纯净度评分:    100.0/100  (优秀)
⚠️  风险等级:     低风险

🔎 检测结果:
  ✅ VPN: 未检测到
  ✅ 代理: 未检测到
  ✅ Tor: 未检测到
  ✅ 数据中心: 未检测到
  ✅ 黑名单: 未检测到

💡 建议:
  ✓ IP 纯净度很高，适合大多数场景使用
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**纯净度评分标准：**
- 90-100 分：优秀，住宅 IP
- 75-89 分：良好
- 50-74 分：一般，可能有限制
- 25-49 分：较差，建议更换
- 0-24 分：很差，强烈建议更换

**检测项说明：**
- VPN: 检测是否使用 VPN 服务
- 代理: 检测是否使用代理服务器
- Tor: 检测是否使用 Tor 网络
- 数据中心: 检测是否为数据中心 IP
- 黑名单: 检测是否在反垃圾邮件黑名单中

## 自定义配置文件

创建 `sites.json` 文件来自定义测试站点：

```json
[
  {
    "Name": "Google",
    "URL": "https://www.google.com"
  },
  {
    "Name": "GitHub",
    "URL": "https://github.com"
  },
  {
    "Name": "Baidu",
    "URL": "https://www.baidu.com"
  },
  {
    "Name": "Taobao",
    "URL": "https://www.taobao.com"
  }
]
```

使用自定义配置：

```bash
netspeed -test -config sites.json
```

## 延迟评级标准

| 延迟范围 | 状态 | 说明 |
|---------|------|------|
| < 200ms | 优秀 | 网络连接极佳 |
| 200-500ms | 良好 | 网络连接正常 |
| 500-1000ms | 一般 | 网络有一定延迟 |
| > 1000ms | 较差 | 网络延迟较高 |
| Timeout | 超时 | 无法连接或被屏蔽 |

## 代理使用场景

### VPN 连接检测

```bash
# 启动 VPN 后检测出口 IP
netspeed -ip

# 测试 VPN 连接质量
netspeed -test
```

### 企业代理环境

```bash
# 使用公司 HTTP 代理
netspeed -test -proxy http://proxy.company.com:8080

# 使用 SOCKS5 代理
netspeed -test -proxy socks5://proxy.company.com:1080
```

### 持续监控代理连接

```bash
# 每分钟检测一次代理连接质量
netspeed -test -proxy socks5://127.0.0.1:1080 -watch 60
```

## 技术实现

### 核心特性

- **并发测试** - 使用 goroutine 并行测试所有网站
- **超时控制** - 每个请求都有精确的超时控制
- **降级策略** - HEAD 请求失败自动降级到 GET 请求
- **故障转移** - IP 检测支持多个 API 自动切换
- **代理支持** - 完整支持 HTTP/HTTPS/SOCKS5 代理协议

### 默认测试站点

- Google, YouTube, Facebook, Twitter, Instagram
- GitHub, Reddit, Wikipedia, Amazon
- Netflix, OpenAI, Telegram

### IP 检测 API

使用以下 API（按优先级，支持故障转移）：

1. ipapi.co
2. ipinfo.io
3. ip-api.com

## 项目结构

```
net-speed/
├── main.go                  # 主程序
├── go.mod                   # Go 模块定义
├── go.sum                   # 依赖校验和
├── sites.example.json       # 示例配置文件
├── README.md                # 项目文档
└── netspeed.exe            # 编译后的可执行文件（Windows）
```

## 依赖

- Go 1.21+
- golang.org/x/net/proxy (SOCKS5 支持)

## 常见问题

### 1. 某些网站总是超时？

可能原因：
- 网站被防火墙屏蔽
- 网络环境不稳定
- 超时时间过短

解决方案：
```bash
# 增加超时时间到 20 秒
netspeed -test -timeout 20
```

### 2. 如何在 VPN 环境下使用？

VPN 通常会自动接管系统网络，无需额外配置：

```bash
# 直接测试即可（自动使用 VPN）
netspeed -test
netspeed -ip
```

### 3. 代理设置不生效？

检查代理 URL 格式：

```bash
# 正确格式
netspeed -test -proxy socks5://127.0.0.1:1080
netspeed -test -proxy http://proxy:8080

# 错误格式（缺少协议）
netspeed -test -proxy 127.0.0.1:1080  # ❌
```

### 4. Windows 下无法运行？

确保：
1. 使用 PowerShell 或 CMD（不是 Git Bash）
2. 文件名是 `netspeed.exe`
3. 已添加到 PATH 或使用完整路径

```powershell
# PowerShell 示例
.\netspeed.exe -test
```

## 开发计划

- [ ] 添加 JSON 输出格式支持
- [ ] 支持 IPv6 测试
- [ ] 添加 DNS 解析时间测试
- [ ] 支持历史记录保存
- [ ] 添加 Web UI 界面
- [ ] 支持导出测试报告（PDF/HTML）

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 相关项目

基于 [NetPulse Chrome Extension](NetPulse/) 的思路开发。
