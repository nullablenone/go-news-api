package article

import (
	"context"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *Article) error
	FindAll(ctx context.Context) ([]Article, error)
	FindByID(ctx context.Context, id uint) (Article, error)
	FindBySlug(ctx context.Context, slug string) (Article, error)
	Update(ctx context.Context, article *Article) error
	Delete(ctx context.Context, id uint) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, article *Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *articleRepository) FindAll(ctx context.Context) ([]Article, error) {
	var articles []Article
	err := r.db.WithContext(ctx).Preload("Author").Find(&articles).Error
	return articles, err
}

func (r *articleRepository) FindByID(ctx context.Context, id uint) (Article, error) {
	var article Article
	err := r.db.WithContext(ctx).Preload("Author").First(&article, id).Error
	return article, err
}

func (r *articleRepository) FindBySlug(ctx context.Context, slug string) (Article, error) {
	var article Article
	err := r.db.WithContext(ctx).Preload("Author").Where("slug = ?", slug).First(&article).Error
	return article, err
}

func (r *articleRepository) Update(ctx context.Context, article *Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *articleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Article{}, id).Error
}
