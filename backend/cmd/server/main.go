package main

import (
	"assetsentinel/internal/config"
	"assetsentinel/internal/handlers"
	"assetsentinel/internal/middleware"
	"assetsentinel/internal/repository"
	"assetsentinel/internal/services"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting server on port %s", cfg.Port)

	db, err := repository.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.NewRepository(db)
	authService := services.NewAuthService(repo, cfg.JWTSecret)
	assetService := services.NewAssetService(repo)

	authHandler := handlers.NewAuthHandler(authService)
	assetHandler := handlers.NewAssetHandler(assetService)

	r := gin.Default()
	r.Use(middleware.CORS())

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

		orgs := api.Group("/organizations")
		orgs.Use(middleware.RequireRole("admin"))
		{
			orgs.GET("", handlers.ListOrganizations(repo))
			orgs.POST("", handlers.CreateOrganization(repo))
			orgs.GET("/:id", handlers.GetOrganization(repo))
		}

		users := api.Group("/users")
		users.Use(middleware.RequireRole("admin"))
		{
			users.GET("", handlers.ListUsers(repo))
		}
	}

	log.Printf("Server starting on http://localhost:%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
