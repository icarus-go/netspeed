package tester

import "time"

// Site 表示要测试的网站
type Site struct {
	Name string `json:"Name"`
	URL  string `json:"URL"`
}

// TestResult 测试结果
type TestResult struct {
	Name    string
	URL     string
	Latency time.Duration
	Status  string
	Success bool
	Error   string
}

// GetStatusByLatency 根据延迟判断状态
func GetStatusByLatency(latency time.Duration) string {
	ms := latency.Milliseconds()
	switch {
	case ms < 200:
		return "优秀"
	case ms < 500:
		return "良好"
	case ms < 1000:
		return "一般"
	default:
		return "较差"
	}
}

// DefaultSites 默认测试网站列表
var DefaultSites = []Site{
	{"Google", "https://www.google.com"},
	{"GitHub", "https://github.com"},
	{"YouTube", "https://www.youtube.com"},
	{"Twitter", "https://twitter.com"},
	{"Facebook", "https://www.facebook.com"},
	{"Instagram", "https://www.instagram.com"},
	{"Reddit", "https://www.reddit.com"},
	{"Netflix", "https://www.netflix.com"},
	{"Wikipedia", "https://www.wikipedia.org"},
	{"Amazon", "https://www.amazon.com"},
	{"OpenAI", "https://www.openai.com"},
	{"Telegram", "https://telegram.org"},
}
