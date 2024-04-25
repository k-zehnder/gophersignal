// Package api contains the core logic for initializing the HTTP server in the GopherSignal application.
// It includes the NewServer function, which is responsible for setting up the HTTP server with configured routes and middleware.
// This setup ensures that the server is ready to handle requests as per the defined routing logic,
// effectively integrating various components like route handlers and the router for seamless operation.

package api

import (
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
	"github.com/k-zehnder/gophersignal/backend/internal/api/router"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

func NewServer(store store.Store) http.Handler {
	// Initialize route handlers using the provided store.
	articlesHandler := routeHandlers.NewArticlesHandler(store)

	// Configure the router with the route handlers.
	router := router.SetupRouter(articlesHandler)

	// Return the configured router as an http.Handler.
	return router
}
