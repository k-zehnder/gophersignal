// Package main is the entry point for the GopherSignal API server, handling
// initialization and startup tasks including configuration setup, database
// connectivity, API routing, and HTTP server launch.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/docs"
	"github.com/k-zehnder/gophersignal/backend/internal/api"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// main sets up and launches the API server.
func main() {
	cfg := config.NewConfig()

	// Set up config, database, API routing, and HTTP server
	docs.SwaggerInfo.Host = cfg.SwaggerHost
	store, err := store.NewMySQLStore(cfg.DataSourceName)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	handler := api.NewServer(store)

	// Start HTTP server
	server := startServer(cfg.ServerAddress, handler)
	defer gracefulShutdown(server)
}

// startServer runs the HTTP server on the specified address.
func startServer(addr string, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	return server
}

// gracefulShutdown cleanly stops the server on system signals.
func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Blocks

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Println("Server exited cleanly")
}
