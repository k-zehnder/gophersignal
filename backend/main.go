package main

import (
	"log"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/pkg/router"
)

func main() {
	// Load environment variables from .env file
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the router
	r := router.SetupRouter()

	// Start the HTTP server with the router
	addr := config.GetEnvVar("SERVER_ADDRESS", "0.0.0.0:8080")
	log.Printf("Server is running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
