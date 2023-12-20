package main

import (
	"fmt"
	"log"
	"os"

	"github.com/k-zehnder/gophersignal/backend/pkg/hackernews"
	"github.com/k-zehnder/gophersignal/backend/pkg/store"
)

func main() {
	// Get DSN from environment variable
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN not set in .env file")
	}

	// Initialize database connection
	dbStore := store.NewDBStore(dsn)

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
}
