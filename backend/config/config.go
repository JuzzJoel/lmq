package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	GeoIPDBPath string
}

// Load retrieves the configuration from environment variables or a .env file.
func Load() (*Config, error) {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	geoPath := os.Getenv("GEOIP_DB_PATH")
	if geoPath == "" {
		geoPath = "data/GeoLite2-City.mmdb"
	}

	dbURL := os.Getenv("DATABASE_URL")
	redisURL := os.Getenv("REDIS_URL")

	if dbURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}
	if redisURL == "" {
		return nil, errors.New("REDIS_URL environment variable is required")
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		RedisURL:    redisURL,
		GeoIPDBPath: geoPath,
	}, nil
}
