package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	DatabaseURL           string
	OpenRouterApiEndpoint string
	OpenRouterApiKey      string
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}

func LoadConfig() *Config {
	// Load .env only if present (for local dev)
	_ = godotenv.Load()
	cfg := &Config{
		Port:                  getEnv("PORT", "4000"),
		DatabaseURL:           getEnv("DATABASE_URL", ""),
		OpenRouterApiEndpoint: getEnv("OPEN_ROUTER_API_ENDPOINT", ""),
		OpenRouterApiKey:      getEnv("OPEN_ROUTER_API_KEY", ""),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return cfg
}
