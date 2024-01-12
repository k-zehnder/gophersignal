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
	dsn := config.GetEnv("MYSQL_DSN", "")
	if dsn == "" {
		log.Fatal("MYSQL_DSN not set in .env file")
	}

	sqlstore, err := store.NewMySQLStore(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database store: %v", err)
	}

	if err := sqlstore.Init(); err != nil {
		log.Fatalf("Failed to initialize database tables: %v", err)
	}

	handler := routeHandlers.NewHandler(sqlstore)
	r := router.SetupRouter(handler)

	// Start the HTTP server
	addr := config.GetEnv("SERVER_ADDRESS", "0.0.0.0:8080")
	log.Printf("Server is running on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
