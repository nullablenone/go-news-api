package article

type ArticleService interface {
	CreateArticle(req CreateArticleRequest) (Article, error)
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
