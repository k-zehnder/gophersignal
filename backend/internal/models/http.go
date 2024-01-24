package models

import "time"

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type ArticleResponse struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Link         string    `json:"link"`
	Summary      string    `json:"summary"`
	Source       string    `json:"source"`
	IsOnHomepage bool      `json:"is_on_homepage"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
