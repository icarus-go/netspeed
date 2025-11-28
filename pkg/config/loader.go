package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/icarus-go/netspeed/pkg/tester"
)

// Loader 配置加载器
type Loader struct{}

// NewLoader 创建新的配置加载器
func NewLoader() *Loader {
	return &Loader{}
}

// LoadSites 加载测试站点
func (l *Loader) LoadSites(configFile string) ([]tester.Site, error) {
	if configFile == "" {
		// 返回默认站点
		return tester.DefaultSites, nil
	}

	// 从配置文件加载
	return l.loadSitesFromFile(configFile)
}

// loadSitesFromFile 从 JSON 配置文件加载站点
func (l *Loader) loadSitesFromFile(filename string) ([]tester.Site, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var sites []tester.Site
	if err := json.Unmarshal(data, &sites); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	if len(sites) == 0 {
		return nil, fmt.Errorf("配置文件中没有站点")
	}

	return sites, nil
}
