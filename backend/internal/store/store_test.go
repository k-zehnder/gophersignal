package store

import (
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

func TestMockStore_GetArticles(t *testing.T) {
	expectedArticles := []*models.Article{{Title: "Test Article 1"}}
	mockstore := NewMockStore(expectedArticles, nil, nil)

	articles, err := mockstore.GetArticles()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(articles) != len(expectedArticles) {
		t.Fatalf("Expected %d articles, got %d", len(expectedArticles), len(articles))
	}
}

func TestMockStore_SaveArticles(t *testing.T) {
	mockstore := NewMockStore(nil, nil, nil)
	articles := []*models.Article{
		{Title: "Test Article 1"},
		{Title: "Test Article 2"},
	}

	err := mockstore.SaveArticles(articles)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockstore.Articles) != 2 {
		t.Fatalf("Expected 2 articles, got %d", len(mockstore.Articles))
	}
}
