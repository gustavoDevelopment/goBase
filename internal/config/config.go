package config

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	AppName         string      `yaml:"application_name"`
	Description     string      `yaml:"description"`
	Version         string      `yaml:"application_version"`
	Uuid            string      `yaml:"entity_uuid"`
	Environment     string      `yaml:"environment"`
	MongoURI        string      `yaml:"-"`
	MongoDB         string      `yaml:"-"`
	MongoCollection string      `yaml:"-"`
	HTTP            HTTPConfig  `yaml:"http"`
	App             AppConfig   `yaml:"app"`
	JSONConfig      *JSONConfig // Embedded JSON configuration
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Port         string `yaml:"port"`
	BasePath     string `yaml:"base_path"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	IdleTimeout  string `yaml:"idle_timeout"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	MongoDB struct {
		URI      string `yaml:"uri"`
		Database string `yaml:"database"`
		Timeout  string `yaml:"timeout"`
	} `yaml:"mongodb"`
	JWTSecret          string `yaml:"jwt_secret"`
	PasswordSaltRounds int    `yaml:"password_salt_rounds"`
	JSONConfigPath     string `yaml:"json_config_path"`
}

// LoadConfig reads configuration from YAML file, environment variables, and JSON config
func LoadConfig(configPath string) (*Config, error) {
	// Load environment variables from .env file if it exists first
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		logger.Log.Info("No .env file found, using system environment variables", zap.Error(err))
	}

	// Load YAML config
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Expand environment variables in the YAML content
	expandedConfig := os.Expand(string(configFile), func(key string) string {
		val := os.Getenv(key)
		return val
	})

	var config Config
	if err := yaml.Unmarshal([]byte(expandedConfig), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	logger.Log.Info("Config loaded successfully", zap.String("config_path", config.App.JSONConfigPath))

	// Load JSON config if path is provided in config
	if config.App.JSONConfigPath != "" {
		jsonConfig, err := LoadJSONConfig(config.App.JSONConfigPath)
		if err != nil {
			logger.Log.Warn("Failed to load JSON config",
				zap.String("path", config.App.JSONConfigPath),
				zap.Error(err))
		} else {
			config.JSONConfig = jsonConfig
		}
	}

	// Set default values if environment variables are not set
	if config.App.MongoDB.URI == "" {
		config.App.MongoDB.URI = "mongodb://localhost:27017"
	}
	if config.App.MongoDB.Database == "" {
		config.App.MongoDB.Database = "business_orchestrator"
	}
	if config.App.MongoDB.Timeout == "" {
		config.App.MongoDB.Timeout = "10s"
	}

	// Map MongoDB configuration from App.MongoDB to top-level fields
	config.MongoURI = config.App.MongoDB.URI
	config.MongoDB = config.App.MongoDB.Database

	// Set default values if not set
	if config.HTTP.BasePath == "" {
		config.HTTP.BasePath = "/api/v1"
	}

	// Clean up base path (ensure it starts with / and doesn't end with /)
	config.HTTP.BasePath = strings.TrimSpace(config.HTTP.BasePath)
	if config.HTTP.BasePath == "" {
		config.HTTP.BasePath = "/api/v1"
	}
	// Ensure it starts with a single slash and doesn't end with a slash
	config.HTTP.BasePath = "/" + strings.Trim(config.HTTP.BasePath, "/")

	return &config, nil
}
