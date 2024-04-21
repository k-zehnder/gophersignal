// Package main for the GopherSignal API server.
// This server application initializes configuration, database connections, and API routing.
// It also sets up a Swagger documentation endpoint and starts an HTTP server for handling API requests.
package main

import (
	"log"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/docs"
	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
	"github.com/k-zehnder/gophersignal/backend/internal/api/router"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// @title GopherSignal API
// @description API server for the GopherSignal application.
// @version 1
// @BasePath /api/v1
// @host gophersignal.com

// main is the entry point of the application.
// It orchestrates the initialization of the application configuration, database store, HTTP router, and starts the server.
func main() {
	// Initialize application configuration.
	appConfig := config.NewConfig()

	// Configure Swagger host based on the application configuration.
	docs.SwaggerInfo.Host = appConfig.SwaggerHost

	// Initialize database store using configuration data.
	sqlstore, err := store.NewMySQLStore(appConfig.DataSourceName)
	if err != nil {
		log.Fatalf("Failed to initialize database store: %v", err)
	}

	// Set up API handlers and router.
	handler := routeHandlers.NewHandler(sqlstore)
	router := router.SetupRouter(handler)

	// Start the HTTP server using the address from the configuration.
	log.Printf("Starting server on %s\n", appConfig.ServerAddress)
	if err := http.ListenAndServe(appConfig.ServerAddress, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
