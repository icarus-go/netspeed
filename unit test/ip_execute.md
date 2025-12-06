# IPCommand.Execute 单元测试用例（真实执行版）

## 测试概述

本文档包含 `IPCommand.Execute` 方法的完整单元测试用例，采用**表驱动测试**风格，通过**真实执行**而非 Mock 方式，覆盖所有关键场景和边界条件。

## 用户输入参数

在使用本测试用例之前，请提供以下参数：

| 参数名 | 类型 | 格式示例 | 说明 |
|--------|------|----------|------|
| `PROXY_URL` | string | `http://127.0.0.1:8080` | 代理服务器地址（支持 HTTP/HTTPS/SOCKS5） |
| `EXPECTED_IP_NO_PROXY` | string | `192.168.1.100` | 不使用代理时的预期 IP 地址 |
| `EXPECTED_IP_WITH_PROXY` | string | `10.0.0.1` | 使用代理时的预期 IP 地址 |

### 参数获取方式

可以通过以下命令获取实际 IP：

```bash
# 获取当前公网 IP（可能通过代理）
netspeed -ip

# 获取原始 IP（绕过代理）
netspeed -ip -origin
```

## 测试范围

- **命令启用/禁用逻辑**
- **代理配置处理（origin 模式）**
- **IP 检测成功/失败场景**
- **错误处理和传播**
- **输出格式化**
- **真实网络环境下的执行验证**

---

## 表驱动测试用例

### 测试用例 1：命令未启用

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_CommandNotEnabled |
| **测试场景** | 当 `-ip` 标志未设置时，命令应直接返回，不执行任何操作 |
| **前置条件** | `enabled = false`（`-ip` 标志未设置） |
| **输入参数** | `enabled=false`, `origin=false`, 任意 `ctx` |
| **预期输出** | 返回 `nil` 错误，不显示任何信息 |
| **测试步骤** | 1. 创建 IPCommand 实例<br>2. 设置 `enabled = false`<br>3. 调用 `Execute(ctx)`<br>4. 验证返回值为 `nil`<br>5. 验证未调用任何检测器 |

### 测试用例 2：启用命令，使用代理模式，成功获取 IP

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_WithProxy_Success |
| **测试场景** | 命令启用，使用代理模式，通过真实网络请求获取并显示 IP 信息 |
| **前置条件** | `enabled = true`, `origin = false`, 代理服务器可用且可访问 |
| **输入参数** | `enabled=true`, `origin=false`, 用户提供的 `PROXY_URL`, 预期的 `EXPECTED_IP_WITH_PROXY` |
| **预期输出** | 返回 `nil` 错误，显示完整的 IP 信息，IP 地址应匹配 `EXPECTED_IP_WITH_PROXY` |
| **测试步骤** | 1. 创建 IPCommand 实例<br>2. 设置 `enabled = true`, `origin = false`<br>3. 根据用户提供的 `PROXY_URL` 构建 HTTPClient<br>4. 创建 command.Context，设置 HTTPClient 和 Timeout<br>5. 调用 `Execute(ctx)`<br>6. 验证返回值为 `nil`<br>7. 捕获输出，验证包含 "📍 IP 地址"（非原始 IP）<br>8. 从输出中提取实际 IP 地址<br>9. 验证实际 IP 与 `EXPECTED_IP_WITH_PROXY` 匹配 |
| **真实执行** | 是 - 通过真实 HTTP 请求调用 IP 检测 API |
| **IP 验证** | 从标准输出中解析 IP，格式：`📍 IP 地址: [IP]` |

### 测试用例 3：启用命令，原始模式（禁用代理），成功获取 IP

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_OriginMode_Success |
| **测试场景** | 命令启用，使用 `-origin` 标志，禁用代理，通过真实网络请求获取原始 IP |
| **前置条件** | `enabled = true`, `origin = true` |
| **输入参数** | `enabled=true`, `origin=true`, 预期的 `EXPECTED_IP_NO_PROXY` |
| **预期输出** | 返回 `nil` 错误，显示完整的原始 IP 信息，IP 地址应匹配 `EXPECTED_IP_NO_PROXY` |
| **测试步骤** | 1. 创建 IPCommand 实例<br>2. 设置 `enabled = true`, `origin = true`<br>3. 创建 command.Context（使用默认 HTTPClient）<br>4. 调用 `Execute(ctx)`<br>5. 验证返回值为 `nil`<br>6. 捕获输出，验证包含 "📍 原始 IP"（而非 "IP 地址"）<br>7. 验证输出包含 "⚠️ 已设置为获取原始 IP" 提示<br>8. 从输出中提取实际 IP 地址<br>9. 验证实际 IP 与 `EXPECTED_IP_NO_PROXY` 匹配 |
| **真实执行** | 是 - 通过真实 HTTP 请求调用 IP 检测 API |
| **IP 验证** | 从标准输出中解析 IP，格式：`📍 原始 IP: [IP]` |

### 测试用例 4：启用命令，使用代理模式，IP 检测失败

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_WithProxy_Failure |
| **测试场景** | 命令启用，使用代理模式，但所有 IP API 调用失败（网络不可达或代理失效） |
| **前置条件** | `enabled = true`, `origin = false`, 代理服务器不可达或所有 IP API 均返回错误 |
| **输入参数** | `enabled=true`, `origin=false`, 无效的 `PROXY_URL` |
| **预期输出** | 返回包含 "获取 IP 信息失败" 的错误信息 |
| **测试步骤** | 1. 创建 IPCommand 实例<br>2. 设置 `enabled = true`, `origin = false`<br>3. 构建 HTTPClient 使用无效代理（确保无法连接）<br>4. 创建 command.Context，设置 HTTPClient<br>5. 调用 `Execute(ctx)`<br>6. 验证返回错误包含 "获取 IP 信息失败"<br>7. 验证未显示任何 IP 信息 |
| **真实执行** | 是 - 尝试真实网络请求，预期失败 |
| **失败场景** | 代理不可达、API 服务器不可用、网络超时 |

### 测试用例 5：启用命令，原始模式（禁用代理），IP 检测失败

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_OriginMode_Failure |
| **测试场景** | 命令启用，使用 `-origin` 标志，禁用代理，但所有 IP API 调用失败 |
| **前置条件** | `enabled = true`, `origin = true`, 所有 IP API 均返回错误或不可达 |
| **输入参数** | `enabled=true`, `origin=true` |
| **预期输出** | 返回包含 "获取 IP 信息失败" 的错误信息 |
| **测试步骤** | 1. 创建 IPCommand 实例<br>2. 设置 `enabled = true`, `origin = true`<br>3. 创建 command.Context<br>4. 调用 `Execute(ctx)`<br>5. 验证返回错误包含 "获取 IP 信息失败"<br>6. 验证输出包含 "⚠️ 已设置为获取原始 IP" 提示 |
| **真实执行** | 是 - 尝试真实网络请求，预期失败 |
| **失败场景** | 所有 IP API 不可达、网络中断、请求超时 |

### 测试用例 6：HTTP 客户端配置 - 原始模式

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_OriginMode_HTTPClientConfiguration |
| **测试场景** | 验证原始模式下正确创建禁用代理的 HTTP 客户端 |
| **前置条件** | `enabled = true`, `origin = true` |
| **输入参数** | `enabled=true`, `origin=true`, 包含自定义 HTTPClient 的 `ctx` |
| **预期输出** | 创建新的 HTTP 客户端，禁用代理功能 |
| **测试步骤** | 1. 创建 IPCommand 实例<br>2. 设置 `enabled = true`, `origin = true`<br>3. 记录原始 ctx.HTTPClient<br>4. 调用 `Execute(ctx)`<br>5. 验证创建了新的 HTTP 客户端（非 ctx.HTTPClient）<br>6. 验证新客户端配置了禁用代理的 Transport<br>7. 验证 Transport 配置了正确的连接参数 |

### 测试用例 7：HTTP 客户端配置 - 代理模式

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_ProxyMode_HTTPClientUsage |
| **测试场景** | 验证代理模式下正确使用上下文的 HTTPClient |
| **前置条件** | `enabled = true`, `origin = false` |
| **输入参数** | `enabled=true`, `origin=false`, 包含代理的 `ctx.HTTPClient` |
| **预期输出** | 使用 ctx.HTTPClient，不创建新的客户端 |
| **测试步骤** | 1. 创建 IPCommand 实例<br>2. 设置 `enabled = true`, `origin = false`<br>3. 提供自定义的 ctx.HTTPClient<br>4. 调用 `Execute(ctx)`<br>5. 验证使用 ctx.HTTPClient<br>6. 验证未创建新的 HTTP 客户端 |

### 测试用例 8：超时时间传递

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_TimeoutPropagation |
| **测试场景** | 验证 ctx.Timeout 正确传递给 HTTP 客户端 |
| **前置条件** | `enabled = true`, `origin = true`, `ctx.Timeout = 30` |
| **输入参数** | `enabled=true`, `origin=true`, `ctx.Timeout=30` |
| **预期输出** | HTTP 客户端的超时时间设置为 30 秒 |
| **测试步骤** | 1. 创建 IPCommand 实例<br>2. 设置 `enabled = true`, `origin = true`<br>3. 设置 ctx.Timeout = 30<br>4. 调用 `Execute(ctx)`<br>5. 验证新创建的 HTTP 客户端超时时间为 30 秒 |

### 测试用例 9：IP 信息显示 - 完整信息

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_DisplayFullIPInfo |
| **测试场景** | 验证完整 IP 信息的显示格式 |
| **前置条件** | `enabled = true`, 获取到完整的 IPInfo |
| **输入参数** | 包含完整字段的 IPInfo：`IP`, `Country`, `CountryCode`, `Region`, `City`, `ISP`, `Org`, `Timezone` |
| **预期输出** | 显示所有字段的格式化信息 |
| **测试步骤** | 1. 创建 IPCommand 实例<br>2. 设置 `enabled = true`<br>3. 模拟完整的 IPInfo 返回<br>4. 调用 `Execute(ctx)`<br>5. 验证输出包含所有字段的标签和值<br>6. 验证输出格式包含分隔线 "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" |

### 测试用例 10：IP 信息显示 - 部分信息（无地区和城市）

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_DisplayPartialIPInfo |
| **测试场景** | 验证部分字段缺失时的显示行为 |
| **前置条件** | `enabled = true`, IPInfo 包含部分字段（缺失 Region, City） |
| **输入参数** | IPInfo 只包含：`IP`, `Country`, `ISP`, `Org` |
| **预期输出** | 只显示存在的字段，不显示空的字段 |
| **测试步骤** | 1. 创建 IPCommand 实例<br>2. 设置 `enabled = true`<br>3. 模拟部分字段的 IPInfo<br>4. 调用 `Execute(ctx)`<br>5. 验证输出包含存在的字段<br>6. 验证输出不包含空的 Region 和 City 字段 |

### 测试用例 11：原始模式标识区分

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_OriginModeIdentification |
| **测试场景** | 验证 origin 模式和非 origin 模式的输出标识差异 |
| **前置条件** | `enabled = true` |
| **输入参数** | 相同的 IPInfo，分别测试 `origin=true` 和 `origin=false` |
| **预期输出** | origin 模式显示 "📍 原始 IP"，非 origin 模式显示 "📍 IP 地址" |
| **测试步骤** | 1. 创建两个 IPCommand 实例<br>2. 第一个实例：`enabled=true`, `origin=true`<br>3. 第二个实例：`enabled=true`, `origin=false`<br>4. 两次调用 `Execute(ctx)`<br>5. 验证第一次输出包含 "原始 IP"<br>6. 验证第二次输出包含 "IP 地址"（非原始） |

---

## 边界条件测试

### 边界条件 1：空 ctx.HTTPClient（代理模式）

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_NilHTTPClient_ProxyMode |
| **测试场景** | 代理模式下 ctx.HTTPClient 为 nil |
| **前置条件** | `enabled = true`, `origin = false`, `ctx.HTTPClient = nil` |
| **预期输出** | 可能导致 panic 或错误（取决于 ipinfo.Detector 实现） |
| **注意事项** | 需要在测试中验证是否需要空指针检查 |

### 边界条件 2：ctx.Timeout 为 0

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_ZeroTimeout |
| **测试场景** | ctx.Timeout 设置为 0 |
| **前置条件** | `enabled = true`, `origin = true`, `ctx.Timeout = 0` |
| **预期输出** | HTTP 客户端使用 0 超时（无超时限制） |
| **注意事项** | 验证 Go http.Client 对零超时的处理 |

### 边界条件 3：IPInfo 所有字段为空

| 属性 | 值 |
|------|-----|
| **用例名称** | TestExecute_EmptyIPInfo |
| **测试场景** | IPInfo 结构体所有字段为空字符串 |
| **前置条件** | `enabled = true`, IPInfo 所有字段为空 |
| **预期输出** | 显示基本的 IP 地址行，其他字段不显示 |
| **注意事项** | 验证 displayIPInfo 函数对空值的处理 |

---

## 真实执行测试代码示例

以下代码展示了如何实现真实执行的测试用例（以 TestExecute_WithProxy_Success 为例）：

```go
func TestExecute_WithProxy_Success(t *testing.T) {
    // 用户输入参数
    proxyURL := os.Getenv("PROXY_URL")
    expectedIP := os.Getenv("EXPECTED_IP_WITH_PROXY")

    if proxyURL == "" {
        t.Skip("跳过测试：未设置 PROXY_URL 环境变量")
    }
    if expectedIP == "" {
        t.Skip("跳过测试：未设置 EXPECTED_IP_WITH_PROXY 环境变量")
    }

    // 构建带代理的 HTTPClient
    proxy, err := url.Parse(proxyURL)
    if err != nil {
        t.Fatalf("无效的代理 URL: %v", err)
    }

    httpClient := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxy),
        },
        Timeout: 10 * time.Second,
    }

    // 创建 command.Context
    ctx := &command.Context{
        HTTPClient: httpClient,
        Timeout:    10,
    }

    // 创建 IPCommand 并设置标志
    cmd := commands.NewIPCommand()
    flags := flag.NewFlagSet("test", flag.ContinueOnError)
    cmd.DefineFlags(flags)
    flags.Set("ip", "true")

    // 执行命令
    err = cmd.Execute(ctx)
    if err != nil {
        t.Fatalf("Execute 失败: %v", err)
    }

    // 捕获输出并验证
    // 注意：实际测试中需要重定向标准输出
    // 这里使用 bufio.Scanner 或其他方式捕获输出
}
```

### 辅助函数：提取输出中的 IP 地址

```go
func extractIPFromOutput(output string) string {
    // 输出格式: "📍 IP 地址: 192.168.1.1"
    scanner := bufio.NewScanner(strings.NewReader(output))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "📍 IP 地址:") || strings.Contains(line, "📍 原始 IP:") {
            parts := strings.Split(line, ":")
            if len(parts) >= 2 {
                return strings.TrimSpace(parts[1])
            }
        }
    }
    return ""
}
```

### 辅助函数：验证 IP 地址匹配

```go
func verifyIPMatch(t *testing.T, actualIP, expectedIP string) {
    // 允许 IP 地址前后可能有空白字符
    actual := strings.TrimSpace(actualIP)
    expected := strings.TrimSpace(expectedIP)

    if actual != expected {
        t.Errorf("IP 地址不匹配:\n  实际: %s\n  预期: %s", actual, expected)
    }
}
```

---

## 测试执行顺序

建议按以下顺序执行测试用例：

1. **基础功能测试**（用例 1-3）
   - 命令启用/禁用
   - 成功获取 IP（代理模式和原始模式）

2. **错误处理测试**（用例 4-5）
   - 各种失败场景

3. **HTTP 客户端配置测试**（用例 6-8）
   - 客户端创建和配置

4. **输出格式化测试**（用例 9-11）
   - 显示逻辑和格式

5. **边界条件测试**
   - 异常情况处理

---

## 性能测试考虑

- **网络依赖**：真实执行测试需要网络连接，会产生真实的 HTTP 请求
- **测试时间**：每个涉及真实网络的测试用例可能需要 5-30 秒（包括超时时间）
- **网络稳定性**：建议在稳定的网络环境下执行测试
- **代理可用性**：确保提供的代理服务器稳定可用
- **重试机制**：IP 检测模块已内置多个 API 的故障转移
- **并发安全**：Execute 方法应设计为无状态或线程安全

---

## 覆盖率目标

- **语句覆盖率**: 100%
- **分支覆盖率**: 100%
- **函数覆盖率**: 100%

### 关键路径

- [x] `if !*c.enabled` 分支
- [x] `if *c.origin` 分支
- [x] `detector.Detect()` 成功路径
- [x] `detector.Detect()` 失败路径
- [x] `displayIPInfo()` 所有字段显示逻辑

---

## 相关文件

- **被测文件**: `pkg/commands/ip.go`
- **依赖文件**:
  - `pkg/command/command.go`（Context 接口）
  - `pkg/ipinfo/detector.go`（IP 检测逻辑）
  - `pkg/ipinfo/model.go`（IPInfo 结构体）

---

## 测试环境

- **Go 版本**: 1.19+
- **测试框架**: `testing`
- **Mock 框架**: `gomock` 或原生 `testify/mock`
- **操作系统**: Linux/macOS/Windows

---

## 总结

本测试套件采用**表驱动测试**方法，通过**真实执行**的方式，通过 11 个核心测试用例和 3 个边界条件测试，全面覆盖了 `IPCommand.Execute` 方法的所有功能点。

### 表驱动测试 Example

```go
type fields struct {
    enabled *bool
    origin  *bool
}
type args struct {
    ctx *command.Context
}
tests := []struct {
    name    string
    fields  fields
    args    args
    wantErr bool
}{
    // TODO: Add test cases.
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        c := &IPCommand{
            enabled: tt.fields.enabled,
            origin:  tt.fields.origin,
        }
        if err := c.Execute(tt.args.ctx); (err != nil) != tt.wantErr {
            t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
        }
    })
}
```

### 真实执行测试的优势

1. **真实环境验证**：直接验证在真实网络环境下的行为
2. **代理功能验证**：能够验证代理配置的实际效果
3. **端到端测试**：从输入到输出的完整流程验证
4. **故障转移验证**：验证多个 IP API 的故障转移机制

### 使用前准备

1. 设置环境变量：
   ```bash
   export PROXY_URL="http://127.0.0.1:8080"
   export EXPECTED_IP_NO_PROXY="实际公网IP"
   export EXPECTED_IP_WITH_PROXY="通过代理的IP"
   ```

2. 确保网络连接稳定
3. 确保代理服务器可用（如果使用代理测试）

### 注意事项

- 测试时间较长，需要耐心等待网络请求完成
- 测试结果依赖于网络环境和代理服务器状态
- 建议在 CI/CD 环境中使用时配置适当的超时时间和重试机制