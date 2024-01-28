// Package store defines an interface and an implementation for article storage and retrieval.
// It includes the Store interface for defining methods to interact with the article data in the database.
package store

import (
	"database/sql"
	"fmt"

	// Import the MySQL driver with a blank identifier to ensure its `init()` function is executed.
	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

// Store interface defines methods for article storage and retrieval.
type Store interface {
	SaveArticles(articles []*models.Article) error
	GetArticles() ([]*models.Article, error)
}

// MySQLStore implements the Store interface using a MySQL database.
type MySQLStore struct {
	db *sql.DB // db represents the connection to the database.
}

// NewMySQLStore establishes a new MySQL database connection.
func NewMySQLStore(dataSourceName string) (*MySQLStore, error) {
	// Attempt to open a database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ping the database to ensure the connection is active and the server is reachable.
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Return a new MySQLStore instance with the established database connection.
	return &MySQLStore{db: db}, nil
}

// SaveArticles handles the addition or update of articles in the database.
func (store *MySQLStore) SaveArticles(articles []*models.Article) error {
	// Start a transaction.
	tx, err := store.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Prepare SQL statement for article insertion or update.
	stmt, err := tx.Prepare(`
        INSERT INTO articles (title, link, content, summary, source, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE title=VALUES(title), link=VALUES(link), content=VALUES(content),
        summary=VALUES(summary), source=VALUES(source), updated_at=VALUES(updated_at);`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement for each article.
	for _, article := range articles {
		_, err := stmt.Exec(article.Title, article.Link, article.Content, article.Summary, article.Source, article.CreatedAt, article.UpdatedAt)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed for article '%s': %w", article.Title, err)
		}
	}

	// Commit the transaction after successful execution.
	return tx.Commit()
}

// GetArticles retrieves the latest 30 articles, sorted by their update timestamp.
func (store *MySQLStore) GetArticles() ([]*models.Article, error) {
	// Query to fetch articles.
	query := "SELECT id, title, link, content, summary, source, created_at, updated_at FROM articles ORDER BY updated_at DESC LIMIT 30;"

	rows, err := store.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Populate articles from query results.
	var articles []*models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Link, &article.Content, &article.Summary, &article.Source, &article.CreatedAt, &article.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, &article)
	}

	// Handle any iteration errors.
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}

	return articles, nil
}
