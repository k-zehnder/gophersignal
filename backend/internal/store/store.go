// Package store defines an interface and an implementation for article storage and retrieval.
package store

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
)

// Store defines methods for article storage and retrieval.
type Store interface {
	SaveArticles(articles []*models.Article) error
	GetArticles(limit, offset int) ([]*models.Article, error)
	GetFilteredArticles(flagged, dead, dupe *bool, limit, offset int) ([]*models.Article, error)
	GetArticlesWithThresholds(limit, offset, minUpvotes, minComments int) ([]*models.Article, error)
	GetArticlesWithThresholdsAndFilters(limit, offset, minUpvotes, minComments int, flagged, dead, dupe *bool) ([]*models.Article, error)
}

// MySQLStore implements Store using a MySQL database.
type MySQLStore struct {
	db *sql.DB
}

// NewMySQLStore creates a new MySQLStore.
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
          article_rank,
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
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, article := range articles {
		_, execErr := stmt.Exec(
			article.Title,
			article.Link,
			article.ArticleRank,
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

// GetArticles retrieves deduplicated articles.
func (store *MySQLStore) GetArticles(limit, offset int) ([]*models.Article, error) {
	query := `
		SELECT a.id, a.title, a.link, a.article_rank, a.content, a.summary, a.source,
		       a.upvotes, a.comment_count, a.comment_link, a.flagged,
		       a.dead, a.dupe, a.created_at, a.updated_at
		FROM articles a
		INNER JOIN (
			SELECT title, MAX(id) AS max_id
			FROM articles
			WHERE summary IS NOT NULL 
			  AND TRIM(summary) != ''
			  AND summary NOT LIKE 'No summary available%'
			  AND flagged = FALSE
			  AND dead = FALSE
			  AND dupe = FALSE
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
			&article.ArticleRank,
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

// GetFilteredArticles retrieves articles with optional filters.
func (store *MySQLStore) GetFilteredArticles(flagged, dead, dupe *bool, limit, offset int) ([]*models.Article, error) {
	innerQuery := `
		SELECT title, MAX(id) AS max_id
		FROM articles
		WHERE 1=1
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

	query := `
		SELECT a.id, a.title, a.link, a.article_rank, a.content, a.summary, a.source,
		       a.upvotes, a.comment_count, a.comment_link, a.flagged,
		       a.dead, a.dupe, a.created_at, a.updated_at
		FROM articles a
		INNER JOIN (
	` + innerQuery + `
		) b ON a.title = b.title AND a.id = b.max_id
		ORDER BY a.id DESC
		LIMIT ? OFFSET ?;
	`
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
			&article.ArticleRank,
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

// GetArticlesWithThresholds retrieves articles using provided minimum upvote and comment thresholds.
func (store *MySQLStore) GetArticlesWithThresholds(limit, offset, minUpvotes, minComments int) ([]*models.Article, error) {
	query := `
		SELECT a.id, a.title, a.link, a.article_rank, a.content, a.summary, a.source,
		       a.upvotes, a.comment_count, a.comment_link, a.flagged,
		       a.dead, a.dupe, a.created_at, a.updated_at
		FROM articles a
		INNER JOIN (
			SELECT title, MAX(id) AS max_id
			FROM articles
			WHERE summary IS NOT NULL 
			  AND TRIM(summary) != ''
			  AND summary NOT LIKE 'No summary available%'
			  AND flagged = FALSE
			  AND dead = FALSE
			  AND dupe = FALSE
			  AND upvotes >= ?
			  AND comment_count >= ?
			GROUP BY title
		) b ON a.title = b.title AND a.id = b.max_id
		ORDER BY a.id DESC
		LIMIT ? OFFSET ?;
	`
	// If minUpvotes or minComments is 0, the query condition ">= 0" includes all articles.
	rows, err := store.db.Query(query, minUpvotes, minComments, limit, offset)
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
			&article.ArticleRank,
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

// GetArticlesWithThresholdsAndFilters retrieves articles that satisfy both threshold
// conditions (minUpvotes and minComments) and additional boolean filters (flagged, dead, dupe).
func (store *MySQLStore) GetArticlesWithThresholdsAndFilters(limit, offset, minUpvotes, minComments int, flagged, dead, dupe *bool) ([]*models.Article, error) {
	innerQuery := `
		SELECT title, MAX(id) AS max_id
		FROM articles
		WHERE summary IS NOT NULL 
		  AND TRIM(summary) != ''
		  AND summary NOT LIKE 'No summary available%'
	`
	// Append boolean filter conditions.
	if flagged != nil {
		innerQuery += " AND flagged = ?"
	} else {
		innerQuery += " AND flagged = FALSE"
	}
	if dead != nil {
		innerQuery += " AND dead = ?"
	} else {
		innerQuery += " AND dead = FALSE"
	}
	if dupe != nil {
		innerQuery += " AND dupe = ?"
	} else {
		innerQuery += " AND dupe = FALSE"
	}

	// Append threshold conditions.
	innerQuery += " AND upvotes >= ? AND comment_count >= ? GROUP BY title"

	// Build the outer query.
	fullQuery := `
		SELECT a.id, a.title, a.link, a.article_rank, a.content, a.summary, a.source,
		       a.upvotes, a.comment_count, a.comment_link, a.flagged,
		       a.dead, a.dupe, a.created_at, a.updated_at
		FROM articles a
		INNER JOIN (
			` + innerQuery + `
		) b ON a.title = b.title AND a.id = b.max_id
		ORDER BY a.id DESC
		LIMIT ? OFFSET ?;
	`

	// Build the arguments list.
	var args []interface{}
	if flagged != nil {
		args = append(args, boolToInt(*flagged))
	} else {
		args = append(args, 0)
	}
	if dead != nil {
		args = append(args, boolToInt(*dead))
	} else {
		args = append(args, 0)
	}
	if dupe != nil {
		args = append(args, boolToInt(*dupe))
	} else {
		args = append(args, 0)
	}
	args = append(args, minUpvotes, minComments, limit, offset)

	rows, err := store.db.Query(fullQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute combined query: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Link,
			&article.ArticleRank,
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

// Convert boolean to integer (1 or 0).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
