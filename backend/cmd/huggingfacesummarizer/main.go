// Package main for the Hugging Face summarizer tool.
// It connects to a MySQL database to fetch articles and uses the Hugging Face API to generate summaries.
package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	// Import the MySQL driver with a blank identifier to ensure its `init()` function is executed.
	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/config"
)

// Constants for Hugging Face API.
const huggingFaceURL = "https://api-inference.huggingface.co/models/facebook/bart-large-cnn"

// HuggingFaceRequest defines the structure expected by the Hugging Face API.
type HuggingFaceRequest struct {
	Inputs string `json:"inputs"`
}

// HuggingFaceResponseItem represents a single response item from the Hugging Face API.
type HuggingFaceResponseItem struct {
	SummaryText string `json:"summary_text"`
}

// main is the entry point of the application. It begins by loading the application configuration,
// and then proceeds to fetch and summarize articles.
func main() {
	// Load application configuration.
	appConfig := config.NewConfig()

	// Check for required configuration.
	if appConfig.DataSourceName == "" || appConfig.HuggingFaceAPIKey == "" {
		log.Fatal("Required configuration(s) missing")
	}

	// Open database connection.
	db := openDatabaseConnection(appConfig.DataSourceName)
	defer db.Close() // Close the database connection after processing.

	// Process articles for summarization.
	processArticles(db, appConfig.HuggingFaceAPIKey)
}

// openDatabaseConnection establishes a connection to the database.
func openDatabaseConnection(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	return db
}

// processArticles processes each article for summarization.
func processArticles(db *sql.DB, apiKey string) {
	query := "SELECT id, content FROM articles WHERE (summary IS NULL OR summary = '');"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Error querying database:", err)
	}
	defer rows.Close() // Close the rows when done processing.

	// Iterate through the articles.
	for rows.Next() {
		// Declare variables for article ID and content.
		var id int
		var content string
		if err := rows.Scan(&id, &content); err != nil {
			log.Fatal("Error scanning database rows:", err)
		}

		// Check if content is empty, and if so, skip processing.
		if content == "" {
			log.Printf("Skipping Article ID %d: content is empty", id)
			continue
		}

		// Summarize the article content.
		summary, err := summarizeContent(apiKey, content)
		if err != nil {
			log.Printf("Error summarizing Article ID %d: %v", id, err)
			continue
		}

		// Update the article summary in the database.
		updateArticleSummary(db, id, summary)
	}

	// Check if there was an error during iteration over database rows.
	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating over database rows:", err)
	}
}

// summarizeContent sends content to the Hugging Face API for summarization.
func summarizeContent(apiKey, content string) (string, error) {
	// Construct the request body with the content to be summarized.
	reqBody := HuggingFaceRequest{Inputs: content}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Create a new POST request to the Hugging Face API endpoint.
	req, err := http.NewRequest("POST", huggingFaceURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}

	// Set necessary headers for authorization and content type.
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Initialize an HTTP client with a specified timeout.
	client := &http.Client{Timeout: 30 * time.Second}

	// Send the request and handle potential retries.
	resp, err := sendRequestWithRetries(client, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // Ensure the response body is closed after processing.

	// Read and parse the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the response to extract the summary text.
	return parseSummaryResponse(body)
}

// sendRequestWithRetries attempts to send an HTTP request with retries.
func sendRequestWithRetries(client *http.Client, req *http.Request) (*http.Response, error) {
	const maxRetries = 3
	const retryWaitTime = 5 * time.Second // Wait time between retries.

	// Attempt to send the request up to the maximum number of retries.
	for retryCount := 0; retryCount < maxRetries; retryCount++ {
		resp, err := client.Do(req)
		if err == nil {
			return resp, nil // Return the response if request is successful.
		}

		// Wait for a specified time before retrying the request.
		time.Sleep(retryWaitTime)
	}
	return nil, fmt.Errorf("max retries reached for request")
}

// parseSummaryResponse parses the response from the Hugging Face API.
func parseSummaryResponse(body []byte) (string, error) {
	var apiResps []HuggingFaceResponseItem
	if err := json.Unmarshal(body, &apiResps); err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	if len(apiResps) > 0 && apiResps[0].SummaryText != "" {
		// Assuming the first item in the array is the desired summary.
		return apiResps[0].SummaryText, nil
	}
	return "", nil
}

// updateArticleSummary updates the database with the provided summary for the given article ID.
func updateArticleSummary(db *sql.DB, id int, summary string) {
	if summary == "" {
		// Output a message if no summary is available for the article.
		fmt.Printf("No summary available for Article ID %d\n", id)
		return
	}

	// Prepare the SQL statement for updating the article's summary.
	stmt, err := db.Prepare("UPDATE articles SET summary = ? WHERE id = ?")
	if err != nil {
		log.Fatal("Error preparing statement for updating summary:", err)
	}
	defer stmt.Close() // Ensure the statement is closed after execution.

	// Execute the SQL statement with the summary and article ID.
	if _, err := stmt.Exec(summary, id); err != nil {
		log.Printf("Error updating database for Article ID %d: %v", id, err)
	} else {
		fmt.Printf("Article ID %d summarized\n", id)
	}
}
