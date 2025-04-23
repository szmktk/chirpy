package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	// DBURL is the database connection URL.
	DBURL string
	// Platform is the running platform identifier (optional).
	Platform string
	// PolkaKey is the API key for Polka webhooks.
	PolkaKey string
	// TokenSecret is the secret used to sign JWTs.
	TokenSecret string
}

// LoadConfig reads environment variables (optionally from a .env file) and returns a Config.
// It returns an error if any required variable is missing.
func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	var missing []string
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		missing = append(missing, "DB_URL")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		missing = append(missing, "POLKA_KEY")
	}
	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		missing = append(missing, "TOKEN_SECRET")
	}
	platform := os.Getenv("PLATFORM")

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}
	return &Config{
		DBURL:       dbURL,
		Platform:    platform,
		PolkaKey:    polkaKey,
		TokenSecret: tokenSecret,
	}, nil
}
