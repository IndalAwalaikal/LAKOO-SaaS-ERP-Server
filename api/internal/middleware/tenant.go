package middleware

import (
	"strings"
	"github.com/gin-gonic/gin"
	"lakoo/backend/internal/repository"
	"lakoo/backend/pkg/response"
)

func TenantResolver(repo repository.TenantRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		
		if strings.Contains(host, ":") {
			parts := strings.Split(host, ":")
			host = parts[0]
		}

		parts := strings.Split(host, ".")
		
		// Bypass resolver for plain localhost (development)
		if host == "localhost" || host == "127.0.0.1" {
			c.Next()
			return
		}

		if len(parts) < 2 || (len(parts) == 2 && parts[1] != "localhost") {
			response.Error(c, 400, "BAD_REQUEST", "Tenant subdomain is missing or invalid")
			c.Abort()
			return
		}

		slug := parts[0]
		tenant, err := repo.GetBySlug(slug)
		if err != nil || tenant == nil {
			response.Error(c, 404, "NOT_FOUND", "Tenant not registered")
			c.Abort()
			return
		}

		c.Set("tenant_slug", slug)
		c.Set("resolved_tenant_id", tenant.ID)
		c.Next()
	}
}
