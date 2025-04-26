// Package models defines the core data structures used within the GopherSignal application.
package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// NullableInt is a custom type for nullable integers.
// swagger:model NullableInt
// x-nullable: true
// type: integer
// format: int64
type NullableInt struct {
	sql.NullInt64 `swaggerignore:"true"`
}

// NewNullableInt returns a new NullableInt with the given value.
func NewNullableInt(value int64) NullableInt {
	return NullableInt{sql.NullInt64{Int64: value, Valid: true}}
}

// MarshalJSON ensures numeric output or null.
func (n NullableInt) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON handles null or numeric input.
func (n *NullableInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(data, &n.Int64)
	n.Valid = (err == nil)
	return err
}

// NullableString is a custom type for nullable strings.
// swagger:model NullableString
// x-nullable: true
// type: string
type NullableString struct {
	sql.NullString `swaggerignore:"true"`
}

// MarshalJSON ensures string or null output.
func (n NullableString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON handles null or string input.
func (n *NullableString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(data, &n.String)
	n.Valid = (err == nil)
	return err
}

// Article represents an article in the system.
type Article struct {
	ID           int            `json:"id"`
	HNID         int            `json:"hn_id"`
	Title        string         `json:"title"`
	Link         string         `json:"link"`
	ArticleRank  int            `json:"article_rank"`
	Content      string         `json:"content"`
	Summary      NullableString `json:"summary"`
	Source       string         `json:"source"`
	CommitHash   string         `json:"commit_hash"`
	ModelName    string         `json:"model_name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Upvotes      NullableInt    `json:"upvotes"`
	CommentCount NullableInt    `json:"comment_count"`
	CommentLink  NullableString `json:"comment_link"`
	Flagged      bool           `json:"flagged"`
	Dead         bool           `json:"dead"`
	Dupe         bool           `json:"dupe"`
}

// NewArticle constructs a new Article.
func NewArticle(
	id int,
	hnID int,
	title, link string,
	articleRank int,
	content, summary, source string,
	commitHash, modelName string,
	createdAt, updatedAt time.Time,
	upvotes int64,
	commentCount int64,
	commentLink string,
	flagged, dead, dupe bool,
) *Article {
	return &Article{
		ID:           id,
		HNID:         hnID,
		Title:        title,
		Link:         link,
		ArticleRank:  articleRank,
		Content:      content,
		Summary:      NullableString{NullString: sql.NullString{String: summary, Valid: summary != ""}},
		Source:       source,
		CommitHash:   commitHash,
		ModelName:    modelName,
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

// ArticlesResponse represents the response for multiple articles.
type ArticlesResponse struct {
	Code       int        `json:"code"`        // HTTP status code
	Status     string     `json:"status"`      // Response status message
	TotalCount int        `json:"total_count"` // Total number of articles
	Articles   []*Article `json:"articles"`    // List of articles
}

// ErrorResponse represents the error response format.
type ErrorResponse struct {
	Code    int    `json:"code"`    // HTTP status code
	Status  string `json:"status"`  // Error status message
	Message string `json:"message"` // Detailed error message
}
