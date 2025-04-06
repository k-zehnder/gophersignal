// Package store provides an interface and implementations for article storage and retrieval.
// This file contains comprehensive tests for the MockStore type, which serves as a testing double for the Store interface.
package store

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

// TestMockStore_GetArticles verifies basic article retrieval functionality with pagination.
func TestMockStore_GetArticles(t *testing.T) {
	// Setup test data and mock store.
	expectedArticles := []*models.Article{
		{Title: "Test Article 1"},
	}
	mockStore := NewMockStore(expectedArticles, nil, nil)

	// Execute GetArticles with full result range.
	articles, err := mockStore.GetArticles(len(expectedArticles), 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate result count matches expectation.
	if len(articles) != len(expectedArticles) {
		t.Fatalf("Expected %d articles, got %d", len(expectedArticles), len(articles))
	}

	// Verify individual article contents.
	for i, article := range articles {
		if article.Title != expectedArticles[i].Title {
			t.Errorf("Expected title %s, got %s", expectedArticles[i].Title, article.Title)
		}
	}
}

// TestMockStore_SaveArticles validates the article persistence mechanism.
func TestMockStore_SaveArticles(t *testing.T) {
	mockStore := NewMockStore(nil, nil, nil)
	articles := []*models.Article{
		{Title: "Test Article 1"},
		{Title: "Test Article 2"},
	}

	err := mockStore.SaveArticles(articles)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockStore.Articles) != len(articles) {
		t.Fatalf("Expected %d articles, got %d", len(articles), len(mockStore.Articles))
	}

	for i, article := range mockStore.Articles {
		if article.Title != articles[i].Title {
			t.Errorf("Expected title %s, got %s", articles[i].Title, article.Title)
		}
	}
}

// TestMockStore_SaveArticles_Error verifies error propagation in save operations.
func TestMockStore_SaveArticles_Error(t *testing.T) {
	expectedErr := errors.New("save error")
	mockStore := NewMockStore(nil, expectedErr, nil)

	err := mockStore.SaveArticles([]*models.Article{{Title: "Test Article"}})
	if err != expectedErr {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}
}

// TestMockStore_GetArticles_Error checks error handling in retrieval operations.
func TestMockStore_GetArticles_Error(t *testing.T) {
	expectedErr := errors.New("get error")
	mockStore := NewMockStore(nil, nil, expectedErr)

	_, err := mockStore.GetArticles(10, 0)
	if err != expectedErr {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}
}

// TestMockStore_GetFilteredArticles_NoFilter verifies unfiltered retrieval.
func TestMockStore_GetFilteredArticles_NoFilter(t *testing.T) {
	articles := []*models.Article{
		{Title: "Article 1", Flagged: false, Dead: false, Dupe: false},
		{Title: "Article 2", Flagged: true, Dead: false, Dupe: false},
		{Title: "Article 3", Flagged: false, Dead: true, Dupe: false},
		{Title: "Article 4", Flagged: false, Dead: false, Dupe: true},
	}
	mockStore := NewMockStore(articles, nil, nil)

	filtered, err := mockStore.GetFilteredArticles(nil, nil, nil, len(articles), 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(filtered) != len(articles) {
		t.Fatalf("Expected %d articles, got %d", len(articles), len(filtered))
	}
}

// TestMockStore_GetFilteredArticles_Flagged validates single-flag filtering.
func TestMockStore_GetFilteredArticles_Flagged(t *testing.T) {
	articles := []*models.Article{
		{Title: "Article 1", Flagged: false, Dead: false, Dupe: false},
		{Title: "Article 2", Flagged: true, Dead: false, Dupe: false},
		{Title: "Article 3", Flagged: true, Dead: true, Dupe: false},
		{Title: "Article 4", Flagged: false, Dead: false, Dupe: true},
	}
	mockStore := NewMockStore(articles, nil, nil)

	flagged := true
	filtered, err := mockStore.GetFilteredArticles(&flagged, nil, nil, len(articles), 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(filtered) != 2 {
		t.Fatalf("Expected 2 articles, got %d", len(filtered))
	}
	for _, a := range filtered {
		if !a.Flagged {
			t.Errorf("Expected flagged article, got %v", a.Title)
		}
	}
}

// TestMockStore_GetFilteredArticles_DeadAndDupe validates combined flag filtering.
func TestMockStore_GetFilteredArticles_DeadAndDupe(t *testing.T) {
	articles := []*models.Article{
		{Title: "Article 1", Flagged: false, Dead: false, Dupe: false},
		{Title: "Article 2", Flagged: true, Dead: true, Dupe: false},
		{Title: "Article 3", Flagged: true, Dead: true, Dupe: true},
		{Title: "Article 4", Flagged: false, Dead: false, Dupe: true},
	}
	mockStore := NewMockStore(articles, nil, nil)

	dead := true
	dupe := true
	filtered, err := mockStore.GetFilteredArticles(nil, &dead, &dupe, len(articles), 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(filtered) != 1 {
		t.Fatalf("Expected 1 article, got %d", len(filtered))
	}
	if filtered[0].Title != "Article 3" {
		t.Errorf("Expected article 3, got %s", filtered[0].Title)
	}
}

// TestMockStore_GetFilteredArticles_Error verifies error handling in filtered retrieval.
func TestMockStore_GetFilteredArticles_Error(t *testing.T) {
	expectedErr := errors.New("get error")
	mockStore := NewMockStore(nil, nil, expectedErr)

	_, err := mockStore.GetFilteredArticles(nil, nil, nil, 10, 0)
	if err != expectedErr {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}
}

// TestMockStore_GetArticlesWithThresholds validates threshold-based filtering.
func TestMockStore_GetArticlesWithThresholds(t *testing.T) {
	articles := []*models.Article{
		{Title: "Article 1", Upvotes: models.NewNullableInt(50), CommentCount: models.NewNullableInt(10)},
		{Title: "Article 2", Upvotes: models.NewNullableInt(30), CommentCount: models.NewNullableInt(5)},
		{Title: "Article 3", Upvotes: models.NewNullableInt(40), CommentCount: models.NewNullableInt(15)},
	}
	mockStore := NewMockStore(articles, nil, nil)

	t.Run("MeetBothThresholds", func(t *testing.T) {
		result, err := mockStore.GetArticlesWithThresholds(10, 0, 40, 10)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("Expected 2 articles, got %d", len(result))
		}
	})

	t.Run("ZeroThresholds", func(t *testing.T) {
		result, err := mockStore.GetArticlesWithThresholds(10, 0, 0, 0)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(result) != len(articles) {
			t.Fatalf("Expected %d articles, got %d", len(articles), len(result))
		}
	})
}

// TestMockStore_GetArticlesWithThresholdsAndFilters validates combined filtering.
func TestMockStore_GetArticlesWithThresholdsAndFilters(t *testing.T) {
	articles := []*models.Article{
		{Title: "Article 1", Upvotes: models.NewNullableInt(50), CommentCount: models.NewNullableInt(10), Flagged: false, Dead: false, Dupe: false},
		{Title: "Article 2", Upvotes: models.NewNullableInt(30), CommentCount: models.NewNullableInt(5), Flagged: true, Dead: false, Dupe: false},
		{Title: "Article 3", Upvotes: models.NewNullableInt(40), CommentCount: models.NewNullableInt(15), Flagged: false, Dead: true, Dupe: false},
		{Title: "Article 4", Upvotes: models.NewNullableInt(35), CommentCount: models.NewNullableInt(20), Flagged: true, Dead: true, Dupe: true},
	}
	mockStore := NewMockStore(articles, nil, nil)

	t.Run("CombinedFilters", func(t *testing.T) {
		flagged := true
		dead := true
		dupe := true
		result, err := mockStore.GetArticlesWithThresholdsAndFilters(10, 0, 30, 5, &flagged, &dead, &dupe)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(result) != 1 || result[0].Title != "Article 4" {
			t.Errorf("Expected Article 4, got %v", result)
		}
	})

	t.Run("PartialFilters", func(t *testing.T) {
		flagged := false
		result, err := mockStore.GetArticlesWithThresholdsAndFilters(10, 0, 40, 10, &flagged, nil, nil)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("Expected 2 articles, got %d", len(result))
		}

		// Verify both expected articles are present
		titles := map[string]bool{
			"Article 1": false,
			"Article 3": false,
		}
		for _, article := range result {
			titles[article.Title] = true
		}

		for title, found := range titles {
			if !found {
				t.Errorf("Expected article %q not found in results", title)
			}
		}
	})
}

// TestMockStore_ThresholdMethods_Error verifies error propagation in threshold methods.
func TestMockStore_ThresholdMethods_Error(t *testing.T) {
	expectedErr := errors.New("get error")
	mockStore := NewMockStore(nil, nil, expectedErr)

	t.Run("ThresholdsError", func(t *testing.T) {
		_, err := mockStore.GetArticlesWithThresholds(10, 0, 0, 0)
		if err != expectedErr {
			t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
		}
	})

	t.Run("CombinedError", func(t *testing.T) {
		_, err := mockStore.GetArticlesWithThresholdsAndFilters(10, 0, 0, 0, nil, nil, nil)
		if err != expectedErr {
			t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
		}
	})
}

// TestMockStore_Pagination validates pagination behavior across methods.
func TestMockStore_Pagination(t *testing.T) {
	articles := make([]*models.Article, 20)
	for i := range articles {
		articles[i] = &models.Article{Title: fmt.Sprintf("Article %d", i+1)}
	}
	mockStore := NewMockStore(articles, nil, nil)

	t.Run("LimitOffset", func(t *testing.T) {
		result, err := mockStore.GetArticles(5, 10)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(result) != 5 {
			t.Fatalf("Expected 5 articles, got %d", len(result))
		}
		if !strings.Contains(result[0].Title, "Article 11") {
			t.Errorf("Expected first article to be Article 11, got %s", result[0].Title)
		}
	})

	t.Run("OffsetExceedsResults", func(t *testing.T) {
		result, err := mockStore.GetArticlesWithThresholds(10, 25, 0, 0)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Fatalf("Expected empty result, got %d articles", len(result))
		}
	})
}
