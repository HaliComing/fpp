package proxylistplus

import (
	"github.com/HaliComing/fpp/extractors/types"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/request"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)

type extractor struct{}

// New returns a extractor.
func New() types.Extractor {
	return &extractor{}
}

func (e *extractor) Key() string {
	return "list.proxylistplus.com"
}

// Extract is the main function to extract the data.
func (e *extractor) Extract() ([]*models.ProxyIP, error) {
	var proxyIPs []*models.ProxyIP
	client := request.HTTPClient{}
	html, err := client.Request("GET", "https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-1", nil,
		request.WithHeader(http.Header{
			"Referer":    {"https://list.proxylistplus.com/update-2"},
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36"},
		}), request.WithTimeout(time.Duration(10)*time.Second)).GetResponse()
	if err != nil {
		return nil, err
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	dom.Find("#page table.bg tbody tr.cells").Each(func(i int, s *goquery.Selection) {
		ip := strings.TrimSpace(s.Find("td:nth-child(2)").Text())
		port := strings.TrimSpace(s.Find("td:nth-child(3)").Text())
		proxyIP := &models.ProxyIP{
			IP:   ip,
			Port: port,
		}
		proxyIPs = append(proxyIPs, proxyIP)
	})
	return proxyIPs, nil
}
