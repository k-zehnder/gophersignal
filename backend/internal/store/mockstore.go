package store

import "github.com/k-zehnder/gophersignal/backend/internal/models"

// MockStore satisfies the Store interface.
type MockStore struct {
	Articles         []*models.Article
	SaveError        error
	GetArticlesError error
	IsOnHomepage     bool
	Limit            int
}

// NewMockStore creates and returns a new instance of MockStore.
func NewMockStore(articles []*models.Article, saveError error, getArticlesError error) *MockStore {
	return &MockStore{
		Articles:         articles,
		SaveError:        saveError,
		GetArticlesError: getArticlesError,
	}
}

func (ms *MockStore) SaveArticles(articles []*models.Article) error {
	if ms.SaveError != nil {
		return ms.SaveError
	}
	ms.Articles = articles
	return nil
}

func (ms *MockStore) GetArticles(isOnHomepage bool, limit int) ([]*models.Article, error) {
	if ms.GetArticlesError != nil {
		return nil, ms.GetArticlesError
	}
	ms.IsOnHomepage = isOnHomepage
	ms.Limit = limit

	if limit > 0 && len(ms.Articles) > limit {
		return ms.Articles[:limit], nil
	}

	return ms.Articles, nil
}
