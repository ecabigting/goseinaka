package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	APIPort            string
	APIKey             string
	JWTSecret          string
	JWTAccessTokenTTL  int
	JWTRefreshTokenTTL int
	DatabaseURL        string
	GoogleClientID     string
	GoogleClientSecret string
	GinMode            string
	LogLevel           string
}

const (
	defaultAPIPort = "8080"
	aTokenTTL      = 15   // 15 mins default
	rTokenTTL      = 3600 // 3600 mins default
)

func Load() (*Config, error) {
	var cfg Config

	err := godotenv.Load() // Load .env file if present
	if err != nil {
		log.Println("INFO: No .env file found or error loading, relying on system environment variables.")
	}

	// Get API Port
	cfg.APIPort = os.Getenv("API_PORT")
	if cfg.APIPort == "" {
		log.Printf(" ..empty APIPort defaulting to: %s", defaultAPIPort)
		cfg.APIPort = defaultAPIPort
	}

	// Get API Key
	cfg.APIKey = os.Getenv("API_KEY")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("empty API Key")
	}

	// GET JWT Secret
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("empty JWT Secret")
	}

	// GET access token time to live
	cfg.JWTAccessTokenTTL, err = strconv.Atoi(os.Getenv("JWT_AccessToken_TTL"))
	if err != nil {
		cfg.JWTAccessTokenTTL = aTokenTTL
	}

	// GET refresh token time to live
	cfg.JWTRefreshTokenTTL, err = strconv.Atoi(os.Getenv("JWT_RefreshToken_TTL"))
	if err != nil {
		cfg.JWTRefreshTokenTTL = rTokenTTL
	}

	// Get Database URL
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("empty Database URL")
	}

	// Get Log level
	cfg.LogLevel = os.Getenv("LOG_LEVEL")
	if cfg.LogLevel == "" {
		log.Printf(" ..log level not set defaulting to info")
		cfg.LogLevel = "Info"
	}

	// Get GIN Mode
	cfg.GinMode = os.Getenv("GIN_MODE")
	if cfg.GinMode == "" {
		log.Printf(" ..GIN Mode not set, setting to default: debug")
		cfg.GinMode = "debug"
	}

	// AUTH
	// GOOGLE
	cfg.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	if cfg.GoogleClientID == "" {
		return nil, fmt.Errorf("google Client ID Not set")
	}

	cfg.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	if cfg.GoogleClientSecret == "" {
		return nil, fmt.Errorf("google Client Secret Not set")
	}

	return &cfg, nil
}
