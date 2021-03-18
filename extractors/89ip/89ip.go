package ip89

import (
	"github.com/HaliComing/fpp/extractors/types"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/request"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type extractor struct{}

// New returns a extractor.
func New() types.Extractor {
	return &extractor{}
}

func (e *extractor) Key() string {
	return "www.89ip.cn"
}

// Extract is the main function to extract the data.
func (e *extractor) Extract() ([]*models.ProxyIP, error) {
	var proxyIPs []*models.ProxyIP
	client := request.HTTPClient{}
	html, err := client.Request("GET", "https://www.89ip.cn/index_1.html", nil,
		request.WithHeader(http.Header{
			"Referer":    {"https://www.89ip.cn/"},
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36"},
		})).GetResponse()
	if err != nil {
		return nil, err
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	dom.Find("div div.layui-form > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
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
