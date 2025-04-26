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
		1,
		"Test Title",
		"https://example.com",
		1,
		"Full content here.",
		"Short summary.",
		"Hacker News",
		"commit123",
		"test-model",
		createdAt,
		updatedAt,
		100,
		50,
		"https://news.ycombinator.com/item?id=1",
		false,
		false,
		true,
	)

	assert.Equal(t, 1, article.ID)
	assert.Equal(t, 1, article.HNID)
	assert.Equal(t, "Test Title", article.Title)
	assert.Equal(t, 1, article.ArticleRank)
	assert.True(t, article.Summary.Valid)
	assert.Equal(t, "Short summary.", article.Summary.String)
	assert.True(t, article.CommentLink.Valid)
	assert.Equal(t, "https://news.ycombinator.com/item?id=1", article.CommentLink.String)
	assert.False(t, article.Flagged)
	assert.False(t, article.Dead)
	assert.True(t, article.Dupe)
	assert.Equal(t, "commit123", article.CommitHash)
	assert.Equal(t, "test-model", article.ModelName)
}

// TestNullableIntJSONMarshaling tests JSON marshaling for NullableInt.
func TestNullableIntJSONMarshaling(t *testing.T) {
	validInt := NullableInt{NullInt64: sql.NullInt64{Int64: 42, Valid: true}}
	data, err := json.Marshal(validInt)
	assert.NoError(t, err)
	assert.JSONEq(t, "42", string(data))

	invalidInt := NullableInt{NullInt64: sql.NullInt64{Valid: false}}
	data, err = json.Marshal(invalidInt)
	assert.NoError(t, err)
	assert.JSONEq(t, "null", string(data))
}

// TestNullableStringJSONMarshaling tests JSON marshaling for NullableString.
func TestNullableStringJSONMarshaling(t *testing.T) {
	validString := NullableString{NullString: sql.NullString{String: "test", Valid: true}}
	data, err := json.Marshal(validString)
	assert.NoError(t, err)
	assert.JSONEq(t, `"test"`, string(data))

	invalidString := NullableString{NullString: sql.NullString{Valid: false}}
	data, err = json.Marshal(invalidString)
	assert.NoError(t, err)
	assert.JSONEq(t, "null", string(data))
}

// TestArticleMarshaling tests marshaling and unmarshaling of an Article.
func TestArticleMarshaling(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)

	article := NewArticle(
		2,
		5,
		"Test Article",
		"https://example.com",
		2,
		"Full content here",
		"Summary here",
		"Example Source",
		"commit456",
		"prod-model",
		createdAt,
		updatedAt,
		123,
		45,
		"https://example.com/comments",
		true,
		false,
		false,
	)

	data, err := json.Marshal(article)
	assert.NoError(t, err)

	var result Article
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, article.ID, result.ID)
	assert.Equal(t, article.HNID, result.HNID)
	assert.Equal(t, article.Title, result.Title)
	assert.Equal(t, article.Link, result.Link)
	assert.Equal(t, article.ArticleRank, result.ArticleRank)
	assert.Equal(t, article.Content, result.Content)
	assert.Equal(t, article.Summary.String, result.Summary.String)
	assert.Equal(t, article.CommitHash, result.CommitHash)
	assert.Equal(t, article.ModelName, result.ModelName)
	assert.Equal(t, article.Upvotes.Int64, result.Upvotes.Int64)
	assert.Equal(t, article.CommentCount.Int64, result.CommentCount.Int64)
	assert.Equal(t, article.CommentLink.String, result.CommentLink.String)
	assert.Equal(t, article.Flagged, result.Flagged)
	assert.Equal(t, article.Dead, result.Dead)
	assert.Equal(t, article.Dupe, result.Dupe)
}
