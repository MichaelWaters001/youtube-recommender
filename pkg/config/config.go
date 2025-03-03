package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config structure to hold application configurations
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	OAuth    OAuthConfig
}

// ServerConfig holds server-related configurations
type ServerConfig struct {
	Port int
}

// DatabaseConfig holds database connection details
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

// OAuthConfig holds Google OAuth credentials
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// AppConfig is the global configuration instance
var AppConfig Config

// LoadConfig initializes Viper to read configuration from `config.toml`
func LoadConfig() {
	viper.SetConfigName("config") // Config file name (without extension)
	viper.SetConfigType("toml")   // File format
	viper.AddConfigPath(".")      // Look for config in the current directory

	// Allow environment variables to override config file values
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	log.Println("Config loaded successfully")
}
