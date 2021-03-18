package controllers

import (
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/conf"
	"github.com/HaliComing/fpp/pkg/serializer"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Code: 0,
		Data: conf.Version,
	})
}

// 统计数据
func Count(c *gin.Context) {
	count, err := models.ProxyCount()
	if err != nil {
		c.JSON(200, serializer.Err(serializer.CodeNotSet, err.Error(), err))
		return
	}
	c.JSON(200, serializer.Response{
		Code: 0,
		Data: map[string]interface{}{
			"Count": count,
		},
	})
}
