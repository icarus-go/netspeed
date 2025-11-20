package ipinfo

// IPInfo IP 地理信息
type IPInfo struct {
	IP          string `json:"ip"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Region      string `json:"region"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	Timezone    string `json:"timezone"`
	Org         string `json:"org"`
}

// Provider IP API 提供商
type Provider struct {
	Name   string
	URL    string
	Format string // "json" or "text"
}

// DefaultProviders 默认 API 提供商列表
var DefaultProviders = []Provider{
	{"ping0.cc", "https://ping0.cc/geo", "text"},         // 优先使用，响应快
	{"ipapi.co", "https://ipapi.co/json/", "json"},
	{"ipinfo.io", "https://ipinfo.io/json", "json"},
	{"ip-api.com", "http://ip-api.com/json/", "json"},
}
