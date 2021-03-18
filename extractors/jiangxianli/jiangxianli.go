package jiangxianli

import (
	"github.com/HaliComing/fpp/extractors/types"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/request"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type extractor struct{}

// New returns a extractor.
func New() types.Extractor {
	return &extractor{}
}

func (e *extractor) Key() string {
	return "ip.jiangxianli.com"
}

// Extract is the main function to extract the data.
func (e *extractor) Extract() ([]*models.ProxyIP, error) {
	var proxyIPs []*models.ProxyIP
	client := request.HTTPClient{}
	html, err := client.Request("GET", "https://ip.jiangxianli.com/?page=1&country=%E4%B8%AD%E5%9B%BD", nil).GetResponse()
	if err != nil {
		return nil, err
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	dom.Find(".ip-tables table tbody tr").Each(func(i int, s *goquery.Selection) {
		ip := strings.TrimSpace(s.Find("td:nth-child(1)").Text())
		port := strings.TrimSpace(s.Find("td:nth-child(2)").Text())
		proxyIP := &models.ProxyIP{
			IP:   ip,
			Port: port,
		}
		proxyIPs = append(proxyIPs, proxyIP)
	})
	return proxyIPs, nil
}
