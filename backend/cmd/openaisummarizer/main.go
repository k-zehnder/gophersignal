// Package main for the OpenAI summarizer tool.
// It connects to a MySQL database to fetch articles and uses the OpenAI API to generate summaries.

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// Import the MySQL driver with a blank identifier to ensure its `init()` function is executed.
	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/config"
)

// Constants for OpenAI API.
const openAIURL = "https://api.openai.com/v1/engines/gpt-3.5-turbo-instruct/completions"

// OpenAIRequest defines the structure expected by the OpenAI API.
type OpenAIRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

// OpenAIResponse represents the response structure from the OpenAI API.
type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

// main is the starting point of the application where it initializes configuration
// and begins the process of summarizing articles.
func main() {
	// Initialize application configuration.
	appConfig := config.NewConfig()

	// Establish database connection.
	db, err := sql.Open("mysql", appConfig.DataSourceName)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}
	defer db.Close()

	// Retrieve articles needing summarization.
	processArticles(db, appConfig.OpenAIAPIKey)
}

// processArticles fetches articles from the database and processes them using the OpenAI API.
func processArticles(db *sql.DB, apiKey string) {
	// SQL query to select articles without summaries.
	query := "SELECT id, content FROM articles WHERE (summary IS NULL OR summary = '');"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Error querying database:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var content string
		if err := rows.Scan(&id, &content); err != nil {
			log.Fatal("Error scanning database rows:", err)
		}

		// Generate summary using OpenAI API.
		summary, err := summarizeContent(apiKey, content)
		if err != nil {
			log.Printf("Error summarizing Article ID %d: %v", id, err)
			continue
		}

		// Update database with the summary.
		updateArticleSummary(db, id, summary)
	}

	// Check for errors after processing rows.
	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating over database rows:", err)
	}
}

// summarizeContent makes a request to OpenAI to summarize the provided content.
func summarizeContent(apiKey, content string) (string, error) {
	// Prepare the request body.
	reqBody := OpenAIRequest{
		Prompt:    fmt.Sprintf("Summarize the following article:\n\n%s", content),
		MaxTokens: 75,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Create the HTTP request.
	req, err := http.NewRequest("POST", openAIURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request and handle response.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for a successful response.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Parse the response body.
	var apiResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", err
	}

	// Check if apiResp.Choices is empty before accessing its elements.
	if len(apiResp.Choices) == 0 {
		log.Printf("No summary available for the following content:\n%s\n", content)
		return "", nil
	}

	// Extract and return the summary.
	return apiResp.Choices[0].Text, nil
}

// updateArticleSummary updates the given article in the database with the provided summary.
func updateArticleSummary(db *sql.DB, id int, summary string) {
	if summary == "" {
		fmt.Printf("No summary available for Article ID %d\n", id)
		return
	}

	// Prepare the SQL statement for updating the article.
	stmt, err := db.Prepare("UPDATE articles SET summary = ? WHERE id = ?")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	// Execute the update.
	if _, err := stmt.Exec(summary, id); err != nil {
		log.Printf("Error updating article ID %d: %v", id, err)
	} else {
		fmt.Printf("Article ID %d summarized\n", id)
	}
}
