package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// InitHTTPClient 初始化 HTTP 客户端，支持代理配置
// 支持通过参数 -proxy 或环境变量 (HTTP_PROXY, HTTPS_PROXY, ALL_PROXY) 设置代理
func InitHTTPClient(proxyURL string, timeout time.Duration) (*http.Client, error) {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	}

	// 如果没有通过参数设置代理，则让 Go 的 http 库自动从环境变量读取代理
	// 当 transport.Proxy 为 nil 时，Go 会自动检查环境变量 HTTP_PROXY, HTTPS_PROXY, ALL_PROXY
	if proxyURL != "" {
		// 解析代理 URL
		parsedURL, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("无效的代理 URL: %v", err)
		}

		// 根据协议类型设置代理
		switch parsedURL.Scheme {
		case "http", "https":
			transport.Proxy = http.ProxyURL(parsedURL)
			fmt.Printf("✓ 已设置 HTTP/HTTPS 代理: %s\n", proxyURL)

		case "socks5":
			// SOCKS5 代理
			dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, nil, proxy.Direct)
			if err != nil {
				return nil, fmt.Errorf("SOCKS5 代理配置失败: %v", err)
			}
			transport.Dial = dialer.Dial
			fmt.Printf("✓ 已设置 SOCKS5 代理: %s\n", proxyURL)

		default:
			return nil, fmt.Errorf("不支持的代理协议: %s (支持 http, https, socks5)", parsedURL.Scheme)
		}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}, nil
}
