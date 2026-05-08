package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"lakoo/backend/pkg/auth"
	"lakoo/backend/pkg/config"
	"lakoo/backend/pkg/response"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenStr := ""

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		}

		// Fallback to Cookie if Header is missing or invalid
		if tokenStr == "" {
			cookie, err := c.Cookie("token")
			if err == nil {
				tokenStr = cookie
			}
		}

		if tokenStr == "" {
			response.Error(c, 401, "UNAUTHORIZED", "Missing authorization token")
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(tokenStr, cfg)
		if err != nil {
			response.Error(c, 401, "UNAUTHORIZED", "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
