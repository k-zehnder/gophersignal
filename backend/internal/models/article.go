package models

import (
	"database/sql"
	"time"
)

type Article struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	Link         string         `json:"link"`
	Content      string         `json:"content"`
	Summary      sql.NullString `json:"summary"`
	Source       string         `json:"source"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	IsOnHomepage bool           `json:"is_on_homepage"`
}

func NewArticle(id int, title, link, content, summary, source string, createdAt, updatedAt time.Time, isOnHomepage bool) *Article {
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
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		IsOnHomepage: isOnHomepage,
	}
}
