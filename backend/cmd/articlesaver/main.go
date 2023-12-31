package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/internal/scraper"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dsn := config.GetEnvVar("SCRAPER_MYSQL_DSN", "") // Hack
	if dsn == "" {
		log.Fatal("SCRAPER_MYSQL_DSN not set in .env file")
	}

	// Initialize the database store
	dbStore, err := store.NewDBStore(dsn)
	if err != nil {
		log.Fatal("Error initializing DBStore:", err)
	}

	// Scrape and save articles
	SaveArticles(dbStore)
}

func SaveArticles(dbStore *store.DBStore) {
	hns := scraper.HackerNewsScraper{}
	articles, err := hns.Scrape()
	if err != nil {
		log.Printf("Error scraping articles: %v\n", err)
		return
	}

	if err := dbStore.SaveArticles(articles); err != nil {
		log.Printf("Error saving articles: %v\n", err)
		return
	}

	for _, article := range articles {
		log.Printf("Saved article: %s - %s\n", article.Title, article.Link)
	}
}
