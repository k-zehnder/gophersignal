// Package main is the entry point for the GopherSignal API server, handling
// initialization and startup tasks including configuration setup, database
// connectivity, API routing, and HTTP server launch.

package main

import (
    "log"
    "net/http"

    "github.com/k-zehnder/gophersignal/backend/config"
    "github.com/k-zehnder/gophersignal/backend/docs"
    "github.com/k-zehnder/gophersignal/backend/internal/api"
    "github.com/k-zehnder/gophersignal/backend/internal/store"
)

// Main function initializes and starts the GopherSignal API server.
func main() {
    // Initialize application configuration from environment variables.
    cfg := config.NewConfig()

    // Configure Swagger documentation host per application settings.
    docs.SwaggerInfo.Host = cfg.SwaggerHost

    // Establish connection to the SQL database.
    sqlStore, err := store.NewMySQLStore(cfg.DataSourceName)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    // Instantiate the server handler with database access.
    handler := api.NewServer(sqlStore)

    // Set up the HTTP server with configured address and handler.
    httpServer := &http.Server{
        Addr:    cfg.ServerAddress,
        Handler: handler,
    }

    // Launch the HTTP server and handle potential start-up errors.
    log.Printf("Server starting on %s\n", cfg.ServerAddress)
    if err := httpServer.ListenAndServe(); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
