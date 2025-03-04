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

// VoteTag handles upvoting or downvoting a tag
func VoteTag(c *gin.Context) {
	var request struct {
		CreatorTagID int `json:"creator_tag_id" binding:"required"`
		VoteType     int `json:"vote_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate vote_type (1 = upvote, -1 = downvote)
	if request.VoteType != 1 && request.VoteType != -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote_type, must be 1 or -1"})
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

	// Store vote in DB using pgxpool
	err := db.VoteTag(ctx, userID.(int), request.CreatorTagID, request.VoteType)
	if err != nil {
		logger.Log.Error("Failed to store vote", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded"})
}

// RemoveVote removes a user's vote from a tag
func RemoveVote(c *gin.Context) {
	creatorTagID, err := strconv.Atoi(c.Param("creator_tag_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid creator tag ID"})
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

	// Remove vote in DB using pgxpool
	err = db.RemoveVote(ctx, userID.(int), creatorTagID)
	if err != nil {
		logger.Log.Error("Failed to remove vote", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote removed"})
}
