package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/HaliComing/fpp/bootstrap"
	"github.com/HaliComing/fpp/extractors"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/conf"
	"github.com/HaliComing/fpp/pkg/util"
	"github.com/HaliComing/fpp/routers"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "c", "conf.ini", "配置文件路径")
	flag.Parse()
	bootstrap.Init(confPath)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(3)

	// 开启API接口
	go func() {
		defer wg.Done()
		RunApi()
	}()

	// 开启Proxy接口
	go func() {
		defer wg.Done()
		RunProxy()
	}()

	proxyIPChan := make(chan *models.ProxyIP, 2000)

	// IP提取器
	ticker := time.NewTicker(time.Duration(conf.SystemConfig.ExtractionInterval) * time.Minute)
	go func(ticker *time.Ticker) {
		defer wg.Done()
		for {
			RunExtractors(proxyIPChan)
			<-ticker.C // 等待时间到达
		}
	}(ticker)

	// IP消化器 启动n个go程
	for i := 0; i < conf.SystemConfig.NumberOfThreads; i++ {
		go RunCheck(proxyIPChan)
	}
	wg.Wait()
}

func RunApi() {
	api := routers.InitRouter()

	// 如果启用了SSL
	if conf.SSLConfig.CertPath != "" {
		go func() {
			util.Log().Info("[API] SSL Listen = %s", conf.SSLConfig.Listen)
			if err := api.RunTLS(conf.SSLConfig.Listen,
				conf.SSLConfig.CertPath, conf.SSLConfig.KeyPath); err != nil {
				util.Log().Error("[API] SSL Listen = %s, Error = %s", conf.SSLConfig.Listen, err)
			}
		}()
	}

	util.Log().Info("[API] Listen = %s", conf.SystemConfig.Listen)
	if err := api.Run(conf.SystemConfig.Listen); err != nil {
		util.Log().Error("[API] Listen = %s, Error = %s", conf.SystemConfig.Listen, err)
	}
}

func RunProxy() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Logger = log.New(os.Stderr, "[proxy] ", log.LstdFlags)
	proxy.Verbose = conf.SystemConfig.Debug
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		println("DoFunc", req.URL.Scheme)
		switchProxy(proxy, []int{models.ProtocolTypeHTTP, models.ProtocolTypeHTTPS})
		return req, nil
	})
	proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		println("HandleConnectFunc")
		switchProxy(proxy, []int{models.ProtocolTypeHTTPS, models.ProtocolTypeALL})
		return goproxy.OkConnect, host
	})
	log.Fatal(http.ListenAndServe(conf.SystemConfig.ProxyListen, proxy))
}

func switchProxy(proxy *goproxy.ProxyHttpServer, protocols []int) {
	random, err := models.ProxyRandom(nil, protocols, nil)
	useProxy := true
	var proxyStr string
	var proxyURL *url.URL
	if err != nil {
		useProxy = false
	} else {
		proxyStr = fmt.Sprintf("http://%s:%s", random.IP, random.Port)
		proxyURL, err = url.Parse(proxyStr)
		if err != nil {
			useProxy = false
		}
	}
	if useProxy {
		proxy.Tr = &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		proxy.ConnectDial = proxy.NewConnectDialToProxy(proxyStr)
	} else {
		proxy.Tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		proxy.ConnectDial = nil
	}
}

func RunExtractors(proxyIPChan chan *models.ProxyIP) {
	proxyIPs := extractors.Extract()
	for _, proxyIP := range proxyIPs {
		proxyIPChan <- proxyIP
	}
}

func RunCheck(proxyIPChan chan *models.ProxyIP) {
	for {
		proxyIP := <-proxyIPChan
		if proxyIP.CheckIP() {
			util.Log().Info("[CheckIP] testIP = %s:%s ，Model = %+v", proxyIP.IP, proxyIP.Port, proxyIP)
			models.Save(proxyIP)
		} else {
			models.Delete(proxyIP)
		}
		if proxyIP.Score > 0 {
			timer := time.NewTimer(time.Duration(conf.SystemConfig.CheckInterval) * time.Minute)
			go func(timer *time.Timer, proxyIP *models.ProxyIP, proxyIPChan chan *models.ProxyIP) {
				<-timer.C
				proxyIPChan <- proxyIP
			}(timer, proxyIP, proxyIPChan)
		}
	}
}
