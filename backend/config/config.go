package config

import (
	"log"
	"os"

	"github.com/k-zehnder/gophersignal/backend/docs"
)

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func Init() string {
	dsn := GetEnv("MYSQL_DSN", "")
	if dsn == "" {
		log.Fatal("MYSQL_DSN not set in .env file")
	}

	env := GetEnv("GO_ENV", "dev")
	if env == "dev" {
		docs.SwaggerInfo.Host = "localhost:8080"
	} else {
		docs.SwaggerInfo.Host = "0.0.0.0:8080"
	}

	return dsn
}
