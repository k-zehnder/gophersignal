package models

type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
	Summary string `json:"summary"`
	Source  string `json:"source"`
}

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type ArticleResponse struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
	Summary string `json:"summary"`
	Source  string `json:"source"`
}

type UpdateArticleRequest struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
