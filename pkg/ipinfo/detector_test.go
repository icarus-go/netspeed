package ipinfo

import (
	"testing"
)

func TestParseTextFormat(t *testing.T) {
	detector := &Detector{}

	tests := []struct {
		name     string
		input    string
		expected IPInfo
	}{
		{
			name: "中文格式",
			input: `14.153.68.158
中国 广东省深圳市福田中国电信
AS4134
CHINANET Guangdong province network`,
			expected: IPInfo{
				IP:      "14.153.68.158",
				Country: "中国",
				Region:  "广东省",
				City:    "深圳市",
				ISP:     "福田中国电信",
				Org:     "AS4134 CHINANET Guangdong province network",
			},
		},
		{
			name: "简单格式",
			input: `203.0.113.45
United States California Los Angeles AT&T Services
AS7018
ATT-INTERNET4`,
			expected: IPInfo{
				IP:      "203.0.113.45",
				Country: "United",
				Region:  "States",
				City:    "California",
				ISP:     "Los Angeles AT&T Services",
				Org:     "Los Angeles AT&T Services",
			},
		},
		{
			name: "最小格式",
			input: `1.2.3.4
China Beijing`,
			expected: IPInfo{
				IP:      "1.2.3.4",
				Country: "China",
				Region:  "Beijing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.parseTextFormat([]byte(tt.input))
			if err != nil {
				t.Fatalf("解析失败: %v", err)
			}

			if result.IP != tt.expected.IP {
				t.Errorf("IP 不匹配: got %s, want %s", result.IP, tt.expected.IP)
			}

			if result.Country != tt.expected.Country {
				t.Errorf("Country 不匹配: got %s, want %s", result.Country, tt.expected.Country)
			}

			if result.Region != tt.expected.Region {
				t.Errorf("Region 不匹配: got %s, want %s", result.Region, tt.expected.Region)
			}

			if result.City != tt.expected.City {
				t.Errorf("City 不匹配: got %s, want %s", result.City, tt.expected.City)
			}

			if tt.expected.ISP != "" && result.ISP != tt.expected.ISP {
				t.Errorf("ISP 不匹配: got %s, want %s", result.ISP, tt.expected.ISP)
			}
		})
	}
}

func TestParseLocationLine(t *testing.T) {
	detector := &Detector{}

	tests := []struct {
		name     string
		line     string
		expected IPInfo
	}{
		{
			name: "中文完整格式",
			line: "中国 广东省深圳市福田中国电信",
			expected: IPInfo{
				Country: "中国",
				Region:  "广东省",
				City:    "深圳市",
				ISP:     "福田中国电信",
				Org:     "福田中国电信",
			},
		},
		{
			name: "中文省市格式",
			line: "中国 北京市朝阳区联通",
			expected: IPInfo{
				Country: "中国",
				City:    "北京市",
				ISP:     "朝阳区联通",
				Org:     "朝阳区联通",
			},
		},
		{
			name: "英文格式",
			line: "United States California Los Angeles Comcast",
			expected: IPInfo{
				Country: "United",
				Region:  "States",
				City:    "California",
				ISP:     "Los Angeles Comcast",
				Org:     "Los Angeles Comcast",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipInfo := &IPInfo{}
			detector.parseLocationLine(tt.line, ipInfo)

			if ipInfo.Country != tt.expected.Country {
				t.Errorf("Country 不匹配: got %s, want %s", ipInfo.Country, tt.expected.Country)
			}

			if ipInfo.Region != tt.expected.Region {
				t.Errorf("Region 不匹配: got %s, want %s", ipInfo.Region, tt.expected.Region)
			}

			if ipInfo.City != tt.expected.City {
				t.Errorf("City 不匹配: got %s, want %s", ipInfo.City, tt.expected.City)
			}

			if tt.expected.ISP != "" && ipInfo.ISP != tt.expected.ISP {
				t.Errorf("ISP 不匹配: got %s, want %s", ipInfo.ISP, tt.expected.ISP)
			}
		})
	}
}

// TestParseJSONFormat 测试 JSON 格式解析
func TestParseJSONFormat(t *testing.T) {
	detector := &Detector{}

	tests := []struct {
		name      string
		jsonInput string
		wantErr   bool
		expected  IPInfo
	}{
		{
			name: "完整 JSON",
			jsonInput: `{
				"ip": "203.0.113.45",
				"country": "United States",
				"countryCode": "US",
				"region": "California",
				"city": "Los Angeles",
				"isp": "AT&T Services",
				"timezone": "America/Los_Angeles",
				"org": "AS7018 AT&T Services Inc."
			}`,
			wantErr: false,
			expected: IPInfo{
				IP:          "203.0.113.45",
				Country:     "United States",
				CountryCode: "US",
				Region:      "California",
				City:        "Los Angeles",
				ISP:         "AT&T Services",
				Timezone:    "America/Los_Angeles",
				Org:         "AS7018 AT&T Services Inc.",
			},
		},
		{
			name: "部分字段 JSON",
			jsonInput: `{
				"ip": "1.2.3.4",
				"country": "China",
				"city": "Beijing"
			}`,
			wantErr: false,
			expected: IPInfo{
				IP:      "1.2.3.4",
				Country: "China",
				City:    "Beijing",
			},
		},
		{
			name:      "无效 JSON",
			jsonInput: `{invalid json`,
			wantErr:   true,
		},
		{
			name:      "空 JSON",
			jsonInput: `{}`,
			wantErr:   false,
			expected:  IPInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.parseJSONFormat([]byte(tt.jsonInput))

			if (err != nil) != tt.wantErr {
				t.Errorf("parseJSONFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if result.IP != tt.expected.IP {
				t.Errorf("IP = %s, want %s", result.IP, tt.expected.IP)
			}
			if result.Country != tt.expected.Country {
				t.Errorf("Country = %s, want %s", result.Country, tt.expected.Country)
			}
			if result.City != tt.expected.City {
				t.Errorf("City = %s, want %s", result.City, tt.expected.City)
			}
		})
	}
}

// TestParseTextFormat_EdgeCases 测试文本格式边界情况
func TestParseTextFormat_EdgeCases(t *testing.T) {
	detector := &Detector{}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "空输入",
			input:   "",
			wantErr: true,
		},
		{
			name:    "仅 IP 地址",
			input:   "1.2.3.4",
			wantErr: true, // 行数不足
		},
		{
			name: "包含空行",
			input: `1.2.3.4

China Beijing

`,
			wantErr: false,
		},
		{
			name: "超长文本",
			input: `1.2.3.4
China Beijing
AS1234
Very Long Organization Name With Many Words That Might Cause Issues If Not Handled Properly`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := detector.parseTextFormat([]byte(tt.input))

			if (err != nil) != tt.wantErr {
				t.Errorf("parseTextFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestContainsIgnoreCase 测试大小写不敏感的字符串包含
func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{"完全匹配", "VPN Service", "VPN", true},
		{"大小写不同", "vpn service", "VPN", true},
		{"部分匹配", "ExpressVPN", "VPN", true},
		{"不匹配", "Residential ISP", "VPN", false},
		{"空字符串", "", "VPN", false},
		{"子串为空", "VPN", "", true},
		{"中文字符", "中国电信", "电信", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsIgnoreCase(tt.s, tt.substr)
			if result != tt.want {
				t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.want)
			}
		})
	}
}

// TestToLower 测试小写转换
func TestToLower(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"VPN", "vpn"},
		{"ExpressVPN", "expressvpn"},
		{"UPPER", "upper"},
		{"MiXeD", "mixed"},
		{"", ""},
		{"123", "123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toLower(tt.input)
			if result != tt.want {
				t.Errorf("toLower(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}
