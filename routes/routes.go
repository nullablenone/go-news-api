package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/config"
	"github.com/nullablenone/go-news-api/internal/domain/article"
	"github.com/nullablenone/go-news-api/internal/domain/middleware"
	"github.com/nullablenone/go-news-api/internal/domain/user"
)

func SetRoutes(env *config.Env, article *article.ArticleHandler, user *user.UserHandler) *gin.Engine {
	router := gin.Default()

	// Public
	router.POST("/register", user.Register)
	router.POST("/login", user.Login)

	// Private (role: user, admin)
	publicArticleRoutes := router.Group("/public")
	publicArticleRoutes.Use(middleware.AuthMiddleware(env.JWTSecret), middleware.RoleMiddleware("user", "admin"))
	{
		publicArticleRoutes.GET("/articles", article.ListArticles)
		publicArticleRoutes.GET("/articles/:slug", article.GetArticleBySlug)
	}

	// Private (role: admin)
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(env.JWTSecret), middleware.RoleMiddleware("admin"))
	{
		adminRoutes.POST("/article", article.CreateArticle)
		adminRoutes.GET("/articles", article.GetArticles)
		adminRoutes.GET("/article/:id", article.GetArticle)
		adminRoutes.PUT("/article/:id", article.UpdateArticle)
		adminRoutes.DELETE("/article/:id", article.DeleteArticle)
	}

	return router
}
