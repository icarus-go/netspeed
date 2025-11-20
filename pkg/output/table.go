package output

import (
	"fmt"
	"time"

	"github.com/icarus-go/net-speed/pkg/tester"
)

// PrintResultsTable 以表格形式输出结果
func PrintResultsTable(results []tester.TestResult) {
	// 表头
	fmt.Println("┌─────────────────┬──────────────┬────────────────────────────┬──────────┐")
	fmt.Printf("│ %-15s │ %-12s │ %-26s │ %-8s │\n", "网站", "延迟", "URL", "状态")
	fmt.Println("├─────────────────┼──────────────┼────────────────────────────┼──────────┤")

	// 数据行
	for _, result := range results {
		var statusIcon string
		var latencyStr string

		if result.Success {
			statusIcon = "✓"
			latencyStr = fmt.Sprintf("%6d ms", result.Latency.Milliseconds())

			// 根据状态着色
			switch result.Status {
			case "优秀":
				statusIcon = "✓"
			case "良好":
				statusIcon = "✓"
			case "一般":
				statusIcon = "⚠"
			case "较差":
				statusIcon = "!"
			}
		} else {
			statusIcon = "✗"
			latencyStr = "   Timeout"
		}

		// 截断 URL
		url := result.URL
		if len(url) > 26 {
			url = url[:23] + "..."
		}

		fmt.Printf("│ %s %-13s │ %-12s │ %-26s │ %s %-6s │\n",
			statusIcon,
			result.Name,
			latencyStr,
			url,
			statusIcon,
			result.Status,
		)
	}

	fmt.Println("└─────────────────┴──────────────┴────────────────────────────┴──────────┘")
}

// PrintSummary 打印统计摘要
func PrintSummary(results []tester.TestResult) {
	var online, total int
	var totalLatency time.Duration
	var minLatency, maxLatency time.Duration
	var minSite, maxSite string

	total = len(results)
	minLatency = time.Hour // 初始值设为很大

	for _, result := range results {
		if result.Success {
			online++
			totalLatency += result.Latency

			if result.Latency < minLatency {
				minLatency = result.Latency
				minSite = result.Name
			}
			if result.Latency > maxLatency {
				maxLatency = result.Latency
				maxSite = result.Name
			}
		}
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 统计摘要")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("在线网站: %d/%d (%.1f%%)\n", online, total, float64(online)/float64(total)*100)

	if online > 0 {
		avgLatency := totalLatency / time.Duration(online)
		fmt.Printf("平均延迟: %d ms\n", avgLatency.Milliseconds())
		fmt.Printf("最低延迟: %d ms (%s)\n", minLatency.Milliseconds(), minSite)
		fmt.Printf("最高延迟: %d ms (%s)\n", maxLatency.Milliseconds(), maxSite)

		// 网络质量评级
		quality := "优秀"
		if avgLatency.Milliseconds() > 500 {
			quality = "一般"
		} else if avgLatency.Milliseconds() > 200 {
			quality = "良好"
		}
		fmt.Printf("网络质量: %s\n", quality)
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
