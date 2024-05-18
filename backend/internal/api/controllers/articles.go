// Package controllers handles HTTP requests for the GopherSignal application.
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// Controller defines the interface for controllers that handle HTTP requests.
// It extends the http.Handler interface by expecting implementation of ServeHTTP.
type Controller interface {
	http.Handler
}

// ArticlesController manages article-related HTTP requests.
type ArticlesController struct {
	Store store.Store // Store provides access to the data layer.
}

// NewArticlesController creates a new ArticlesController with the provided store.
func NewArticlesController(store store.Store) *ArticlesController {
	return &ArticlesController{Store: store}
}

// ServeHTTP routes HTTP requests to the appropriate controller methods.
func (h *ArticlesController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.GetArticles(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetArticles handles the HTTP request to retrieve articles.
// @Summary Get articles
// @Description Retrieve a list of articles from the database.
// @Tags Articles
// @Accept json
// @Produce json
// @Success 200 {object} models.ArticlesResponse "Articles data"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /articles [get]
func (h *ArticlesController) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.Store.GetArticles()
	if err != nil {
		h.jsonErrorResponse(w, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	h.jsonResponse(w, models.ArticlesResponse{
		Code:       http.StatusOK,
		Status:     "success",
		TotalCount: len(articles),
		Articles:   articles,
	}, http.StatusOK)
}

// jsonResponse sends a JSON response with the specified data and status code.
func (h *ArticlesController) jsonResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// jsonErrorResponse sends a JSON error response with the specified data and status code.
func (h *ArticlesController) jsonErrorResponse(w http.ResponseWriter, response models.ErrorResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
