package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MichaelWaters001/youtube-recommender/backend/internal/db"
	"github.com/MichaelWaters001/youtube-recommender/backend/pkg/config"
	"github.com/MichaelWaters001/youtube-recommender/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

// AddCreator handles adding a new YouTube creator by @handle
func AddCreator(c *gin.Context) {
	var request struct {
		YouTubeHandle string `json:"youtube_handle" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Fetch channel details from YouTube API using @handle
	channelID, name, description, err := fetchYouTubeChannelDetails(request.YouTubeHandle)
	if err != nil {
		logger.Log.Error("Failed to fetch YouTube channel details", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch channel details"})
		return
	}

	// Store creator in DB using pgxpool
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	creatorID, err := db.AddCreator(ctx, request.YouTubeHandle, channelID, name, description)
	if err != nil {
		logger.Log.Error("Failed to store creator", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store creator"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":             creatorID,
		"youtube_handle": request.YouTubeHandle,
		"youtube_id":     channelID,
		"name":           name,
		"description":    description,
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

// Fetch YouTube channel details using YouTube API with @handle support
func fetchYouTubeChannelDetails(youtubeHandle string) (string, string, string, error) {
	apiKey := config.AppConfig.YouTube.APIKey

	// Remove '@' from handle if present
	if youtubeHandle[0] == '@' {
		youtubeHandle = youtubeHandle[1:]
	}

	// Step 1: Get channel ID from @handle
	searchURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?part=snippet&type=channel&q=%s&key=%s",
		youtubeHandle, apiKey,
	)

	resp, err := http.Get(searchURL)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	var searchResult struct {
		Items []struct {
			ID struct {
				ChannelID string `json:"channelId"`
			} `json:"id"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return "", "", "", err
	}

	if len(searchResult.Items) == 0 {
		return "", "", "", fmt.Errorf("channel not found for handle: @%s", youtubeHandle)
	}

	channelID := searchResult.Items[0].ID.ChannelID

	// Step 2: Get channel details using Channel ID
	detailsURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/channels?part=snippet&id=%s&key=%s",
		channelID, apiKey,
	)

	resp, err = http.Get(detailsURL)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	var detailsResult struct {
		Items []struct {
			Snippet struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&detailsResult); err != nil {
		return "", "", "", err
	}

	if len(detailsResult.Items) == 0 {
		return "", "", "", fmt.Errorf("failed to fetch details for channel ID: %s", channelID)
	}

	return channelID, detailsResult.Items[0].Snippet.Title, detailsResult.Items[0].Snippet.Description, nil
}
