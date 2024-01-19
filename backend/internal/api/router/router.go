package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/k-zehnder/gophersignal/internal/api/controller"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(articlesController controller.ArticleController) *gin.Engine {
	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Welcome to the API")
	})

	// CORS middleware setup
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "https://gophersignal.com"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(corsConfig))

	baseRouter := router.Group("/api/v1")
	articlesRouter := baseRouter.Group("/articles")

	// Setup routes for ArticleController
	articlesRouter.GET("", articlesController.GetAll)
	articlesRouter.POST("", articlesController.Create)

	return router
}
