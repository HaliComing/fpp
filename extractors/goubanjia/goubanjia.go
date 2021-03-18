package goubanjia

import (
	"github.com/HaliComing/fpp/extractors/types"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/request"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

type extractor struct{}

// New returns a extractor.
func New() types.Extractor {
	return &extractor{}
}

func (e *extractor) Key() string {
	return "www.goubanjia.com"
}

// Extract is the main function to extract the data.
func (e *extractor) Extract() ([]*models.ProxyIP, error) {
	var proxyIPs []*models.ProxyIP
	client := request.HTTPClient{}
	html, err := client.Request("GET", "http://www.goubanjia.com/", nil,
		request.WithHeader(http.Header{
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36"},
		})).GetResponse()
	if err != nil {
		return nil, err
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	dom.Find(".services div table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// 获取IP
		s.Find("td:nth-child(1)").ChildrenFiltered("p").Remove()
		ip := strings.TrimSpace(s.Find("td:nth-child(1)").Text())
		ipport := strings.Split(ip, ":")
		// 获取端口
		attr, exists := s.Find("td:nth-child(1)").Find(".port").Attr("class")
		port := ""
		if exists {
			portClass := strings.Split(attr, " ")
			for _, ch := range portClass[1] {
				port = port + strconv.Itoa(strings.Index("ABCDEFGHIZ", string(ch)))
			}
		}
		atoi, err2 := strconv.Atoi(port)
		if err2 != nil {
			return
		}
		proxyIP := &models.ProxyIP{
			IP:   ipport[0],
			Port: strconv.Itoa(atoi >> 3),
		}
		proxyIPs = append(proxyIPs, proxyIP)
	})
	return proxyIPs, nil
}
