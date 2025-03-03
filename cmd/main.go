package main

import (
	"fmt"

	"github.com/MichaelWaters001/youtube-recommender/internal/auth"
	"github.com/MichaelWaters001/youtube-recommender/internal/db"
	"github.com/MichaelWaters001/youtube-recommender/internal/handlers"
	"github.com/MichaelWaters001/youtube-recommender/pkg/config"
	"github.com/MichaelWaters001/youtube-recommender/pkg/logger"
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

	// Protected routes (require JWT)
	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		protected.POST("/creators", handlers.AddCreator)
		protected.GET("/creators/:id", handlers.GetCreator)
		protected.POST("/creators/:id/tags", handlers.AddTag)
		protected.GET("/creators/:id/tags", handlers.GetTags)
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
