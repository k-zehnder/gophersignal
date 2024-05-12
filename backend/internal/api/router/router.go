// Package router configures the HTTP router for the GopherSignal API,
// including CORS settings and Swagger documentation.
package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/k-zehnder/gophersignal/backend/internal/api/controllers"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRouter initializes and returns a configured mux.Router.
func SetupRouter(articlesController *controllers.ArticlesController) *mux.Router {
	r := mux.NewRouter()

	// Setup CORS for various development and production environments
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{
			"http://localhost:3000",        // Local frontend dev server
			"http://localhost:8080",        // Local dev server
			"https://gophersignal.com",     // Production frontend
			"https://www.gophersignal.com", // Production frontend with www
		}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	r.Use(cors)

	// Setup API v1 routes, including the GET endpoint for articles
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.Handle("/articles", articlesController)

	// Enable Swagger documentation at '/swagger'
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	return r
}
