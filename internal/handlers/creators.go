package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MichaelWaters001/youtube-recommender/internal/db"
	"github.com/MichaelWaters001/youtube-recommender/pkg/logger"
	"github.com/gin-gonic/gin"
)

// YouTube API Key (should be in config)
const youtubeAPIKey = "YOUR_YOUTUBE_API_KEY"

// AddCreator handles adding a new YouTube creator by ID
func AddCreator(c *gin.Context) {
	var request struct {
		YouTubeID string `json:"youtube_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Fetch creator details from YouTube API
	name, description, err := fetchYouTubeChannelDetails(request.YouTubeID)
	if err != nil {
		logger.Log.Error("Failed to fetch YouTube channel details", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch channel details"})
		return
	}

	// Store creator in DB using pgxpool
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	creatorID, err := db.AddCreator(ctx, request.YouTubeID, name, description)
	if err != nil {
		logger.Log.Error("Failed to store creator", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store creator"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          creatorID,
		"youtube_id":  request.YouTubeID,
		"name":        name,
		"description": description,
	})
}

// GetCreator retrieves a creator's details
func GetCreator(c *gin.Context) {
	creatorID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	creator, err := db.GetCreator(ctx, creatorID)
	if err != nil {
		logger.Log.Error("Failed to fetch creator", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Creator not found"})
		return
	}

	c.JSON(http.StatusOK, creator)
}

// Fetch YouTube channel details using YouTube API
func fetchYouTubeChannelDetails(youtubeID string) (string, string, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=snippet&id=%s&key=%s", youtubeID, youtubeAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result struct {
		Items []struct {
			Snippet struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if len(result.Items) == 0 {
		return "", "", fmt.Errorf("channel not found")
	}

	return result.Items[0].Snippet.Title, result.Items[0].Snippet.Description, nil
}
