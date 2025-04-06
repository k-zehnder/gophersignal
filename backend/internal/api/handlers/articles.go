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
// @Summary Get filtered articles
// @Description Retrieve paginated articles with optional filters and thresholds
// @Tags Articles
// @Accept  json
// @Produce  json
// @Param   flagged       query   boolean  false  "Filter by flagged status"
// @Param   dead          query   boolean  false  "Filter by dead status"
// @Param   dupe          query   boolean  false  "Filter by duplicate status"
// @Param   limit         query   integer  false  "Results per page (max 100)"  default(30) minimum(1) maximum(100)
// @Param   offset        query   integer  false  "Pagination offset"            default(0) minimum(0)
// @Param   min_upvotes   query   integer  false  "Minimum upvotes threshold"    default(0) minimum(0) format(int64)
// @Param   min_comments  query   integer  false  "Minimum comments threshold"   default(0) minimum(0) format(int64)
// @Success 200 {object} models.ArticlesResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /articles [get]
func (h *ArticlesHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Parse optional boolean filters.
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

	// Extract optional threshold parameters.
	// A value of 0 means no threshold filtering.
	minUpvotes := 0
	minComments := 0
	if upStr := q.Get("min_upvotes"); upStr != "" {
		if up, err := strconv.Atoi(upStr); err == nil {
			minUpvotes = up
		}
	}
	if commStr := q.Get("min_comments"); commStr != "" {
		if comm, err := strconv.Atoi(commStr); err == nil {
			minComments = comm
		}
	}

	var (
		articles []*models.Article
		err      error
	)

	// Determine which store method to call based on provided filters.
	if (minUpvotes > 0 || minComments > 0) && (flagged != nil || dead != nil || dupe != nil) {
		articles, err = h.Store.GetArticlesWithThresholdsAndFilters(limit, offset, minUpvotes, minComments, flagged, dead, dupe)
	} else if minUpvotes > 0 || minComments > 0 {
		articles, err = h.Store.GetArticlesWithThresholds(limit, offset, minUpvotes, minComments)
	} else if flagged != nil || dead != nil || dupe != nil {
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
