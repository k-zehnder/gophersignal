// Package handlers contains unit tests for HTTP handler functions.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// TestGetArticles_Success tests the GetArticles handler for a successful response.
func TestGetArticles_Success(t *testing.T) {
	// Set up a mock store with predefined data to simulate database interactions.
	// The mock store should implement GetArticles(limit, offset int)
	mockStore := store.NewMockStore([]*models.Article{
		{ID: 1, Title: "Test Article 1"},
	}, nil, nil)

	// Initialize the handler with the mock store.
	handler := NewArticlesHandler(mockStore)

	// Create a new HTTP GET request for the articles endpoint.
	req := httptest.NewRequest("GET", "/dummy-url", nil)
	rr := httptest.NewRecorder()

	// Serve the HTTP request and record the response.
	handler.ServeHTTP(rr, req)

	// Check if the response status code is as expected (200 OK).
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestGetArticles_Error tests the GetArticles handler for an error scenario.
func TestGetArticles_Error(t *testing.T) {
	// Set up a mock store to simulate an error scenario.
	mockStore := store.NewMockStore(nil, nil, errors.New("database error"))

	// Initialize the handler with the mock store.
	handler := NewArticlesHandler(mockStore)

	// Create a new HTTP GET request for the articles endpoint.
	req := httptest.NewRequest("GET", "/dummy-url", nil)
	rr := httptest.NewRecorder()

	// Serve the HTTP request and record the response.
	handler.ServeHTTP(rr, req)

	// Check if the response status code is as expected (500 Internal Server Error).
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

// TestServeHTTP_MethodNotAllowed tests the ServeHTTP method for non-GET requests.
func TestServeHTTP_MethodNotAllowed(t *testing.T) {
	// Initialize the handler with a mock store (can be nil as it's not used in this test).
	handler := NewArticlesHandler(nil)

	// Create a new HTTP POST request (or any non-GET request).
	req := httptest.NewRequest("POST", "/dummy-url", nil)
	rr := httptest.NewRecorder()

	// Serve the HTTP request and record the response.
	handler.ServeHTTP(rr, req)

	// Check if the response status code is as expected (405 Method Not Allowed).
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

// TestGetArticles_WithQueryParams tests the GetArticles handler when query parameters are provided.
func TestGetArticles_WithQueryParams(t *testing.T) {
	// Prepare a list of articles with varying boolean fields.
	articles := []*models.Article{
		{ID: 1, Title: "Article 1", Flagged: false, Dead: false, Dupe: false},
		{ID: 2, Title: "Article 2", Flagged: true, Dead: false, Dupe: false},
		{ID: 3, Title: "Article 3", Flagged: true, Dead: true, Dupe: false},
		{ID: 4, Title: "Article 4", Flagged: false, Dead: false, Dupe: true},
	}
	mockStore := store.NewMockStore(articles, nil, nil)
	handler := NewArticlesHandler(mockStore)

	// Case 1: flagged=true (should return only articles where Flagged == true, i.e. Articles 2 and 3).
	req := httptest.NewRequest("GET", "/dummy-url?flagged=true", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	// Decode the JSON response.
	var resp models.ArticlesResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.TotalCount != 2 {
		t.Errorf("Expected 2 articles for flagged=true, got %d", resp.TotalCount)
	}

	// Case 2: dead=true and dupe=true (should return 0 articles, as no article satisfies both).
	req = httptest.NewRequest("GET", "/dummy-url?dead=true&dupe=true", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.TotalCount != 0 {
		t.Errorf("Expected 0 articles for dead=true and dupe=true, got %d", resp.TotalCount)
	}
}

// TestGetFilteredArticles_Success tests the endpoint when query parameters are provided.
// Since our handler uses a single endpoint (/articles) for both default and filtered queries,
// this test simulates a filtered request.
func TestGetFilteredArticles_Success(t *testing.T) {
	// Prepare a list of articles with varying boolean fields.
	articles := []*models.Article{
		{ID: 1, Title: "Article 1", Flagged: false, Dead: false, Dupe: false},
		{ID: 2, Title: "Article 2", Flagged: true, Dead: false, Dupe: false},
		{ID: 3, Title: "Article 3", Flagged: true, Dead: true, Dupe: false},
		{ID: 4, Title: "Article 4", Flagged: false, Dead: false, Dupe: true},
	}
	mockStore := store.NewMockStore(articles, nil, nil)
	handler := NewArticlesHandler(mockStore)

	// In this test we provide flagged=true.
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
	// In our test data, Articles 2 and 3 have Flagged == true.
	if resp.TotalCount != 2 {
		t.Errorf("Expected 2 articles for flagged=true, got %d", resp.TotalCount)
	}

	// Also test pagination: for example, use limit=2 and offset=0.
	req = httptest.NewRequest("GET", "/dummy-url?limit=2&offset=0", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	// Since our mock store has 4 articles, default GetArticles with limit=30 would return all 4,
	// but with limit=2 we expect only 2 articles.
	if resp.TotalCount != 2 {
		t.Errorf("Expected 2 articles for limit=2, offset=0, got %d", resp.TotalCount)
	}
}
