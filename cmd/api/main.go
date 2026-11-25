package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/internal/domain/article"
	"github.com/nullablenone/go-news-api/internal/domain/user"
	"github.com/nullablenone/go-news-api/internal/routes"
	"github.com/nullablenone/go-news-api/internal/shared/config"
)

func setGinMode(appEnv string) {
	switch strings.ToLower(appEnv) {
	case "production", "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Inisialisasi Config Gagal: %v", err)
	}

	setGinMode(cfg.App.Env)

	db, err := config.ConnectPostgre(cfg.DB)
	if err != nil {
		log.Fatalf("Inisialisasi Postgre Gagal: %v", err)
	}

	rdb, err := config.ConnectRedis(cfg.Redis)
	if err != nil {
		log.Printf("Peringatan: Inisialisasi Redis gagal (%v). Aplikasi berjalan tanpa cache Redis.", err)
	}

	if err = db.AutoMigrate(article.Article{}, user.User{}, user.RefreshToken{}); err != nil {
		log.Fatalf("Gagal membuat migrasi: %v", err)
	}

	// Seeder
	user.RunAdminSeeder(db)
	// article.RunArticleSeeder(db)

	// Article Wiring
	articleRepo := article.NewArticleRepository(db)
	articleService := article.NewArticleService(articleRepo, rdb)
	articleHandler := article.NewArticleHandler(articleService)

	// User Wiring
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo, cfg)
	userHandler := user.NewUserHandler(userService)

	router := routes.SetRoutes(cfg, articleHandler, userHandler)
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("Menjalankan %s (%s) pada %s", cfg.App.Name, cfg.App.Env, addr)
	if err = router.Run(addr); err != nil {
		log.Fatalf("failed to start HTTP server on %s: %v", addr, err)
	}

}
