// Package store provides an interface and implementations for article storage and retrieval.
// It includes tests for the MockStore type, which simulates a testing double for the Store interface.

package store

import (
	"errors"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

// TestMockStore_GetArticles validates the retrieval of articles from MockStore.
// It ensures the method returns the correct set of articles and checks for the absence of errors.
func TestMockStore_GetArticles(t *testing.T) {
	// Prepare expected articles for the mock store.
	expectedArticles := []*models.Article{{Title: "Test Article 1"}}
	mockStore := NewMockStore(expectedArticles, nil, nil)

	// Call GetArticles to retrieve articles from the mock store.
	articles, err := mockStore.GetArticles()

	// Ensure no error is returned.
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

// TestMockStore_SaveArticles checks the functionality of adding articles to MockStore.
// It tests whether the articles are correctly saved in the store.
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

	// Ensure that no error is returned from SaveArticles.
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

// TestMockStore_SaveArticles_Error tests the error handling of the SaveArticles method.
// It ensures that the method returns a predefined error when it is set in the mock store.
func TestMockStore_SaveArticles_Error(t *testing.T) {
	// Create a predefined error to be returned by the mock store.
	expectedErr := errors.New("save error")
	mockStore := NewMockStore(nil, expectedErr, nil)

	// Call SaveArticles, which should trigger the predefined error.
	err := mockStore.SaveArticles([]*models.Article{{Title: "Test Article"}})

	// Check if the returned error matches the expected error.
	if err != expectedErr {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}
}

// TestMockStore_GetArticles_Error tests the error handling of the GetArticles method.
// It checks that the method returns a predefined error when set in the mock store.
func TestMockStore_GetArticles_Error(t *testing.T) {
	// Set a predefined error to be returned by the mock store during GetArticles.
	expectedErr := errors.New("get error")
	mockStore := NewMockStore(nil, nil, expectedErr)

	// Attempt to retrieve articles, expecting the predefined error to be returned.
	_, err := mockStore.GetArticles()

	// Verify if the actual error matches the expected error.
	if err != expectedErr {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}
}
