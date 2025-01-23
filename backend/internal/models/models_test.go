package models

import (
	"testing"
	"time"
)

// TestNewArticle verifies that the NewArticle constructor initializes fields correctly.
func TestNewArticle(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now()

	article := NewArticle(
		1,
		"Test Title",
		"https://example.com",
		"Full content here.",
		"Short summary.",
		"Hacker News",
		createdAt,
		updatedAt,
		100,
		50,
		"https://news.ycombinator.com/item?id=1",
	)

	if article.ID != 1 {
		t.Errorf("Expected ID 1, got %d", article.ID)
	}
	if article.Title != "Test Title" {
		t.Errorf("Expected Title 'Test Title', got '%s'", article.Title)
	}
	if !article.Summary.Valid || article.Summary.String != "Short summary." {
		t.Errorf("Expected Summary 'Short summary.', got '%v'", article.Summary)
	}
	if !article.CommentLink.Valid || article.CommentLink.String != "https://news.ycombinator.com/item?id=1" {
		t.Errorf("Expected CommentLink 'https://news.ycombinator.com/item?id=1', got '%v'", article.CommentLink)
	}
}
