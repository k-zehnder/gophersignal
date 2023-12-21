package main

import (
	"log"
	"os"

	"github.com/k-zehnder/gophersignal/backend/pkg/hackernews"
	"github.com/k-zehnder/gophersignal/backend/pkg/store"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
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
		log.Fatal(err)
	}

	dbStore.SaveArticles(articles)

	for _, article := range articles {
		log.Printf("Saved article: %s - %s\n", article.Title, article.Link)
	}
}