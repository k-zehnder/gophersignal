// Package routeHandlers provides specialized HTTP handler functions for the GopherSignal application routes.
// Includes a struct 'ArticlesHandler' to handle article-related database interactions.

package routeHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// Handler interface for HTTP handlers.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// ArticlesHandler struct holds a reference to the store to interact with the database.
type ArticlesHandler struct {
	Store store.Store
}

// NewArticlesHandler initializes a new ArticlesHandler with the given store.
func NewArticlesHandler(store store.Store) *ArticlesHandler {
	return &ArticlesHandler{Store: store}
}

// ServeHTTP implements the HTTP request routing for the ArticlesHandler.
func (h *ArticlesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Route different HTTP requests to appropriate methods
	if r.Method == "GET" {
		h.GetArticles(w, r)
	} else {
		// Respond with Method Not Allowed for non-GET requests
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
func (h *ArticlesHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.Store.GetArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cache headers for the response for 600 seconds (i.e., 10 minutes).
	h.setCacheHeaders(w, 600)

	// Send a JSON response with the retrieved articles and HTTP status OK (200).
	h.jsonResponse(w, articles, http.StatusOK)

}

// jsonResponse sends a JSON response with the given data and HTTP status code.
func (h *ArticlesHandler) jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// setCacheHeaders sets cache-related headers in the HTTP response.
func (h *ArticlesHandler) setCacheHeaders(w http.ResponseWriter, maxAgeSeconds int) {
	cacheControl := fmt.Sprintf("public, max-age=%d", maxAgeSeconds)
	w.Header().Set("Cache-Control", cacheControl)
}
