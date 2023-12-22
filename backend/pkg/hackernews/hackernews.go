package hackernews

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/k-zehnder/gophersignal/backend/pkg/models"
)

const maxContentLength = 10000

type HackerNewsScraper struct{}

func (hns *HackerNewsScraper) Scrape() ([]*models.Article, error) {
	var articles []*models.Article
	c := colly.NewCollector(
		colly.AllowedDomains("news.ycombinator.com"),
	)

	c.OnHTML("tr.athing", func(e *colly.HTMLElement) {
		title := e.ChildText("td.title > span.titleline > a")
		link := e.ChildAttr("td.title > span.titleline > a", "href")

		if title != "" && link != "" {
			// Check if the URL has a supported protocol (e.g., http or https)
			if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
				content, err := fetchArticleContent(link)
				if err != nil {
					fmt.Printf("Failed to fetch content for %s: %v\n", link, err)
					return
				}
				if len(content) > maxContentLength {
					content = content[:maxContentLength] + "..." // Truncate content
				}
				article := models.NewArticle(0, title, link, content, "", "Hacker News", time.Now())
				articles = append(articles, article)
			} else {
				fmt.Printf("Skipping unsupported protocol for URL: %s\n", link)
			}
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

func fetchArticleContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to fetch content from %s. Status code: %d", url, resp.StatusCode)
	}

	// Read the body of the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	// Extract and concatenate text from certain elements
	var textContent strings.Builder
	doc.Find("article, p, h1, h2, h3").Each(func(i int, s *goquery.Selection) {
		textContent.WriteString(s.Text())
		textContent.WriteString("\n\n") // Add some spacing between elements
	})

	return textContent.String(), nil
}
