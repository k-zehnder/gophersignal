package models

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestNullableIntJSONMarshaling(t *testing.T) {
	// Valid integer
	validInt := NullableInt{NullInt64: sql.NullInt64{Int64: 42, Valid: true}}
	data, err := json.Marshal(validInt)
	assert.NoError(t, err)
	assert.JSONEq(t, "42", string(data))

	// Null value
	invalidInt := NullableInt{NullInt64: sql.NullInt64{Valid: false}}
	data, err = json.Marshal(invalidInt)
	assert.NoError(t, err)
	assert.JSONEq(t, "null", string(data))
}

func TestNullableStringJSONMarshaling(t *testing.T) {
	// Valid string
	validString := NullableString{NullString: sql.NullString{String: "test", Valid: true}}
	data, err := json.Marshal(validString)
	assert.NoError(t, err)
	assert.JSONEq(t, `"test"`, string(data))

	// Null value
	invalidString := NullableString{NullString: sql.NullString{Valid: false}}
	data, err = json.Marshal(invalidString)
	assert.NoError(t, err)
	assert.JSONEq(t, "null", string(data))
}

func TestArticleMarshaling(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)

	article := NewArticle(
		1,
		"Test Article",
		"https://example.com",
		"Full content here",
		"Summary here",
		"Example Source",
		createdAt,
		updatedAt,
		123,
		45,
		"https://example.com/comments",
	)

	data, err := json.Marshal(article)
	assert.NoError(t, err)

	var result Article
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, article.ID, result.ID)
	assert.Equal(t, article.Title, result.Title)
	assert.Equal(t, article.Link, result.Link)
	assert.Equal(t, article.Content, result.Content)
	assert.Equal(t, article.Summary.String, result.Summary.String)
	assert.Equal(t, article.Upvotes.Int64, result.Upvotes.Int64)
	assert.Equal(t, article.CommentCount.Int64, result.CommentCount.Int64)
	assert.Equal(t, article.CommentLink.String, result.CommentLink.String)
}
