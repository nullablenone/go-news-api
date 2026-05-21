package article

type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Slug    string `json:"slug" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
}

type UpdateArticleRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Slug    string `json:"slug" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
}
