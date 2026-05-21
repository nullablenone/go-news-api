package main

import (
	"log"

	"github.com/nullablenone/go-news-api/config"
	"github.com/nullablenone/go-news-api/internal/domain/article"
	"github.com/nullablenone/go-news-api/routes"
)

func main() {
	db, err := config.ConnectPostgre()
	if err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(article.Article{}); err != nil {
		log.Fatal(err)
	}

	articleRepo := article.NewArticleRepository(db)
	articleService := article.NewArticleService(articleRepo)
	articleHandler := article.NewArticleHandler(articleService)

	router := routes.SetRoutes(articleHandler)
	if err = router.Run(":3333"); err != nil {
		log.Fatalf("failed to start HTTP server on :3333: %v", err)
	}

}
