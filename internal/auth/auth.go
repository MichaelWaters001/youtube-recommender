package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MichaelWaters001/youtube-recommender/internal/db"
	"github.com/MichaelWaters001/youtube-recommender/pkg/config"
	"github.com/MichaelWaters001/youtube-recommender/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConfig *oauth2.Config
var jwtSecret = []byte("your-secret-key") // Move this to config

// Initialize Google OAuth config
func InitAuth() {
	oauthConfig = &oauth2.Config{
		ClientID:     config.AppConfig.OAuth.ClientID,
		ClientSecret: config.AppConfig.OAuth.ClientSecret,
		RedirectURL:  config.AppConfig.OAuth.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// Generate a random key for session encryption (Deprecated with JWT)
func generateRandomKey() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Handle Google Login Redirect
func GoogleLogin(c *gin.Context) {
	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Handle Google OAuth Callback and Return JWT
func GoogleCallback(c *gin.Context) {
	ctx := context.Background()
	code := c.Query("code")

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		logger.Log.Error("OAuth exchange error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := oauthConfig.Client(ctx, token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logger.Log.Error("Failed to get user info", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer userInfoResp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		logger.Log.Error("Failed to decode user info", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Store user in DB if not exists
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	userID, err := db.EnsureUser(ctx, userInfo.Email)
	if err != nil {
		logger.Log.Error("Failed to ensure user in DB", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ensure user in DB"})
		return
	}

	// Generate JWT token
	jwtToken, err := generateJWT(userID)
	if err != nil {
		logger.Log.Error("Failed to generate JWT", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

// Logout handler (JWT is stateless, so no real logout is needed)
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful. Delete JWT on client side."})
}

// GenerateJWT creates a JWT for authentication
func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
