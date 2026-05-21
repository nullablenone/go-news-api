package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/internal/domain/article"
	"github.com/nullablenone/go-news-api/internal/domain/user"
)

func SetRoutes(article *article.ArticleHandler, user *user.UserHandler) *gin.Engine {
	router := gin.Default()

	// Auth Endpoints
	router.POST("/register", user.Register)
	router.POST("/login", user.Login)

	// Artikel - Private (Sementara belum dipasang middleware)
	router.POST("/article", article.CreateArticle)
	router.GET("/articles", article.GetArticles)
	router.GET("/article/:id", article.GetArticle)
	router.PUT("/article/:id", article.UpdateArticle)
	router.DELETE("/article/:id", article.DeleteArticle)

	return router
}
