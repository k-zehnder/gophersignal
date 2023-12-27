package router

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/pkg/myhandlers"
	"github.com/k-zehnder/gophersignal/backend/pkg/store"
)

func SetupRouter() *mux.Router {
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

	// Initialize the database store
	dsn := config.GetEnvVar("MYSQL_DSN", "")
	dbStore := store.NewDBStore(dsn)

	// Setup API routes and subroutes
	setupAPIRoutes(r, dbStore)

	return r
}

func setupAPIRoutes(r *mux.Router, dbStore *store.DBStore) {
	// API Version 1
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Setup a route for handling /api/v1/articles
	v1.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.GetArticlesHandler(w, r, dbStore)
	}).Methods("GET")

	// Add more routes and subroutes here
}
