package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/k-zehnder/gophersignal/backend/pkg/models"
)

const maxContentLength = 10000

func main() {
	fmt.Println("Scraping Started")
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://news.ycombinator.com/"},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			hackerNewsParse(g, r)
		},
		UserAgent: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	}).Start()
	fmt.Println("Scraping Finished")
}

func hackerNewsParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("tr.athing").Each(func(i int, s *goquery.Selection) {
		title := s.Find("td.title > a.storylink").Text()
		link, exists := s.Find("td.title > a.storylink").Attr("href")

		if !exists {
			fmt.Println("No link found for article")
			return
		}

		content, err := fetchArticleContent(link)
		if err != nil {
			fmt.Printf("Failed to fetch content for %s: %v\n", link, err)
			return
		}
		if len(content) > maxContentLength {
			content = content[:maxContentLength] + "..."
		}

		article := models.NewArticle(0, title, link, content, "", "Hacker News", time.Now())
		fmt.Printf("Saved article: %s - %s\n", article.Title, article.Link)
	})
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
		return "", fmt.Errorf("Failed to fetch content from %s. Status code: %d", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	var textContent strings.Builder
	doc.Find("article, p, h1, h2, h3").Each(func(i int, s *goquery.Selection) {
		textContent.WriteString(s.Text())
		textContent.WriteString("\n\n")
	})

	return textContent.String(), nil
}
