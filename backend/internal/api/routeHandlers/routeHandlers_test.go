// Package routeHandlers provides HTTP handler functions for the routes of the GopherSignal application.
package routeHandlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

// TestGetArticles_Success tests the GetArticles handler for a successful scenario.
func TestGetArticles_Success(t *testing.T) {
	// Set up a mock store with predefined data to simulate database interactions.
	mockstore := &store.MockStore{
		Articles: []*models.Article{{Title: "Test Article 1"}},
	}

	// Initialize the handler with the mock store.
	handler := NewHandler(mockstore)

	// Set up the router and associate the handler with the GET method for the articles endpoint.
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")

	// Create a new HTTP GET request for the articles endpoint.
	req, err := http.NewRequest("GET", "/api/v1/articles", nil)
	if err != nil {
		t.Fatal(err) // Fail the test if request creation fails.
	}

	// Use httptest to record the HTTP response.
	rr := httptest.NewRecorder()

	// Serve the HTTP request and record the response.
	r.ServeHTTP(rr, req)

	// Check if the response status code is as expected (200 OK).
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestGetArticles_StoreError tests the GetArticles handler for a scenario where the store returns an error.
func TestGetArticles_StoreError(t *testing.T) {
	// Set up a mock store to simulate an error scenario.
	mockStore := store.NewMockStore(nil, nil, errors.New("store error"))

	// Initialize the handler with the mock store.
	handler := NewHandler(mockStore)

	// Create a new HTTP GET request for the articles endpoint.
	req, err := http.NewRequest("GET", "/api/v1/articles", nil)
	if err != nil {
		t.Fatal(err) // Fail the test if request creation fails.
	}

	// Use httptest to record the HTTP response.
	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.GetArticles)

	// Serve the HTTP request and record the response.
	handlerFunc.ServeHTTP(rr, req)

	// Check if the response status code is as expected (500 Internal Server Error) for an error scenario.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check if the response body contains the expected error message.
	expected := "store error\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
