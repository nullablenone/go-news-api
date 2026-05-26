package article

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type ArticleService interface {
	CreateArticle(authorID uint, req CreateArticleRequest) (Article, error)
	GetAllArticles() ([]Article, error)
	GetArticleByID(id uint) (Article, error)
	GetArticleBySlug(slug string) (Article, error)
	UpdateArticle(id uint, req UpdateArticleRequest) (Article, error)
	DeleteArticle(id uint) error
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

func (s *articleService) CreateArticle(authorID uint, req CreateArticleRequest) (Article, error) {
	article := Article{
		Title:    req.Title,
		Slug:     req.Slug,
		Summary:  req.Summary,
		Category: req.Category,
		ReadTime: req.ReadTime,
		AuthorID: authorID,
		Content:  req.Content,
	}

	err := s.repo.Create(&article)
	if err != nil {
		return Article{}, err
	}

	res, err := s.repo.FindByID(article.ID)
	if err == nil {
		s.clearArticleCache("")
	}

	return res, err
}

func (s *articleService) GetAllArticles() ([]Article, error) {
	ctx := context.Background()
	var articles []Article

	// Panggilan cache menjadi sangat ringkas dan bersih
	if s.getCache(ctx, cacheKeyAllArticles, &articles) {
		log.Println("====== [REDIS] Cache Hit: Mengambil semua artikel dari Redis ======")
		return articles, nil
	}

	log.Println("====== [DB] Cache Miss: Mengambil semua artikel dari PostgreSQL ======")
	articles, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKeyAllArticles, articles)
	return articles, nil
}

func (s *articleService) GetArticleBySlug(slug string) (Article, error) {
	ctx := context.Background()
	cacheKey := cacheKeyArticlePrefix + slug
	var article Article

	if s.getCache(ctx, cacheKey, &article) {
		log.Printf("====== [REDIS] Cache Hit: Mengambil artikel slug [%s] dari Redis ======", slug)
		return article, nil
	}

	log.Printf("====== [DB] Cache Miss: Mengambil artikel slug [%s] dari PostgreSQL ======", slug)
	article, err := s.repo.FindBySlug(slug)
	if err != nil {
		return Article{}, err
	}

	s.setCache(ctx, cacheKey, article)
	return article, nil
}

func (s *articleService) GetArticleByID(id uint) (Article, error) {
	return s.repo.FindByID(id)
}

func (s *articleService) UpdateArticle(id uint, req UpdateArticleRequest) (Article, error) {
	article, err := s.repo.FindByID(id)
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

	err = s.repo.Update(&article)
	if err != nil {
		return Article{}, err
	}

	s.clearArticleCache(oldSlug)
	if oldSlug != req.Slug {
		s.clearArticleCache(req.Slug)
	}

	return article, nil
}

func (s *articleService) DeleteArticle(id uint) error {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	s.clearArticleCache(article.Slug)
	return nil
}
