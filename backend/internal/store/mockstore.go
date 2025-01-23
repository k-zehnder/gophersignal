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
func (ms *MockStore) GetArticles() ([]*models.Article, error) {
	if ms.GetAllError != nil {
		return nil, ms.GetAllError
	}
	return ms.Articles, nil
}
