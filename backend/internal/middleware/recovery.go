package middleware

import (
	"apihub/pkg/logger"
	"apihub/pkg/response"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: " + err.(string))
				response.Error(c, 500, "Internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
