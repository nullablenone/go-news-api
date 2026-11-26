package article

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	cacheTTL              = 1 * time.Hour
	cacheTimeout          = 1 * time.Second
	cacheKeyAllArticles   = "articles:all"
	cacheKeyArticlePrefix = "article:slug:"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, authorID uint, req CreateArticleRequest) (Article, error)
	GetAllArticles(ctx context.Context) ([]Article, error)
	GetArticleByID(ctx context.Context, id uint) (Article, error)
	GetArticleBySlug(ctx context.Context, slug string) (Article, error)
	UpdateArticle(ctx context.Context, id uint, req UpdateArticleRequest) (Article, error)
	DeleteArticle(ctx context.Context, id uint) error
}

type articleService struct {
	repo ArticleRepository
	rdb  *redis.Client
}

func NewArticleService(repo ArticleRepository, rdb *redis.Client) ArticleService {
	return &articleService{
		repo: repo,
		rdb:  rdb,
	}
}

func (s *articleService) CreateArticle(ctx context.Context, authorID uint, req CreateArticleRequest) (Article, error) {
	article := Article{
		Title:    req.Title,
		Slug:     req.Slug,
		Summary:  req.Summary,
		Category: req.Category,
		ReadTime: req.ReadTime,
		AuthorID: authorID,
		Content:  req.Content,
	}

	err := s.repo.Create(ctx, &article)
	if err != nil {
		return Article{}, err
	}

	res, err := s.repo.FindByID(ctx, article.ID)
	if err == nil {
		s.clearArticleCache(ctx, "")
	}

	return res, err
}

func (s *articleService) GetAllArticles(ctx context.Context) ([]Article, error) {
	var articles []Article

	if s.getCache(ctx, cacheKeyAllArticles, &articles) {
		log.Println("====== [REDIS] Cache Hit: fetching all articles from Redis ======")
		return articles, nil
	}

	log.Println("====== [DB] Cache Miss: fetching all articles from PostgreSQL ======")
	articles, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKeyAllArticles, articles)
	return articles, nil
}

func (s *articleService) GetArticleBySlug(ctx context.Context, slug string) (Article, error) {
	cacheKey := cacheKeyArticlePrefix + slug
	var article Article

	if s.getCache(ctx, cacheKey, &article) {
		log.Printf("====== [REDIS] Cache Hit: fetching article slug [%s] from Redis ======", slug)
		return article, nil
	}

	log.Printf("====== [DB] Cache Miss: fetching article slug [%s] from PostgreSQL ======", slug)
	article, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return Article{}, err
	}

	s.setCache(ctx, cacheKey, article)
	return article, nil
}

func (s *articleService) GetArticleByID(ctx context.Context, id uint) (Article, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *articleService) UpdateArticle(ctx context.Context, id uint, req UpdateArticleRequest) (Article, error) {
	article, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Article{}, err
	}

	oldSlug := article.Slug

	article.Title = req.Title
	article.Slug = req.Slug
	article.Summary = req.Summary
	article.Category = req.Category
	article.ReadTime = req.ReadTime
	article.Content = req.Content

	err = s.repo.Update(ctx, &article)
	if err != nil {
		return Article{}, err
	}

	s.clearArticleCache(ctx, oldSlug)
	if oldSlug != req.Slug {
		s.clearArticleCache(ctx, req.Slug)
	}

	return article, nil
}

func (s *articleService) DeleteArticle(ctx context.Context, id uint) error {
	article, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	s.clearArticleCache(ctx, article.Slug)
	return nil
}

func (s *articleService) getCache(ctx context.Context, key string, dest interface{}) bool {
	if s.rdb == nil {
		return false
	}
	cacheCtx, cancel := context.WithTimeout(ctx, cacheTimeout)
	defer cancel()

	cachedData, err := s.rdb.Get(cacheCtx, key).Result()
	if err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(cachedData), dest); err != nil {
		return false
	}
	return true
}

func (s *articleService) setCache(ctx context.Context, key string, data interface{}) {
	if s.rdb == nil {
		return
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	cacheCtx, cancel := context.WithTimeout(ctx, cacheTimeout)
	defer cancel()

	if err := s.rdb.Set(cacheCtx, key, jsonData, cacheTTL).Err(); err != nil {
		log.Printf("warning: failed to set redis cache [%s]: %v", key, err)
	}
}

func (s *articleService) clearArticleCache(ctx context.Context, slug string) {
	if s.rdb == nil {
		return
	}
	cacheCtx, cancel := context.WithTimeout(ctx, cacheTimeout)
	defer cancel()

	keys := []string{cacheKeyAllArticles}
	if slug != "" {
		keys = append(keys, cacheKeyArticlePrefix+slug)
	}
	if err := s.rdb.Del(cacheCtx, keys...).Err(); err != nil {
		log.Printf("warning: failed to clear redis cache %v: %v", keys, err)
	}
}
