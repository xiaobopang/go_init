package routers

/*
 * @Script: routers.go
 * @Author: pangxiaobo
 * @Email: 10846295@qq.com
 * @Create At: 2018-11-27 18:19:27
 * @Last Modified By: pangxiaobo
 * @Last Modified At: 2018-11-27 18:19:27
 * @Description: This is description.
 */

import (
	"github.com/gin-gonic/gin"
	"github.com/go_init/controllers"
	"github.com/go_init/middleware"
	"net/http"
)

var indexCtl = new(controllers.IndexController)
var testCtl = new(controllers.TestController)

func SetupRouter() *gin.Engine {

	router := gin.Default()
	router.Use(gin.Recovery())
	//router.Use(gin.Logger())

	router.GET("/", indexCtl.Welcome)
	router.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.unclepang.com/")
	})
	router.NoRoute(indexCtl.Handle404)

	v1 := router.Group("/v1")
	v1.Use(middleware.CORS(middleware.CORSOptions{}))
	{
		v1.GET("/test", testCtl.GetNick)
	}

	v2 := router.Group("/v2")
	v2.Use(middleware.CORS(middleware.CORSOptions{}))
	{
		v2.GET("/user", testCtl.GetUser)
		v2.POST("/user", testCtl.AddUser)
		v2.DELETE("/user", testCtl.DelUser)
		v2.PATCH("/user", testCtl.UptUser)
	}

	return router
}
