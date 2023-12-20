package hackernews

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/k-zehnder/gophersignal/backend/pkg/models"
)

type HackerNewsScraper struct{}

func (hns *HackerNewsScraper) Scrape() ([]*models.Article, error) {
	var articles []*models.Article
	c := colly.NewCollector(
		colly.AllowedDomains("news.ycombinator.com"),
	)

	// Use the 'span.titleline' to be more specific and avoid selecting other 'td.title' elements.
	c.OnHTML("tr.athing", func(e *colly.HTMLElement) {
		title := e.ChildText("td.title > span.titleline > a")
		link := e.ChildAttr("td.title > span.titleline > a", "href")

		if title != "" && link != "" {
			article := models.NewArticle(title, link, "Hacker News")
			articles = append(articles, article)
		} else {
			fmt.Printf("No title or link found for element: %v\n", e)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit("https://news.ycombinator.com/")
	if err != nil {
		return nil, err
	}

	c.Wait()

	return articles, nil
}
