package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/internal/domain/article"
	"github.com/nullablenone/go-news-api/internal/domain/user"
	"github.com/nullablenone/go-news-api/internal/shared/config"
	"github.com/nullablenone/go-news-api/internal/shared/middleware"
	"github.com/nullablenone/go-news-api/internal/shared/response"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetRoutes(cfg *config.Config, db *gorm.DB, rdb *redis.Client, article *article.ArticleHandler, user *user.UserHandler) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	router.Use(cors.New(corsConfig))

	// Docs
	registerDocsRoutes(router)

	// Root Information
	router.GET("/", func(c *gin.Context) {
		response.Success(c, http.StatusOK, cfg.App.Name+" is running", gin.H{
			"app": cfg.App.Name,
			"env": cfg.App.Env,
		})
	})

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		if db == nil {
			response.Error(c, http.StatusInternalServerError, "Database handle is unavailable", "INTERNAL_ERROR")
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Database handle is unavailable", "INTERNAL_ERROR")
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			response.Error(c, http.StatusInternalServerError, "Database ping failed: "+err.Error(), "INTERNAL_ERROR")
			return
		}

		redisStatus := "disabled"
		if rdb != nil {
			if err := rdb.Ping(ctx).Err(); err != nil {
				redisStatus = "disconnected"
			} else {
				redisStatus = "connected"
			}
		}

		response.Success(c, http.StatusOK, "System is healthy.", gin.H{
			"status":   "UP",
			"database": "connected",
			"redis":    redisStatus,
			"time":     time.Now().Format(time.RFC3339),
		})
	})

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
