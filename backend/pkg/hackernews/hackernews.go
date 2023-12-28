package hackernews

import (
	"context"
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
			fmt.Printf("Found article: %s - %s\n", title, link) // Log found article
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
				fmt.Printf("Saved article: %s - %s\nContent: %s\n", article.Title, article.Link, content) // Log the saved article with content
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
	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to fetch content from %s. Status code: %d", url, resp.StatusCode)
	}

	// Read the entire body of the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	// Extract text from the entire body
	textContent := doc.Find("body").Text()

	// Remove excess whitespace and blank lines
	cleanedText := removeExtraWhitespace(textContent)

	return cleanedText, nil
}

// Removes extra whitespace and blank lines from a string
func removeExtraWhitespace(text string) string {
	// Split the text into lines
	lines := strings.Split(text, "\n")

	var cleanedLines []string
	for _, line := range lines {
		// Trim whitespace from each line
		trimmedLine := strings.TrimSpace(line)

		// Add the line if it is not empty
		if trimmedLine != "" {
			cleanedLines = append(cleanedLines, trimmedLine)
		}
	}

	// Join the cleaned lines back into a single string
	return strings.Join(cleanedLines, "\n")
}
