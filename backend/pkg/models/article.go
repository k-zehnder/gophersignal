package models

import "time"

// Article represents a generic article structure
type Article struct {
	Title     string
	Link      string
	Source    string
	ScrapedAt time.Time
}
