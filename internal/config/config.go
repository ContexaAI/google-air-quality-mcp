package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	APIKey        string
	MCPServerName string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default/env vars")
	}

	return &Config{
		Port:          getEnv("PORT", "8080"),
		APIKey:        getEnv("API_KEY", ""),
		MCPServerName: getEnv("MCP_SERVER_NAME", "google-air-quality-mcp"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
