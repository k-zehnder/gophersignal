// Package store provides an interface and implementations for article storage and retrieval.
// It includes a MockStore type, which serves as a testing double for the Store interface.
package store

import (
	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

// MockStore serves as a testing double for the Store interface.
type MockStore struct {
	Articles    []*models.Article
	SaveError   error
	GetAllError error
}

// NewMockStore initializes a MockStore with predefined articles and potential errors.
func NewMockStore(articles []*models.Article, saveError error, getAllError error) *MockStore {
	return &MockStore{
		Articles:    articles,
		SaveError:   saveError,
		GetAllError: getAllError,
	}
}

// SaveArticles simulates storing articles, returning a predefined error if set.
func (ms *MockStore) SaveArticles(articles []*models.Article) error {
	if ms.SaveError != nil {
		return ms.SaveError
	}
	ms.Articles = articles
	return nil
}

// GetArticles simulates fetching articles, returning a predefined error if set.
func (ms *MockStore) GetArticles(limit, offset int) ([]*models.Article, error) {
	if ms.GetAllError != nil {
		return nil, ms.GetAllError
	}
	if offset >= len(ms.Articles) {
		return []*models.Article{}, nil
	}
	end := offset + limit
	if end > len(ms.Articles) {
		end = len(ms.Articles)
	}
	return ms.Articles[offset:end], nil
}

// GetFilteredArticles simulates fetching articles based on optional filter criteria.
func (ms *MockStore) GetFilteredArticles(flagged, dead, dupe *bool, limit, offset int) ([]*models.Article, error) {
	if ms.GetAllError != nil {
		return nil, ms.GetAllError
	}

	var filtered []*models.Article
	for _, article := range ms.Articles {
		if flagged != nil && article.Flagged != *flagged {
			continue
		}
		if dead != nil && article.Dead != *dead {
			continue
		}
		if dupe != nil && article.Dupe != *dupe {
			continue
		}
		filtered = append(filtered, article)
	}
	if offset >= len(filtered) {
		return []*models.Article{}, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], nil
}

// GetArticlesWithThresholds implements threshold-based filtering.
// It converts models.NullableInt values to int64 for comparison.
// If minUpvotes or minComments is 0, that metric is not filtered.
func (ms *MockStore) GetArticlesWithThresholds(limit, offset, minUpvotes, minComments int) ([]*models.Article, error) {
	if ms.GetAllError != nil {
		return nil, ms.GetAllError
	}

	var filtered []*models.Article
	for _, article := range ms.Articles {
		// Convert Upvotes and CommentCount from NullableInt to int64.
		upvotes := int64(0)
		if article.Upvotes.Valid {
			upvotes = article.Upvotes.Int64
		}
		comments := int64(0)
		if article.CommentCount.Valid {
			comments = article.CommentCount.Int64
		}

		if upvotes >= int64(minUpvotes) && comments >= int64(minComments) {
			filtered = append(filtered, article)
		}
	}

	if offset >= len(filtered) {
		return []*models.Article{}, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], nil
}

// GetArticlesWithThresholdsAndFilters implements combined filtering:
// both threshold conditions (minUpvotes and minComments) and additional boolean filters.
func (ms *MockStore) GetArticlesWithThresholdsAndFilters(limit, offset, minUpvotes, minComments int, flagged, dead, dupe *bool) ([]*models.Article, error) {
	if ms.GetAllError != nil {
		return nil, ms.GetAllError
	}

	var filtered []*models.Article
	for _, article := range ms.Articles {
		// Apply boolean filters.
		if flagged != nil && article.Flagged != *flagged {
			continue
		}
		if dead != nil && article.Dead != *dead {
			continue
		}
		if dupe != nil && article.Dupe != *dupe {
			continue
		}

		// Convert Upvotes and CommentCount from NullableInt to int64.
		upvotes := int64(0)
		if article.Upvotes.Valid {
			upvotes = article.Upvotes.Int64
		}
		comments := int64(0)
		if article.CommentCount.Valid {
			comments = article.CommentCount.Int64
		}

		if upvotes >= int64(minUpvotes) && comments >= int64(minComments) {
			filtered = append(filtered, article)
		}
	}

	if offset >= len(filtered) {
		return []*models.Article{}, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], nil
}
