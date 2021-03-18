package routers

import (
	"github.com/HaliComing/fpp/middleware"
	"github.com/HaliComing/fpp/pkg/util"
	"github.com/HaliComing/fpp/routers/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// InitCORS 初始化跨域配置
func InitCORS(router *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	router.Use(cors.New(config))
}

// InitRouter 初始化主机模式路由
func InitRouter() *gin.Engine {
	util.Log().Info("[API] Start API, Init Router.")

	r := gin.Default()

	/*
		静态资源
	*/
	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/api/"})))
	r.Use(middleware.FrontendFileHandler())

	v1 := r.Group("/api/v1")

	// 跨域相关
	InitCORS(r)

	v1.Use(middleware.TokenHandler())
	/*
		路由
	*/
	{
		// 全局相关
		site := v1.Group("site")
		{
			// 测试用路由
			site.GET("ping", controllers.Ping)
			// 统计数据
			site.GET("count", controllers.Count)
		}

		// IP代理池
		proxy := v1.Group("proxy")
		{
			// 随机获取一个IP
			proxy.GET("random", controllers.ProxyRandom)
			// 获取全部IP
			proxy.GET("all", controllers.ProxyAll)
			// 删除IP
			proxy.GET("delete", controllers.ProxyDelete)
		}
	}
	return r
}
