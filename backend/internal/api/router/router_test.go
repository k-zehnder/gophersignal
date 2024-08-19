// Package router contains the unit tests for verifying the router configuration and route handling in the GopherSignal application.
package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/api/handlers"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// TestRouter_ArticlesRoute tests the articles route in the router.
func TestRouter_ArticlesRoute(t *testing.T) {
	// Set up a mock store with no articles to simulate database interaction.
	mockStore := store.NewMockStore([]*models.Article{}, nil, nil)

	// Initialize the ArticlesHandler with the mock store.
	articlesHandler := handlers.NewArticlesHandler(mockStore)

	// Set up the router.
	router := SetupRouter(articlesHandler)

	// Create a new HTTP request to test the articles route.
	req := httptest.NewRequest("GET", "/api/v1/articles", nil)
	rr := httptest.NewRecorder()

	// Serve the request using the router and record the response.
	router.ServeHTTP(rr, req)

	// Check if the response status code is as expected (200 OK).
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
