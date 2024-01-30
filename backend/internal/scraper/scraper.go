// Package scraper provides a Hacker News scraper for fetching articles from the website.
// It includes the HackerNewsScraper struct, Scrape method for scraping articles, and supporting functions.
package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

// maxContentLength defines the maximum length of content that can be fetched from an article.
const maxContentLength = 10000

// HackerNewsScraper is a struct for scraping articles from Hacker News.
type HackerNewsScraper struct{}

// Scrape performs the scraping of articles from Hacker News. It returns a slice of article pointers or an error.
func (hns *HackerNewsScraper) Scrape() ([]*models.Article, error) {
	var articles []*models.Article

	// Initialize a new Colly collector for scraping.
	c := colly.NewCollector(
		colly.AllowedDomains("news.ycombinator.com"),
	)

	// OnHTML defines what the scraper should do when it encounters specific HTML elements.
	c.OnHTML("tr.athing", func(e *colly.HTMLElement) {
		// Extract title and link from the HTML element.
		title := e.ChildText("td.title > span.titleline > a")
		link := e.ChildAttr("td.title > span.titleline > a", "href")

		// Ensure title and link are not empty and fetch the article content.
		if title != "" && link != "" {
			fmt.Printf("Found article: %s - %s\n", title, link)

			// Skip non-http(s) links.
			if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
				content, err := fetchArticleContent(link)
				if err != nil {
					fmt.Printf("Failed to fetch content for %s: %v\n", link, err)
					return
				}

				// Truncate content to a maximum length.
				if len(content) > maxContentLength {
					content = content[:maxContentLength] + "..."
				}

				// Create and store the article.
				now := time.Now()
				article := models.NewArticle(0, title, link, content, "", "Hacker News", now, now)
				articles = append(articles, article)
			} else {
				fmt.Printf("Skipping unsupported protocol for URL: %s\n", link)
			}
		}
	})

	// Log every request made.
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Log any errors encountered during the scraping.
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping.
	err := c.Visit("https://news.ycombinator.com/")
	if err != nil {
		return nil, err
	}

	// Reverse the order of articles before returning.
	// This is done to ensure that newer articles have higher ID values,
	// while older articles have lower ID values. Consequently, when querying
	// in descending order by ID (ORDER BY id DESC), the most recent articles
	// will be returned first.
	reverseArticles(articles)

	// Wait for scraping to be completed.
	c.Wait()
	return articles, nil
}

// fetchArticleContent fetches and processes the content of an article given its URL.
func fetchArticleContent(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create an HTTP request for the article URL.
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	// Send the HTTP request and get the response.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response status code is not OK (200).
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch content from %s. Status code: %d", url, resp.StatusCode)
	}

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Detect and convert the content to UTF-8 encoding.
	utf8Body, err := detectAndConvertToUTF8(body)
	if err != nil {
		fmt.Printf("Error converting content to UTF-8 for %s: %v\n", url, err)
		return "", nil
	}

	// Check if the converted content is valid UTF-8.
	if !utf8.ValidString(utf8Body) {
		fmt.Printf("Invalid UTF-8 content for %s\n", url)
		return "", nil
	}

	// Parse the HTML content using goquery.
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(utf8Body))
	if err != nil {
		return "", err
	}

	// Extract text content from the HTML document.
	textContent := doc.Find("body").Text()

	// Remove extra whitespace and return the cleaned text content.
	cleanedText := removeExtraWhitespace(textContent)

	return cleanedText, nil
}

// detectAndConvertToUTF8 detects the character encoding of content and converts it to UTF-8.
func detectAndConvertToUTF8(content []byte) (string, error) {
	r := bytes.NewReader(content)
	e, _, _ := charset.DetermineEncoding(content, "")
	utf8Reader := transform.NewReader(r, e.NewDecoder())

	// Read and convert the content to UTF-8 encoding.
	transformedContent, err := io.ReadAll(utf8Reader)
	if err != nil {
		return "", err
	}
	return string(transformedContent), nil
}

// removeExtraWhitespace removes extra whitespace from a text.
func removeExtraWhitespace(text string) string {
	// Split text into lines
	lines := strings.Split(text, "\n")
	var cleanedLines []string

	// Remove extra whitespace from each line
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			cleanedLines = append(cleanedLines, trimmedLine)
		}
	}

	// Join the cleaned lines and return the result
	return strings.Join(cleanedLines, "\n")
}

// reverseArticles reverses the order of articles in a slice.
func reverseArticles(articles []*models.Article) {
	// Iterate through half of the slice and swap elements to reverse the order
	for i := 0; i < len(articles)/2; i++ {
		j := len(articles) - i - 1
		articles[i], articles[j] = articles[j], articles[i]
	}
}
