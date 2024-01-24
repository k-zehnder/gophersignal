package routeHandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

func TestGetArticles_Success(t *testing.T) {
	// Set up mock store with predefined data where IsOnHomepage is true.
	mockstore := &store.MockStore{
		Articles: []*models.Article{{Title: "Test Article 1", IsOnHomepage: true}},
	}

	// Initialize handler with mock store.
	handler := NewHandler(mockstore)

	// Create a new request using the router with is_on_homepage=true.
	req, err := http.NewRequest("GET", "/api/v1/articles?is_on_homepage=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record HTTP response using httptest.
	rr := httptest.NewRecorder()

	// Create a router and serve the request.
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")
	router.ServeHTTP(rr, req)

	// Validate response status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	// Unmarshal response body into a slice of articles.
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
		GetArticlesError: errors.New("store error"),
	}

	// Initialize handler with mock store.
	handler := NewHandler(mockstore)

	// Create a new request using the router with is_on_homepage=true.
	req, err := http.NewRequest("GET", "/api/v1/articles?is_on_homepage=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record HTTP response using httptest.
	rr := httptest.NewRecorder()

	// Create a router and serve the request.
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")
	router.ServeHTTP(rr, req)

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

func TestGetArticles_NegativeCase(t *testing.T) {
	// Set up mock store with predefined data where IsOnHomepage is false.
	mockstore := &store.MockStore{
		Articles: []*models.Article{{Title: "Test Article 2", IsOnHomepage: false}},
	}

	// Initialize handler with mock store.
	handler := NewHandler(mockstore)

	// Create a new request using the router with is_on_homepage=false.
	req, err := http.NewRequest("GET", "/api/v1/articles?is_on_homepage=false", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record HTTP response using httptest.
	rr := httptest.NewRecorder()

	// Create a router and serve the request.
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")
	router.ServeHTTP(rr, req)

	// Validate response status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	// Unmarshal response body into a slice of articles.
	var response struct {
		Data []models.Article `json:"data"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Failed to unmarshal response body:", err)
	}

	// Check if the response contains the expected data.
	if len(response.Data) != 1 || response.Data[0].Title != "Test Article 2" {
		t.Errorf("Unexpected response body: got %v", response.Data)
	}
}

func TestHandler_JsonResponse(t *testing.T) {
	// Initialize a new HTTP response recorder.
	rr := httptest.NewRecorder()

	// Create a sample data structure.
	data := models.Response{Code: http.StatusOK, Status: "success", Data: "Sample data"}

	// Create a new Handler.
	handler := &Handler{}

	// Call the jsonResponse function.
	handler.jsonResponse(rr, data, http.StatusOK)

	// Validate response status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	// Validate response body.
	var response models.Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Failed to unmarshal response body:", err)
	}

	// Check if the response matches the expected data.
	if response != data {
		t.Errorf("Unexpected response body: got %v, expected %v", response, data)
	}
}

func TestGetArticles_ValidQueryParam(t *testing.T) {
	// Set up mock store with predefined data.
	mockstore := &store.MockStore{
		Articles: []*models.Article{{Title: "Test Article", IsOnHomepage: true}},
	}

	// Initialize handler with mock store.
	handler := NewHandler(mockstore)

	// Create a new request with valid query parameter.
	req, err := http.NewRequest("GET", "/api/v1/articles?is_on_homepage=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record HTTP response using httptest.
	rr := httptest.NewRecorder()

	// Create a router and serve the request.
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")
	router.ServeHTTP(rr, req)

	// Validate response status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}
}

func TestGetArticles_InvalidQueryParam(t *testing.T) {
	// Initialize handler with mock store.
	handler := NewHandler(&store.MockStore{})

	// Create a new request with an invalid query parameter value.
	req, err := http.NewRequest("GET", "/api/v1/articles?is_on_homepage=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record HTTP response using httptest.
	rr := httptest.NewRecorder()

	// Create a router and serve the request.
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")
	router.ServeHTTP(rr, req)

	// Validate response status code.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, status)
	}
}

func TestGetArticles_WithLimitParam(t *testing.T) {
	// Set up mock store with predefined data.
	mockstore := &store.MockStore{
		Articles: []*models.Article{
			{Title: "Test Article 1", IsOnHomepage: true},
			{Title: "Test Article 2", IsOnHomepage: true},
		},
	}

	// Initialize handler with mock store.
	handler := NewHandler(mockstore)

	// Create a new request with a valid limit query parameter.
	req, err := http.NewRequest("GET", "/api/v1/articles?limit=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record HTTP response using httptest.
	rr := httptest.NewRecorder()

	// Create a router and serve the request.
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")
	router.ServeHTTP(rr, req)

	// Validate response status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	// Unmarshal response body into a slice of articles.
	var response struct {
		Data []models.Article `json:"data"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Failed to unmarshal response body:", err)
	}

	// Check if the response contains only 1 article as set by the limit.
	if len(response.Data) != 1 {
		t.Errorf("Expected 1 article, got %v", len(response.Data))
	}
}

func TestGetArticles_InvalidLimitParam(t *testing.T) {
	// Set up mock store (can be empty because this test is for the limit logic).
	mockstore := &store.MockStore{}

	// Initialize handler with mock store.
	handler := NewHandler(mockstore)

	// Test scenarios with different invalid limit values.
	testCases := []string{
		"/api/v1/articles?limit=-1",  // Negative limit
		"/api/v1/articles?limit=0",   // Zero limit
		"/api/v1/articles?limit=101", // Limit greater than 100
	}

	for _, tc := range testCases {
		// Create a new request for each test case.
		req, err := http.NewRequest("GET", tc, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Record HTTP response using httptest.
		rr := httptest.NewRecorder()

		// Create a router and serve the request.
		router := mux.NewRouter()
		router.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")
		router.ServeHTTP(rr, req)

		// Validate response status code.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status code %v for %v, got %v", http.StatusOK, tc, status)
		}
	}
}

func TestGetArticles_ExceedsMaxLimit(t *testing.T) {
	// Set up a mock store with more than 100 articles.
	var mockArticles []*models.Article
	for i := 0; i < 105; i++ {
		mockArticles = append(mockArticles, &models.Article{Title: fmt.Sprintf("Test Article %d", i), IsOnHomepage: true})
	}
	mockstore := &store.MockStore{
		Articles: mockArticles,
	}

	// Initialize the handler with the mock store.
	handler := NewHandler(mockstore)

	// Create a new request with a limit query parameter exceeding 100.
	req, err := http.NewRequest("GET", "/api/v1/articles?limit=150", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the HTTP response using httptest.
	rr := httptest.NewRecorder()

	// Create a router and serve the request.
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")
	router.ServeHTTP(rr, req)

	// Validate the response status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	// Unmarshal the response body into a slice of articles.
	var response struct {
		Data []models.Article `json:"data"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Failed to unmarshal response body:", err)
	}

	// Check if the number of articles in the response does not exceed 100.
	if len(response.Data) > 100 {
		t.Errorf("Expected a maximum of 100 articles, got %v", len(response.Data))
	}
}
