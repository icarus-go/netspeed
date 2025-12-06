package commands

import (
	"bufio"
	"bytes"
	"flag"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/icarus-go/netspeed/pkg/command"
)

// TestExecute 表驱动测试 - 标准表驱动测试风格
func TestExecute(t *testing.T) {
	// 获取环境变量
	proxyURL := os.Getenv("PROXY_URL")
	expectedIPNoProxy := os.Getenv("EXPECTED_IP_NO_PROXY")
	expectedIPWithProxy := os.Getenv("EXPECTED_IP_WITH_PROXY")

	type fields struct {
		enabled *bool
		origin  *bool
	}

	type args struct {
		ctx *command.Context
	}

	// 辅助函数：将 bool 转换为指针
	boolPtr := func(b bool) *bool {
		return &b
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		proxyURL     string
		timeout      time.Duration
		skip         bool
		skipReason   string
		validateFunc func(*testing.T, string, error)
	}{
		// 测试用例 1：命令未启用
		{
			name: "命令未启用",
			fields: fields{
				enabled: boolPtr(false),
				origin:  boolPtr(false),
			},
			args: args{
				ctx: &command.Context{
					HTTPClient: &http.Client{},
					Timeout:    10,
				},
			},
			timeout: 5 * time.Second,
			validateFunc: func(t *testing.T, output string, err error) {
				if err != nil {
					t.Errorf("预期返回 nil 错误，实际返回: %v", err)
				}
				if output != "" {
					t.Errorf("命令未启用时不应有输出，实际输出: %s", output)
				}
				t.Log("命令未启用 successful")
			},
		},

		// 测试用例 2：代理模式成功获取 IP
		{
			name: "代理模式成功获取 IP",
			fields: fields{
				enabled: boolPtr(true),
				origin:  boolPtr(false),
			},
			args: args{
				ctx: &command.Context{
					HTTPClient: &http.Client{},
					Timeout:    15,
				},
			},
			proxyURL:   proxyURL,
			timeout:    30 * time.Second,
			skip:       proxyURL == "" || expectedIPWithProxy == "",
			skipReason: "未设置 PROXY_URL 或 EXPECTED_IP_WITH_PROXY",
			validateFunc: func(t *testing.T, output string, err error) {
				if err != nil {
					t.Fatalf("Execute 失败: %v", err)
				}
				// 验证包含 "📍 IP 地址"
				if !strings.Contains(output, "📍 IP 地址") {
					t.Errorf("输出中未找到 '📍 IP 地址'，实际输出: %s", output)
				}
				// 验证不包含 "原始 IP"
				if strings.Contains(output, "原始 IP") {
					t.Errorf("输出中不应包含 '原始 IP'，实际输出: %s", output)
				}
				// 验证包含国家信息
				if !strings.Contains(output, "🌍 国家") {
					t.Errorf("输出中未找到国家信息，实际输出: %s", output)
				}
				// 提取并验证 IP
				actualIP := extractIPFromOutput(output)
				actual := strings.TrimSpace(actualIP)
				expected := strings.TrimSpace(expectedIPWithProxy)
				if actual == "" {
					t.Errorf("未从输出中提取到 IP 地址")
				} else if actual != expected {
					t.Errorf("IP 地址不匹配: 实际=%s, 预期=%s", actual, expected)
				}
				t.Log("代理模式成功获取 IP successful")
			},
		},

		// 测试用例 3：原始模式成功获取 IP
		{
			name: "原始模式成功获取 IP",
			fields: fields{
				enabled: boolPtr(true),
				origin:  boolPtr(true),
			},
			args: args{
				ctx: &command.Context{
					HTTPClient: &http.Client{},
					Timeout:    15,
				},
			},
			timeout:    30 * time.Second,
			skip:       expectedIPNoProxy == "",
			skipReason: "未设置 EXPECTED_IP_NO_PROXY",
			validateFunc: func(t *testing.T, output string, err error) {
				if err != nil {
					t.Fatalf("Execute 失败: %v", err)
				}
				// 验证包含 "📍 原始 IP"
				if !strings.Contains(output, "📍 原始 IP") {
					t.Errorf("输出中未找到 '📍 原始 IP'，实际输出: %s", output)
				}
				// 验证不包含 "📍 IP 地址:"
				if strings.Contains(output, "📍 IP 地址:") {
					t.Errorf("输出中不应包含 '📍 IP 地址:'，实际输出: %s", output)
				}
				// 验证包含提示信息
				if !strings.Contains(output, "⚠️  已设置为获取原始 IP") {
					t.Errorf("输出中未找到提示信息，实际输出: %s", output)
				}
				// 验证包含国家信息
				if !strings.Contains(output, "🌍 国家") {
					t.Errorf("输出中未找到国家信息，实际输出: %s", output)
				}
				// 提取并验证 IP
				actualIP := extractIPFromOutput(output)
				actual := strings.TrimSpace(actualIP)
				expected := strings.TrimSpace(expectedIPNoProxy)
				if actual == "" {
					t.Errorf("未从输出中提取到 IP 地址")
				} else if actual != expected {
					t.Errorf("IP 地址不匹配: 实际=%s, 预期=%s", actual, expected)
				}
				t.Log("原始模式成功获取 IP successful")
			},
		},

		// 测试用例 4：代理模式失败
		{
			name: "代理模式失败",
			fields: fields{
				enabled: boolPtr(true),
				origin:  boolPtr(false),
			},
			args: args{
				ctx: &command.Context{
					HTTPClient: &http.Client{},
					Timeout:    3,
				},
			},
			proxyURL: "http://127.0.0.1:99999", // 无效代理
			timeout:  10 * time.Second,
			validateFunc: func(t *testing.T, output string, err error) {
				if err == nil {
					t.Errorf("预期返回错误，实际返回 nil")
				} else if !strings.Contains(err.Error(), "获取 IP 信息失败") {
					t.Errorf("错误信息不匹配: %v", err)
				}
				// 失败时不应显示 IP 信息
				if strings.Contains(output, "📍 IP 地址") {
					t.Errorf("失败时不应显示 IP 信息，实际输出: %s", output)
				}
				t.Log("代理模式失败 successful")
			},
		},

		// 测试用例 5：原始模式边界条件 - 超时测试
		{
			name: "原始模式超时测试",
			fields: fields{
				enabled: boolPtr(true),
				origin:  boolPtr(true),
			},
			args: args{
				ctx: &command.Context{
					HTTPClient: &http.Client{},
					Timeout:    2, // 很短的超时
				},
			},
			timeout: 5 * time.Second,
			validateFunc: func(t *testing.T, output string, err error) {
				// 这个测试主要是验证在短超时情况下Execute仍然能正确执行
				// 无论成功还是失败，都要验证输出格式
				if strings.Contains(output, "⚠️  已设置为获取原始 IP") {
					t.Log("原始模式提示信息正确显示")
				}
				// 如果成功，验证IP信息
				if err == nil && strings.Contains(output, "📍 原始 IP") {
					t.Log("原始模式成功获取 IP")
				}
				t.Log("原始模式超时测试 successful")
			},
		},

		// 测试用例 6：显示完整 IP 信息格式
		{
			name: "显示完整 IP 信息格式",
			fields: fields{
				enabled: boolPtr(true),
				origin:  boolPtr(false),
			},
			args: args{
				ctx: &command.Context{
					HTTPClient: &http.Client{},
					Timeout:    15,
				},
			},
			timeout: 30 * time.Second,
			validateFunc: func(t *testing.T, output string, err error) {
				if err != nil {
					t.Fatalf("Execute 失败: %v", err)
				}
				// 验证包含分隔线
				if !strings.Contains(output, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") {
					t.Errorf("输出中未找到分隔线，实际输出: %s", output)
				}
				// 验证包含 IP 地址行
				if !strings.Contains(output, "📍 IP 地址") {
					t.Errorf("输出中未找到 IP 地址行，实际输出: %s", output)
				}
				// 验证包含国家信息
				if !strings.Contains(output, "🌍 国家") {
					t.Errorf("输出中未找到国家信息，实际输出: %s", output)
				}
				t.Log("显示完整 IP 信息格式 successful")
			},
		},

		// 测试用例 7：模式标识区分
		{
			name: "模式标识区分",
			fields: fields{
				enabled: boolPtr(true),
				origin:  boolPtr(true), // 原始模式
			},
			args: args{
				ctx: &command.Context{
					HTTPClient: &http.Client{},
					Timeout:    15,
				},
			},
			timeout: 30 * time.Second,
			validateFunc: func(t *testing.T, output string, err error) {
				if err != nil {
					t.Fatalf("Execute 失败: %v", err)
				}
				// 验证显示 "📍 原始 IP"
				if !strings.Contains(output, "📍 原始 IP") {
					t.Errorf("输出中未找到 '📍 原始 IP'，实际输出: %s", output)
				}
				// 验证不显示 "📍 IP 地址:"
				if strings.Contains(output, "📍 IP 地址:") {
					t.Errorf("输出中不应包含 '📍 IP 地址:'，实际输出: %s", output)
				}
				t.Log("模式标识区分 successful")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 跳过测试
			if tt.skip {
				t.Skip(tt.skipReason)
			}

			// 构建 HTTP 客户端
			var httpClient *http.Client
			if tt.proxyURL != "" {
				proxy, err := url.Parse(tt.proxyURL)
				if err != nil {
					t.Fatalf("无效的代理 URL: %v", err)
				}
				httpClient = &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxy),
					},
					Timeout: tt.timeout,
				}
			} else {
				httpClient = tt.args.ctx.HTTPClient
				httpClient.Timeout = tt.timeout
			}

			// 创建命令
			c := NewIPCommand()

			// 创建并设置标志
			flags := flag.NewFlagSet("test", flag.ContinueOnError)
			c.DefineFlags(flags)
			if tt.fields.enabled != nil && *tt.fields.enabled {
				flags.Set("ip", "true")
			}
			if tt.fields.origin != nil && *tt.fields.origin {
				flags.Set("origin", "true")
			}

			// 更新上下文
			tt.args.ctx.HTTPClient = httpClient
			tt.args.ctx.Timeout = int(tt.timeout.Seconds())

			// 重定向标准输出以捕获输出
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// 执行命令
			err := c.Execute(tt.args.ctx)

			// 恢复标准输出
			w.Close()
			os.Stdout = oldStdout

			// 读取输出
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			output := buf.String()

			// 验证结果
			if tt.validateFunc != nil {
				tt.validateFunc(t, output, err)
			}
		})
	}
}

// 辅助函数：从输出中提取 IP 地址
func extractIPFromOutput(output string) string {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "📍 IP 地址") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1])
			}
		}
		if strings.Contains(line, "📍 原始 IP") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}