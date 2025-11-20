package ipinfo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Detector IP 检测器
type Detector struct {
	client    *http.Client
	providers []Provider
}

// NewDetector 创建新的 IP 检测器
func NewDetector(client *http.Client) *Detector {
	return &Detector{
		client:    client,
		providers: DefaultProviders,
	}
}

// Detect 检测 IP 信息（带故障转移）
func (d *Detector) Detect() (*IPInfo, error) {
	var lastErr error
	for _, provider := range d.providers {
		ipInfo, err := d.fetchFromProvider(provider)
		if err == nil {
			return ipInfo, nil
		}
		lastErr = err
		fmt.Printf("⚠️  %s 失败，尝试下一个...\n", provider.Name)
	}

	return nil, fmt.Errorf("所有 IP API 都失败: %v", lastErr)
}

// DetectScore 检测 IP 纯净度（简化版，使用已有的 IP 信息）
func (d *Detector) DetectScore() (*IPScore, error) {
	// 先获取基本 IP 信息
	info, err := d.Detect()
	if err != nil {
		return nil, fmt.Errorf("获取 IP 信息失败: %v", err)
	}

	score := &IPScore{
		IP:     info.IP,
		ASN:    "", // 从基础 API 可能获取不到
		ASNOrg: info.Org,
	}

	// 计算纯净度分数
	score.Score = d.calculatePuritySimple(info)
	score.RiskLevel = score.GetRiskLevel()

	// 基于 ISP/Org 名称判断特征
	d.analyzeIPCharacteristics(info, score)

	return score, nil
}

// calculatePuritySimple 简化的纯净度计算
func (d *Detector) calculatePuritySimple(info *IPInfo) float64 {
	score := 100.0

	// 检查 ISP/组织名称中的关键词
	org := info.Org
	if org == "" {
		org = info.ISP
	}

	if org != "" {
		// VPN/代理关键词检测
		vpnKeywords := []string{"VPN", "Proxy", "Datacenter", "Hosting", "Cloud", "Virtual", "Server"}
		for _, keyword := range vpnKeywords {
			if containsIgnoreCase(org, keyword) {
				score -= 30
				break
			}
		}

		// 云服务商检测
		cloudProviders := []string{"Amazon", "Google Cloud", "Microsoft Azure", "DigitalOcean",
			"Linode", "Vultr", "Hetzner", "OVH", "Cloudflare"}
		for _, provider := range cloudProviders {
			if containsIgnoreCase(org, provider) {
				score -= 20
				break
			}
		}

		// 住宅 ISP 加分
		residentialISPs := []string{"Comcast", "AT&T", "Verizon", "China Telecom", "China Unicom",
			"China Mobile", "Chinanet", "Telekom", "Orange", "Vodafone", "BT"}
		for _, isp := range residentialISPs {
			if containsIgnoreCase(org, isp) {
				score += 10
				break
			}
		}
	}

	// 确保分数在 0-100 范围内
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// analyzeIPCharacteristics 分析 IP 特征
func (d *Detector) analyzeIPCharacteristics(info *IPInfo, score *IPScore) {
	org := info.Org
	if org == "" {
		org = info.ISP
	}

	// 检测各种特征
	score.IsVPN = containsIgnoreCase(org, "VPN")
	score.IsProxy = containsIgnoreCase(org, "Proxy")
	score.IsTor = containsIgnoreCase(org, "Tor")
	score.IsDatacenter = containsIgnoreCase(org, "Datacenter") ||
		containsIgnoreCase(org, "Hosting") ||
		containsIgnoreCase(org, "Cloud")
}

func containsIgnoreCase(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && contains(s, substr))
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// fetchFromProvider 从指定 API 获取 IP 信息
func (d *Detector) fetchFromProvider(provider Provider) (*IPInfo, error) {
	resp, err := d.client.Get(provider.URL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("服务器返回错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 根据格式类型解析
	switch provider.Format {
	case "text":
		return d.parseTextFormat(body)
	case "json":
		return d.parseJSONFormat(body)
	default:
		return nil, fmt.Errorf("不支持的格式: %s", provider.Format)
	}
}

// parseJSONFormat 解析 JSON 格式的响应
func (d *Detector) parseJSONFormat(body []byte) (*IPInfo, error) {
	var ipInfo IPInfo
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}
	return &ipInfo, nil
}

// parseTextFormat 解析 ping0.cc/geo 的文本格式
// 格式示例:
// 14.153.68.158                          // IP 地址
// 中国 广东省深圳市福田中国电信               // 地区信息
// AS4134                                  // ASN 编号
// CHINANET Guangdong province network     // 网络运营商完整名称
func (d *Detector) parseTextFormat(body []byte) (*IPInfo, error) {
	scanner := bufio.NewScanner(bytes.NewReader(body))
	lines := []string{}

	// 读取所有非空行
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) < 2 {
		return nil, fmt.Errorf("文本格式不正确，行数不足")
	}

	ipInfo := &IPInfo{
		IP: lines[0], // 第一行: IP 地址
	}

	// 第二行: 地区信息，格式如 "中国 广东省深圳市福田中国电信"
	if len(lines) >= 2 {
		d.parseLocationLine(lines[1], ipInfo)
	}

	// 第三行: ASN 编号（如果存在），格式如 "AS4134"
	if len(lines) >= 3 && strings.HasPrefix(lines[2], "AS") {
		ipInfo.Org = lines[2] // 暂存 ASN 到 Org 字段
	}

	// 第四行: 网络运营商完整名称（如果存在），格式如 "CHINANET Guangdong province network"
	if len(lines) >= 4 {
		// 如果有 ASN，将其与运营商名称组合
		if ipInfo.Org != "" && strings.HasPrefix(ipInfo.Org, "AS") {
			ipInfo.Org = ipInfo.Org + " " + lines[3]
		} else {
			ipInfo.Org = lines[3]
		}
	}

	return ipInfo, nil
}

// parseLocationLine 解析地区信息行
// 格式: "中国 广东省深圳市福田中国电信"
// 或: "United States California Los Angeles AT&T"
func (d *Detector) parseLocationLine(line string, ipInfo *IPInfo) {
	// 原始完整字符串
	remaining := line

	// 1. 提取国家（第一个空格之前，或包含"省"/"市"之前的部分）
	spaceIdx := strings.Index(remaining, " ")
	if spaceIdx > 0 {
		// 检查第一个词是否包含省/市（说明是单字国家名）
		firstWord := remaining[:spaceIdx]
		if !strings.Contains(firstWord, "省") && !strings.Contains(firstWord, "市") {
			ipInfo.Country = firstWord
			remaining = strings.TrimSpace(remaining[spaceIdx+1:])
		}
	}

	// 如果没有空格，整行就是国家
	if ipInfo.Country == "" {
		ipInfo.Country = remaining
		return
	}

	// 2. 提取省份（包含"省"的部分）
	if idx := strings.Index(remaining, "省"); idx > 0 {
		ipInfo.Region = remaining[:idx+len("省")]
		remaining = remaining[idx+len("省"):]
	}

	// 3. 提取城市（包含"市"的部分）
	if idx := strings.Index(remaining, "市"); idx > 0 {
		ipInfo.City = remaining[:idx+len("市")]
		remaining = remaining[idx+len("市"):]
	}

	// 4. 剩余部分作为 ISP/运营商
	remaining = strings.TrimSpace(remaining)
	if remaining != "" {
		ipInfo.ISP = remaining
		ipInfo.Org = remaining
	}

	// 5. 如果没有提取到省/市（英文格式），则按空格分割
	if ipInfo.Region == "" && ipInfo.City == "" {
		// 重新解析，按空格分割
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			// parts[0] 已经是 Country
			if len(parts) >= 2 {
				ipInfo.Region = parts[1]
			}
			if len(parts) >= 3 {
				ipInfo.City = parts[2]
			}
			if len(parts) >= 4 {
				ipInfo.ISP = strings.Join(parts[3:], " ")
				ipInfo.Org = ipInfo.ISP
			}
		}
	}
}
