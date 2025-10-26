package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	HTTPAddr          string
	TSServerURL       string
	TSAPIToken        string
	TSVirtualServerID string
	LogLevel          string
}

// Load reads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	cfg := &Config{
		HTTPAddr:          getEnv("HTTP_ADDR", ":8080"),
		TSServerURL:       getEnv("TS_SERVER_URL", ""),
		TSAPIToken:        getEnv("TS_API_TOKEN", ""),
		TSVirtualServerID: getEnv("TS_VIRTUAL_SERVER_ID", ""),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
	}

	// Validate critical settings if needed
	// For now, we just log the config
	return cfg, nil
}

// getEnv reads an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// String returns a string representation of the config (without sensitive data)
func (c *Config) String() string {
	return fmt.Sprintf("HTTPAddr=%s, LogLevel=%s", c.HTTPAddr, c.LogLevel)
}
