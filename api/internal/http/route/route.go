package route

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	
	"lakoo/backend/internal/http/handler"
	"lakoo/backend/internal/middleware"
	"lakoo/backend/internal/repository"
	"lakoo/backend/pkg/config"
)

type RouterParams struct {
	Engine              *gin.Engine
	TenantHandler       *handler.TenantHandler
	ProductHandler      *handler.ProductHandler
	CustomerHandler     *handler.CustomerHandler
	SaleHandler         *handler.SaleHandler
	InventoryHandler    *handler.InventoryHandler
	FinanceHandler      *handler.FinanceHandler
	MediaHandler        *handler.MediaHandler
	NotificationHandler *handler.NotificationHandler
	Config              *config.Config
	RedisClient         *redis.Client
	TenantRepo          repository.TenantRepository
}

func RegisterRoutes(p RouterParams) {
	// Global Middlewares
	p.Engine.Use(middleware.CORSMiddleware())
	p.Engine.Use(middleware.SecurityMiddleware())
	p.Engine.Use(middleware.TenantResolver(p.TenantRepo))
	
	v1 := p.Engine.Group("/api/v1")
	
	// Open Routes
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up", "service": "api"})
	})

	auth := v1.Group("/auth")
	{
		auth.POST("/register", p.TenantHandler.Register)
		auth.POST("/login", middleware.RateLimitMiddleware(p.RedisClient, 5, 15*time.Minute), p.TenantHandler.Login)
		auth.POST("/forgot-password", p.TenantHandler.ForgotPassword)
		auth.POST("/reset-password", p.TenantHandler.ResetPassword)
		auth.POST("/logout", p.TenantHandler.Logout)
	}

	// Protected Routes
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware(p.Config))
	{
		protected.GET("/notifications", p.NotificationHandler.GetNotifications)

		tenantGroup := protected.Group("/tenant")
		{
			tenantGroup.PUT("", p.TenantHandler.UpdateProfile)
			tenantGroup.PUT("/password", p.TenantHandler.ChangePassword)
			tenantGroup.GET("/users", p.TenantHandler.ListStaff)
			tenantGroup.POST("/users", middleware.RequireRole("owner"), p.TenantHandler.AddStaff)
			tenantGroup.DELETE("/users/:id", middleware.RequireRole("owner"), p.TenantHandler.RemoveStaff)
		}
		
		products := protected.Group("/products")
		{
			products.POST("", p.ProductHandler.Create)
			products.PUT("/:id", p.ProductHandler.Update)
			products.DELETE("/:id", p.ProductHandler.Delete)
			products.GET("", p.ProductHandler.List)
			products.GET("/:id", p.ProductHandler.Get)
		}

		customers := protected.Group("/customers")
		{
			customers.POST("", p.CustomerHandler.Create)
			customers.PUT("/:id", p.CustomerHandler.Update)
			customers.DELETE("/:id", p.CustomerHandler.Delete)
			customers.GET("", p.CustomerHandler.List)
			customers.GET("/:id", p.CustomerHandler.Get)
		}

		sales := protected.Group("/sales")
		{
			sales.POST("", p.SaleHandler.Create)
			sales.GET("", p.SaleHandler.List)
			sales.GET("/trend", p.SaleHandler.GetTrend)
		}

		inventory := protected.Group("/inventory")
		{
			inventory.POST("/adjust", p.InventoryHandler.Adjust)
			inventory.GET("/:productId/history", p.InventoryHandler.History)
		}

		finance := protected.Group("/finance")
		finance.Use(middleware.RequireRole("owner", "manager"))
		{
			finance.POST("", p.FinanceHandler.Record)
			finance.DELETE("/:id", p.FinanceHandler.Delete)
			finance.GET("", p.FinanceHandler.List)
			finance.GET("/summary", p.FinanceHandler.Summary)
		}

		media := protected.Group("/media")
		{
			media.POST("/upload", p.MediaHandler.Upload)
		}
	}
}
