// Package main is the entry point for the GopherSignal API server.
// @title GopherSignal API
// @description API server for the GopherSignal application.
// @version 1
// @BasePath /api/v1
package main

import (
	"log"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/docs"
	"github.com/k-zehnder/gophersignal/backend/internal/api/server"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// main initializes and launches the API server.
func main() {
	// Load server configuration
	cfg := config.NewConfig()

	// Configure Swagger and initialize the database
	docs.SwaggerInfo.Host = cfg.SwaggerHost
	store, err := store.NewMySQLStore(cfg.DataSourceName)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// Create and start the HTTP server
	srv := server.StartServer(cfg.ServerAddress, server.NewServer(store))
	defer server.GracefulShutdown(srv)
}
