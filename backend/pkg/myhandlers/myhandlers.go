package myhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/pkg/store"
)

func GetArticlesHandler(w http.ResponseWriter, r *http.Request, dbStore *store.DBStore) {
	articles, err := dbStore.GetAllArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}
