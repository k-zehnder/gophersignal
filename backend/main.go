package main

import (
	"log"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/config"

	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
	"github.com/k-zehnder/gophersignal/backend/internal/api/router"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// @title GopherSignal API
// @description This is the GopherSignal API server.
// @version 1
// @BasePath /api/v1
// @host gophersignal.com
func main() {
	// Initialize configuration
	dsn := config.Init()

	// Initialize database store
	sqlstore, err := store.NewMySQLStore(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database store: %v", err)
	}

	// Initialize API handler and router
	handler := routeHandlers.NewHandler(sqlstore)
	r := router.SetupRouter(handler)

	// Start the HTTP server
	addr := config.GetEnv("SERVER_ADDRESS", "0.0.0.0:8080")
	log.Printf("Server is running on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
