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
