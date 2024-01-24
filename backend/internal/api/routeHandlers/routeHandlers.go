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

// @Summary Get articles
// @Description Retrieve a list of articles from the database.
// @Tags Articles
// @Accept json
// @Produce json
// @Param is_on_homepage query boolean false "Filter by is_on_homepage" default(true)
// @Param limit query int false "Maximum number of articles to return" default(100)
// @Success 200 {array} models.ArticleResponse
// @Failure 400 {object} models.Response{data=string} "Invalid Query Parameter"
// @Failure 500 {object} models.Response{data=string} "Internal Server Error"
// @Router /articles [get]
func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	isOnHomepageBool := true
	limit := 30 // Changed default limit to 30

	isOnHomepageParam, ok := r.URL.Query()["is_on_homepage"]
	if ok && len(isOnHomepageParam[0]) > 0 {
		var err error
		isOnHomepageBool, err = strconv.ParseBool(isOnHomepageParam[0])
		if err != nil {
			h.jsonResponse(w, models.Response{Code: http.StatusBadRequest, Status: "error", Data: "Invalid is_on_homepage parameter"}, http.StatusBadRequest)
			return
		}
	}

	limitParam, ok := r.URL.Query()["limit"]
	if ok && len(limitParam[0]) > 0 {
		var err error
		limit, err = strconv.Atoi(limitParam[0])
		if err != nil || limit <= 0 || limit > 100 {
			limit = 100 // Enforce a maximum limit of 100
		}
	}

	articles, err := h.Store.GetArticles(isOnHomepageBool, limit)
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
