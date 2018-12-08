package controllers

/*
 * @Script: index.go
 * @Author: pangxiaobo
 * @Email: 10846295@qq.com
 * @Create At: 2018-11-08 20:07:59
 * @Last Modified By: pangxiaobo
 * @Last Modified At: 2018-11-09 11:46:02
 * @Description: This is description.
 */

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type IndexController struct{}

func (i *IndexController) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"msg":       "Welcome.",
		"data":      nil,
		"timestamp": time.Now().Unix(),
	})

	return
}

func (i *IndexController) Handle404(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"code":      404,
		"msg":       "Page is Not Found.",
		"data":      nil,
		"timestamp": time.Now().Unix(),
	})
	return
}
