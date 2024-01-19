package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/k-zehnder/gophersignal/internal/models"
	"github.com/k-zehnder/gophersignal/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestArticleController_GetAll(t *testing.T) {
	// Create a mock store for testing
	mockStore := &store.MockStore{
		GetAllArticlesFunc: func() []models.Article {
			return []models.Article{
				{ID: 1, Title: "Article 1"},
				{ID: 2, Title: "Article 2"},
			}
		},
	}

	// Create an instance of the controller
	controller := NewController(mockStore)

	// Create a mock Gin context for the test
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/articles", controller.GetAll)
	req, _ := http.NewRequest("GET", "/articles", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check the response status code and content
	assert.Equal(t, http.StatusOK, resp.Code)
	var response models.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "success", response.Status)
	assert.NotNil(t, response.Data)
}

func TestArticleController_Create(t *testing.T) {
	// Create a mock store for testing
	mockStore := &store.MockStore{
		SaveArticleFunc: func(article models.Article) error {
			return nil
		},
	}

	// Create an instance of the controller
	controller := NewController(mockStore)

	// Create a mock Gin context for the test with JSON request
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/articles", controller.Create)
	reqBody := `{"title": "Test Article", "content": "Test Content"}`
	req, _ := http.NewRequest("POST", "/articles", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check the response status code and content
	assert.Equal(t, http.StatusCreated, resp.Code)
	var response models.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Article created successfully", response.Data)
}
