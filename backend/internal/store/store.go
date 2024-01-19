package store

import (
	"github.com/k-zehnder/gophersignal/internal/models"
	"gorm.io/gorm"
)

// ArticleStore is an interface for storing and retrieving articles.
type ArticleStore interface {
	GetAllArticles() []models.Article
	SaveArticle(article models.Article) error
}

// Store is the implementation of the ArticleStore interface using a GORM database.
type Store struct {
	Db *gorm.DB
}

// NewStore creates a new Store instance.
func NewStore(db *gorm.DB) *Store {
	return &Store{Db: db}
}

// GetAllArticles retrieves all articles from the database.
func (s *Store) GetAllArticles() []models.Article {
	var articles []models.Article
	result := s.Db.Find(&articles)
	if result.Error != nil {
		return nil
	}
	return articles
}

// SaveArticle saves an article to the database.
func (s *Store) SaveArticle(article models.Article) error {
	result := s.Db.Create(&article)
	return result.Error
}
