package routeHandlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"github.com/k-zehnder/gophersignal/backend/internal/store"
)

func TestGetArticlesHandler(t *testing.T) {
	// Create a mock store with some dummy data
	mockstore := &store.MockStore{
		Articles: []*models.Article{{Title: "Test Article 1"}},
	}

	// Create a new handler with the mock store
	handler := NewHandler(mockstore)

	// Create a new router and register the handler
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/articles", handler.GetArticles).Methods("GET")

	// Create a new HTTP request to the endpoint
	req, err := http.NewRequest("GET", "/api/v1/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body, headers, etc next
}
