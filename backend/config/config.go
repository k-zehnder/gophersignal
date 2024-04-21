// Package config handles the configuration management for the GopherSignal application.
// It provides functionalities to initialize the application configuration and retrieve environment variables,
// with support for default values.
package config

import (
	"os"
)

// AppConfig represents the application's configuration.
type AppConfig struct {
	DataSourceName    string
	Environment       string
	ServerAddress     string
	SwaggerHost       string
	HuggingFaceAPIKey string
	OpenAIAPIKey      string
}

// NewConfig initializes and returns a new AppConfig with default values obtained from environment variables.
func NewConfig() *AppConfig {
	return &AppConfig{
		DataSourceName:    GetEnv("MYSQL_DSN", "default_dsn"),
		Environment:       GetEnv("GO_ENV", "dev"),
		ServerAddress:     GetEnv("SERVER_ADDRESS", "0.0.0.0:8080"),
		SwaggerHost:       GetDefaultSwaggerHost(GetEnv("GO_ENV", "dev")), // Pass the environment directly
		HuggingFaceAPIKey: GetEnv("HUGGING_FACE_API_KEY", ""),
		OpenAIAPIKey:      GetEnv("OPEN_AI_API_KEY", ""),
	}
}

// GetEnv retrieves the value of an environment variable or returns a fallback value if the variable is not set.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// GetDefaultSwaggerHost returns the default Swagger host based on the environment.
func GetDefaultSwaggerHost(env string) string {
	switch env {
	case "dev":
		return "localhost:8080"
	default:
		return "gophersignal.com"
	}
}
