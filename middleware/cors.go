package middleware

/*
 * @Script: cors.go
 * @Author: pangxiaobo
 * @Email: 10846295@qq.com
 * @Create At: 2018-11-07 10:31:21
 * @Last Modified By: pangxiaobo
 * @Last Modified At: 2018-11-07 10:56:25
 * @Description: This is description.
 */

import "github.com/gin-gonic/gin"

type CORSOptions struct {
	Origin string
}

// CORS middleware from https://github.com/gin-gonic/gin/issues/29#issuecomment-89132826
func CORS(options CORSOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1") // allow any origin domain
		if options.Origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", options.Origin)
		}
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
