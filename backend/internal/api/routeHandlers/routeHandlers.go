package routeHandlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
// @Description Retrieve a list of articles from the database.
// @Tags Articles
// @Accept json
// @Produce json
// @Param is_on_homepage query boolean false "Filter by is_on_homepage (default: false)"
// @Success 200 {array} models.ArticleResponse
// @Failure 500 {object} models.Response{data=string}
// @Router /articles [get]
func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	isOnHomepage := r.URL.Query().Get("is_on_homepage")

	isOnHomepageBool, err := strconv.ParseBool(isOnHomepage)
	if err != nil {
		h.jsonResponse(w, models.Response{Code: http.StatusBadRequest, Status: "error", Data: "Invalid is_on_homepage parameter"}, http.StatusBadRequest)
		return
	}

	articles, err := h.Store.GetArticles(isOnHomepageBool)
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
