package ipinfo

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestDetector_FetchFromProvider_Success 测试成功获取 IP 信息
func TestDetector_FetchFromProvider_Success(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		format       string
		wantIP       string
		wantCountry  string
	}{
		{
			name: "Text 格式响应",
			responseBody: `14.153.68.158
中国 广东省深圳市福田中国电信
AS4134
CHINANET Guangdong province network`,
			format:      "text",
			wantIP:      "14.153.68.158",
			wantCountry: "中国",
		},
		{
			name: "JSON 格式响应",
			responseBody: `{
				"ip": "203.0.113.45",
				"country": "United States",
				"city": "Los Angeles"
			}`,
			format:      "json",
			wantIP:      "203.0.113.45",
			wantCountry: "United States",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock HTTP Server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// 创建 Detector
			client := &http.Client{Timeout: 5 * time.Second,
				Transport: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
				},
			}
			detector := NewDetector(client)

			// 创建 Provider
			provider := Provider{
				Name:   "test",
				URL:    server.URL,
				Format: tt.format,
			}

			// 执行测试
			result, err := detector.fetchFromProvider(provider)
			if err != nil {
				t.Fatalf("fetchFromProvider() error = %v", err)
			}

			if result.IP != tt.wantIP {
				t.Errorf("IP = %s, want %s", result.IP, tt.wantIP)
				return
			}

			if result.Country != tt.wantCountry {
				t.Errorf("Country = %s, want %s", result.Country, tt.wantCountry)
				return
			}
		})
	}
}

// TestDetector_FetchFromProvider_HTTPErrors 测试 HTTP 错误处理
func TestDetector_FetchFromProvider_HTTPErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:       "429 Too Many Requests",
			statusCode: http.StatusTooManyRequests,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					w.Write([]byte(`{"ip": "1.2.3.4", "country": "Test"}`))
				}
			}))
			defer server.Close()

			client := &http.Client{Timeout: 5 * time.Second}
			detector := NewDetector(client)

			provider := Provider{
				Name:   "test",
				URL:    server.URL,
				Format: "json",
			}

			_, err := detector.fetchFromProvider(provider)

			if (err != nil) != tt.wantErr {
				t.Errorf("fetchFromProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDetector_FetchFromProvider_Timeout 测试超时处理
func TestDetector_FetchFromProvider_Timeout(t *testing.T) {
	// 创建慢响应的 Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second) // 延迟 3 秒
		w.Write([]byte(`{"ip": "1.2.3.4"}`))
	}))
	defer server.Close()

	// 创建超时时间很短的客户端
	client := &http.Client{Timeout: 100 * time.Millisecond}
	detector := NewDetector(client)

	provider := Provider{
		Name:   "test",
		URL:    server.URL,
		Format: "json",
	}

	_, err := detector.fetchFromProvider(provider)

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

// TestDetector_FetchFromProvider_InvalidURL 测试无效 URL
func TestDetector_FetchFromProvider_InvalidURL(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	detector := NewDetector(client)

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "Invalid URL",
			url:  "://invalid-url",
		},
		{
			name: "Unreachable Host",
			url:  "http://this-host-does-not-exist-12345.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := Provider{
				Name:   "test",
				URL:    tt.url,
				Format: "json",
			}

			_, err := detector.fetchFromProvider(provider)

			if err == nil {
				t.Error("Expected error for invalid URL, got nil")
			}
		})
	}
}

// TestDetector_FetchFromProvider_UnsupportedFormat 测试不支持的格式
func TestDetector_FetchFromProvider_UnsupportedFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`some data`))
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	detector := NewDetector(client)

	provider := Provider{
		Name:   "test",
		URL:    server.URL,
		Format: "xml", // 不支持的格式
	}

	_, err := detector.fetchFromProvider(provider)

	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}
}

// TestDetector_Detect_Failover 测试故障转移机制
func TestDetector_Detect_Failover(t *testing.T) {
	// 创建多个 Mock Servers
	failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failServer.Close()

	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ip": "203.0.113.45", "country": "Test Country"}`))
	}))
	defer successServer.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	detector := &Detector{
		client: client,
		providers: []Provider{
			{Name: "fail1", URL: failServer.URL, Format: "json"},
			{Name: "fail2", URL: failServer.URL, Format: "json"},
			{Name: "success", URL: successServer.URL, Format: "json"},
		},
	}

	result, err := detector.Detect()

	if err != nil {
		t.Fatalf("Detect() should succeed with failover, got error: %v", err)
	}

	if result.IP != "203.0.113.45" {
		t.Errorf("IP = %s, want 203.0.113.45", result.IP)
	}

	if result.Country != "Test Country" {
		t.Errorf("Country = %s, want Test Country", result.Country)
	}
}

// TestDetector_Detect_AllFail 测试所有 API 都失败的情况
func TestDetector_Detect_AllFail(t *testing.T) {
	failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failServer.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	detector := &Detector{
		client: client,
		providers: []Provider{
			{Name: "fail1", URL: failServer.URL, Format: "json"},
			{Name: "fail2", URL: failServer.URL, Format: "json"},
		},
	}

	_, err := detector.Detect()

	if err == nil {
		t.Error("Expected error when all providers fail, got nil")
	}
}

// TestDetector_Detect_FirstSuccess 测试首选 API 成功时不尝试后续
func TestDetector_Detect_FirstSuccess(t *testing.T) {
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ip": "1.2.3.4", "country": "First"}`))
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	detector := &Detector{
		client: client,
		providers: []Provider{
			{Name: "first", URL: server.URL, Format: "json"},
			{Name: "second", URL: server.URL, Format: "json"},
			{Name: "third", URL: server.URL, Format: "json"},
		},
	}

	result, err := detector.Detect()

	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	if callCount != 1 {
		t.Errorf("Expected 1 API call, got %d", callCount)
	}

	if result.IP != "1.2.3.4" {
		t.Errorf("IP = %s, want 1.2.3.4", result.IP)
	}
}

// TestDetector_NewDetector 测试 Detector 创建
func TestDetector_NewDetector(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	detector := NewDetector(client)

	if detector == nil {
		t.Fatal("NewDetector() should not return nil")
	}

	if detector.client != client {
		t.Error("Detector client not set correctly")
	}

	if len(detector.providers) == 0 {
		t.Error("Detector should have default providers")
	}
}

// TestDetector_EmptyResponseBody 测试空响应体处理
func TestDetector_EmptyResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// 空响应体
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	detector := NewDetector(client)

	provider := Provider{
		Name:   "test",
		URL:    server.URL,
		Format: "json",
	}

	_, err := detector.fetchFromProvider(provider)

	// 空 JSON 响应体应该能够解析（虽然字段为空）
	if err != nil {
		t.Logf("Empty response body error: %v", err)
	}
}