package middleware

import (
	"errors"
	"github.com/HaliComing/fpp/pkg/conf"
	"github.com/HaliComing/fpp/pkg/serializer"
	"github.com/gin-gonic/gin"
)

// TokenHandler Token处理
func TokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		// 成功跳过
		if token == conf.SystemConfig.Token {
			c.Next()
			return
		}

		err := errors.New("token is error")
		c.JSON(200, serializer.Err(
			serializer.CodeNoPermissionErr, err.Error(), err))
		c.Abort()
	}
}
