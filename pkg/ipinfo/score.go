package ipinfo

// IPScore IP 纯净度评分信息
type IPScore struct {
	IP           string  `json:"ip"`
	Score        float64 `json:"score"`        // 0-100 分数，越高越纯净
	IsVPN        bool    `json:"isVPN"`        // 是否是 VPN
	IsProxy      bool    `json:"isProxy"`      // 是否是代理
	IsTor        bool    `json:"isTor"`        // 是否是 Tor
	IsDatacenter bool    `json:"isDatacenter"` // 是否是数据中心 IP
	IsBlacklisted bool   `json:"isBlacklisted"` // 是否在黑名单
	RiskLevel    string  `json:"riskLevel"`    // 风险等级: low, medium, high
	ASN          string  `json:"asn"`          // AS 号
	ASNOrg       string  `json:"asnOrg"`       // AS 组织
}

// IPScoreProvider IP 纯净度检测提供商
type IPScoreProvider struct {
	Name string
	URL  string
}

// GetRiskLevel 根据分数获取风险等级
func (s *IPScore) GetRiskLevel() string {
	switch {
	case s.Score >= 80:
		return "低风险"
	case s.Score >= 50:
		return "中风险"
	default:
		return "高风险"
	}
}

// GetQualityDescription 获取质量描述
func (s *IPScore) GetQualityDescription() string {
	switch {
	case s.Score >= 90:
		return "优秀"
	case s.Score >= 75:
		return "良好"
	case s.Score >= 50:
		return "一般"
	case s.Score >= 25:
		return "较差"
	default:
		return "很差"
	}
}
