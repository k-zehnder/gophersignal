// Package main for the article saving tool.
// It initializes a database connection, scrapes articles from HackerNews,
// and saves them into a MySQL database.
package main

import (
	"log"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/internal/scraper"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// main is the entry point for the command-line tool that connects to a MySQL database, retrieves articles,
// and saves them into the database.
func main() {
	// Initialize the application configuration.
	appConfig := config.NewConfig()

	// Validate the database connection string.
	if appConfig.DataSourceName == "" {
		log.Fatal("MYSQL_DSN not set in .env file")
	}

	// Connect to the MySQL database.
	sqlstore, err := store.NewMySQLStore(appConfig.DataSourceName)
	if err != nil {
		log.Fatal("Error initializing DBStore:", err)
	}

	// Scrape and save articles into the database.
	SaveArticles(sqlstore)
}

// SaveArticles handles scraping articles from HackerNews and storing them in the database.
// It reports on the success or failure of these tasks.
func SaveArticles(sqlstore *store.MySQLStore) {
	// Initialize the scraper for HackerNews.
	hns := scraper.HackerNewsScraper{}

	// Scrape articles from HackerNews.
	articles, err := hns.Scrape()
	if err != nil {
		log.Printf("Error scraping articles: %v\n", err)
		return
	}

	// Save the scraped articles to the MySQL database.
	if err := sqlstore.SaveArticles(articles); err != nil {
		log.Printf("Error saving articles: %v\n", err)
		return
	}

	// Log the successful saving of each article.
	for _, article := range articles {
		log.Printf("Saved article: %s - %s\n", article.Title, article.Link)
	}
}
