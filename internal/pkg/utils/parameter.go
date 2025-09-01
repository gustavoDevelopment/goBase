package utils

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/config"
)

// GetParam searches for a parameter by name in the JSON configuration
// Returns the value and true if found, empty string and false otherwise
func GetParam(name string) (string, bool) {
	jsonConfig := config.GetJSONConfig()
	if jsonConfig == nil {
		return "", false
	}

	// Search in Params
	for _, param := range jsonConfig.Params {
		if param.Name == name {
			return param.Value, true
		}
	}

	// Search in IntegrationPaths (if not found in Params)
	for _, path := range jsonConfig.IntegrationPaths {
		if path.Name == name {
			return path.Value, true
		}
	}

	// Search in Certificates (if needed)
	for _, cert := range jsonConfig.Certificates {
		if cert.Name == name {
			return cert.Value, true
		}
	}

	return "", false
}

// GetParamOrDefault returns the parameter value if found, otherwise returns the default value
func GetParamOrDefault(name, defaultValue string) string {
	if value, found := GetParam(name); found {
		return value
	}
	return defaultValue
}

// GetIntegrationPath is a helper to get integration path by name
func GetIntegrationPath(name string) (string, bool) {
	jsonConfig := config.GetJSONConfig()
	if jsonConfig == nil {
		return "", false
	}

	for _, path := range jsonConfig.IntegrationPaths {
		if path.Name == name {
			return path.Value, true
		}
	}
	return "", false
}

// GetCertificate is a helper to get certificate by name
func GetCertificate(name string) (string, bool) {
	jsonConfig := config.GetJSONConfig()
	if jsonConfig == nil {
		return "", false
	}

	for _, cert := range jsonConfig.Certificates {
		if cert.Name == name {
			return cert.Value, true
		}
	}
	return "", false
}
