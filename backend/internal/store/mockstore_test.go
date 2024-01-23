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

	articles, err := mockstore.GetArticles(true)
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

	if _, err := mockstore.GetArticles(true); err != expectedError {
		t.Fatalf("Expected error '%v', got '%v'", expectedError, err)
	}
}
