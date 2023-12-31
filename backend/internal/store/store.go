package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

// Store defines methods for storing and retrieving articles in a repository.
type Store interface {
	SaveArticles(articles []*models.Article) error
	GetAllArticles() ([]*models.Article, error)
}

// DBStore wraps the SQL database connection.
type DBStore struct {
	db *sql.DB
}

// NewDBStore initializes a new DBStore instance.
func NewDBStore(dataSourceName string) (*DBStore, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DBStore{db: db}, nil
}

// SaveArticles updates articles in the database.
func (store *DBStore) SaveArticles(articles []*models.Article) error {
	// Reset IsOnHomepage for all articles
	if _, err := store.db.Exec("UPDATE articles SET isOnHomepage = FALSE"); err != nil {
		return fmt.Errorf("failed to reset articles isOnHomepage status: %w", err)
	}

	// Update articles with the new data
	for _, article := range articles {
		if _, err := store.db.Exec("INSERT INTO articles (title, link, content, source, isOnHomepage) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE title=VALUES(title), link=VALUES(link), content=VALUES(content), source=VALUES(source), isOnHomepage=VALUES(isOnHomepage)",
			article.Title, article.Link, article.Content, article.Source, article.IsOnHomepage); err != nil {
			return fmt.Errorf("failed to save article: %s, error: %w", article.Title, err)
		}
	}

	return nil
}

// GetAllArticles retrieves all articles from the database.
func (store *DBStore) GetAllArticles() ([]*models.Article, error) {
	rows, err := store.db.Query("SELECT id, title, link, content, summary, source, scraped_at, isOnHomepage FROM articles")
	if err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Link, &article.Content, &article.Summary, &article.Source, &article.ScrapedAt, &article.IsOnHomepage); err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, &article)
	}

	return articles, nil
}
