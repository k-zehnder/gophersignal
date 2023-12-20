package main

import (
	"fmt"
	"log"

	"github.com/k-zehnder/gophersignal/backend/pkg/hackernews"
)

func main() {
	hns := hackernews.HackerNewsScraper{}
	articles, err := hns.Scrape()
	if err != nil {
		log.Fatal(err)
	}
	for i, article := range articles {
		msg := fmt.Sprintf("[%d] %s - %s", i, article.Title, article.Link)
		fmt.Println(msg)
	}
}
