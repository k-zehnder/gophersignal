package routehandlers

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
	var articles []*models.Article
	err = json.Unmarshal(rr.Body.Bytes(), &articles)
	if err != nil {
		t.Fatal("Failed to unmarshal response body:", err)
	}

	// Check if the response contains the expected data.
	if len(articles) != 1 || articles[0].Title != "Test Article 1" {
		t.Errorf("Unexpected response body: got %v", articles)
	}
}

func TestGetArticles_StoreError(t *testing.T) {
	// Set up mock store to return an error.
	mockStore := store.NewMockStore(nil, nil, errors.New("store error"))
	handler := NewHandler(mockStore)

	// Simulate HTTP request to API endpoint.
	req, err := http.NewRequest("GET", "/api/v1/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record HTTP response.
	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.GetArticles)

	// Serve request and capture response.
	handlerFunc.ServeHTTP(rr, req)

	// Validate response status code for error scenario.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %v, got %v", http.StatusInternalServerError, status)
	}

	// Validate response body for error message.
	expected := "store error\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected response body %v, got %v", expected, rr.Body.String())
	}
}
