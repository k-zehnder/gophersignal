package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/k-zehnder/gophersignal/backend/internal/models"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

const maxContentLength = 10000

// Scraper defines the interface for a scraper
type Scraper interface {
	Scrape() ([]*models.Article, error)
}

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
			fmt.Printf("Found article: %s - %s\n", title, link)
			if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
				content, err := fetchArticleContent(link)
				if err != nil {
					fmt.Printf("Failed to fetch content for %s: %v\n", link, err)
					return
				}
				if len(content) > maxContentLength {
					content = content[:maxContentLength] + "..."
				}
				now := time.Now()
				article := models.NewArticle(0, title, link, content, "", "Hacker News", now, now, true)
				articles = append(articles, article)
				fmt.Printf("Saved article: %s - %s\nContent: %s\n", article.Title, article.Link, content)
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
		return "", fmt.Errorf("failed to fetch content from %s. Status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	utf8Body, err := detectAndConvertToUTF8(body)
	if err != nil {
		return "", fmt.Errorf("failed to convert content to UTF-8 for %s: %v", url, err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(utf8Body))
	if err != nil {
		return "", err
	}

	textContent := doc.Find("body").Text()
	cleanedText := removeExtraWhitespace(textContent)

	return cleanedText, nil
}

func detectAndConvertToUTF8(content []byte) (string, error) {
	r := bytes.NewReader(content)
	e, _, _ := charset.DetermineEncoding(content, "")
	utf8Reader := transform.NewReader(r, e.NewDecoder())

	transformedContent, err := io.ReadAll(utf8Reader)
	if err != nil {
		return "", err
	}
	return string(transformedContent), nil
}

func removeExtraWhitespace(text string) string {
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			cleanedLines = append(cleanedLines, trimmedLine)
		}
	}
	return strings.Join(cleanedLines, "\n")
}
