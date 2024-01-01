package store

import (
    "testing"
    "github.com/k-zehnder/gophersignal/backend/internal/models"
)

// TestGetAllArticles verifies that GetArticles method retrieves the correct number of articles.
func TestGetAllArticles(t *testing.T) {
    expectedArticles := []*models.Article{{Title: "Test Article 1"}}
    mockstore := NewMockStore(expectedArticles, nil, nil)

    articles, err := mockstore.GetArticles()
    if err != nil {
        t.Fatalf("GetAllArticles() error = %v, wantErr %v", err, nil)
    }

    if len(articles) != len(expectedArticles) {
        t.Fatalf("Expected %d articles, got %d", len(expectedArticles), len(articles))
    }
}

// TestSaveArticles checks the SaveArticles method for correct article storage.
func TestSaveArticles(t *testing.T) {
    mockstore := NewMockStore(nil, nil, nil)
    articles := []*models.Article{
        {Title: "Test Article 1"},
        {Title: "Test Article 2"},
    }

    err := mockstore.SaveArticles(articles)
    if err != nil {
        t.Fatalf("SaveArticles() error = %v, wantErr nil", err)
    }

    if len(mockstore.Articles) != 2 {
        t.Fatalf("Expected 2 articles, got %d", len(mockstore.Articles))
    }
}
