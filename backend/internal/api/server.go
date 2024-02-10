// Package api contains the core HTTP server initialization and testing logic for the GopherSignal application.
// It includes the NewServer function which sets up the HTTP server with routes and middleware,
// and a unit test for NewServer which ensures correct server behavior and route handling using mock data.

package api

import (
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
	"github.com/k-zehnder/gophersignal/backend/internal/api/router"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

func NewServer(cfg *config.AppConfig, store store.Store) http.Handler {
	// Instantiate route handlers
	articlesHandler := routeHandlers.NewArticlesHandler(store)

	// Set up the router
	router := router.SetupRouter(articlesHandler)

	// Return the configured router as an http.Handler
	return router
}
