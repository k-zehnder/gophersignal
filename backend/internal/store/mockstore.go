package store

import "github.com/k-zehnder/gophersignal/backend/internal/models"

type MockStore struct {
	Articles    []*models.Article
	SaveError   error
	GetAllError error
}

func (ms *MockStore) SaveArticles(articles []*models.Article) error {
	if ms.SaveError != nil {
		return ms.SaveError
	}
	ms.Articles = articles
	return nil
}

func (ms *MockStore) GetAllArticles() ([]*models.Article, error) {
	if ms.GetAllError != nil {
		return nil, ms.GetAllError
	}
	return ms.Articles, nil
}
