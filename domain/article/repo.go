package article

import "gorm.io/gorm"

type ArticleRepository interface {
	Create(article *Article) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(article *Article) error {
	return r.db.Create(article).Error
}
