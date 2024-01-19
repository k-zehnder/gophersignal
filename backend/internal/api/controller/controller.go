package controller

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k-zehnder/gophersignal/internal/models"
	"github.com/k-zehnder/gophersignal/internal/store"
)

// ArticleController struct manages the HTTP layer for articles.
type ArticleController struct {
	store store.ArticleStore
}

// NewController creates a new instance of ArticleController.
func NewController(store store.ArticleStore) ArticleController {
	return ArticleController{
		store: store,
	}
}

func (controller *ArticleController) GetAll(ctx *gin.Context) {
	articles := controller.store.GetAllArticles()
	if articles == nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Code:   http.StatusOK,
		Status: "success",
		Data:   articles,
	})
}

func (controller *ArticleController) Create(ctx *gin.Context) {
	var req models.CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   err.Error(),
		})
		return
	}

	article := models.Article{
		ID:           0,
		Title:        req.Title,
		Link:         "",
		Content:      req.Content,
		Summary:      sql.NullString{String: "", Valid: false},
		Source:       "",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
		IsOnHomepage: false,
	}

	if err := controller.store.SaveArticle(article); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.Response{
		Code:   http.StatusCreated,
		Status: "success",
		Data:   "Article created successfully",
	})
}
