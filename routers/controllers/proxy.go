package controllers

import (
	"errors"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/serializer"
	"github.com/HaliComing/fpp/pkg/util"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// 随机获取一个IP
func ProxyRandom(c *gin.Context) {
	anonymous, _ := strconv.Atoi(c.Query("anonymous"))
	protocol, _ := strconv.Atoi(c.Query("protocol"))
	country := c.Query("country")
	random, err := models.ProxyRandom(
		util.IfIntArray(anonymous != 0, []int{anonymous}, nil),
		util.IfIntArray(protocol != 0, []int{protocol}, nil),
		util.IfStringArray(country != "", []string{country}, nil))
	if err != nil {
		c.JSON(200, serializer.Err(serializer.CodeNotSet, err.Error(), err))
		return
	}
	c.JSON(200, serializer.Response{
		Code: 0,
		Data: random,
	})
}

// 获取全部IP
func ProxyAll(c *gin.Context) {
	anonymous, _ := strconv.Atoi(c.Query("anonymous"))
	protocol, _ := strconv.Atoi(c.Query("protocol"))
	country := c.Query("country")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	all, err := models.ProxyAll(
		util.IfIntArray(anonymous != 0, []int{anonymous}, nil),
		util.IfIntArray(protocol != 0, []int{protocol}, nil),
		util.IfStringArray(country != "", []string{country}, nil), page, pageSize)
	if err != nil {
		c.JSON(200, serializer.Err(serializer.CodeNotSet, err.Error(), err))
		return
	}
	c.JSON(200, serializer.Response{
		Code: 0,
		Data: all,
	})
}

// 删除IP
func ProxyDelete(c *gin.Context) {
	ip := c.Query("ip")
	port := c.Query("port")
	proxy := c.Query("proxy")
	if strings.Index(proxy, ":") != -1 {
		split := strings.Split(proxy, ":")
		ip = split[0]
		port = split[1]
	}
	if ip == "" || port == "" {
		err := errors.New("param is error")
		c.JSON(200, serializer.Err(serializer.CodeNotSet, err.Error(), err))
		return
	}
	models.ProxyDelete(ip, port)
	c.JSON(200, serializer.Response{
		Code: 0,
	})
}
