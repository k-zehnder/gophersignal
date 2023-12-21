package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/k-zehnder/gophersignal/backend/pkg/myhandlers"
	"github.com/k-zehnder/gophersignal/backend/pkg/store"
)

func main() {
	// Get DSN from environment variable
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN not set in .env file")
	}

	// Initialize database connection
	dbStore := store.NewDBStore(dsn)

	// Create a new mux.Router
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

	// Apply CORS middleware to your router
	r.Use(cors)

	// Define API routes here

	// Setup a route for handling /articles
	r.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.GetArticlesHandler(w, r, dbStore)
	}).Methods("GET")

	// Start the HTTP server with your router
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
