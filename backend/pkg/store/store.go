package store

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/pkg/models"
)

type DBStore struct {
	db *sql.DB
}

func NewDBStore(dataSourceName string) *DBStore {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &DBStore{db: db}
}

func (store *DBStore) SaveArticles(articles []*models.Article) {
	for _, article := range articles {
		_, err := store.db.Exec("INSERT INTO articles (title, link, content, source) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE title=VALUES(title), link=VALUES(link), content=VALUES(content), source=VALUES(source)",
			article.Title, article.Link, article.Content, article.Source)

		if err != nil {
			log.Printf("Failed to save article: %s\nError: %v\n", article.Title, err)
		}
	}
}

func (store *DBStore) GetAllArticles() ([]*models.Article, error) {
	rows, err := store.db.Query("SELECT title, link, content, summary, source, scraped_at FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.Title, &article.Link, &article.Content, &article.Summary, &article.Source, &article.ScrapedAt); err != nil {
			return nil, err
		}
		articles = append(articles, &article)
	}
	return articles, nil
}
