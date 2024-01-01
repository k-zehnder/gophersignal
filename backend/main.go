package main

import (
	"log"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
	"github.com/k-zehnder/gophersignal/backend/internal/api/router"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

func main() {
	// Load environment variables from the .env file.
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Retrieve the MySQL data source name (DSN) from environment variables
	dsn := config.GetEnvVar("MYSQL_DSN", "")
	if dsn == "" {
		log.Fatal("MYSQL_DSN not set in .env file")
	}

	// Initialize the database store with the provided DSN
	dbStore, err := store.NewDBStore(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database store: %v", err)
	}

	// Initialize the database tables
	if err := dbStore.Init(); err != nil {
		log.Fatalf("Failed to initialize database tables: %v", err)
	}

	// Create a new handler for routing, injecting the database store dependency
	handler := routeHandlers.NewHandler(dbStore)

	// Set up the HTTP router with the handler
	r := router.SetupRouter(handler)

	// Determine the server address from environment variables and start the HTTP server
	addr := config.GetEnvVar("SERVER_ADDRESS", "0.0.0.0:8080")
	log.Printf("Server is running on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
