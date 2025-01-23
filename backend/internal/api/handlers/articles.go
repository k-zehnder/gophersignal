// Package handlers handles HTTP requests for the GopherSignal application.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// Handler defines the interface for handlers that manage HTTP requests.
type Handler interface {
	http.Handler
}

// ArticlesHandler manages article-related HTTP requests.
type ArticlesHandler struct {
	Store store.Store // Store provides access to the data layer.
}

// NewArticlesHandler creates a new ArticlesHandler with the provided store.
func NewArticlesHandler(store store.Store) *ArticlesHandler {
	return &ArticlesHandler{Store: store}
}

// ServeHTTP routes HTTP requests to the appropriate handler methods.
func (h *ArticlesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.GetArticles(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetArticles handles the HTTP request to retrieve articles.
func (h *ArticlesHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.Store.GetArticles()
	if err != nil {
		h.jsonErrorResponse(w, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	h.setCacheHeaders(w, 900) // Set cache headers for 15 minutes.
	h.jsonResponse(w, models.ArticlesResponse{
		Code:       http.StatusOK,
		Status:     "success",
		TotalCount: len(articles),
		Articles:   articles,
	}, http.StatusOK)
}

// setCacheHeaders adds caching headers to the HTTP response.
func (h *ArticlesHandler) setCacheHeaders(w http.ResponseWriter, maxAgeSeconds int) {
	cacheControl := fmt.Sprintf("public, max-age=%d", maxAgeSeconds)
	w.Header().Set("Cache-Control", cacheControl)
}

// jsonResponse sends a JSON response with the specified data and status code.
func (h *ArticlesHandler) jsonResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// jsonErrorResponse sends a JSON error response with the specified data and status code.
func (h *ArticlesHandler) jsonErrorResponse(w http.ResponseWriter, response models.ErrorResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
