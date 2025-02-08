// Package store defines an interface and an implementation for article storage and retrieval.
package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

// Store interface defines methods for article storage and retrieval.
type Store interface {
	SaveArticles(articles []*models.Article) error
	GetArticles(limit, offset int) ([]*models.Article, error)
	GetFilteredArticles(flagged, dead, dupe *bool, limit, offset int) ([]*models.Article, error)
}

// MySQLStore implements the Store interface using a MySQL database.
type MySQLStore struct {
	db *sql.DB
}

// NewMySQLStore establishes a new MySQL database connection.
func NewMySQLStore(dataSourceName string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &MySQLStore{db: db}, nil
}

// SaveArticles inserts articles into the database.
// Note: The "ON DUPLICATE KEY UPDATE" clause has been removed.
func (store *MySQLStore) SaveArticles(articles []*models.Article) error {
	stmt, err := store.db.Prepare(`
        INSERT INTO articles (
          title,
          link,
          content,
          summary,
          source,
          upvotes,
          comment_count,
          comment_link,
          flagged,
          dead,
          dupe,
          created_at,
          updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, article := range articles {
		_, execErr := stmt.Exec(
			article.Title,
			article.Link,
			article.Content,
			article.Summary,
			article.Source,
			article.Upvotes,
			article.CommentCount,
			article.CommentLink,
			article.Flagged,
			article.Dead,
			article.Dupe,
			article.CreatedAt,
			article.UpdatedAt,
		)
		if execErr != nil {
			fmt.Printf("Failed for article '%s': %v\n", article.Title, execErr)
			continue
		}
	}
	return nil
}

// GetArticles retrieves the latest articles (by id descending) that have a non-empty summary.
func (store *MySQLStore) GetArticles(limit, offset int) ([]*models.Article, error) {
	query := `
		SELECT id, title, link, content, summary, source,
		       upvotes, comment_count, comment_link, flagged,
		       dead, dupe, created_at, updated_at
		FROM articles
		WHERE summary IS NOT NULL AND TRIM(summary) != ''
		ORDER BY id DESC
		LIMIT ? OFFSET ?;
	`

	rows, err := store.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Link,
			&article.Content,
			&article.Summary,
			&article.Source,
			&article.Upvotes,
			&article.CommentCount,
			&article.CommentLink,
			&article.Flagged,
			&article.Dead,
			&article.Dupe,
			&article.CreatedAt,
			&article.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, &article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}
	return articles, nil
}

// GetFilteredArticles retrieves the latest filtered articles (by id descending)
// and applies optional filters for flagged, dead, and dupe statuses.
func (store *MySQLStore) GetFilteredArticles(flagged, dead, dupe *bool, limit, offset int) ([]*models.Article, error) {
	// Base query.
	query := `
		SELECT id, title, link, content, summary, source,
		       upvotes, comment_count, comment_link, flagged,
		       dead, dupe, created_at, updated_at
		FROM articles
		WHERE 1
	`
	var args []interface{}

	// Build filtering conditions.
	if flagged != nil {
		var flaggedVal int
		if *flagged {
			flaggedVal = 1
		} else {
			flaggedVal = 0
		}
		query += " AND flagged = ?"
		args = append(args, flaggedVal)
	}
	if dead != nil {
		var deadVal int
		if *dead {
			deadVal = 1
		} else {
			deadVal = 0
		}
		query += " AND dead = ?"
		args = append(args, deadVal)
	}
	if dupe != nil {
		var dupeVal int
		if *dupe {
			dupeVal = 1
		} else {
			dupeVal = 0
		}
		query += " AND dupe = ?"
		args = append(args, dupeVal)
	}

	// Append ordering and pagination.
	query += " ORDER BY id DESC LIMIT ? OFFSET ?;"
	args = append(args, limit, offset)

	// Execute the query.
	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute filtered query: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Link,
			&article.Content,
			&article.Summary,
			&article.Source,
			&article.Upvotes,
			&article.CommentCount,
			&article.CommentLink,
			&article.Flagged,
			&article.Dead,
			&article.Dupe,
			&article.CreatedAt,
			&article.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, &article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}
	return articles, nil
}
