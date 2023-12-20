package scraper

import "github.com/k-zehnder/gophersignal/backend/pkg/models"

// Scraper defines the interface for a scraper
type Scraper interface {
	Scrape() (*[]models.Article, error)
}
