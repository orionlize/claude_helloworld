package main

import (
	"log"
	"os"

	"apihub/internal/config"
	"apihub/internal/handler"
	"apihub/internal/middleware"
	"apihub/pkg/database"
	"apihub/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Connect to database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

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
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize handlers
	h := handler.New(db, cfg)

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.Auth(cfg.JWT.Secret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", h.GetCurrentUser)
				users.PUT("/me", h.UpdateCurrentUser)
			}

			// Project routes
			projects := protected.Group("/projects")
			{
				projects.GET("", h.ListProjects)
				projects.POST("", h.CreateProject)
				projects.GET("/:id", h.GetProject)
				projects.PUT("/:id", h.UpdateProject)
				projects.DELETE("/:id", h.DeleteProject)
			}

			// API Collection routes
			collections := protected.Group("/projects/:project_id/collections")
			{
				collections.GET("", h.ListCollections)
				collections.POST("", h.CreateCollection)
				collections.GET("/:id", h.GetCollection)
				collections.PUT("/:id", h.UpdateCollection)
				collections.DELETE("/:id", h.DeleteCollection)
			}

			// API Endpoint routes
			endpoints := protected.Group("/collections/:collection_id/endpoints")
			{
				endpoints.GET("", h.ListEndpoints)
				endpoints.POST("", h.CreateEndpoint)
				endpoints.GET("/:id", h.GetEndpoint)
				endpoints.PUT("/:id", h.UpdateEndpoint)
				endpoints.DELETE("/:id", h.DeleteEndpoint)
			}

			// Environment routes
			environments := protected.Group("/projects/:project_id/environments")
			{
				environments.GET("", h.ListEnvironments)
				environments.POST("", h.CreateEnvironment)
				environments.GET("/:id", h.GetEnvironment)
				environments.PUT("/:id", h.UpdateEnvironment)
				environments.DELETE("/:id", h.DeleteEnvironment)
			}

			// Request testing routes
			api.GET("/test/request", h.SendRequest)
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
