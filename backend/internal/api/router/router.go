// Package router configures the HTTP router for the GopherSignal API,
// including CORS settings and Swagger documentation.
package router

import (
	"net/http"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/k-zehnder/gophersignal/backend/internal/api/handlers"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter creates an http.Handler with configured routes and handlers.
func NewRouter(store store.Store) http.Handler {
	articlesHandler := handlers.NewArticlesHandler(store)
	return SetupRouter(articlesHandler)
}

// SetupRouter initializes and returns a configured mux.Router.
func SetupRouter(articlesHandler *handlers.ArticlesHandler) *mux.Router {
	r := mux.NewRouter()

	// Setup CORS for various development and production environments
	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{
			"http://localhost:3000",        // Local frontend dev server
			"http://localhost:8080",        // Local dev server
			"https://gophersignal.com",     // Production frontend
			"https://www.gophersignal.com", // Production frontend with www
		}),
		gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	r.Use(cors)

	// Setup API v1 routes, including the GET endpoint for articles
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.Handle("/articles", articlesHandler)

	// Enable Swagger documentation at '/swagger'
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	return r
}
