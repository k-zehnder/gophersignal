package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

// Store interface defines the operations for storing and retrieving articles.
type Store interface {
	Init() error
	SaveArticles(articles []*models.Article) error
	GetArticles() ([]*models.Article, error)
}

// DBStore implements the Store interface using a SQL database.
type DBStore struct {
	db *sql.DB
}

// NewDBStore establishes a connection to the SQL database and returns a DBStore.
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

// Init sets up the necessary database tables, particularly 'articles'.
func (store *DBStore) Init() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS articles (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		link VARCHAR(512) NOT NULL,
		content TEXT,
		summary TEXT,
		source VARCHAR(100) NOT NULL,
		scraped_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		isOnHomepage BOOLEAN,
		UNIQUE KEY unique_article (title, link)
	);`
	_, err := store.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create articles table: %w", err)
	}
	return nil
}

// SaveArticles updates or adds new articles in the database.
func (store *DBStore) SaveArticles(articles []*models.Article) error {
	if _, err := store.db.Exec("UPDATE articles SET isOnHomepage = FALSE"); err != nil {
		return fmt.Errorf("failed to reset articles isOnHomepage status: %w", err)
	}

	for _, article := range articles {
		if _, err := store.db.Exec("INSERT INTO articles (title, link, content, source, isOnHomepage) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE title=VALUES(title), link=VALUES(link), content=VALUES(content), source=VALUES(source), isOnHomepage=VALUES(isOnHomepage)",
			article.Title, article.Link, article.Content, article.Source, article.IsOnHomepage); err != nil {
			return fmt.Errorf("failed to save article: %s, error: %w", article.Title, err)
		}
	}

	return nil
}

// GetArticles retrieves all the articles from the database.
func (store *DBStore) GetArticles() ([]*models.Article, error) {
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
