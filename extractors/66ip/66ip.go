package ip66

import (
	"github.com/HaliComing/fpp/extractors/types"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/request"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strings"
)

type extractor struct{}

// New returns a extractor.
func New() types.Extractor {
	return &extractor{}
}

func (e *extractor) Key() string {
	return "www.66ip.cn"
}

// Extract is the main function to extract the data.
func (e *extractor) Extract() ([]*models.ProxyIP, error) {
	var proxyIPs []*models.ProxyIP
	ips, err := ip66()
	if err != nil {
		return nil, err
	}
	proxyIPs = append(proxyIPs, ips...)
	ips2, err := ip66mo()
	if err != nil {
		return nil, err
	}
	proxyIPs = append(proxyIPs, ips2...)
	return proxyIPs, nil
}

func ip66() ([]*models.ProxyIP, error) {
	var proxyIPs []*models.ProxyIP
	client := request.HTTPClient{}
	html, err := client.Request("GET", "http://www.66ip.cn/", nil,
		request.WithHeader(http.Header{
			"Referer":    {"http://www.66ip.cn/"},
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36"},
		})).GetResponse()
	if err != nil {
		return nil, err
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	dom.Find("#main > div table > tbody > tr:not(:first-child)").Each(func(i int, s *goquery.Selection) {
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

func ip66mo() ([]*models.ProxyIP, error) {
	var proxyIPs []*models.ProxyIP
	client := request.HTTPClient{}
	html, err := client.Request("GET", "http://www.66ip.cn/mo.php", nil,
		request.WithHeader(http.Header{
			"Referer":    {"http://www.66ip.cn/"},
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36"},
		})).GetResponse()
	if err != nil {
		return nil, err
	}
	reg, err := regexp.Compile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d{1,5})`)
	if err != nil {
		return nil, err
	}
	matchOk := reg.MatchString(html)
	if !matchOk {
		return proxyIPs, nil
	}
	submatchs := reg.FindAllStringSubmatch(html, -1)
	for _, submatch := range submatchs {
		proxyIP := &models.ProxyIP{
			IP:   submatch[1],
			Port: submatch[2],
		}
		proxyIPs = append(proxyIPs, proxyIP)
	}
	return proxyIPs, nil
}
