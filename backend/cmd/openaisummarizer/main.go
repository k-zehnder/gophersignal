package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/config"
)

const openAIURL = "https://api.openai.com/v1/engines/text-davinci-003/completions"

// const openAIURL = "https://api.openai.com/v1/engines/text-ada-001/completions"

type OpenAIRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	dsn := config.GetEnv("SCRAPER_MYSQL_DSN", "")
	if dsn == "" {
		log.Fatal("SCRAPER_MYSQL_DSN not set in .env file")
	}

	apiKey := config.GetEnv("OPEN_AI_API_KEY", "")

	if apiKey == "" {
		log.Fatal("OPEN_AI_API_KEY not set in .env file")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, content FROM articles WHERE summary = '' AND is_on_homepage = true")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var content string
		if err := rows.Scan(&id, &content); err != nil {
			log.Fatal(err)
		}

		// Prepare the API request to summarize the content
		reqBody := OpenAIRequest{
			Prompt:    fmt.Sprintf("Summarize the following article scraped from hackernews.com:\n\n%s", content),
			MaxTokens: 75,
		}
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", openAIURL, bytes.NewBuffer(reqBytes))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var apiResp OpenAIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			log.Fatal(err)
		}

		stmt, err := db.Prepare("UPDATE articles SET summary = ? WHERE id = ?")
		if err != nil {
			log.Fatal("Error preparing statement:", err)
		}
		defer stmt.Close()

		if len(apiResp.Choices) > 0 {
			summary := apiResp.Choices[0].Text
			_, err := stmt.Exec(summary, id)
			if err != nil {
				log.Fatal("Error updating article:", err)
			}
			fmt.Printf("Article ID %d summarized\n", id)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
