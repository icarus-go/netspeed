package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
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

	// 检查是否通过参数设置了代理
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
	} else {
		// 检查环境变量代理设置
		httpProxy := getEnvProxy("HTTP_PROXY", "http_proxy")
		httpsProxy := getEnvProxy("HTTPS_PROXY", "https_proxy")
		allProxy := getEnvProxy("ALL_PROXY", "all_proxy")

		if httpProxy != "" || httpsProxy != "" || allProxy != "" {
			// 使用环境变量代理，让 Go 的 http 库自动读取
			// transport.Proxy 为 nil 时，Go 会自动检查环境变量 HTTP_PROXY, HTTPS_PROXY, ALL_PROXY
			fmt.Println("✓ 已自动检测到环境变量代理配置:")
			if httpProxy != "" {
				fmt.Printf("  - HTTP_PROXY: %s\n", httpProxy)
			}
			if httpsProxy != "" {
				fmt.Printf("  - HTTPS_PROXY: %s\n", httpsProxy)
			}
			if allProxy != "" {
				fmt.Printf("  - ALL_PROXY: %s\n", allProxy)
			}
			transport.Proxy = http.ProxyFromEnvironment
		}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}, nil
}

// getEnvProxy 获取环境变量代理，优先检查大写，再检查小写
func getEnvProxy(upper, lower string) string {
	if val := os.Getenv(upper); val != "" {
		return val
	}
	return os.Getenv(lower)
}
