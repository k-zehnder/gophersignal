// Package server includes the unit test for the NewServer function of the GopherSignal application.
package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// TestNewServer validates the server's behavior by simulating HTTP requests with mock data
// and verifying the responses, thus confirming the accuracy and reliability of the API's outputs.
func TestNewServer(t *testing.T) {
	// Initialize your test articles here
	mockArticles := []*models.Article{
		{
			ID:      1,
			Title:   "Test Article 1",
			Content: "Content of Test Article 1",
		},
		{
			ID:      2,
			Title:   "Test Article 2",
			Content: "Content of Test Article 2",
		},
	}

	// Create a MockStore with the mock data
	mockStore := store.NewMockStore(mockArticles, nil, nil)

	// Initialize the server with the mock store
	handler := NewServer(mockStore)

	// Adjust the request URL to include the API prefix
	req, err := http.NewRequest("GET", "/api/v1/articles", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check for the expected status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse the response body to check if it returns the correct articles
	var articlesResponse models.ArticlesResponse
	err = json.Unmarshal(rr.Body.Bytes(), &articlesResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Assert the content of the response
	if len(articlesResponse.Articles) != len(mockArticles) {
		t.Errorf("Expected %d articles, got %d", len(mockArticles), len(articlesResponse.Articles))
	}
	for i, article := range articlesResponse.Articles {
		if article.ID != mockArticles[i].ID || article.Title != mockArticles[i].Title {
			t.Errorf("Expected article %v, got %v", mockArticles[i], article)
		}
	}
}
