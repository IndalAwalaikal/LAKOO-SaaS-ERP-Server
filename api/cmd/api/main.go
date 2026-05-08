package main

import (
	"log"

	"github.com/gin-gonic/gin"
	
	"lakoo/backend/pkg/config"
	"lakoo/backend/pkg/database"
	"lakoo/backend/internal/repository"
	"lakoo/backend/internal/service"
	"lakoo/backend/internal/http/handler"
	"lakoo/backend/internal/http/route"
	"lakoo/backend/pkg/storage"
)

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database Connections
	db := database.NewMySQLConnection(cfg)
	defer db.Close()

	// Run Auto Migrations
	database.RunMigrations(db)
	redisClient := database.NewRedisClient(cfg)

	// 3. Initialize Local Storage
	localSvc := storage.NewLocalStorage(cfg.StoragePath)

	// 4. Initialize Repositories
	tenantRepo := repository.NewTenantRepository(db)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	saleRepo := repository.NewSaleRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	financeRepo := repository.NewFinanceRepository(db)

	// 5. Initialize UseCases
	tenantUcase := service.NewTenantService(tenantRepo, userRepo, redisClient, cfg)
	productService := service.NewProductService(productRepo, financeRepo)
	saleService := service.NewSaleService(saleRepo, productRepo, financeRepo)
	customerService := service.NewCustomerService(customerRepo)
	inventoryService := service.NewInventoryService(inventoryRepo, productRepo, financeRepo)
	financeService := service.NewFinanceService(financeRepo)
	notificationService := service.NewNotificationService(productRepo, saleRepo)

	// 6. Initialize Handlers
	tenantHnd := handler.NewTenantHandler(tenantUcase)
	productHnd := handler.NewProductHandler(productService)
	saleHnd := handler.NewSaleHandler(saleService)
	customerHnd := handler.NewCustomerHandler(customerService)
	inventoryHnd := handler.NewInventoryHandler(inventoryService)
	financeHnd := handler.NewFinanceHandler(financeService)
	notificationHnd := handler.NewNotificationHandler(notificationService)
	mediaHnd := handler.NewMediaHandler(localSvc)

	// 7. Setup Gin Router
	r := gin.Default()
	
	// Serve static files from storage directory
	r.Static("/storage", cfg.StoragePath)
	
	// 8. Register Routes securely in delivery layer
	route.RegisterRoutes(route.RouterParams{
		Engine:           r,
		TenantHandler:    tenantHnd,
		ProductHandler:   productHnd,
		SaleHandler:      saleHnd,
		CustomerHandler:  customerHnd,
		InventoryHandler:    inventoryHnd,
		FinanceHandler:      financeHnd,
		MediaHandler:        mediaHnd,
		NotificationHandler: notificationHnd,
		Config:              cfg,
		RedisClient:         redisClient,
		TenantRepo:          tenantRepo,
	})

	// 8. Start Server
	log.Printf("Starting Lakoo API server on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
