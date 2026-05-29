package main

import (
	"log"

	"github.com/nullablenone/go-news-api/config"
	"github.com/nullablenone/go-news-api/internal/domain/article"
	"github.com/nullablenone/go-news-api/internal/domain/user"
	"github.com/nullablenone/go-news-api/routes"
)

func main() {

	env, err := config.NewEnv()
	if err != nil {
		log.Fatalf("Inisialisasi Env file Gagal: %v", err)
	}

	db, err := config.ConnectPostgre(env)
	if err != nil {
		log.Fatalf("Inisialisasi Postgre Gagal: %v", err)
	}

	rdb, err := config.ConnectRedis(env)
	if err != nil {
		log.Fatalf("Inisialisasi Redis Gagal: %v", err)
	}

	if err = db.AutoMigrate(article.Article{}, user.User{}); err != nil {
		log.Fatalf("Gagal membuat migrasi: %v", err)
	}

	// Seeder
	user.RunAdminSeeder(db)

	// Article Wiring
	articleRepo := article.NewArticleRepository(db)
	articleService := article.NewArticleService(articleRepo, rdb)
	articleHandler := article.NewArticleHandler(articleService)

	// User Wiring
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo, env)
	userHandler := user.NewUserHandler(userService)

	router := routes.SetRoutes(env, articleHandler, userHandler)
	if err = router.Run(":8888"); err != nil {
		log.Fatalf("failed to start HTTP server on :8888: %v", err)
	}

}
