package extractors

import (
	ip66 "github.com/HaliComing/fpp/extractors/66ip"
	yip7 "github.com/HaliComing/fpp/extractors/7yip"
	ip89 "github.com/HaliComing/fpp/extractors/89ip"
	"github.com/HaliComing/fpp/extractors/goubanjia"
	"github.com/HaliComing/fpp/extractors/ip3366"
	"github.com/HaliComing/fpp/extractors/jiangxianli"
	"github.com/HaliComing/fpp/extractors/kuaidaili"
	"github.com/HaliComing/fpp/extractors/proxylistplus"
	"github.com/HaliComing/fpp/extractors/seofangfa"
	"github.com/HaliComing/fpp/extractors/sudaili"
	"github.com/HaliComing/fpp/extractors/types"
	"github.com/HaliComing/fpp/extractors/xiladaili"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/request"
	"github.com/HaliComing/fpp/pkg/util"
	"time"
)

var extractors []types.Extractor

func init() {
	extractors = []types.Extractor{
		yip7.New(),
		ip66.New(),
		ip89.New(),
		goubanjia.New(),
		ip3366.New(),
		jiangxianli.New(),
		kuaidaili.New(),
		proxylistplus.New(),
		seofangfa.New(),
		sudaili.New(),
		xiladaili.New(),
	}
}

// 提取器主要函数
func Extract() []*models.ProxyIP {
	ip, err := request.RequestIP(models.TestUrlHttps)
	if err == nil {
		models.LocalIp = ip
	}
	var result []*models.ProxyIP
	util.Log().Info("[Extractor] ExtractorNumber = %d", len(extractors))
	for _, extractor := range extractors {
		proxyIPs, err := extractor.Extract()
		if err != nil {
			util.Log().Error("[Extractor] Extractor = %s, Error = %s", extractor.Key(), err)
			continue
		}
		if proxyIPs == nil || len(proxyIPs) == 0 {
			util.Log().Info("[Extractor] Extractor = %s, IPNumber = 0", extractor.Key())
			continue
		}
		for _, proxyIP := range proxyIPs {
			if proxyIP == nil || proxyIP.IP == "" || proxyIP.Port == "" {
				continue
			}
			proxyIP.Score = 10
			proxyIP.Source = extractor.Key()
			proxyIP.CreateTime = time.Now().Format("2006-01-02 15:04:05")
			proxyIP.FailedCount = 0
			result = append(result, proxyIP)
			util.Log().Info("[Extractor] Extractor = %s, Proxy = %s:%s", extractor.Key(), proxyIP.IP, proxyIP.Port)
		}
	}
	util.Log().Info("[Extractor] Extractor = ALL, IPNumber = %d", len(result))
	return result
}
