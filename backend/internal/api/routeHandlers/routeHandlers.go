package routeHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

type Handler struct {
	Store store.Store
}

// NewHandler creates and returns a new instance of Handler
func NewHandler(store store.Store) *Handler {
	return &Handler{Store: store}
}

func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.Store.GetArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return articles as JSON response
	h.jsonResponse(w, articles, http.StatusOK)
}

// Utility function to send JSON responses
func (h *Handler) jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
