package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/internal/domain/article"
	"github.com/nullablenone/go-news-api/internal/domain/user"
	"github.com/nullablenone/go-news-api/internal/shared/config"
	"github.com/nullablenone/go-news-api/internal/shared/middleware"
)

func SetRoutes(cfg *config.Config, article *article.ArticleHandler, user *user.UserHandler) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true 
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	router.Use(cors.New(corsConfig))

	// Public
	router.POST("/register", user.Register)
	router.POST("/login", user.Login)

	// Auth (Refresh Token)
	authRoutes := router.Group("/auth")
	authRoutes.Use(middleware.RefreshTokenMiddleware(cfg.JWT.Secret))
	{
		authRoutes.POST("/refresh", user.RefreshToken)
	}

	publicArticleRoutes := router.Group("/public")
	publicArticleRoutes.GET("/articles", article.ListArticles)
	publicArticleRoutes.GET("/articles/:slug", article.GetArticleBySlug)

	// Private (role: admin)
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(cfg.JWT.Secret), middleware.RoleMiddleware("admin"))
	{
		adminRoutes.POST("/article", article.CreateArticle)
		adminRoutes.GET("/articles", article.GetArticles)
		adminRoutes.GET("/article/:id", article.GetArticle)
		adminRoutes.PUT("/article/:id", article.UpdateArticle)
		adminRoutes.DELETE("/article/:id", article.DeleteArticle)
	}

	return router
}
