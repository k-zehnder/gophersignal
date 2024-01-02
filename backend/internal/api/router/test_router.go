package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/api/routeHandlers"
)

func TestRouter_ArticlesRoute(t *testing.T) {
	// Initialize mock handler.
	handler := &routeHandlers.Handler{}

	// Setup router with the handler.
	router := SetupRouter(handler)

	// Simulate a GET request.
	req := httptest.NewRequest("GET", "/api/v1/articles", nil)
	rr := httptest.NewRecorder()

	// Serve the request.
	router.ServeHTTP(rr, req)

	// Assert the response status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	// Validate existence of the route.
	if route := router.GetRoute("/api/v1/articles"); route == nil {
		t.Error("Expected /api/v1/articles route to exist")
	}
}
