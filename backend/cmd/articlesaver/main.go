package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/k-zehnder/gophersignal/backend/pkg/hackernews"
	"github.com/k-zehnder/gophersignal/backend/pkg/store"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("SCRAPER_MYSQL_DSN") // Hack
	if dsn == "" {
		log.Fatal("MYSQL_DSN not set in .env file")
	}

	dbStore := store.NewDBStore(dsn)
	SaveArticles(dbStore)
}

func SaveArticles(dbStore *store.DBStore) {
	hns := hackernews.HackerNewsScraper{}
	articles, err := hns.Scrape()
	if err != nil {
		log.Printf("Error scraping articles: %v\n", err)
		return // Don't terminate the program, just return
	}

	dbStore.SaveArticles(articles)

	for _, article := range articles {
		log.Printf("Saved article: %s - %s\n", article.Title, article.Link)
	}
}
