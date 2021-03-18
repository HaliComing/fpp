package bootstrap

import (
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/conf"
	"github.com/gin-gonic/gin"
)

// Init 初始化启动
func Init(path string) {
	InitApplication()
	conf.Init(path)
	// Debug 关闭时，切换为生产模式
	if !conf.SystemConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	models.Init()
	InitStatic()
}
