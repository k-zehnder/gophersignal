package store

import "github.com/k-zehnder/gophersignal/internal/models"

type MockStore struct {
	GetAllArticlesFunc func() []models.Article
	SaveArticleFunc    func(article models.Article) error
}

func (m *MockStore) GetAllArticles() []models.Article {
	if m.GetAllArticlesFunc != nil {
		return m.GetAllArticlesFunc()
	}
	return nil
}

func (m *MockStore) SaveArticle(article models.Article) error {
	if m.SaveArticleFunc != nil {
		return m.SaveArticleFunc(article)
	}
	return nil
}
