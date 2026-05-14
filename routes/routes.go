package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/domain/article"
)

func SetRoutes(article *article.ArticleHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/article", article.CreateArticle)

	return router
}
