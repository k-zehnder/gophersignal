// Package models defines the core data structures used within the GopherSignal application.
package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// NullableInt is a custom type for nullable integers to improve Swagger compatibility.
type NullableInt struct {
	sql.NullInt64
}

// NewNullableInt returns a new NullableInt with the given value marked as valid.
func NewNullableInt(value int64) NullableInt {
	return NullableInt{sql.NullInt64{Int64: value, Valid: true}}
}

// MarshalJSON customizes JSON output for NullableInt.
func (n NullableInt) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON customizes JSON input for NullableInt.
func (n *NullableInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(data, &n.Int64)
	n.Valid = (err == nil)
	return err
}

// NullableString is a custom type for nullable strings to improve Swagger compatibility.
type NullableString struct {
	sql.NullString
}

// MarshalJSON customizes JSON output for NullableString.
func (n NullableString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON customizes JSON input for NullableString.
func (n *NullableString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(data, &n.String)
	n.Valid = (err == nil)
	return err
}

// Article represents the structure for an article in the system.
type Article struct {
	ID           int            `json:"id"`                                  // Unique identifier for the article.
	Title        string         `json:"title"`                               // Title of the article.
	Link         string         `json:"link"`                                // URL link to the article.
	ArticleRank  int            `json:"article_rank"`                        // Article rank extracted from the source.
	Content      string         `json:"content"`                             // Full content of the article.
	Summary      NullableString `json:"summary" swaggertype:"string"`        // Summary of the article, nullable.
	Source       string         `json:"source"`                              // Source from where the article was fetched.
	CreatedAt    time.Time      `json:"created_at"`                          // Timestamp when the article was created.
	UpdatedAt    time.Time      `json:"updated_at"`                          // Timestamp when the article was last updated.
	Upvotes      NullableInt    `json:"upvotes" swaggertype:"integer"`       // Upvote count.
	CommentCount NullableInt    `json:"comment_count" swaggertype:"integer"` // Number of comments.
	CommentLink  NullableString `json:"comment_link" swaggertype:"string"`   // Link to the comment thread (if any).
	Flagged      bool           `json:"flagged"`                             // Whether the article is flagged.
	Dead         bool           `json:"dead"`                                // Whether the article is dead.
	Dupe         bool           `json:"dupe"`                                // Whether the article is marked as duplicate.
}

// NewArticle is a constructor for creating a new Article instance.
func NewArticle(
	id int,
	title, link string,
	articleRank int,
	content, summary, source string,
	createdAt, updatedAt time.Time,
	upvotes int64,
	commentCount int64,
	commentLink string,
	flagged, dead, dupe bool,
) *Article {
	return &Article{
		ID:           id,
		Title:        title,
		Link:         link,
		ArticleRank:  articleRank,
		Content:      content,
		Summary:      NullableString{NullString: sql.NullString{String: summary, Valid: summary != ""}},
		Source:       source,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Upvotes:      NullableInt{NullInt64: sql.NullInt64{Int64: upvotes, Valid: true}},
		CommentCount: NullableInt{NullInt64: sql.NullInt64{Int64: commentCount, Valid: true}},
		CommentLink:  NullableString{NullString: sql.NullString{String: commentLink, Valid: commentLink != ""}},
		Flagged:      flagged,
		Dead:         dead,
		Dupe:         dupe,
	}
}

// ArticlesResponse represents the response structure for a list of articles.
type ArticlesResponse struct {
	Code       int        `json:"code"`        // HTTP status code.
	Status     string     `json:"status"`      // Status message.
	TotalCount int        `json:"total_count"` // Total count of articles.
	Articles   []*Article `json:"articles"`    // List of articles.
}

// ErrorResponse represents the format for API error responses.
type ErrorResponse struct {
	Code    int    `json:"code"`    // HTTP status code.
	Status  string `json:"status"`  // Error status message.
	Message string `json:"message"` // Detailed error message.
}
