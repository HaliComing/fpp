package types

import "github.com/HaliComing/fpp/models"

// 提取器
type Extractor interface {
	// 提取器名称 一般为域名
	Key() string
	// 提取器主要函数
	Extract() ([]*models.ProxyIP, error)
}
