// Package controllers contains unit tests for HTTP controller functions.
package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// TestGetArticles_Success tests the GetArticles controller for a successful response.
func TestGetArticles_Success(t *testing.T) {
	// Set up a mock store with predefined data to simulate database interactions.
	mockStore := store.NewMockStore([]*models.Article{{Title: "Test Article 1"}}, nil, nil)

	// Initialize the controller with the mock store.
	controller := NewArticlesController(mockStore)

	// Create a new HTTP GET request for the articles endpoint.
	req := httptest.NewRequest("GET", "/dummy-url", nil)
	rr := httptest.NewRecorder()

	// Serve the HTTP request and record the response.
	controller.ServeHTTP(rr, req)

	// Check if the response status code is as expected (200 OK).
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("controller returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestGetArticles_Error tests the GetArticles controller for an error scenario.
func TestGetArticles_Error(t *testing.T) {
	// Set up a mock store to simulate an error scenario.
	mockStore := store.NewMockStore(nil, nil, errors.New("database error"))

	// Initialize the controller with the mock store.
	controller := NewArticlesController(mockStore)

	// Create a new HTTP GET request for the articles endpoint.
	req := httptest.NewRequest("GET", "/dummy-url", nil)
	rr := httptest.NewRecorder()

	// Serve the HTTP request and record the response.
	controller.ServeHTTP(rr, req)

	// Check if the response status code is as expected (500 Internal Server Error).
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("controller returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

// TestServeHTTP_MethodNotAllowed tests the ServeHTTP method for non-GET requests.
func TestServeHTTP_MethodNotAllowed(t *testing.T) {
	// Initialize the controller with a mock store (can be nil as it's not used in this test).
	controller := NewArticlesController(nil)

	// Create a new HTTP POST request (or any non-GET request).
	req := httptest.NewRequest("POST", "/dummy-url", nil)
	rr := httptest.NewRecorder()

	// Serve the HTTP request and record the response.
	controller.ServeHTTP(rr, req)

	// Check if the response status code is as expected (405 Method Not Allowed).
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("controller returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}
