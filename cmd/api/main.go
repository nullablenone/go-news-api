package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nullablenone/go-news-api/internal/domain/article"
	"github.com/nullablenone/go-news-api/internal/domain/user"
	"github.com/nullablenone/go-news-api/internal/routes"
	"github.com/nullablenone/go-news-api/internal/shared/config"
	"github.com/nullablenone/go-news-api/internal/shared/database"
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
		log.Fatalf("failed to load config: %v", err)
	}

	setGinMode(cfg.App.Env)

	db, err := database.ConnectPostgre(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	rdb, err := database.ConnectRedis(cfg.Redis)
	if err != nil {
		log.Printf("warning: failed to connect to Redis (%v). Running without Redis cache", err)
	}

	if err = db.AutoMigrate(article.Article{}, user.User{}, user.RefreshToken{}); err != nil {
		log.Fatalf("failed to run database migration: %v", err)
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

	router := routes.SetRoutes(cfg, db, rdb, articleHandler, userHandler)
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("starting %s (%s) on %s", cfg.App.Name, cfg.App.Env, addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start HTTP server on %s: %v", addr, err)
		}
	}()

	<-sigCtx.Done()
	log.Println("shutting down server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("warning: forced HTTP server shutdown: %v", err)
	}

	if rdb != nil {
		if err := rdb.Close(); err != nil {
			log.Printf("warning: failed to close Redis connection: %v", err)
		}
	}

	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("warning: failed to close PostgreSQL connection: %v", err)
		}
	}

	log.Println("server exited cleanly")
}
