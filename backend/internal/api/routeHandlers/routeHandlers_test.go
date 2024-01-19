package routeHandlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

func TestGetArticles_Success(t *testing.T) {
	// Set up mock store with predefined data.
	mockstore := &store.MockStore{
		Articles: []*models.Article{{Title: "Test Article 1"}},
	}

	// Initialize handler with mock store.
	handler := NewHandler(mockstore)

	// Configure router with API endpoint.
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")

	// Simulate HTTP request to API endpoint.
	req, err := http.NewRequest("GET", "/api/v1/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record HTTP response.
	rr := httptest.NewRecorder()

	// Serve request and capture response.
	r.ServeHTTP(rr, req)

	// Validate response status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	// Unmarshal response body into a slice of pointers to articles.
	var response struct {
		Data []models.Article `json:"data"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Failed to unmarshal response body:", err)
	}

	// Check if the response contains the expected data.
	if len(response.Data) != 1 || response.Data[0].Title != "Test Article 1" {
		t.Errorf("Unexpected response body: got %v", response.Data)
	}
}

func TestGetArticles_StoreError(t *testing.T) {
	// Set up mock store to return an error when GetArticles is called.
	mockstore := &store.MockStore{
		GetAllError: errors.New("store error"),
	}

	// Initialize handler with mock store.
	handler := NewHandler(mockstore)

	// Configure router with API endpoint.
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")

	// Simulate HTTP request to API endpoint.
	req, err := http.NewRequest("GET", "/api/v1/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record HTTP response.
	rr := httptest.NewRecorder()

	// Serve request and capture response.
	r.ServeHTTP(rr, req)

	// Validate response status code.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %v, got %v", http.StatusInternalServerError, status)
	}

	// Unmarshal response body into a struct to check for the error message.
	var response models.Response
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Failed to unmarshal response body:", err)
	}

	// Check if the response contains the expected error message.
	expectedError := "Failed to retrieve articles from the store"
	if response.Data != expectedError {
		t.Errorf("Expected error message %v, got %v", expectedError, response.Data)
	}
}
