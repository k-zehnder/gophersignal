package store

import (
	"errors"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

// TestMockStore_GetArticles verifies that the mock store returns the expected articles.
func TestMockStore_GetArticles(t *testing.T) {
	expectedArticles := []*models.Article{{Title: "Test Article 1"}}
	mockstore := NewMockStore(expectedArticles, nil, nil)

	articles, err := mockstore.GetArticles(true, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(articles) != len(expectedArticles) {
		t.Fatalf("Expected %d articles, got %d", len(expectedArticles), len(articles))
	}
}

// TestMockStore_SaveArticles checks that the mock store successfully saves articles.
func TestMockStore_SaveArticles(t *testing.T) {
	mockstore := NewMockStore(nil, nil, nil)
	articles := []*models.Article{
		{Title: "Test Article 1"},
		{Title: "Test Article 2"},
	}

	if err := mockstore.SaveArticles(articles); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(mockstore.Articles) != len(articles) {
		t.Fatalf("Expected %d articles saved, got %d", len(articles), len(mockstore.Articles))
	}
}

// TestMockStore_SaveArticles_Error ensures that the mock store returns an error when saving articles fails.
func TestMockStore_SaveArticles_Error(t *testing.T) {
	expectedError := errors.New("save error")
	mockstore := NewMockStore(nil, expectedError, nil)
	articles := []*models.Article{{Title: "Test Article Error"}}

	if err := mockstore.SaveArticles(articles); err != expectedError {
		t.Fatalf("Expected error '%v', got '%v'", expectedError, err)
	}
}

// TestMockStore_GetArticles_Error confirms that the mock store returns an error when article retrieval fails.
func TestMockStore_GetArticles_Error(t *testing.T) {
	expectedError := errors.New("get all error")
	mockstore := NewMockStore(nil, nil, expectedError)

	if _, err := mockstore.GetArticles(true, 10); err != expectedError {
		t.Fatalf("Expected error '%v', got '%v'", expectedError, err)
	}
}

func TestMockStore_GetArticles_LimitExceeded(t *testing.T) {
	// Create a list of expected articles with more items than the specified limit.
	expectedArticles := []*models.Article{
		{Title: "Test Article 1"},
		{Title: "Test Article 2"},
		{Title: "Test Article 3"},
	}
	mockstore := NewMockStore(expectedArticles, nil, nil)

	// Specify the limit to be less than the number of expected articles.
	limit := 2

	articles, err := mockstore.GetArticles(true, limit)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Ensure that the number of returned articles matches the specified limit.
	if len(articles) != limit {
		t.Fatalf("Expected %d articles, got %d", limit, len(articles))
	}

	// Check if the returned articles match the first 'limit' articles in the expected list.
	for i := 0; i < limit; i++ {
		if articles[i].Title != expectedArticles[i].Title {
			t.Fatalf("Expected article '%s', got '%s'", expectedArticles[i].Title, articles[i].Title)
		}
	}
}
