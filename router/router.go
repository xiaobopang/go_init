package router

/*
 * @Script: routers.go
 * @Author: pangxiaobo
 * @Email: 10846295@qq.com
 * @Create At: 2018-11-27 18:19:27
 * @Last Modified By: pangxiaobo
 * @Last Modified At: 2018-12-12 14:25:18
 * @Description: This is description.
 */

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiaobopang/go_init/controller"
	"github.com/xiaobopang/go_init/middleware"
)

var indexCtl = new(controller.IndexController)
var testCtl = new(controller.TestController)
var wsCtl = new(controller.WsController)
var mqCtl = new(controller.MqController)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	//router.Use(gin.Logger())

	router.GET("/", indexCtl.Welcome)
	router.NoRoute(indexCtl.Handle404)

	// 简单的路由组: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/redis", testCtl.RedisTest) //redis 测试

		v1.POST("/exchange", func(c *gin.Context) {
			mqCtl.ExchangeHandler(c.Writer, c.Request)
		})
		v1.POST("/queue/bind", func(c *gin.Context) {
			mqCtl.QueueBindHandler(c.Writer, c.Request)
		})
		v1.GET("/queue", func(c *gin.Context) {
			mqCtl.QueueHandler(c.Writer, c.Request)
		}) //consume queue
		v1.POST("/queue", func(c *gin.Context) {
			mqCtl.QueueHandler(c.Writer, c.Request)
		}) //declare queue
		v1.DELETE("/queue", func(c *gin.Context) {
			mqCtl.QueueHandler(c.Writer, c.Request)
		}) //delete queue
		v1.POST("/publish", func(c *gin.Context) {
			mqCtl.PublishHandler(c.Writer, c.Request)
		})
		v1.GET("/ws", func(c *gin.Context) {
			wsCtl.WsHandler(c.Writer, c.Request)
		})

		v1.GET("/get_token", testCtl.GetToken)
	}

	router.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.unclepang.com/")
	})

	v2 := router.Group("/v2")
	v2.Use(middleware.CORS(middleware.CORSOptions{}))
	{
		v2.GET("/user", testCtl.GetUser)
		v2.GET("/es", testCtl.ES)
		v2.POST("/user", testCtl.AddUser)
		v2.DELETE("/user", testCtl.DelUser)
		v2.PATCH("/user", testCtl.UptUser)
	}

	return router
}
