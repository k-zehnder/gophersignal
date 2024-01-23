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
