package main

import (
	"fmt"
	"time"

	article "github.com/k-zehnder/gophersignal/backend/pkg/models"
)

func main() {
	article0 := &article.Article{
		Title:     "Some article title",
		Link:      "https://example.com",
		Source:    "hackernews",
		ScrapedAt: time.Now(),
	}
	fmt.Printf("[x] article0: %v\n", article0)
}
