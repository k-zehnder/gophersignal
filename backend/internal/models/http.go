// Package models defines the core data structures used within the GopherSignal application.

package models

import "time"

// Response is the standard format for API responses.
// swagger:response Response
type Response struct {
	Code   int         `json:"code"`           // The HTTP status code of the response (default: 200)
	Status string      `json:"status"`         // The status message accompanying the code (default: "success")
	Data   interface{} `json:"data,omitempty"` // The data payload of the response
}

// ArticleResponse represents an article with detailed information.
// swagger:response ArticleResponse
type ArticleResponse struct {
	ID        int       `json:"id"`         // Unique identifier of the article (default: 0), example: 1
	Title     string    `json:"title"`      // Title of the article (default: ""), example: "Sample Title"
	Content   string    `json:"content"`    // Full content of the article (default: ""), example: "Sample content..."
	Link      string    `json:"link"`       // External link to the article (default: ""), example: "https://example.com"
	Summary   string    `json:"summary"`    // Brief summary of the article (default: ""), example: "This is a sample summary."
	Source    string    `json:"source"`     // Source from where the article was obtained (default: ""), example: "Sample Source"
	CreatedAt time.Time `json:"created_at"` // Timestamp of when the article was created (default: current time), example: "2022-01-01T12:00:00Z"
	UpdatedAt time.Time `json:"updated_at"` // Timestamp of the last update to the article (default: current time), example: "2022-01-01T12:30:00Z"
}
