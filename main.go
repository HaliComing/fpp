package main

import (
	"flag"
	"github.com/HaliComing/fpp/bootstrap"
	"github.com/HaliComing/fpp/extractors"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/conf"
	"github.com/HaliComing/fpp/pkg/util"
	"github.com/HaliComing/fpp/routers"
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
	wg.Add(2)

	// 开启API接口
	go func() {
		defer wg.Done()
		RunApi()
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
