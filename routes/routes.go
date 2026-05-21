package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/domain/article"
)

func SetRoutes(article *article.ArticleHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/article", article.CreateArticle)
	router.GET("/articles", article.GetArticles)
	router.GET("/article/:id", article.GetArticle)
	router.PUT("/article/:id", article.UpdateArticle)
	router.DELETE("/article/:id", article.DeleteArticle)

	return router
}
