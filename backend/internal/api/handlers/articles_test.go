// Package handlers contains unit tests for HTTP handler functions.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/config"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// TestGetArticles_Success tests the GetArticles handler for a successful response.
func TestGetArticles_Success(t *testing.T) {
	// Set up a mock store with predefined data.
	mockStore := store.NewMockStore([]*models.Article{
		{ID: 1, Title: "Test Article 1"},
	}, nil, nil)

	cfg := config.NewConfig()
	handler := NewArticlesHandler(mockStore, cfg)

	req := httptest.NewRequest("GET", "/dummy-url", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestGetArticles_Error tests the GetArticles handler for an error scenario.
func TestGetArticles_Error(t *testing.T) {
	// Set up a mock store to simulate an error scenario.
	mockStore := store.NewMockStore(nil, nil, errors.New("database error"))

	cfg := config.NewConfig()
	handler := NewArticlesHandler(mockStore, cfg)

	req := httptest.NewRequest("GET", "/dummy-url", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

// TestServeHTTP_MethodNotAllowed tests the ServeHTTP method for non-GET requests.
func TestServeHTTP_MethodNotAllowed(t *testing.T) {
	cfg := config.NewConfig()
	handler := NewArticlesHandler(nil, cfg)

	req := httptest.NewRequest("POST", "/dummy-url", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

// TestGetArticles_WithQueryParams tests the GetArticles handler when boolean query parameters are provided.
func TestGetArticles_WithQueryParams(t *testing.T) {
	// Prepare articles with varying boolean fields.
	articles := []*models.Article{
		{ID: 1, Title: "Article 1", Flagged: false, Dead: false, Dupe: false},
		{ID: 2, Title: "Article 2", Flagged: true, Dead: false, Dupe: false},
		{ID: 3, Title: "Article 3", Flagged: true, Dead: true, Dupe: false},
		{ID: 4, Title: "Article 4", Flagged: false, Dead: false, Dupe: true},
	}
	mockStore := store.NewMockStore(articles, nil, nil)
	cfg := config.NewConfig()
	handler := NewArticlesHandler(mockStore, cfg)

	// Test flagged=true.
	req := httptest.NewRequest("GET", "/dummy-url?flagged=true", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	var resp models.ArticlesResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	// Expecting Articles 2 and 3.
	if resp.TotalCount != 2 {
		t.Errorf("Expected 2 articles for flagged=true, got %d", resp.TotalCount)
	}

	// Test pagination: limit=2, offset=0.
	req = httptest.NewRequest("GET", "/dummy-url?limit=2&offset=0", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.TotalCount != 2 {
		t.Errorf("Expected 2 articles for limit=2, offset=0, got %d", resp.TotalCount)
	}
}

// TestGetArticles_WithThresholdFilters tests the GetArticles handler when threshold query parameters are provided.
func TestGetArticles_WithThresholdFilters(t *testing.T) {
	// Prepare articles with varying upvotes and comment counts.
	articles := []*models.Article{
		{
			ID: 1, Title: "Article 1",
			Upvotes:      models.NewNullableInt(5),
			CommentCount: models.NewNullableInt(2),
			Flagged:      false, Dead: false, Dupe: false,
		},
		{
			ID: 2, Title: "Article 2",
			Upvotes:      models.NewNullableInt(10),
			CommentCount: models.NewNullableInt(5),
			Flagged:      false, Dead: false, Dupe: false,
		},
		{
			ID: 3, Title: "Article 3",
			Upvotes:      models.NewNullableInt(15),
			CommentCount: models.NewNullableInt(8),
			Flagged:      false, Dead: false, Dupe: false,
		},
		{
			ID: 4, Title: "Article 4",
			Upvotes:      models.NewNullableInt(3),
			CommentCount: models.NewNullableInt(1),
			Flagged:      false, Dead: false, Dupe: false,
		},
	}
	mockStore := store.NewMockStore(articles, nil, nil)
	cfg := config.NewConfig()
	handler := NewArticlesHandler(mockStore, cfg)

	// Use thresholds that should filter out articles 1 and 4.
	req := httptest.NewRequest("GET", "/dummy-url?min_upvotes=10&min_comments=5", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}
	var resp models.ArticlesResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	// Expecting Articles 2 and 3.
	if resp.TotalCount != 2 {
		t.Errorf("Expected 2 articles for min_upvotes=10 and min_comments=5, got %d", resp.TotalCount)
	}
}

// TestGetArticles_WithCombinedFilters tests when both thresholds and boolean filters are provided.
func TestGetArticles_WithCombinedFilters(t *testing.T) {
	// Prepare articles with varying fields.
	articles := []*models.Article{
		{
			ID: 1, Title: "Article 1",
			Upvotes:      models.NewNullableInt(10),
			CommentCount: models.NewNullableInt(5),
			Flagged:      true, Dead: false, Dupe: false,
		},
		{
			ID: 2, Title: "Article 2",
			Upvotes:      models.NewNullableInt(20),
			CommentCount: models.NewNullableInt(10),
			Flagged:      true, Dead: false, Dupe: false,
		},
		{
			ID: 3, Title: "Article 3",
			Upvotes:      models.NewNullableInt(5),
			CommentCount: models.NewNullableInt(2),
			Flagged:      true, Dead: false, Dupe: false,
		},
		{
			ID: 4, Title: "Article 4",
			Upvotes:      models.NewNullableInt(25),
			CommentCount: models.NewNullableInt(15),
			Flagged:      false, Dead: false, Dupe: false,
		},
	}
	mockStore := store.NewMockStore(articles, nil, nil)
	cfg := config.NewConfig()
	handler := NewArticlesHandler(mockStore, cfg)

	// Provide both threshold and boolean filter (flagged=true).
	req := httptest.NewRequest("GET", "/dummy-url?flagged=true&min_upvotes=10&min_comments=5", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}
	var resp models.ArticlesResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	// Expected: Only articles that are flagged true and meet thresholds: Articles 1 and 2.
	if resp.TotalCount != 2 {
		t.Errorf("Expected 2 articles for combined filters, got %d", resp.TotalCount)
	}
}
