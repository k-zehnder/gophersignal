package store

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

// Store interface defines the operations for storing and retrieving articles.
type Store interface {
	SaveArticles(articles []*models.Article) error
	GetArticles() ([]*models.Article, error)
}

// MySQLStore implements the Store interface using a MySQL database.
type MySQLStore struct {
	db *sql.DB
}

// NewMySQLStore establishes a connection to the MySQL database and returns a MySQLStore.
func NewMySQLStore(dataSourceName string) (*MySQLStore, error) {
    // Separate the DSN into base DSN (without the database name) and the database name
    parts := strings.Split(dataSourceName, "/")
    baseDSN := parts[0] + "/"
    dbName := strings.Split(parts[1], "?")[0]

    // Connect to MySQL without specifying the database
    db, err := sql.Open("mysql", baseDSN)
    if err != nil {
        return nil, fmt.Errorf("failed to open database connection: %w", err)
    }

    // Ping the database to ensure the connection is valid
    err = db.Ping()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MySQL server: %w", err)
    }

    // Create the database if it doesn't exist
    _, err = db.Exec("CREATE DATABASE IF NOT EXISTS gophersignal")
    if err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to create database '%s': %w", dbName, err)
    }

    // Now connect to the new database
    db.Close()
    db, err = sql.Open("mysql", dataSourceName)
    if err != nil {
        return nil, fmt.Errorf("failed to open database with name '%s': %w", dbName, err)
    }

    // Ping again to check the new connection
    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to connect to database '%s': %w", dbName, err)
    }

    // Check if the articles table exists and create it if not
    var tableExists bool
    err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = 'articles'", dbName).Scan(&tableExists)
    if err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to check if articles table exists: %w", err)
    }

    if !tableExists {
        _, err = db.Exec(`CREATE TABLE IF NOT EXISTS articles (
            id INT AUTO_INCREMENT PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            link VARCHAR(512) NOT NULL,
            content TEXT,
            summary VARCHAR(2000),
            source VARCHAR(100) NOT NULL,
            is_on_homepage BOOLEAN,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        );`)
        if err != nil {
            db.Close()
            return nil, fmt.Errorf("failed to create articles table: %w", err)
        }

        // Insert a default article
        _, err = db.Exec("INSERT INTO articles (title, link, content, summary, source, is_on_homepage) VALUES (?, ?, ?, ?, ?, ?)",
            "Default Title", "http://default.link", "Default Content", "Default Summary", "Default Source", true)
        if err != nil {
            db.Close()
            return nil, fmt.Errorf("failed to insert default article: %w", err)
        }
    }

    // Return the MySQLStore instance with the new connection
    return &MySQLStore{db: db}, nil
}

func (store *MySQLStore) SaveArticles(articles []*models.Article) error {
	// Begin a transaction
	tx, err := store.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Reset is_on_homepage for all articles
	if _, err = tx.Exec("UPDATE articles SET is_on_homepage = FALSE"); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to reset is_on_homepage: %w", err)
	}

	// Insert or update articles
	for _, article := range articles {
		_, err := tx.Exec("INSERT INTO articles (title, link, content, summary, source, created_at, updated_at, is_on_homepage) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE title=VALUES(title), link=VALUES(link), content=VALUES(content), summary=VALUES(summary), source=VALUES(source), updated_at=VALUES(updated_at), is_on_homepage=VALUES(is_on_homepage);",
			article.Title, article.Link, article.Content, article.Summary, article.Source, article.CreatedAt, article.UpdatedAt, article.IsOnHomepage)
		if err != nil {
			fmt.Printf("Error saving article: %s, error: %v\n", article.Title, err)
			continue
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func (store *MySQLStore) GetArticles() ([]*models.Article, error) {
	rows, err := store.db.Query("SELECT id, title, link, content, summary, source, created_at, updated_at, is_on_homepage FROM articles WHERE is_on_homepage = TRUE;")
	if err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Link, &article.Content, &article.Summary, &article.Source, &article.CreatedAt, &article.UpdatedAt, &article.IsOnHomepage); err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, &article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}

	return articles, nil
}
