package store

import "github.com/k-zehnder/gophersignal/backend/internal/models"

// MockStore satisfies the Store interface.
type MockStore struct {
	Articles    []*models.Article
	SaveError   error
	GetAllError error
}

// NewMockStore creates and returns a new instance of MockStore.
func NewMockStore(articles []*models.Article, saveError error, getAllError error) *MockStore {
	return &MockStore{
		Articles:    articles,
		SaveError:   saveError,
		GetAllError: getAllError,
	}
}

func (ms *MockStore) Init() error {
	return nil // No-op for the mock implementation of Store
}

func (ms *MockStore) SaveArticles(articles []*models.Article) error {
	if ms.SaveError != nil {
		return ms.SaveError
	}
	ms.Articles = articles
	return nil
}

func (ms *MockStore) GetArticles() ([]*models.Article, error) {
	if ms.GetAllError != nil {
		return nil, ms.GetAllError
	}
	return ms.Articles, nil
}
