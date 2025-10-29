package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	HTTPAddr string                  `yaml:"http_addr"`
	LogLevel string                  `yaml:"log_level"`
	Servers  map[string]ServerConfig `yaml:"servers"`
}

// ServerConfig holds configuration for a single TeamSpeak server
type ServerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// Load reads configuration from a YAML file
// Falls back to default configuration if file doesn't exist
func Load() (*Config, error) {
	configPath := getEnv("TS_CONFIG_FILE", "config.yaml")

	// Try to read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default configuration if file doesn't exist
			return getDefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for missing fields
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	// Apply defaults for each server
	for name, server := range cfg.Servers {
		if server.Port == 0 {
			server.Port = 10011
		}
		cfg.Servers[name] = server
	}

	return cfg, nil
}

// getDefaultConfig returns a default configuration
func getDefaultConfig() *Config {
	return &Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Servers:  make(map[string]ServerConfig),
	}
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
	return fmt.Sprintf("HTTPAddr=%s, LogLevel=%s, Servers=%d", c.HTTPAddr, c.LogLevel, len(c.Servers))
}

// GetServer returns the configuration for a specific server by name
func (c *Config) GetServer(name string) (ServerConfig, bool) {
	server, ok := c.Servers[name]
	return server, ok
}
