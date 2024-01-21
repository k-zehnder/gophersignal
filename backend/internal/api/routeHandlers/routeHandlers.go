package routeHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

type Handler struct {
	Store store.Store
}

func NewHandler(store store.Store) *Handler {
	return &Handler{Store: store}
}

// GetArticles godoc
// @Summary Get articles
// @Description Retrieve a list of articles from the store
// @Tags articles
// @Accept json
// @Produce json
// @Success 200 {array} models.Response{data=[]models.ArticleResponse}
// @Failure 500 {object} models.Response{data=string}
// @Router /articles [get]
func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.Store.GetArticles()
	if err != nil {
		h.jsonResponse(w, models.Response{Code: http.StatusInternalServerError, Status: "error", Data: "Failed to retrieve articles from the store"}, http.StatusInternalServerError)
		return
	}

	h.jsonResponse(w, models.Response{Code: http.StatusOK, Status: "success", Data: articles}, http.StatusOK)
}

func (h *Handler) jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}