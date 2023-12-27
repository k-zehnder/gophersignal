package myhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/pkg/store"
)

func GetArticlesHandler(w http.ResponseWriter, r *http.Request, store *store.DBStore) {
	articles, err := store.GetAllArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return articles as JSON response
	jsonResponse(w, articles, http.StatusOK)
}

// Utility function to send JSON responses
func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
