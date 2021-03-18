package xiladaili

import (
	"github.com/HaliComing/fpp/extractors/types"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/request"
	"net/http"
	"regexp"
)

type extractor struct{}

// New returns a extractor.
func New() types.Extractor {
	return &extractor{}
}

func (e *extractor) Key() string {
	return "www.xiladaili.com"
}

// Extract is the main function to extract the data.
func (e *extractor) Extract() ([]*models.ProxyIP, error) {
	var proxyIPs []*models.ProxyIP
	client := request.HTTPClient{}
	html, err := client.Request("GET", "http://www.xiladaili.com/", nil,
		request.WithHeader(http.Header{
			"Referer":    {"http://www.xiladaili.com/"},
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
