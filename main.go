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
		log.Fatal(err)
	}

	db, err := config.ConnectPostgre(env)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(article.Article{}, user.User{}); err != nil {
		log.Fatal(err)
	}

	// Article Wiring
	articleRepo := article.NewArticleRepository(db)
	articleService := article.NewArticleService(articleRepo)
	articleHandler := article.NewArticleHandler(articleService)

	// User Wiring
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo, env)
	userHandler := user.NewUserHandler(userService)

	router := routes.SetRoutes(articleHandler, userHandler)
	if err = router.Run(":3333"); err != nil {
		log.Fatalf("failed to start HTTP server on :3333: %v", err)
	}

}
