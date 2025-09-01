package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// JSONConfig represents the structure of the JSON configuration file
type JSONConfig struct {
	IntegrationPaths []IntegrationPath `json:"integrationPaths"`
	Certificates    []Certificate     `json:"certificates"`
	Params          []Parameter       `json:"params"`
}

// IntegrationPath represents a path configuration
// Example: {"name": "path.name", "value": "/some/path"}
type IntegrationPath struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Certificate represents a certificate configuration
// Example: {"name": "cert.name", "value": "cert-data"}
type Certificate struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Parameter represents a parameter configuration
// Example: {"name": "mongo.tenant.schema.identification", "value": "orio484001"}
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var (
	instance *JSONConfig
	once     sync.Once
)

// LoadJSONConfig loads the JSON configuration from the specified file path
// and stores it in a singleton instance
func LoadJSONConfig(filePath string) (*JSONConfig, error) {
	var loadErr error
	once.Do(func() {
		if filePath == "" {
			instance = &JSONConfig{}
			return
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			loadErr = fmt.Errorf("error reading JSON config file: %w", err)
			return
		}

		var config JSONConfig
		if err := json.Unmarshal(data, &config); err != nil {
			loadErr = fmt.Errorf("error unmarshaling JSON config: %w", err)
			return
		}

		instance = &config
	})

	if loadErr != nil {
		return nil, loadErr
	}

	return instance, nil
}

// GetJSONConfig returns the loaded JSON configuration
// Returns nil if not loaded yet
func GetJSONConfig() *JSONConfig {
	return instance
}

// GetJSONConfigPath gets the JSON config file path from environment variables
func GetJSONConfigPath() string {
	return os.Getenv("JSON_CONFIG_PATH")
}
