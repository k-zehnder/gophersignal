package models

import "time"

// Response is the standard format for API responses.
// swagger:models.Response
type Response struct {
	Code   int         `json:"code"`           // The HTTP status code of the response (default: 200)
	Status string      `json:"status"`         // The status message accompanying the code (default: "success")
	Data   interface{} `json:"data,omitempty"` // The data payload of the response
}

// ArticleResponse represents an article with detailed information.
// swagger:ArticleResponse
type ArticleResponse struct {
	ID           int       `json:"id"`             // Unique identifier of the article (default: 0)
	Title        string    `json:"title"`          // Title of the article (default: "")
	Content      string    `json:"content"`        // Full content of the article (default: "")
	Link         string    `json:"link"`           // External link to the article (default: "")
	Summary      string    `json:"summary"`        // Brief summary of the article (default: "")
	Source       string    `json:"source"`         // Source from where the article was obtained (default: "")
	IsOnHomepage bool      `json:"is_on_homepage"` // Flag indicating if the article is displayed on the homepage (default: false)
	CreatedAt    time.Time `json:"created_at"`     // Timestamp of when the article was created (default: current time)
	UpdatedAt    time.Time `json:"updated_at"`     // Timestamp of the last update to the article (default: current time)
}
