package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/config"
	"github.com/nullablenone/go-news-api/internal/domain/article"
	"github.com/nullablenone/go-news-api/internal/domain/user"
	"github.com/nullablenone/go-news-api/internal/middleware"
)

func SetRoutes(env *config.Env, article *article.ArticleHandler, user *user.UserHandler) *gin.Engine {
	router := gin.Default()

	// --- SETUP CORS ---
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Mengizinkan semua domain (termasuk localhost frontend)
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	// Authorization agar pengiriman JWT tidak diblokir
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	router.Use(cors.New(corsConfig))

	// Public
	router.POST("/register", user.Register)
	router.POST("/login", user.Login)

	publicArticleRoutes := router.Group("/public")
	publicArticleRoutes.GET("/articles", article.ListArticles)
	publicArticleRoutes.GET("/articles/:slug", article.GetArticleBySlug)

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
