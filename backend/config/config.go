// Package config handles configuration management.
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/k-zehnder/gophersignal/backend/docs"
)

// AppConfig represents the application's configuration.
type AppConfig struct {
	DataSourceName    string // Database connection string
	Environment       string // Application environment (e.g., "development", "production")
	ServerAddress     string // Address on which the server should listen
	SwaggerHost       string // Host for Swagger documentation
	HuggingFaceAPIKey string // API key for Hugging Face service
}

// NewConfig initializes and returns a new AppConfig, loading environment variables from .env file with defaults if not present.
func NewConfig() *AppConfig {
	// Load environment variables
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Failed to load environment variables: %v", err)
	}

	cfg := &AppConfig{
		DataSourceName:    GetDataSourceName(),
		Environment:       GetEnv("GO_ENV", "development"),
		ServerAddress:     GetEnv("SERVER_ADDRESS", "0.0.0.0:8080"),
		SwaggerHost:       GetDefaultSwaggerHost(GetEnv("GO_ENV", "development")),
		HuggingFaceAPIKey: GetEnv("HUGGING_FACE_API_KEY", ""),
	}

	// Configure Swagger host
	docs.SwaggerInfo.Host = cfg.SwaggerHost

	return cfg
}

// GetEnv retrieves the value of an environment variable or returns a fallback if it doesn't exist.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// GetDefaultSwaggerHost returns the default Swagger host based on the environment.
func GetDefaultSwaggerHost(env string) string {
	switch env {
	case "development":
		return "localhost:8080"
	default:
		return "gophersignal.com"
	}
}

// GetDataSourceName constructs the MYSQL_DSN from individual environment variables.
func GetDataSourceName() string {
	user := GetEnv("MYSQL_USER", "user")
	password := GetEnv("MYSQL_PASSWORD", "password")
	host := GetEnv("MYSQL_HOST", "mysql")
	port := GetEnv("MYSQL_PORT", "3306")
	database := GetEnv("MYSQL_DATABASE", "airaccidentdata")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)
}
