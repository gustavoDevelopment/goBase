package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	MongoURI        string
	MongoDB         string
	MongoCollection string
	Timeout         time.Duration
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (*Config, error) {
	// Try to load from .env file, but don't fail if it doesn't exist
	_ = godotenv.Load(path)

	timeout, err := time.ParseDuration(getEnv("TIMEOUT", "10s"))
	if err != nil {
		log.Printf("Invalid timeout value, using default: %v", err)
		timeout = 10 * time.Second
	}

	config := &Config{
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:         getEnv("MONGO_DATABASE", "test"),
		MongoCollection: getEnv("MONGO_COLLECTION", "users"),
		Timeout:         timeout,
	}

	return config, nil
}

// Helper function to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
