package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"lakoo/backend/pkg/response"
)

// RequireRole enforces role-based access control by verifying the "role" claim
// previously set into the context by the AuthMiddleware.
// Usage: router.Use(middleware.RequireRole("ADMIN", "OWNER"))
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Error(c, 403, "FORBIDDEN", "Akses Ditolak: Kredensial tidak memuat spesifikasi peran.")
			c.Abort()
			return
		}

		userRole, ok := roleVal.(string)
		if !ok || userRole == "" {
			response.Error(c, 403, "FORBIDDEN", "Akses Ditolak: Format data peran tidak dikenali.")
			c.Abort()
			return
		}

		hasAccess := false
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			log.Printf("RBAC Blocked: Context Role '%s' attempted to breach guards requiring %v", userRole, allowedRoles)
			response.Error(c, 403, "FORBIDDEN", "Akses Ditolak: Anda tidak memiliki otoritas (RBAC) untuk mengeksekusi layanan administrasi ini.")
			c.Abort()
			return
		}

		// Proceed to handler explicitly validated
		c.Next()
	}
}
