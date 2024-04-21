// Package routeHandlers provides HTTP handler functions for the routes of the GopherSignal application.
// It defines a Handler struct that holds a reference to the store for database interactions.
package routeHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// Handler struct holds a reference to the store to interact with the database.
type Handler struct {
	Store store.Store
}

// NewHandler initializes a new Handler with the given store.
func NewHandler(store store.Store) *Handler {
	return &Handler{Store: store}
}

// GetArticles handles the HTTP request to retrieve articles. It fetches articles from the store and sends a JSON response.
// @Summary Get articles
// @Description Retrieve a list of articles from the database.
// @Tags Articles
// @Accept json
// @Produce json
// @Success 200 {array} models.ArticleResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /articles [get]
func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.Store.GetArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cache-related headers for the response with a 30-minute cache duration
	setCacheHeaders(w, 1800)

	// Return articles as JSON response
	h.jsonResponse(w, articles, http.StatusOK)
}

// jsonResponse sends a JSON response with the given data and HTTP status code.
func (h *Handler) jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// setCacheHeaders sets cache-related headers in the HTTP response with the specified max age in seconds.
func setCacheHeaders(w http.ResponseWriter, maxAgeSeconds int) {
	// Set Cache-Control header to enable caching for the specified duration
	cacheControl := fmt.Sprintf("public, max-age=%d", maxAgeSeconds)
	w.Header().Set("Cache-Control", cacheControl)
}
