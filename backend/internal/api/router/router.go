// Package router provides a router configuration for the GoPhersignal API service.
// It initializes a Gorilla Mux router with configured routes, CORS settings, and Swagger documentation.
package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRouter initializes and returns a new mux.Router with configured routes and middleware.
func SetupRouter(handler *routeHandlers.Handler) *mux.Router {
	r := mux.NewRouter()

	// Set up CORS (Cross-Origin Resource Sharing) with specified settings.
	// This allows requests from the defined origins and with specified methods and headers.
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{
			"http://localhost:3000",    // Allow local frontend development server
			"http://localhost:8080",    // Allow local development server
			"https://gophersignal.com", // Allow production frontend
		}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Apply CORS middleware to the router.
	r.Use(cors)

	// Configure API routes and their subroutes.
	setupAPIRoutes(r, handler)

	// Serve Swagger-generated API documentation on the /swagger endpoint.
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	return r
}

// setupAPIRoutes configures API routes with their respective handlers.
func setupAPIRoutes(r *mux.Router, handler *routeHandlers.Handler) {
	// Create a subrouter for version 1 of the API.
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Define the route for fetching articles and associate it with the GetArticles handler.
	// This route responds to GET requests on the /api/v1/articles endpoint.
	route := v1.HandleFunc("/articles", handler.GetArticles).Methods("GET")

	// Assign a name to the route for easier identification and debugging.
	route.Name("GetArticles")
}
