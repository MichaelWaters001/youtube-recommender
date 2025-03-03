package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/MichaelWaters001/youtube-recommender/internal/db"
	"github.com/MichaelWaters001/youtube-recommender/pkg/logger"
	"github.com/gin-gonic/gin"
)

// SearchCreators finds creators based on a tag
func SearchCreators(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag query parameter is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	creators, err := db.SearchCreatorsByTag(ctx, tag)
	if err != nil {
		logger.Log.Error("Failed to search creators", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search creators"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"creators": creators})
}
