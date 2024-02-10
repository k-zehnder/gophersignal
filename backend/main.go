// Package main is the entry point of the GopherSignal API server. It is responsible for the overall
// initialization and startup of the server. This includes setting up the application configuration,
// establishing database connectivity, defining API routes via the api package, and launching the HTTP server.

package main

import (
	"log"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/docs"
	"github.com/k-zehnder/gophersignal/backend/internal/api"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// @title GopherSignal API
// @description API server for the GopherSignal application.
// @version 1
// @BasePath /api/v1
// @host gophersignal.com

// main function: Initializes and starts the GopherSignal API server.
func main() {
	// Initialize application configuration.
	cfg := config.NewConfig()

	// Set Swagger documentation host using the application configuration.
	docs.SwaggerInfo.Host = cfg.SwaggerHost

	// Establish database connection.
	sqlStore, err := store.NewMySQLStore(cfg.DataSourceName)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Create the server handler using the store interface.
	handler := api.NewServer(cfg, sqlStore)

	// Configure the HTTP server
	httpServer := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: handler,
	}

	// Start the HTTP server
	log.Printf("Server starting on %s\n", cfg.ServerAddress)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server start-up failed: %v", err)
	}
}
