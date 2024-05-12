// Package server sets up the web server and routing logic for the GopherSignal application.
package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k-zehnder/gophersignal/backend/internal/api/controllers"
	"github.com/k-zehnder/gophersignal/backend/internal/api/router"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// NewServer creates an http.Handler with configured routes and controllers.
func NewServer(store store.Store) http.Handler {
	articlesController := controllers.NewArticlesController(store)
	return router.SetupRouter(articlesController)
}

// StartServer launches an HTTP server on the specified address.
func StartServer(addr string, handler http.Handler) *http.Server {
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

// GracefulShutdown handles server shutdown on receiving interrupt or termination signals.
func GracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Blocks until a signal is received

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
}
