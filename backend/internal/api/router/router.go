// Package router provides the configuration for the HTTP router of the GopherSignal API.
// It focuses on initializing a Gorilla Mux router with configured routes, CORS settings, and Swagger documentation support.

package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRouter initializes and returns a new mux.Router with configured routes and middleware.
func SetupRouter(articlesHandler *routeHandlers.ArticlesHandler) *mux.Router {
	// Initialize a new Gorilla Mux router.
	r := mux.NewRouter()

	// Set up CORS (Cross-Origin Resource Sharing) with specified settings.
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{
			"http://localhost:3000",    // Allow local frontend development server
			"http://localhost:8080",    // Allow local development server
			"https://gophersignal.com", // Allow production frontend
			"https://www.gophersignal.com",  // Allow production frontend with www
		}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Apply CORS middleware to the router.
	r.Use(cors)

	// Create a subrouter for version 1 of the API.
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	// Associate API routes with their respective handlers.
	apiRouter.Handle("/articles", articlesHandler)

	// Serve Swagger documentation.
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Return the configured router as a mux.Router.
	return r
}
