package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HaliComing/fpp/pkg/request"
	"github.com/HaliComing/fpp/pkg/util"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var LocalIp = ""

// 代理协议类型
const (
	// HTTP代理
	ProtocolTypeHTTP int = 1
	// HTTPS代理
	ProtocolTypeHTTPS int = 2
	// HTTP/HTTPS代理
	ProtocolTypeALL int = 3
)

// 代理匿名类型
const (
	// 透明代理
	AnonymousTypeT int = 1
	// 匿名代理
	AnonymousTypeN int = 2
	// 混淆代理
	AnonymousTypeH int = 3
	// 高匿代理
	AnonymousTypeG int = 4
)

// Proxy 代理IP模型
type ProxyIP struct {
	IP          string `gorm:"primary_key"` // IP
	Port        string `gorm:"primary_key"` // Port端口
	Protocol    int    // 代理协议类型 HTTP:1 HTTPS:2 HTTP/HTTPS:3
	Anonymous   int    // 匿名类型 透明:1 普通匿名:2 欺骗匿名:3 高匿:4
	Country     string // 国家
	Province    string // 省
	Attribution string // 归属 全路径
	ISP         string // 运营商
	Score       int    // 分数 默认10分
	Source      string // 来源
	Speed       int64  // 连接速度 单位ms
	CreateTime  string // 创建时间
	LastTime    string // 最后检测时间
	FailedCount int    // 连续失败次数 扣分时采用 Score - (2 * FailedCount + 1)
}

func (proxyIP *ProxyIP) String() string {
	return "{IP:" + proxyIP.IP +
		" " + "Port:" + proxyIP.Port +
		" " + "Protocol:" + strconv.Itoa(proxyIP.Protocol) +
		" " + "Anonymous:" + strconv.Itoa(proxyIP.Anonymous) +
		" " + "Speed:" + strconv.FormatInt(proxyIP.Speed, 10) +
		" " + "FailedCount:" + strconv.Itoa(proxyIP.FailedCount) +
		" " + "Country:" + proxyIP.Country +
		" " + "Source:" + proxyIP.Source +
		"}"
}

const (
	PrefixProtocolHttp  = "http://"
	PrefixProtocolHttps = "https://"
	TestUrlHttp         = "http://httpbin.org/get?show_env=1"
	TestUrlHttps        = "https://httpbin.org/get?show_env=1"
)

func (proxyIP *ProxyIP) CheckIP() bool {
	testIp := fmt.Sprintf("%s:%s", proxyIP.IP, proxyIP.Port)
	begin1 := time.Now()
	http, originHttp, errHttp := request.RequestProxy(TestUrlHttp, PrefixProtocolHttp+testIp)
	speedHttp := time.Now().Sub(begin1).Nanoseconds() / 1000 / 1000 //ms
	begin2 := time.Now()
	https, originHttps, errHttps := request.RequestProxy(TestUrlHttps, PrefixProtocolHttps+testIp)
	speedHttps := time.Now().Sub(begin2).Nanoseconds() / 1000 / 1000 //ms
	if http && https {
		proxyIP.Protocol = ProtocolTypeALL
		proxyIP.Speed = (speedHttp + speedHttps) / 2
		proxyIP.Anonymous = proxyIP.GetAnonymous(originHttps)
	} else if http {
		proxyIP.Protocol = ProtocolTypeHTTP
		proxyIP.Speed = speedHttp
		proxyIP.Anonymous = proxyIP.GetAnonymous(originHttp)
	} else if https {
		proxyIP.Protocol = ProtocolTypeHTTPS
		proxyIP.Speed = speedHttps
		proxyIP.Anonymous = proxyIP.GetAnonymous(originHttps)
	}
	proxyIP.LastTime = time.Now().Format("2006-01-02 15:04:05")
	if http || https {
		proxyIP.FailedCount = 0
		proxyIP.GetIpInfo()
		return true
	}
	proxyIP.Score = proxyIP.Score - (2*proxyIP.FailedCount + 1)
	proxyIP.FailedCount = proxyIP.FailedCount + 1
	if errHttp != nil || errHttps != nil {
		util.Log().Warning("[CheckIP] testIP = %s, ErrorHttp = %s, ErrorHttps = %s", testIp, errHttp, errHttps)
	}
	return false
}

type pos struct {
	Ct     string `json:"ct"`
	Prov   string `json:"prov"`
	City   string `json:"city"`
	Area   string `json:"area"`
	Idc    string `json:"idc"`
	Yunyin string `json:"yunyin"`
	Net    string `json:"net"`
}

// TODO ip_c_list数组值判断
func (proxyIP *ProxyIP) GetIpInfo() {
	html := getHtml("https://www.ip138.com/iplookup.asp?ip=" + proxyIP.IP + "&action=2")
	if html == "" {
		return
	}
	html = util.ConvertToString(html, "GBK", "UTF-8")
	jsonStr := between(html, "\"ip_c_list\":[", "]")
	if jsonStr == "" {
		return
	}
	jsonStr = "[" + jsonStr + "]"
	var pos []pos
	err := json.Unmarshal([]byte(jsonStr), &pos)
	if err == nil {
		proxyIP.Country = strings.TrimSpace(pos[len(pos)-1].Ct)
		proxyIP.Province = strings.TrimSpace(pos[len(pos)-1].Prov)
		proxyIP.Attribution = strings.TrimSpace(pos[len(pos)-1].City + " " + pos[len(pos)-1].Area + " " + pos[len(pos)-1].Idc)
		proxyIP.ISP = strings.TrimSpace(pos[len(pos)-1].Yunyin + " " + pos[len(pos)-1].Net)
	}
}

func between(str, starting, ending string) string {
	s := strings.Index(str, starting)
	if s < 0 {
		return ""
	}
	s += len(starting)
	e := strings.Index(str[s:], ending)
	if e < 0 {
		return ""
	}
	return str[s : s+e]
}

func getHtml(url string) string {
	client := request.HTTPClient{}
	res, err := client.Request("GET", url, nil,
		request.WithHeader(http.Header{
			"Referer":    {"https://www.ip138.com/"},
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36"},
		})).GetResponse()
	if err != nil {
		util.Log().Warning("[getHtml] Error = %s", err)
		return ""
	}
	return res
}

func (proxyIP *ProxyIP) GetAnonymous(origin string) int {
	if LocalIp == "" {
		return 0
	}
	if strings.Index(origin, ",") != -1 {
		if strings.Index(origin, LocalIp) != -1 {
			return AnonymousTypeN
		} else {
			return AnonymousTypeH
		}
	} else {
		if origin == proxyIP.IP {
			return AnonymousTypeG
		} else if origin == LocalIp {
			return AnonymousTypeT
		} else {
			return AnonymousTypeH
		}
	}
}

func Exist(proxyIP *ProxyIP) bool {
	var count int64
	DB.Model(&ProxyIP{}).Where("IP = ?", proxyIP.IP).Where("Port = ?", proxyIP.Port).Count(&count)
	if count == 0 {
		return false
	} else {
		return true
	}
}

func Save(proxyIP *ProxyIP) {
	DB.Save(proxyIP)
}

func Delete(proxyIP *ProxyIP) {
	DB.Delete(proxyIP)
}

func ProxyDelete(ip, port string) {
	proxyIP := &ProxyIP{
		IP:   ip,
		Port: port,
	}
	DB.Delete(proxyIP)
}

func ProxyCount() (int64, error) {
	var count int64
	result := DB.Model(&ProxyIP{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func ProxyRandomCount(anonymous, protocols []int, countries []string) (int64, error) {
	var count int64
	db := DB.Model(&ProxyIP{})
	if anonymous != nil && len(anonymous) != 0 {
		db = db.Where("Anonymous in (?)", anonymous)
	}
	if protocols != nil && len(protocols) != 0 {
		db = db.Where("Protocol in (?)", protocols)
	}
	if countries != nil && len(countries) != 0 {
		db = db.Where("Country in (?)", countries)
	}
	result := db.Count(&count)
	if result.Error != nil {
		return count, result.Error
	}
	if count == 0 {
		return count, errors.New("proxy ip number is 0")
	}
	return count, nil
}

func ProxyRandom(anonymous, protocols []int, countries []string) (ProxyIP, error) {
	count, err := ProxyRandomCount(anonymous, protocols, countries)
	if err != nil {
		return ProxyIP{}, err
	}
	var proxyIP ProxyIP
	db := DB
	if anonymous != nil && len(anonymous) != 0 {
		db = db.Where("Anonymous in (?)", anonymous)
	}
	if protocols != nil && len(protocols) != 0 {
		db = db.Where("Protocol in (?)", protocols)
	}
	if countries != nil && len(countries) != 0 {
		db = db.Where("Country in (?)", countries)
	}
	rand.Seed(time.Now().UnixNano())
	randInt64 := rand.Int63n(count)
	result := db.Limit(1).Offset(randInt64).Find(&proxyIP)
	if result.Error != nil {
		return ProxyIP{}, result.Error
	}
	return proxyIP, nil
}

func ProxyAll(anonymous, protocols []int, countries []string, page, pageSize int) ([]ProxyIP, error) {
	var proxyIPs []ProxyIP
	db := DB
	if anonymous != nil && len(anonymous) != 0 {
		db = db.Where("Anonymous in (?)", anonymous)
	}
	if protocols != nil && len(protocols) != 0 {
		db = db.Where("Protocol in (?)", protocols)
	}
	if countries != nil && len(countries) != 0 {
		db = db.Where("Country in (?)", countries)
	}
	result := db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&proxyIPs)
	if result.Error != nil {
		return nil, result.Error
	}
	return proxyIPs, nil
}
