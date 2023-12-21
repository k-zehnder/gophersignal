package main

import (
	// "fmt"

	"log"
	"net/http"
	"os"

	// "os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/k-zehnder/gophersignal/backend/myhandlers"
	// "github.com/k-zehnder/gophersignal/backend/pkg/hackernews"
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

	/*
	hns := hackernews.HackerNewsScraper{}
	articles, err := hns.Scrape()
	if err != nil {
		log.Fatal(err)
	}

	// Save articles to database
	dbStore.SaveArticles(articles)

	for i, article := range articles {
		fmt.Printf("[%d] %s - %s\n", i, article.Title, article.Link)
	}
	*/

	// Create a new mux.Router
	r := mux.NewRouter()

	// Enable CORS
	cors := handlers.CORS(
	    handlers.AllowedOrigins([]string{
		"http://localhost:3000", // Add other origins as needed
		"https://gophersignal.com", // Add your domain here
	    }),
	    handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
	    handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Apply CORS middleware to your router
	r.Use(cors)

	// Define your API routes here

	// Example: Setup a route for handling /articles
	r.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.GetArticlesHandler(w, r, dbStore)
	}).Methods("GET")

	// Start the HTTP server with your router
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
