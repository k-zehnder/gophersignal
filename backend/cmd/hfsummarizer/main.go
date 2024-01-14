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

	_ "github.com/go-sql-driver/mysql"
	"github.com/k-zehnder/gophersignal/backend/config"
)

const huggingFaceURL = "https://api-inference.huggingface.co/models/facebook/bart-large-cnn"

type HuggingFaceRequest struct {
	Inputs string `json:"inputs"`
}

type HuggingFaceResponseItem struct {
	SummaryText string `json:"summary_text"`
}

func main() {
	dsn := config.GetEnv("SCRAPER_MYSQL_DSN", "") 
	if dsn == "" {
		log.Fatal("SCRAPER_MYSQL_DSN not set in .env file")
	}

	apiKey := config.GetEnv("HUGGING_FACE_API_KEY", "")
	if apiKey == "" {
		log.Fatal("HUGGING_FACE_API_KEY not set in .env file")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, content FROM articles WHERE (summary IS NULL OR summary = '') AND is_on_homepage = TRUE;")
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

		if content == "" {
			log.Printf("Skipping Article ID %d: content is empty", id)
			continue
		}

		reqBody := HuggingFaceRequest{
			Inputs: content,
		}
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			log.Printf("Error marshaling request body for Article ID %d: %v", id, err)
			continue
		}

		req, err := http.NewRequest("POST", huggingFaceURL, bytes.NewBuffer(reqBytes))
		if err != nil {
			log.Printf("Error creating HTTP request for Article ID %d: %v", id, err)
			continue
		}

		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		var resp *http.Response
		var retryCount int

		for retryCount < 3 { // Max 3 retries for a failed request
			resp, err = client.Do(req)
			if err != nil {
				log.Printf("Error sending request for Article ID %d (Retry %d/3): %v", id, retryCount+1, err)
				retryCount++
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}

		if resp == nil {
			log.Printf("Max retries reached for Article ID %d. Skipping...", id)
			continue
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body for Article ID %d: %v", id, err)
			continue
		}

		if !isValidJSONArray(body) {
			log.Printf("Invalid JSON response for Article ID %d: %s", id, string(body))
			continue
		}

		if isServerWarmingUp(body) {
			estimatedTime := getEstimatedTime(body)
			log.Printf("Server is warming up for Article ID %d. Retrying in %.2f seconds...", id, estimatedTime)
			time.Sleep(time.Duration(estimatedTime) * time.Second)
			continue
		}

		var apiResps []HuggingFaceResponseItem
		if err := json.Unmarshal(body, &apiResps); err != nil {
			log.Printf("Error parsing response for Article ID %d: %v", id, err)
			log.Printf("Response body for Article ID %d: %s", id, string(body))
			continue
		}

		if len(apiResps) > 0 {
			summary := apiResps[0].SummaryText
			_, err := db.Exec("UPDATE articles SET summary = ? WHERE id = ?", summary, id)
			if err != nil {
				log.Printf("Error updating database for Article ID %d: %v", id, err)
				continue
			}
			fmt.Printf("Article ID %d summarized\n", id)
		} else {
			fmt.Printf("No summary received for Article ID %d\n", id)
		}
	}
}

func isValidJSONArray(body []byte) bool {
	var jsonArray []interface{}
	if err := json.Unmarshal(body, &jsonArray); err != nil {
		return false
	}
	return true
}

func getEstimatedTime(body []byte) float64 {
	var response struct {
		Error         string  `json:"error"`
		EstimatedTime float64 `json:"estimated_time"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return 0
	}
	return response.EstimatedTime
}

func isServerWarmingUp(body []byte) bool {
	var response struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &response); err == nil && response.Error == "Model facebook/bart-large-cnn is currently loading" {
		return true
	}
	return false
}
