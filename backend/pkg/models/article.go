package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Article represents a generic article structure
type Article struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	Link         string         `json:"link"`
	Content      string         `json:"content"`
	Summary      sql.NullString `json:"summary"`
	Source       string         `json:"source"`
	ScrapedAt    Timestamp      `json:"scrapedAt"`
	IsOnHomepage bool           `json:"isOnHomepage"`
}

// NewArticle creates a new Article instance with a Summary field
func NewArticle(id int, title, link, content, summary, source string, scrapedAt time.Time, isOnHomepage bool) *Article {
	var summaryNullString sql.NullString
	if summary != "" {
		summaryNullString = sql.NullString{String: summary, Valid: true}
	} else {
		summaryNullString = sql.NullString{Valid: false}
	}

	return &Article{
		ID:           id,
		Title:        title,
		Link:         link,
		Content:      content,
		Summary:      summaryNullString,
		Source:       source,
		ScrapedAt:    Timestamp{Time: scrapedAt},
		IsOnHomepage: isOnHomepage,
	}
}

// Timestamp is a custom type for handling time.Time in MySQL
type Timestamp struct {
	time.Time
}

// Scan implements the sql.Scanner interface for Timestamp
func (t *Timestamp) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}
	if bytes, ok := value.([]byte); ok {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(bytes))
		if err != nil {
			return err
		}
		t.Time = parsedTime
		return nil
	}
	return fmt.Errorf("unable to convert %v to Timestamp", value)
}
