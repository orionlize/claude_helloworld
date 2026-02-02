package main

import (
	"log"
	"os"
	"time"

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

	if cfg.Database.Mode == "memory" {
		logger.Info("Using in-memory store")
		str = store.NewMemoryStore()
	} else {
		logger.Info("Using PostgreSQL database")
		// Connect to database
		db, err = database.Connect(cfg.Database)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()
		// Use DatabaseStore instead of direct db access
		str = store.NewDatabaseStore(db)
	}

	// Initialize Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware - allow all origins in development
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  cfg.Environment != "production",
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Global middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "mode": cfg.Database.Mode})
	})

	// Initialize handlers - pass db for SendRequest which needs direct HTTP client
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

			// Project routes
			protected.GET("/projects", h.ListProjects)
			protected.POST("/projects", h.CreateProject)

			// Collections (nested under project)
			protected.GET("/projects/:pid/collections", h.ListCollections)
			protected.POST("/projects/:pid/collections", h.CreateCollection)
			protected.GET("/collections/:cid", h.GetCollection)
			protected.PUT("/collections/:cid", h.UpdateCollection)
			protected.DELETE("/collections/:cid", h.DeleteCollection)

			// Environments (nested under project)
			protected.GET("/projects/:pid/environments", h.ListEnvironments)
			protected.POST("/projects/:pid/environments", h.CreateEnvironment)
			protected.GET("/environments/:eid", h.GetEnvironment)
			protected.PUT("/environments/:eid", h.UpdateEnvironment)
			protected.DELETE("/environments/:eid", h.DeleteEnvironment)

			// Endpoints (nested under collection)
			protected.GET("/collections/:cid/endpoints", h.ListEndpoints)
			protected.POST("/collections/:cid/endpoints", h.CreateEndpoint)
			protected.GET("/endpoints/:epid", h.GetEndpoint)
			protected.PUT("/endpoints/:epid", h.UpdateEndpoint)
			protected.DELETE("/endpoints/:epid", h.DeleteEndpoint)

			// Test request - Send HTTP request for testing
			protected.POST("/test/request", h.SendRequest)

			// Documentation routes
			protected.GET("/projects/:pid/docs", h.GenerateDocs)
			protected.GET("/projects/:pid/docs/export", h.ExportPostman)

			// Project detail routes (must be last)
			protected.GET("/projects/:pid", h.GetProject)
			protected.PUT("/projects/:pid", h.UpdateProject)
			protected.DELETE("/projects/:pid", h.DeleteProject)
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
