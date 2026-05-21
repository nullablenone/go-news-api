package article

type ArticleService interface {
	CreateArticle(req CreateArticleRequest) (Article, error)
	GetAllArticles() ([]Article, error)
	GetArticleByID(id uint) (Article, error)
	UpdateArticle(id uint, req UpdateArticleRequest) (Article, error)
	DeleteArticle(id uint) error
}

type articleService struct {
	repo ArticleRepository
}

func NewArticleService(repo ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (s *articleService) CreateArticle(req CreateArticleRequest) (Article, error) {
	article := Article{
		Title:   req.Title,
		Slug:    req.Slug,
		Content: req.Content,
	}

	err := s.repo.Create(&article)
	if err != nil {
		return Article{}, err
	}

	return article, nil
}

func (s *articleService) GetAllArticles() ([]Article, error) {
	return s.repo.FindAll()
}

func (s *articleService) GetArticleByID(id uint) (Article, error) {
	return s.repo.FindByID(id)
}

func (s *articleService) UpdateArticle(id uint, req UpdateArticleRequest) (Article, error) {
	// Cek apakah artikel ada
	article, err := s.repo.FindByID(id)
	if err != nil {
		return Article{}, err
	}

	// Update data field-nya
	article.Title = req.Title
	article.Slug = req.Slug
	article.Content = req.Content

	err = s.repo.Update(&article)
	if err != nil {
		return Article{}, err
	}

	return article, nil
}

func (s *articleService) DeleteArticle(id uint) error {
	// Cek terlebih dahulu apakah artikel ada sebelum dihapus
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
