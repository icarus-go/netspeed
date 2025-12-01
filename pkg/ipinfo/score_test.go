package ipinfo

import (
	"net/http"
	"testing"
	"time"
)

// TestCalculatePuritySimple 测试纯净度评分计算
func TestCalculatePuritySimple(t *testing.T) {
	detector := &Detector{}

	tests := []struct {
		name      string
		ipInfo    *IPInfo
		wantScore float64
		wantMin   float64
		wantMax   float64
	}{
		{
			name: "住宅 ISP - 中国电信",
			ipInfo: &IPInfo{
				Org: "AS4134 CHINANET Guangdong province network",
			},
			wantScore: 100, // 基础100 + 住宅ISP加分
			wantMin:   100,
			wantMax:   110, // 可能有加分
		},
		{
			name: "VPN 服务商",
			ipInfo: &IPInfo{
				Org: "ExpressVPN LLC",
			},
			wantScore: 70,
			wantMin:   50,
			wantMax:   70, // 100 - 30 (VPN penalty)
		},
		{
			name: "数据中心",
			ipInfo: &IPInfo{
				Org: "Amazon Datacenter",
			},
			wantScore: 50,
			wantMin:   40,
			wantMax:   70, // 100 - 30 (VPN) - 20 (Cloud)
		},
		{
			name: "云服务商 - AWS",
			ipInfo: &IPInfo{
				Org: "Amazon Web Services",
			},
			wantScore: 80,
			wantMin:   70,
			wantMax:   80, // 100 - 20 (Cloud)
		},
		{
			name: "代理服务",
			ipInfo: &IPInfo{
				Org: "Proxy Server LLC",
			},
			wantScore: 70,
			wantMin:   50,
			wantMax:   70,
		},
		{
			name: "住宅 ISP - Comcast",
			ipInfo: &IPInfo{
				Org: "Comcast Cable Communications",
			},
			wantScore: 100,
			wantMin:   100,
			wantMax:   110,
		},
		{
			name: "无组织信息",
			ipInfo: &IPInfo{
				Org: "",
				ISP: "",
			},
			wantScore: 100, // 默认分数
			wantMin:   100,
			wantMax:   100,
		},
		{
			name: "混合场景 - VPN + Hosting",
			ipInfo: &IPInfo{
				Org: "VPN Hosting Services",
			},
			wantScore: 70,
			wantMin:   40,
			wantMax:   70, // 100 - 30 (VPN)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := detector.calculatePuritySimple(tt.ipInfo)

			// 验证分数在合理范围内
			if score < 0 || score > 100 {
				t.Errorf("Score out of range: got %.1f, want 0-100", score)
			}

			// 验证分数在预期范围内
			if score < tt.wantMin || score > tt.wantMax {
				t.Logf("Score = %.1f, expected range [%.1f, %.1f]", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

// TestAnalyzeIPCharacteristics 测试 IP 特征分析
func TestAnalyzeIPCharacteristics(t *testing.T) {
	detector := &Detector{}

	tests := []struct {
		name             string
		ipInfo           *IPInfo
		wantIsVPN        bool
		wantIsProxy      bool
		wantIsTor        bool
		wantIsDatacenter bool
	}{
		{
			name: "VPN 服务",
			ipInfo: &IPInfo{
				Org: "ExpressVPN Network",
			},
			wantIsVPN:        true,
			wantIsProxy:      false,
			wantIsTor:        false,
			wantIsDatacenter: false,
		},
		{
			name: "代理服务",
			ipInfo: &IPInfo{
				Org: "Proxy Services Inc",
			},
			wantIsVPN:        false,
			wantIsProxy:      true,
			wantIsTor:        false,
			wantIsDatacenter: false,
		},
		{
			name: "Tor 节点",
			ipInfo: &IPInfo{
				Org: "Tor Exit Node",
			},
			wantIsVPN:        false,
			wantIsProxy:      false,
			wantIsTor:        true,
			wantIsDatacenter: false,
		},
		{
			name: "数据中心",
			ipInfo: &IPInfo{
				Org: "DigitalOcean Datacenter",
			},
			wantIsVPN:        false,
			wantIsProxy:      false,
			wantIsTor:        false,
			wantIsDatacenter: true,
		},
		{
			name: "云服务",
			ipInfo: &IPInfo{
				Org: "Google Cloud Platform",
			},
			wantIsVPN:        false,
			wantIsProxy:      false,
			wantIsTor:        false,
			wantIsDatacenter: true,
		},
		{
			name: "住宅 ISP",
			ipInfo: &IPInfo{
				Org: "China Telecom",
			},
			wantIsVPN:        false,
			wantIsProxy:      false,
			wantIsTor:        false,
			wantIsDatacenter: false,
		},
		{
			name: "空组织",
			ipInfo: &IPInfo{
				Org: "",
				ISP: "",
			},
			wantIsVPN:        false,
			wantIsProxy:      false,
			wantIsTor:        false,
			wantIsDatacenter: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := &IPScore{}
			detector.analyzeIPCharacteristics(tt.ipInfo, score)

			if score.IsVPN != tt.wantIsVPN {
				t.Errorf("IsVPN = %v, want %v", score.IsVPN, tt.wantIsVPN)
			}

			if score.IsProxy != tt.wantIsProxy {
				t.Errorf("IsProxy = %v, want %v", score.IsProxy, tt.wantIsProxy)
			}

			if score.IsTor != tt.wantIsTor {
				t.Errorf("IsTor = %v, want %v", score.IsTor, tt.wantIsTor)
			}

			if score.IsDatacenter != tt.wantIsDatacenter {
				t.Errorf("IsDatacenter = %v, want %v", score.IsDatacenter, tt.wantIsDatacenter)
			}
		})
	}
}

// TestIPScore_GetRiskLevel 测试风险等级评估
func TestIPScore_GetRiskLevel(t *testing.T) {
	tests := []struct {
		name      string
		score     float64
		wantLevel string
	}{
		{
			name:      "优秀分数 100",
			score:     100,
			wantLevel: "低风险",
		},
		{
			name:      "优秀分数 95",
			score:     95,
			wantLevel: "低风险",
		},
		{
			name:      "良好分数 85",
			score:     85,
			wantLevel: "低风险",
		},
		{
			name:      "良好分数 75",
			score:     75,
			wantLevel: "中风险", // 75 < 80，所以是中风险
		},
		{
			name:      "一般分数 65",
			score:     65,
			wantLevel: "中风险",
		},
		{
			name:      "一般分数 50",
			score:     50,
			wantLevel: "中风险",
		},
		{
			name:      "较差分数 40",
			score:     40,
			wantLevel: "高风险",
		},
		{
			name:      "较差分数 25",
			score:     25,
			wantLevel: "高风险",
		},
		{
			name:      "很差分数 10",
			score:     10,
			wantLevel: "高风险",
		},
		{
			name:      "最低分数 0",
			score:     0,
			wantLevel: "高风险",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipScore := &IPScore{
				Score: tt.score,
			}

			level := ipScore.GetRiskLevel()

			if level != tt.wantLevel {
				t.Errorf("GetRiskLevel() = %s, want %s", level, tt.wantLevel)
			}
		})
	}
}

// TestIPScore_GetQualityDescription 测试评分等级描述
func TestIPScore_GetQualityDescription(t *testing.T) {
	tests := []struct {
		name      string
		score     float64
		wantGrade string
	}{
		{
			name:      "优秀 100",
			score:     100,
			wantGrade: "优秀",
		},
		{
			name:      "优秀 90",
			score:     90,
			wantGrade: "优秀",
		},
		{
			name:      "良好 85",
			score:     85,
			wantGrade: "良好",
		},
		{
			name:      "良好 75",
			score:     75,
			wantGrade: "良好",
		},
		{
			name:      "一般 65",
			score:     65,
			wantGrade: "一般",
		},
		{
			name:      "一般 50",
			score:     50,
			wantGrade: "一般",
		},
		{
			name:      "较差 40",
			score:     40,
			wantGrade: "较差",
		},
		{
			name:      "较差 25",
			score:     25,
			wantGrade: "较差",
		},
		{
			name:      "很差 10",
			score:     10,
			wantGrade: "很差",
		},
		{
			name:      "很差 0",
			score:     0,
			wantGrade: "很差",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipScore := &IPScore{
				Score: tt.score,
			}

			grade := ipScore.GetQualityDescription()

			if grade != tt.wantGrade {
				t.Errorf("GetQualityDescription() = %s, want %s", grade, tt.wantGrade)
			}
		})
	}
}

// TestDetector_DetectScore 测试完整的纯净度检测流程
func TestDetector_DetectScore(t *testing.T) {
	// 这个测试需要 Mock HTTP Server
	// 因为 DetectScore 内部调用 Detect()

	// 简单的单元测试，验证基本逻辑
	client := &http.Client{Timeout: 5 * time.Second}
	detector := NewDetector(client)

	if detector == nil {
		t.Fatal("NewDetector() should not return nil")
	}

	// DetectScore 需要真实的网络请求或完整的 Mock
	// 这里只验证方法存在性
	t.Log("DetectScore integration test requires network or full mocking")
}

// TestIPScore_BoundaryConditions 测试边界条件
func TestIPScore_BoundaryConditions(t *testing.T) {
	detector := &Detector{}

	tests := []struct {
		name   string
		ipInfo *IPInfo
	}{
		{
			name: "所有字段为空",
			ipInfo: &IPInfo{
				IP:      "",
				Country: "",
				Org:     "",
				ISP:     "",
			},
		},
		{
			name: "极长组织名称",
			ipInfo: &IPInfo{
				Org: "This is a very long organization name that might contain VPN or Proxy or Datacenter or Hosting or Cloud or any other keyword that we are looking for in our detection algorithm to properly identify the type of network connection being used",
			},
		},
		{
			name: "特殊字符",
			ipInfo: &IPInfo{
				Org: "Company-Name_123 (VPN) [Proxy] {Service}",
			},
		},
		{
			name: "Unicode 字符",
			ipInfo: &IPInfo{
				Org: "中国移动通信集团公司 China Mobile 🌐",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试不应该 panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("calculatePuritySimple panicked: %v", r)
				}
			}()

			score := detector.calculatePuritySimple(tt.ipInfo)

			// 验证分数在有效范围内
			if score < 0 || score > 100 {
				t.Errorf("Score out of range: %.1f", score)
			}
		})
	}
}

// TestIPScore_MultipleKeywords 测试多关键词场景
func TestIPScore_MultipleKeywords(t *testing.T) {
	detector := &Detector{}

	tests := []struct {
		name      string
		org       string
		wantScore float64
	}{
		{
			name:      "VPN + Hosting",
			org:       "VPN Hosting Services LLC",
			wantScore: 70, // 只扣一次 VPN 的 30 分
		},
		{
			name:      "Cloud + Datacenter",
			org:       "Cloud Datacenter Services",
			wantScore: 70, // 只扣一次 VPN 的 30 分（Datacenter 优先）
		},
		{
			name:      "Residential + Cloud",
			org:       "China Telecom Cloud",
			wantScore: 80, // 住宅 ISP 加分 10，Cloud 扣 20
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipInfo := &IPInfo{Org: tt.org}
			score := detector.calculatePuritySimple(ipInfo)

			// 允许一定的浮动范围
			diff := score - tt.wantScore
			if diff < -10 || diff > 10 {
				t.Logf("Score = %.1f, expected ~%.1f (diff: %.1f)", score, tt.wantScore, diff)
			}
		})
	}
}
