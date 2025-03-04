package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MichaelWaters001/youtube-recommender/backend/internal/db"
	"github.com/MichaelWaters001/youtube-recommender/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

// AddTag handles adding a tag to a creator
func AddTag(c *gin.Context) {
	creatorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid creator ID"})
		return
	}

	var request struct {
		TagName string `json:"tag_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Store tag in DB using pgxpool
	tagID, err := db.AddTag(ctx, creatorID, request.TagName, userID.(int))
	if err != nil {
		logger.Log.Error("Failed to store tag", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store tag"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tag_id":     tagID,
		"creator_id": creatorID,
		"tag_name":   request.TagName,
		"user_id":    userID,
	})
}

// GetTags fetches all tags for a given creator
func GetTags(c *gin.Context) {
	creatorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid creator ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tags, err := db.GetTags(ctx, creatorID)
	if err != nil {
		logger.Log.Error("Failed to fetch tags", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
