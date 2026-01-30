package middleware

import (
	"strings"

	"apihub/pkg/auth"
	"apihub/pkg/response"

	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 401, "Authorization header required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, 401, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := auth.ValidateToken(token, secret)
		if err != nil {
			response.Error(c, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
