// Package router provides unit tests for the router configuration and route handling.
package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// TestRouter_ArticlesRoute tests the articles route in the router.
func TestRouter_ArticlesRoute(t *testing.T) {
	// Set up a mock store with no articles to simulate database interaction.
	mockStore := store.NewMockStore([]*models.Article{}, nil, nil)

	// Initialize the Handler with the mock store.
	handler := &routeHandlers.Handler{
		Store: mockStore,
	}

	// Set up the router with the initialized handler.
	router := SetupRouter(handler)

	// Create a new HTTP request to test the articles route.
	// This request simulates a GET request to the articles route.
	req := httptest.NewRequest("GET", "/api/v1/articles", nil)
	rr := httptest.NewRecorder()

	// Serve the request using the router and record the response.
	router.ServeHTTP(rr, req)

	// Check if the response status code is as expected (200 OK).
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Validate that the 'GetArticles' route exists in the router.
	// This ensures that the route is correctly registered and can be retrieved.
	route := router.GetRoute("GetArticles")
	if route == nil {
		t.Error("Expected GetArticles route to exist")
	} else {
		// Log the route name for verification.
		t.Logf("Route found: %v", route.GetName())
	}
}
