package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/k-zehnder/gophersignal/backend/internal/api/docs"
	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRouter(handler *routeHandlers.Handler) *mux.Router {
	r := mux.NewRouter()

	// Enable CORS
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{
			"http://localhost:3000",
			"http://localhost:8080",
			"https://gophersignal.com",
		}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Apply CORS middleware to the router
	r.Use(cors)

	// Setup API routes and subroutes
	setupAPIRoutes(r, handler)

	// Serve Swagger UI documentation at /swagger/index.html
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	return r
}

// Setup API routes and subroutes
func setupAPIRoutes(r *mux.Router, handler *routeHandlers.Handler) {
	v1 := r.PathPrefix("/api/v1").Subrouter()
	route := v1.HandleFunc("/articles", handler.GetArticles).Methods("GET")
	route.Name("GetArticles")
}
