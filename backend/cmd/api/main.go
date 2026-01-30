package main

import (
	"log"
	"os"

	"apihub/internal/config"
	"apihub/internal/handler"
	"apihub/internal/middleware"
	"apihub/pkg/database"
	"apihub/pkg/logger"
	"apihub/pkg/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Initialize store based on DB mode
	var db *pgxpool.Pool
	var str store.Store
	var dbErr error

	if cfg.Database.Mode == "memory" {
		logger.Info("Using in-memory store")
		str = store.NewMemoryStore()
	} else {
		logger.Info("Using PostgreSQL database")
		// Connect to database
		db, dbErr = database.Connect(cfg.Database)
		if dbErr != nil {
			log.Fatal("Failed to connect to database:", dbErr)
		}
		defer db.Close()
	}

	// Initialize Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Global middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "mode": cfg.Database.Mode})
	})

	// Initialize handlers
	h := handler.New(db, cfg, str)

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.Auth(cfg.JWT.Secret))
		{
			// User routes
			protected.GET("/users/me", h.GetCurrentUser)

			// Project nested routes (must be defined before /projects/:id)
			protected.GET("/projects/:project_id/collections", h.ListCollections)
			protected.POST("/projects/:project_id/collections", h.CreateCollection)
			protected.GET("/projects/:project_id/collections/:id", h.GetCollection)
			protected.PUT("/projects/:project_id/collections/:id", h.UpdateCollection)
			protected.DELETE("/projects/:project_id/collections/:id", h.DeleteCollection)

			protected.GET("/projects/:project_id/environments", h.ListEnvironments)
			protected.POST("/projects/:project_id/environments", h.CreateEnvironment)
			protected.GET("/projects/:project_id/environments/:id", h.GetEnvironment)
			protected.PUT("/projects/:project_id/environments/:id", h.UpdateEnvironment)
			protected.DELETE("/projects/:project_id/environments/:id", h.DeleteEnvironment)

			protected.GET("/collections/:collection_id/endpoints", h.ListEndpoints)
			protected.POST("/collections/:collection_id/endpoints", h.CreateEndpoint)
			protected.GET("/collections/:collection_id/endpoints/:id", h.GetEndpoint)
			protected.PUT("/collections/:collection_id/endpoints/:id", h.UpdateEndpoint)
			protected.DELETE("/collections/:collection_id/endpoints/:id", h.DeleteEndpoint)

			// Project routes (must be last - wildcard routes)
			protected.GET("/projects", h.ListProjects)
			protected.POST("/projects", h.CreateProject)
			protected.GET("/projects/:id", h.GetProject)
			protected.PUT("/projects/:id", h.UpdateProject)
			protected.DELETE("/projects/:id", h.DeleteProject)
		}
	}

	// Start server
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting server on port " + port)
	if err := r.Run(":" + port); err != nil {
		logger.Error("Failed to start server: " + err.Error())
		os.Exit(1)
	}
}
