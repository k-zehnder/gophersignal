package models

import "time"

// Article represents a generic article structure
type Article struct {
	Title     string
	Link      string
	Source    string
	ScrapedAt time.Time
}

// NewArticle creates a new Article instance
func NewArticle(title, link, source string) *Article {
	return &Article{
		Title:     title,
		Link:      link,
		Source:    source,
		ScrapedAt: time.Now(),
	}
}
