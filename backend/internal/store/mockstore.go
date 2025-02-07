// Package store provides an interface and implementations for article storage and retrieval.
// It includes a MockStore type, which serves as a testing double for the Store interface.
package store

import "github.com/k-zehnder/gophersignal/backend/internal/models"

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
// It applies pagination (limit, offset) to the stored articles.
func (ms *MockStore) GetArticles(limit, offset int) ([]*models.Article, error) {
	if ms.GetAllError != nil {
		return nil, ms.GetAllError
	}
	// Simple pagination: slice the Articles slice.
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
// Each filter is provided as a pointer to a bool. If a filter is nil, that field is not filtered.
// It applies pagination (limit, offset) and returns at most 'limit' articles starting at 'offset'.
func (ms *MockStore) GetFilteredArticles(flagged, dead, dupe *bool, limit, offset int) ([]*models.Article, error) {
	if ms.GetAllError != nil {
		return nil, ms.GetAllError
	}

	var filtered []*models.Article
	for _, article := range ms.Articles {
		// Apply the flagged filter if provided.
		if flagged != nil && article.Flagged != *flagged {
			continue
		}
		// Apply the dead filter if provided.
		if dead != nil && article.Dead != *dead {
			continue
		}
		// Apply the dupe filter if provided.
		if dupe != nil && article.Dupe != *dupe {
			continue
		}
		filtered = append(filtered, article)
	}
	// Apply pagination to the filtered results.
	if offset >= len(filtered) {
		return []*models.Article{}, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], nil
}
