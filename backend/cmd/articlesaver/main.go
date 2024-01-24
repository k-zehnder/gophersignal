package main

import (
	"log"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/internal/scraper"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

func main() {
	dsn := config.GetEnv("MYSQL_DSN", "") 
	if dsn == "" {
		log.Fatal("MYSQL_DSN not set in .env file")
	}

	sqlstore, err := store.NewMySQLStore(dsn)
	if err != nil {
		log.Fatal("Error initializing DBStore:", err)
	}

	SaveArticles(sqlstore)
}

func SaveArticles(sqlstore *store.MySQLStore) {
	hns := scraper.HackerNewsScraper{}
	articles, err := hns.Scrape()
	if err != nil {
		log.Printf("Error scraping articles: %v\n", err)
		return
	}

	if err := sqlstore.SaveArticles(articles); err != nil {
		log.Printf("Error saving articles: %v\n", err)
		return
	}

	for _, article := range articles {
		log.Printf("Saved article: %s - %s\n", article.Title, article.Link)
	}
}
