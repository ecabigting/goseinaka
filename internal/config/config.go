package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIPort            string
	APIKey             string
	JWTSecret          string
	DatabaseURL        string
	GoogleClientID     string
	GoogleClientSecret string
	GinMode            string
	LogLevel           string
}

const (
	defaultAPIPort = "8080"
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
