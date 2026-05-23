package article

type ArticleService interface {
	CreateArticle(authorID uint, req CreateArticleRequest) (Article, error)
	GetAllArticles() ([]Article, error)
	GetArticleByID(id uint) (Article, error)
	GetArticleBySlug(slug string) (Article, error) // Tambahan untuk Public API
	UpdateArticle(id uint, req UpdateArticleRequest) (Article, error)
	DeleteArticle(id uint) error
}

type articleService struct {
	repo ArticleRepository
}

func NewArticleService(repo ArticleRepository) ArticleService {
	return &articleService{repo: repo}
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

	// Ambil ulang data artikel agar object Relasi Author ter-load sempurna untuk respon
	return s.repo.FindByID(article.ID)
}

func (s *articleService) GetAllArticles() ([]Article, error) {
	return s.repo.FindAll()
}

func (s *articleService) GetArticleByID(id uint) (Article, error) {
	return s.repo.FindByID(id)
}

func (s *articleService) GetArticleBySlug(slug string) (Article, error) {
	return s.repo.FindBySlug(slug)
}

func (s *articleService) UpdateArticle(id uint, req UpdateArticleRequest) (Article, error) {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return Article{}, err
	}

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

	return article, nil
}

func (s *articleService) DeleteArticle(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
