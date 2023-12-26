package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// const openAIURL = "https://api.openai.com/v1/engines/text-davinci-003/completions"

const openAIURL = "https://api.openai.com/v1/engines/text-curie-001/completions"

// const openAIURL = "https://api.openai.com/v1/engines/text-babbage-001/completions"

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
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("OPEN_AI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPEN_AI_API_KEY not set in .env file")
	}

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN not set in .env file")
	}

	// Connect to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Fetch articles from the database
	rows, err := db.Query("SELECT id, content FROM articles WHERE summary IS NULL")
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
			Prompt:    fmt.Sprintf("Summarize the following text in about 50 words: %s", content),
			MaxTokens: 100, // Adjustable
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

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var apiResp OpenAIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			log.Fatal(err)
		}

		// Update the article with the summary
		if len(apiResp.Choices) > 0 {
			summary := apiResp.Choices[0].Text
			_, err := db.Exec("UPDATE articles SET summary = ? WHERE id = ?", summary, id)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Article ID %d summarized\n", id)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
