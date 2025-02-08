// Package store defines an interface and an implementation for article storage and retrieval.
package store

import (
	"database/sql"
	"fmt"
	"strings"

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

// GetArticles retrieves the latest articles (by id descending) that have a non-empty summary,
// and deduplicates them by title (only the article with the highest id for each title is returned).
func (store *MySQLStore) GetArticles(limit, offset int) ([]*models.Article, error) {
	query := `
		SELECT a.id, a.title, a.link, a.content, a.summary, a.source,
		       a.upvotes, a.comment_count, a.comment_link, a.flagged,
		       a.dead, a.dupe, a.created_at, a.updated_at
		FROM articles a
		INNER JOIN (
			SELECT title, MAX(id) AS max_id
			FROM articles
			WHERE summary IS NOT NULL AND TRIM(summary) != ''
			GROUP BY title
		) b ON a.title = b.title AND a.id = b.max_id
		ORDER BY a.id DESC
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

// GetFilteredArticles retrieves the latest filtered articles (by id descending),
// applies optional filters for flagged, dead, and dupe statuses,
// ensures that only articles with non-empty summaries are considered,
// and deduplicates by title.
func (store *MySQLStore) GetFilteredArticles(flagged, dead, dupe *bool, limit, offset int) ([]*models.Article, error) {
	// Build the inner subquery that deduplicates by title.
	innerQuery := `
		SELECT title, MAX(id) AS max_id
		FROM articles
		WHERE summary IS NOT NULL AND TRIM(summary) != ''
	`
	var conditions []string
	var args []interface{}

	if flagged != nil {
		var flaggedVal int
		if *flagged {
			flaggedVal = 1
		} else {
			flaggedVal = 0
		}
		conditions = append(conditions, "flagged = ?")
		args = append(args, flaggedVal)
	}
	if dead != nil {
		var deadVal int
		if *dead {
			deadVal = 1
		} else {
			deadVal = 0
		}
		conditions = append(conditions, "dead = ?")
		args = append(args, deadVal)
	}
	if dupe != nil {
		var dupeVal int
		if *dupe {
			dupeVal = 1
		} else {
			dupeVal = 0
		}
		conditions = append(conditions, "dupe = ?")
		args = append(args, dupeVal)
	}

	if len(conditions) > 0 {
		innerQuery += " AND " + strings.Join(conditions, " AND ")
	}
	innerQuery += " GROUP BY title"

	// Build the outer query that joins the deduplicated inner query.
	query := `
		SELECT a.id, a.title, a.link, a.content, a.summary, a.source,
		       a.upvotes, a.comment_count, a.comment_link, a.flagged,
		       a.dead, a.dupe, a.created_at, a.updated_at
		FROM articles a
		INNER JOIN (
	` + innerQuery + `
		) b ON a.title = b.title AND a.id = b.max_id
		ORDER BY a.id DESC
		LIMIT ? OFFSET ?;
	`
	// Append pagination parameters.
	args = append(args, limit, offset)

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
