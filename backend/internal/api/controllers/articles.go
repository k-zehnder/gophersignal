// Package controllers handles HTTP requests for the GopherSignal application.
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// Controller defines the interface for controllers that handle HTTP requests.
// It extends the http.Handler interface by expecting implementation of ServeHTTP.
type Controller interface {
	http.Handler
}

// ArticlesController manages article-related HTTP requests.
type ArticlesController struct {
	Store store.Store
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
// @Success 200 {array} models.ArticleResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /articles [get]
func (h *ArticlesController) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.Store.GetArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.setCacheHeaders(w, 900) // Set cache headers for 15 minutes.
	h.jsonResponse(w, articles, http.StatusOK)
}

// jsonResponse sends a JSON response with the specified data and status code.
func (h *ArticlesController) jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// setCacheHeaders adds caching headers to the HTTP response.
func (h *ArticlesController) setCacheHeaders(w http.ResponseWriter, maxAgeSeconds int) {
	cacheControl := fmt.Sprintf("public, max-age=%d", maxAgeSeconds)
	w.Header().Set("Cache-Control", cacheControl)
}
