// Package models defines the data models used in the system.
// It includes the Article struct and a constructor function for creating Article instances.

package models

import (
	"database/sql"
	"time"
)

// Article represents the structure for an article in the system.
type Article struct {
	ID        int            `json:"id"`         // Unique identifier for the article.
	Title     string         `json:"title"`      // Title of the article.
	Link      string         `json:"link"`       // URL link to the article.
	Content   string         `json:"content"`    // Full content of the article.
	Summary   sql.NullString `json:"summary"`    // Summary of the article, nullable.
	Source    string         `json:"source"`     // Source from where the article was fetched.
	CreatedAt time.Time      `json:"created_at"` // Timestamp when the article was created.
	UpdatedAt time.Time      `json:"updated_at"` // Timestamp when the article was last updated.
}

// NewArticle is a constructor for creating a new Article instance.
func NewArticle(id int, title, link, content, summary, source string, createdAt, updatedAt time.Time) *Article {
	var summaryNullString sql.NullString
	// Check if a summary is provided, and create a sql.NullString accordingly.
	if summary != "" {
		summaryNullString = sql.NullString{String: summary, Valid: true}
	} else {
		summaryNullString = sql.NullString{Valid: false}
	}

	return &Article{
		ID:        id,
		Title:     title,
		Link:      link,
		Content:   content,
		Summary:   summaryNullString,
		Source:    source,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
