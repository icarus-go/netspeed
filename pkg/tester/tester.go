package tester

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Tester 网站测试器
type Tester struct {
	client  *http.Client
	timeout time.Duration
}

// NewTester 创建新的测试器
func NewTester(client *http.Client, timeout time.Duration) *Tester {
	return &Tester{
		client:  client,
		timeout: timeout,
	}
}

// TestAll 并行测试所有网站
func (t *Tester) TestAll(sites []Site) []TestResult {
	var wg sync.WaitGroup
	results := make([]TestResult, len(sites))

	for i, site := range sites {
		wg.Add(1)
		go func(idx int, s Site) {
			defer wg.Done()
			results[idx] = t.TestSite(s)
		}(i, site)
	}

	wg.Wait()
	return results
}

// TestSite 测试单个网站的延迟
func (t *Tester) TestSite(site Site) TestResult {
	result := TestResult{
		Name:    site.Name,
		URL:     site.URL,
		Success: false,
	}

	// 创建带超时的请求
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", site.URL+"/favicon.ico", nil)
	if err != nil {
		result.Error = err.Error()
		result.Status = "错误"
		return result
	}

	// 记录开始时间
	start := time.Now()

	// 发送请求
	resp, err := t.client.Do(req)
	latency := time.Since(start)

	if err != nil {
		// 降级：尝试 GET 首页
		result = t.fallbackTest(site, ctx)
		if !result.Success {
			result.Status = "超时"
		}
		return result
	}
	defer resp.Body.Close()

	result.Latency = latency
	result.Success = true
	result.Status = GetStatusByLatency(latency)

	return result
}

// fallbackTest 降级测试（使用 GET 请求）
func (t *Tester) fallbackTest(site Site, ctx context.Context) TestResult {
	result := TestResult{
		Name:    site.Name,
		URL:     site.URL,
		Success: false,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", site.URL, nil)
	if err != nil {
		result.Error = err.Error()
		result.Status = "错误"
		return result
	}

	start := time.Now()
	resp, err := t.client.Do(req)
	latency := time.Since(start)

	if err != nil {
		result.Error = err.Error()
		result.Status = "超时"
		return result
	}
	defer resp.Body.Close()

	result.Latency = latency
	result.Success = true
	result.Status = GetStatusByLatency(latency)

	return result
}
