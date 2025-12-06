package commands

import (
	"fmt"
	"os"
	"testing"
)

// setup 在测试执行前设置环境变量
func setup() {
	fmt.Println("🔧 正在设置测试环境...")

	// 在这里设置您的环境变量
	// 这些值可以根据需要修改

	// 代理服务器地址（支持 HTTP/HTTPS/SOCKS5）
	os.Setenv("PROXY_URL", "http://127.0.0.1:7897")
	fmt.Println("  ✅ PROXY_URL = http://127.0.0.1:7897")

	// 不使用代理时的预期 IP 地址
	os.Setenv("EXPECTED_IP_NO_PROXY", "223.73.211.200")
	fmt.Println("  ✅ EXPECTED_IP_NO_PROXY = 223.73.211.200")

	// 使用代理时的预期 IP 地址
	os.Setenv("EXPECTED_IP_WITH_PROXY", "154.9.30.64")
	fmt.Println("  ✅ EXPECTED_IP_WITH_PROXY = 154.9.30.64")

	fmt.Println("✅ 环境变量设置完成")
}

// teardown 在测试执行后清理环境变量（可选）
func teardown() {
	fmt.Println("\n🧹 正在清理测试环境...")

	// 可以选择清理环境变量（可选）
	// 如果保留这些变量，可能影响其他测试
	os.Unsetenv("PROXY_URL")
	os.Unsetenv("EXPECTED_IP_NO_PROXY")
	os.Unsetenv("EXPECTED_IP_WITH_PROXY")

	fmt.Println("✅ 清理完成")
}

// TestMain 是特殊命名的函数，用于协调整个测试套件的设置和清理
// 当存在 TestMain 时，Go 会调用它而不是直接运行测试函数
func TestMain(m *testing.M) {
	// 执行测试前的设置
	setup()

	// 运行所有测试
	// m.Run() 返回一个退出代码，0 表示成功
	code := m.Run()

	// 执行测试后的清理
	teardown()

	// 使用退出代码退出程序
	os.Exit(code)
}
