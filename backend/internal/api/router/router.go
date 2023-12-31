package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
)

func SetupRouter(handler *routeHandlers.Handler) *mux.Router {
	r := mux.NewRouter()

	// Enable CORS
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{
			"http://localhost:3000",
			"https://gophersignal.com",
		}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Apply CORS middleware to the router
	r.Use(cors)

	// Setup API routes and subroutes
	setupAPIRoutes(r, handler)

	return r
}

func setupAPIRoutes(r *mux.Router, handler *routeHandlers.Handler) {
	v1 := r.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/articles", handler.GetArticles).Methods("GET")
}