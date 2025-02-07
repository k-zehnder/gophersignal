// Package store provides an interface and implementations for article storage and retrieval.
// It includes tests for the MockStore type, which simulates a testing double for the Store interface.

package store

import (
	"errors"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

func TestMockStore_GetArticles(t *testing.T) {
	// Prepare expected articles for the mock store.
	expectedArticles := []*models.Article{
		{Title: "Test Article 1"},
	}
	mockStore := NewMockStore(expectedArticles, nil, nil)

	// Call GetArticles with pagination parameters (limit, offset).
	articles, err := mockStore.GetArticles(len(expectedArticles), 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that the number of articles returned matches the expected number.
	if len(articles) != len(expectedArticles) {
		t.Fatalf("Expected %d articles, got %d", len(expectedArticles), len(articles))
	}

	// Check if the returned articles match the expected ones.
	for i, article := range articles {
		if article.Title != expectedArticles[i].Title {
			t.Errorf("Expected title %s, got %s", expectedArticles[i].Title, article.Title)
		}
	}
}

func TestMockStore_SaveArticles(t *testing.T) {
	// Initialize a mock store without any pre-existing articles.
	mockStore := NewMockStore(nil, nil, nil)

	// Define articles to be saved in the mock store.
	articles := []*models.Article{
		{Title: "Test Article 1"},
		{Title: "Test Article 2"},
	}

	// Attempt to save articles in the mock store.
	err := mockStore.SaveArticles(articles)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that the mock store now contains the expected number of articles.
	if len(mockStore.Articles) != len(articles) {
		t.Fatalf("Expected %d articles, got %d", len(articles), len(mockStore.Articles))
	}

	// Confirm that each saved article matches the corresponding input article.
	for i, article := range mockStore.Articles {
		if article.Title != articles[i].Title {
			t.Errorf("Expected title %s, got %s", articles[i].Title, article.Title)
		}
	}
}

func TestMockStore_SaveArticles_Error(t *testing.T) {
	// Create a predefined error to be returned by the mock store.
	expectedErr := errors.New("save error")
	mockStore := NewMockStore(nil, expectedErr, nil)

	// Call SaveArticles, which should trigger the predefined error.
	err := mockStore.SaveArticles([]*models.Article{{Title: "Test Article"}})
	if err != expectedErr {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}
}

func TestMockStore_GetArticles_Error(t *testing.T) {
	// Set a predefined error to be returned by the mock store during GetArticles.
	expectedErr := errors.New("get error")
	mockStore := NewMockStore(nil, nil, expectedErr)

	// Attempt to retrieve articles, expecting the predefined error to be returned.
	_, err := mockStore.GetArticles(10, 0)
	if err != expectedErr {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}
}

func TestMockStore_GetFilteredArticles_NoFilter(t *testing.T) {
	// Prepare a list of articles with various flagged/dead/dupe values.
	articles := []*models.Article{
		{Title: "Article 1", Flagged: false, Dead: false, Dupe: false},
		{Title: "Article 2", Flagged: true, Dead: false, Dupe: false},
		{Title: "Article 3", Flagged: false, Dead: true, Dupe: false},
		{Title: "Article 4", Flagged: false, Dead: false, Dupe: true},
	}
	mockStore := NewMockStore(articles, nil, nil)

	// When no filters are applied, expect all articles.
	filtered, err := mockStore.GetFilteredArticles(nil, nil, nil, len(articles), 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(filtered) != len(articles) {
		t.Fatalf("Expected %d articles, got %d", len(articles), len(filtered))
	}
}

func TestMockStore_GetFilteredArticles_Flagged(t *testing.T) {
	// Prepare a list of articles.
	articles := []*models.Article{
		{Title: "Article 1", Flagged: false, Dead: false, Dupe: false},
		{Title: "Article 2", Flagged: true, Dead: false, Dupe: false},
		{Title: "Article 3", Flagged: true, Dead: true, Dupe: false},
		{Title: "Article 4", Flagged: false, Dead: false, Dupe: true},
	}
	mockStore := NewMockStore(articles, nil, nil)

	// Apply a filter for flagged = true.
	flagged := true
	filtered, err := mockStore.GetFilteredArticles(&flagged, nil, nil, len(articles), 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Expect Article 2 and Article 3 to match.
	if len(filtered) != 2 {
		t.Fatalf("Expected 2 articles, got %d", len(filtered))
	}
	for _, a := range filtered {
		if !a.Flagged {
			t.Errorf("Expected flagged article, got %v", a.Title)
		}
	}
}

func TestMockStore_GetFilteredArticles_DeadAndDupe(t *testing.T) {
	// Prepare a list of articles.
	articles := []*models.Article{
		{Title: "Article 1", Flagged: false, Dead: false, Dupe: false},
		{Title: "Article 2", Flagged: true, Dead: true, Dupe: false},
		{Title: "Article 3", Flagged: true, Dead: true, Dupe: true},
		{Title: "Article 4", Flagged: false, Dead: false, Dupe: true},
	}
	mockStore := NewMockStore(articles, nil, nil)

	// Apply filters: dead = true and dupe = true.
	dead := true
	dupe := true
	filtered, err := mockStore.GetFilteredArticles(nil, &dead, &dupe, len(articles), 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Only Article 3 meets both criteria.
	if len(filtered) != 1 {
		t.Fatalf("Expected 1 article, got %d", len(filtered))
	}
	if filtered[0].Title != "Article 3" {
		t.Errorf("Expected article 3, got %s", filtered[0].Title)
	}
}

func TestMockStore_GetFilteredArticles_Error(t *testing.T) {
	// Set a predefined error to be returned by GetFilteredArticles.
	expectedErr := errors.New("get error")
	mockStore := NewMockStore(nil, nil, expectedErr)

	_, err := mockStore.GetFilteredArticles(nil, nil, nil, 10, 0)
	if err != expectedErr {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}
}
