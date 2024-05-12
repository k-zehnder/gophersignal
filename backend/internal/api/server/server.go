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

	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
	"github.com/k-zehnder/gophersignal/backend/internal/api/router"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// NewServer constructs an http.Handler with necessary route handlers and configurations.
func NewServer(store store.Store) http.Handler {
	articlesHandler := routeHandlers.NewArticlesHandler(store)
	return router.SetupRouter(articlesHandler)
}

// StartServer launches an HTTP server on the specified address with the given request handler.
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

// GracefulShutdown waits for interrupt or termination signals and shuts down the server gracefully.
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
