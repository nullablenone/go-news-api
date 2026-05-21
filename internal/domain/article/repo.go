package article

import "gorm.io/gorm"

type ArticleRepository interface {
	Create(article *Article) error
	FindAll() ([]Article, error)
	FindByID(id uint) (Article, error)
	Update(article *Article) error
	Delete(id uint) error
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

func (r *articleRepository) FindAll() ([]Article, error) {
	var articles []Article
	err := r.db.Find(&articles).Error
	return articles, err
}

func (r *articleRepository) FindByID(id uint) (Article, error) {
	var article Article
	err := r.db.First(&article, id).Error
	return article, err
}

func (r *articleRepository) Update(article *Article) error {
	return r.db.Save(article).Error
}

func (r *articleRepository) Delete(id uint) error {
	// GORM otomatis melakukan Soft Delete karena ada gorm.DeletedAt di model
	return r.db.Delete(&Article{}, id).Error
}
