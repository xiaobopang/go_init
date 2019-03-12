package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xiaobopang/go_init/lib"
	"github.com/xiaobopang/go_init/model"
)

func jwtAbort(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"status":  "error",
		"message": msg,
	})
	c.Abort()
}

func JWTMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		fmt.Println(authHeader)
		if authHeader == "" {
			jwtAbort(c, "Authorization Failed.")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			jwtAbort(c, "Authorization Failed.")
			return
		}

		claims, err := lib.ParseToken(parts[1])
		if err != nil {
			jwtAbort(c, "无效的Token")
			return
		}

		if time.Now().Unix() > claims.ExpiresAt {
			jwtAbort(c, "Token已过期")
			return
		}

		user := model.User{}
		db.First(&user, claims.UserID)

		if user.ID != claims.UserID {
			jwtAbort(c, "无效的Token")
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
