// Package handlers handles HTTP requests for the GopherSignal application.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// ArticlesHandler manages article-related HTTP requests.
type ArticlesHandler struct {
	Store  store.Store       // Store provides access to the data layer.
	Config *config.AppConfig // Config provides application configuration.
}

// NewArticlesHandler creates a new ArticlesHandler with the provided store and configuration.
func NewArticlesHandler(store store.Store, cfg *config.AppConfig) *ArticlesHandler {
	return &ArticlesHandler{
		Store:  store,
		Config: cfg,
	}
}

// ServeHTTP routes HTTP requests to the appropriate handler methods.
func (h *ArticlesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetArticles(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetArticles handles the HTTP request to retrieve articles.
//
// @Summary Get articles
// @Description Retrieve a list of articles from the database. Optional query parameters
// can be provided to filter results by flagged, dead, and dupe statuses. Additionally,
// pagination parameters `limit` and `offset` are supported.
// @Tags Articles
// @Accept json
// @Produce json
// @Param flagged query bool false "Filter by flagged status"
// @Param dead query bool false "Filter by dead status"
// @Param dupe query bool false "Filter by duplicate status"
// @Param limit query int false "Limit the number of articles returned (default is 30)" default(30)
// @Param offset query int false "Offset for pagination (default is 0)" default(0)
// @Success 200 {object} models.ArticlesResponse "Articles data"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /articles [get]
func (h *ArticlesHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var flagged *bool
	if flaggedStr := q.Get("flagged"); flaggedStr != "" {
		f, err := strconv.ParseBool(flaggedStr)
		if err != nil {
			h.jsonErrorResponse(w, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: fmt.Sprintf("Invalid flagged parameter: %s", flaggedStr),
			}, http.StatusBadRequest)
			return
		}
		flagged = &f
	}

	var dead *bool
	if deadStr := q.Get("dead"); deadStr != "" {
		d, err := strconv.ParseBool(deadStr)
		if err != nil {
			h.jsonErrorResponse(w, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: fmt.Sprintf("Invalid dead parameter: %s", deadStr),
			}, http.StatusBadRequest)
			return
		}
		dead = &d
	}

	var dupe *bool
	if dupeStr := q.Get("dupe"); dupeStr != "" {
		d, err := strconv.ParseBool(dupeStr)
		if err != nil {
			h.jsonErrorResponse(w, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: fmt.Sprintf("Invalid dupe parameter: %s", dupeStr),
			}, http.StatusBadRequest)
			return
		}
		dupe = &d
	}

	// Extract pagination parameters.
	limit := 30
	offset := 0
	if lStr := q.Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil {
			limit = l
		}
	}
	if oStr := q.Get("offset"); oStr != "" {
		if o, err := strconv.Atoi(oStr); err == nil {
			offset = o
		}
	}

	var (
		articles []*models.Article
		err      error
	)
	// If any filtering query parameter is provided, use the filtered query.
	// Otherwise, use the default query.
	if flagged != nil || dead != nil || dupe != nil {
		articles, err = h.Store.GetFilteredArticles(flagged, dead, dupe, limit, offset)
	} else {
		articles, err = h.Store.GetArticles(limit, offset)
	}

	if err != nil {
		h.jsonErrorResponse(w, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	h.setCacheHeaders(w, h.Config.CacheMaxAge) 
	h.jsonResponse(w, models.ArticlesResponse{
		Code:       http.StatusOK,
		Status:     "success",
		TotalCount: len(articles),
		Articles:   articles,
	}, http.StatusOK)
}

func (h *ArticlesHandler) setCacheHeaders(w http.ResponseWriter, maxAgeSeconds int) {
	cacheControl := fmt.Sprintf("public, max-age=%d", maxAgeSeconds)
	w.Header().Set("Cache-Control", cacheControl)
}

func (h *ArticlesHandler) jsonResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *ArticlesHandler) jsonErrorResponse(w http.ResponseWriter, response models.ErrorResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
