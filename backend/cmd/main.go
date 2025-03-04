package main

import (
	"fmt"

	"github.com/MichaelWaters001/youtube-recommender/backend/internal/auth"
	"github.com/MichaelWaters001/youtube-recommender/backend/internal/db"
	"github.com/MichaelWaters001/youtube-recommender/backend/internal/handlers"
	"github.com/MichaelWaters001/youtube-recommender/backend/pkg/config"
	"github.com/MichaelWaters001/youtube-recommender/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger.InitLogger()

	// Load config
	config.LoadConfig()

	// Initialize database
	if err := db.InitDB(); err != nil {
		panic(fmt.Sprintf("Database initialization failed: %v", err))
	}

	// Initialize Google OAuth
	auth.InitAuth()

	// Create Gin router
	r := gin.Default()

	// Public routes
	r.GET("/auth/google", auth.GoogleLogin)
	r.GET("/auth/google/callback", auth.GoogleCallback)
	r.POST("/auth/logout", auth.Logout)

	// Public routes for viewing information
	r.GET("/creators/:id", handlers.GetCreator)
	r.GET("/creators/:id/tags", handlers.GetTags)
	r.GET("/search", handlers.SearchCreators) // Allow public searching

	// Protected routes (require JWT for adding/modifying data)
	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/creators", handlers.AddCreator)
		protected.POST("/creators/:id/tags", handlers.AddTag)
		protected.POST("/votes", handlers.VoteTag)
		protected.DELETE("/votes/:creator_tag_id", handlers.RemoveVote)
	}

	// Start server
	port := config.AppConfig.Server.Port
	logger.Log.Info("Starting server", "port", port)

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Log.Error("Server failed", "error", err)
	}
}
