package main

import (
	"assetsentinel/internal/config"
	"assetsentinel/internal/handlers"
	"assetsentinel/internal/middleware"
	"assetsentinel/internal/repository"
	"assetsentinel/internal/services"
	"assetsentinel/internal/websocket"
	"assetsentinel/internal/worker"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting server on port %s", cfg.Port)
	log.Printf("Database path: %s", cfg.DBPath)

	db, err := repository.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	wsHub := websocket.NewHub()
	go wsHub.Run()

	repo := repository.NewRepository(db)
	authService := services.NewAuthService(repo, cfg.JWTSecret)
	assetService := services.NewAssetService(repo)
	maintenanceService := services.NewMaintenanceService(repo, wsHub)
	workOrderService := services.NewWorkOrderService(repo, wsHub)
	inventoryService := services.NewInventoryService(repo, wsHub)
	depreciationService := services.NewDepreciationService(repo)

	authHandler := handlers.NewAuthHandler(authService)
	assetHandler := handlers.NewAssetHandler(assetService)
	maintenanceHandler := handlers.NewMaintenanceHandler(maintenanceService)
	workOrderHandler := handlers.NewWorkOrderHandler(workOrderService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	depreciationHandler := handlers.NewDepreciationHandler(depreciationService)

	scheduler := worker.NewScheduler(repo, wsHub)
	go scheduler.Start()
	defer scheduler.Stop()

	r := gin.Default()

	r.Use(middleware.CORS())

	r.GET("/ws", wsHub.HandleWebSocket)

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/dashboard", handlers.GetDashboard(repo))

		assets := api.Group("/assets")
		{
			assets.GET("", assetHandler.List)
			assets.POST("", middleware.RequireRole("admin", "maintenance_manager"), assetHandler.Create)
			assets.GET("/:id", assetHandler.Get)
			assets.PUT("/:id", middleware.RequireRole("admin", "maintenance_manager"), assetHandler.Update)
			assets.DELETE("/:id", middleware.RequireRole("admin"), assetHandler.Delete)
		}

		maintenance := api.Group("/maintenance-plans")
		{
			maintenance.GET("", maintenanceHandler.List)
			maintenance.POST("", middleware.RequireRole("admin", "maintenance_manager"), maintenanceHandler.Create)
			maintenance.GET("/:id", maintenanceHandler.Get)
			maintenance.PUT("/:id", middleware.RequireRole("admin", "maintenance_manager"), maintenanceHandler.Update)
			maintenance.DELETE("/:id", middleware.RequireRole("admin"), maintenanceHandler.Delete)
		}

		workOrders := api.Group("/work-orders")
		{
			workOrders.GET("", workOrderHandler.List)
			workOrders.POST("", middleware.RequireRole("admin", "maintenance_manager"), workOrderHandler.Create)
			workOrders.GET("/:id", workOrderHandler.Get)
			workOrders.PUT("/:id", workOrderHandler.Update)
			workOrders.DELETE("/:id", middleware.RequireRole("admin"), workOrderHandler.Delete)
		}

		inventory := api.Group("/inventory")
		{
			inventory.GET("", inventoryHandler.List)
			inventory.POST("", middleware.RequireRole("admin", "maintenance_manager"), inventoryHandler.Create)
			inventory.GET("/:id", inventoryHandler.Get)
			inventory.PUT("/:id", middleware.RequireRole("admin", "maintenance_manager"), inventoryHandler.Update)
			inventory.DELETE("/:id", middleware.RequireRole("admin"), inventoryHandler.Delete)
		}

		reports := api.Group("/reports")
		{
			reports.GET("/depreciation/:asset_id", depreciationHandler.GetAssetDepreciation)
			reports.GET("/costs/:asset_id", depreciationHandler.GetAssetCosts)
			reports.GET("/costs", depreciationHandler.GetAllCosts)
		}

		audit := api.Group("/audit")
		{
			audit.GET("", handlers.GetAuditLogs(repo))
		}

		orgs := api.Group("/organizations")
		orgs.Use(middleware.RequireRole("admin"))
		{
			orgs.GET("", handlers.ListOrganizations(repo))
			orgs.POST("", handlers.CreateOrganization(repo))
			orgs.GET("/:id", handlers.GetOrganization(repo))
			orgs.PUT("/:id", handlers.UpdateOrganization(repo))
		}

		users := api.Group("/users")
		users.Use(middleware.RequireRole("admin"))
		{
			users.GET("", handlers.ListUsers(repo))
			users.POST("", handlers.CreateUser(repo))
			users.GET("/:id", handlers.GetUser(repo))
			users.PUT("/:id", handlers.UpdateUser(repo))
			users.DELETE("/:id", handlers.DeleteUser(repo))
		}
	}

	log.Printf("Server starting on http://localhost:%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
