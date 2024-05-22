// Package main is the entry point for the GopherSignal API server.
// @title GopherSignal API
// @description API server for the GopherSignal application.
// @version 1
// @BasePath /api/v1
package main

import (
	"log"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/internal/api/server"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// main initializes and launches the API server.
func main() {
	// Load server configuration
	cfg := config.NewConfig()

	// Initialize the database store
	store, err := store.NewMySQLStore(cfg.DataSourceName)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// Create the router
	router := server.NewServer(store)

	// Start the HTTP server
	srv := server.StartServer(cfg.ServerAddress, router)
	defer server.GracefulShutdown(srv)
}
