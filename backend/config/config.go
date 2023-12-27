package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	// Load environment variables from .env file
	err := godotenv.Load()
	return err
}

func GetEnvVar(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
