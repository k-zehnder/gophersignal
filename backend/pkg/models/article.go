package models

import "time"

// {allPostsData.map(({ id, date, category, title, summary }) => (

// Article represents a generic article structure
type Article struct {
	Title     string
	Link      string
	Source    string
	ScrapedAt time.Time
}

// NewArticle creates a new Article instance
func NewArticle() {}
